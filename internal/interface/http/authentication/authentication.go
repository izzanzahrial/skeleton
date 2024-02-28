package authentication

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"os"

	"github.com/izzanzahrial/skeleton/internal/interface/http/auth0"
	"github.com/izzanzahrial/skeleton/internal/model"
	"github.com/izzanzahrial/skeleton/pkg/token"
	"github.com/jackc/pgx/v5"
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
	slog    *slog.Logger
}

func NewHandler(service authService, auth0 *auth0.Authenticator, slog *slog.Logger) *Handler {
	return &Handler{service: service, auht0: auth0, slog: slog}
}

func (h *Handler) Login(c echo.Context) error {
	var request LoginReq
	if err := c.Bind(&request); err != nil {
		h.slog.Error("fail to bind request", slog.String("error", err.Error()))
		return echo.ErrBadRequest
	}

	if err := c.Validate(&request); err != nil {
		h.slog.Error("fail to validate request", slog.String("error", err.Error()))
		return c.JSON(http.StatusBadRequest, err)
	}

	user, err := h.service.GetuserByEmailOrUsername(context.Background(), request.Email, request.Username, request.Password)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return c.JSON(http.StatusNotFound, errors.New("user not found"))
		}
		return echo.ErrInternalServerError
	}

	tkn, err := token.NewJWT(user.ID, model.Roles(user.Role))
	if err != nil {
		h.slog.Error("failed to create token", slog.String("error", err.Error()))
		return echo.ErrInternalServerError
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
		h.slog.Error("failed to exchange google oauth token", slog.String("error", err.Error()))
		return echo.ErrBadGateway
	}

	client := googleOauthConfig.Client(ctx, tkn)
	response, err := client.Get(os.Getenv("GOOGLE_OAUTH_API_URL"))
	if err != nil {
		h.slog.Error("failed to get user data from google oauth api", slog.String("error", err.Error()))
		return echo.ErrBadGateway
	}
	defer response.Body.Close()

	var userInfo map[string]interface{}
	err = json.NewDecoder(response.Body).Decode(&userInfo)
	if err != nil {
		h.slog.Error("failed to decode response body", slog.String("error", err.Error()))
		return echo.ErrInternalServerError
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
		return echo.ErrInternalServerError
	}

	jwtToken, err := token.NewJWT(newUser.ID, model.Roles(newUser.Role))
	if err != nil {
		h.slog.Error("failed to create token", slog.String("error", err.Error()))
		return echo.ErrInternalServerError
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
		h.slog.Error("failed to refresh token", slog.String("error", err.Error()))
		return echo.ErrBadGateway
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
		return echo.ErrBadGateway
	}

	token, err := h.auht0.Exchange(context.Background(), c.QueryParam("code"))
	if err != nil {
		h.slog.Error("failed to exchange token", slog.String("error", err.Error()))
		return echo.ErrBadGateway
	}

	idToken, err := h.auht0.VerifyIDToken(context.Background(), token)
	if err != nil {
		h.slog.Error("failed to verify token", slog.String("error", err.Error()))
		return echo.ErrBadGateway
	}

	var profile map[string]interface{}
	if err := idToken.Claims(&profile); err != nil {
		h.slog.Error("failed to unmarshal claims into profile", slog.String("error", err.Error()))
		return echo.ErrBadGateway
	}

	// TODO: save the user from profile to database

	// TODO: generate jwt token from user data

	// TODO: redirect to somewhere else
	return c.JSON(http.StatusOK, echo.Map{"token": token, "profile": profile})
}
