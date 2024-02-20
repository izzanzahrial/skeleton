package authentication

import (
	"context"
	"encoding/json"
	"net/http"

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
}

func NewHandler(service authService) *Handler {
	return &Handler{service: service}
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
	ClientID:     "870910343709-gb2o11j3vvp6npv32ef4c3ecphussifh.apps.googleusercontent.com",
	ClientSecret: "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",
	RedirectURL:  "http://localhost:8080/api/v1/callback",
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
	if state != "temp" {
		return echo.ErrBadGateway
	}

	code := c.FormValue("code")
	token, err := googleOauthConfig.Exchange(ctx, code)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadGateway, err)
	}

	client := googleOauthConfig.Client(ctx, token)
	response, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
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
		RefreshToken: token.RefreshToken,
	}

	newUser, err := h.service.CreateOrCheckGoogleUser(ctx, *user)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadGateway, err)
	}

	// TODO: should be redirected to somewhere else
	return c.JSON(http.StatusCreated, newUser)
}

func (h *Handler) RefreshToken(c echo.Context) error {
	// TODO: get refresh token from database
	// refreshToken := c.Get("refresh_token").(string)
	// fmt.Println("refresh_token: ", refreshToken)
	refreshToken := c.FormValue("refresh")

	newtoken := googleOauthConfig.TokenSource(context.Background(), &oauth2.Token{RefreshToken: refreshToken})
	token, err := newtoken.Token()
	if err != nil {
		return echo.NewHTTPError(http.StatusBadGateway, err)
	}
	return c.JSON(http.StatusOK, token)
}
