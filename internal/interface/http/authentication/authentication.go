package authentication

import (
	"context"
	"encoding/json"
	"net/http"
	"os"

	"github.com/izzanzahrial/skeleton/internal/interface/http/auth0"
	"github.com/izzanzahrial/skeleton/internal/model"
	"github.com/izzanzahrial/skeleton/pkg/token"
	"github.com/labstack/echo/v4"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type authService interface {
	GetuserByEmailOrUsername(ctx context.Context, email, username, password string) (model.User, error)
	CreateOrCheckGoogleUser(ctx context.Context, user model.User) (model.User, error)
}

type Handler struct {
	service authService
	auht0   *auth0.Authenticator
}

func NewHandler(service authService, auth0 *auth0.Authenticator) *Handler {
	return &Handler{service: service, auht0: auth0}
}

func (h *Handler) Login(c echo.Context) error {
	ctx := c.Request().Context()

	email := c.FormValue("email")
	username := c.FormValue("username")
	password := c.FormValue("password")

	user, err := h.service.GetuserByEmailOrUsername(ctx, email, username, password)
	if err != nil {
		return c.JSON(http.StatusNotFound, err.Error())
	}

	tkn, err := token.NewJWT(user.ID, model.Roles(user.Role))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusFound, echo.Map{"user": user, "token": tkn})
}

var googleOauthConfig = &oauth2.Config{
	ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
	ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
	RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL"),
	// Learn more: https://developers.google.com/identity/protocols/oauth2/scopes
	Scopes:   []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
	Endpoint: google.Endpoint,
}

func (h *Handler) LoginGoogleOAuth(c echo.Context) error {
	// TODO: generate state
	state := "temp"
	url := googleOauthConfig.AuthCodeURL(state)
	// Add access_type=offline to get refresh token when user first tries to login
	// refrence: https://stackoverflow.com/questions/10827920/not-receiving-google-oauth-refresh-token
	return c.Redirect(http.StatusTemporaryRedirect, url+"&access_type=offline")
}

func (h *Handler) Callback(c echo.Context) error {
	ctx := c.Request().Context()

	state := c.FormValue("state")
	// TODO: get the generate state
	if state != "temp" {
		return echo.ErrBadGateway
	}

	code := c.FormValue("code")
	tkn, err := googleOauthConfig.Exchange(ctx, code)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadGateway, err)
	}

	client := googleOauthConfig.Client(ctx, tkn)
	response, err := client.Get(os.Getenv("GOOGLE_OAUTH_API_URL"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadGateway, err)
	}
	defer response.Body.Close()

	var userInfo map[string]interface{}
	err = json.NewDecoder(response.Body).Decode(&userInfo)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadGateway, err)
	}

	user := &model.User{
		Email:        userInfo["email"].(string),
		FirstName:    userInfo["given_name"].(string),
		LastName:     userInfo["name"].(string),
		PictureUrl:   userInfo["picture"].(string),
		RefreshToken: tkn.RefreshToken,
	}

	newUser, err := h.service.CreateOrCheckGoogleUser(ctx, *user)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadGateway, err)
	}

	jwtToken, err := token.NewJWT(newUser.ID, model.Roles(newUser.Role))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadGateway, err)
	}

	// TODO: should be redirected to somewhere else
	return c.JSON(http.StatusCreated, echo.Map{"user": user, "token": jwtToken})
}

func (h *Handler) RefreshToken(c echo.Context) error {
	// TODO: get refresh token from database
	refreshToken := c.FormValue("refresh")

	// TODO: add check for different origin
	// native origin

	// google origin
	newtoken := googleOauthConfig.TokenSource(context.Background(), &oauth2.Token{RefreshToken: refreshToken})
	token, err := newtoken.Token()
	if err != nil {
		return echo.NewHTTPError(http.StatusBadGateway, err)
	}

	// TODO: generate new jwt
	return c.JSON(http.StatusOK, token)
}

func (h *Handler) LoginAuth0(c echo.Context) error {
	// TODO: generate state
	return c.Redirect(http.StatusTemporaryRedirect, h.auht0.AuthCodeURL("temp"))
}

func (h *Handler) CallbackAuth0(c echo.Context) error {
	// TODO: get the generated state
	if c.QueryParam("state") != "temp" {
		return c.JSON(http.StatusBadRequest, "invalid state parameter")
	}

	token, err := h.auht0.Exchange(context.Background(), c.QueryParam("code"))
	if err != nil {
		return c.JSON(http.StatusUnauthorized, "failed to exchange an authorization code for a token")
	}

	idToken, err := h.auht0.VerifyIDToken(context.Background(), token)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "failed to verify the ID token")
	}

	var profile map[string]interface{}
	if err := idToken.Claims(&profile); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	// TODO: save the user from profile to database

	// TODO: generate jwt token from user data

	// TODO: redirect to somewhere else
	return c.JSON(http.StatusOK, echo.Map{"token": token, "profile": profile})
}
