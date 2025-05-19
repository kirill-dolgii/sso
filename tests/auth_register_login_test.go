package tests

import (
	"sso/tests/suite"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/golang-jwt/jwt/v5"
	ssov1 "github.com/kirill-dolgii/protos/gen/go/sso"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	emptyAppID     = 0
	appID          = 1
	appSecret      = "test-secret"
	passDefaultLen = 10
)

func TestRegisterLogin_Login_HappyPath(t *testing.T) {
	ctx, st := suite.New(t)

	email := gofakeit.Email()
	pass := randomPassword()

	regResponce, err := st.ApiClient.Register(ctx, &ssov1.RegisterRequest{
		Email:    email,
		Password: pass,
	})

	require.NoError(t, err)
	assert.NotEmpty(t, regResponce.GetUserId())

	loginResponse, _ := st.ApiClient.Login(ctx, &ssov1.LoginRequest{
		Email:    email,
		Password: pass,
		AppId:    appID,
	})

	loginTime := time.Now().UTC()

	tokenParsed, err := jwt.Parse(loginResponse.Token, func(token *jwt.Token) (interface{}, error) {
		return []byte(appSecret), nil
	})
	claims, ok := tokenParsed.Claims.(jwt.MapClaims)
	assert.True(t, ok)

	assert.Equal(t, regResponce.GetUserId(), int64(claims["uid"].(float64)))
	assert.Equal(t, email, claims["email"])
	assert.Equal(t, appID, int(claims["app_id"].(float64)))

	const deltaSeconds = 1
	assert.InDelta(t, loginTime.Add(st.Cfg.TokenTTL).Unix(), int(claims["exp"].(float64)), deltaSeconds)
}

func randomPassword() string {
	return gofakeit.Password(true, true, true, true, false, passDefaultLen)
}
