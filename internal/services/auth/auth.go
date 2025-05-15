package auth

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sso/internal/domain/models"
	"sso/internal/lib/jwt"
	"sso/internal/storage"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidAppId       = errors.New("invalid app id")
	ErrUserExists         = errors.New("user exists")
)

type UserSaver interface {
	SaveUser(
		ctx context.Context,
		email string,
		passHash []byte,
	) (userID int64, err error)
}

type UserProvider interface {
	User(
		ctx context.Context,
		email string,
	) (user models.User, err error)
	IsAdmin(
		ctx context.Context,
		userID int64,
	) (isAdmin bool, err error)
}

type AppProvider interface {
	App(
		ctx context.Context,
		appID int,
	) (app models.App, err error)
}

type Auth struct {
	Logger       *slog.Logger
	userSaver    UserSaver
	userProvider UserProvider
	appProvider  AppProvider
	tokenTTL     time.Duration
}

func New(
	log *slog.Logger,
	userSaver UserSaver,
	userProvider UserProvider,
	appProvider AppProvider,
	tokenTTL time.Duration,
) *Auth {
	return &Auth{
		Logger:       log,
		userSaver:    userSaver,
		userProvider: userProvider,
		appProvider:  appProvider,
		tokenTTL:     tokenTTL,
	}
}

func (a *Auth) Login(
	ctx context.Context,
	email string,
	password string,
	appId int,
) (string, error) {
	const op = "auth.Login"
	log := a.Logger.With(
		slog.String("op", op),
		slog.String("email", email),
	)

	log.Info("attempting to login user")

	user, err := a.userProvider.User(ctx, email)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			log.Warn("user not found",
				slog.Attr{Key: "error", Value: slog.StringValue(err.Error())})
			return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		}
		log.Warn("failed to get user",
			slog.Attr{Key: "error", Value: slog.StringValue(err.Error())})
		return "", fmt.Errorf("%s: %w", op, err)
	}

	if err := bcrypt.CompareHashAndPassword(user.PassHash, []byte(password)); err != nil {
		log.Warn("invalid credentials",
			slog.Attr{Key: "error", Value: slog.StringValue(err.Error())})
		return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}

	app, err := a.appProvider.App(ctx, appId)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	log.Info("user login successful")

	token, err := jwt.NewToken(user, app, a.tokenTTL)
	if err != nil {
		log.Error("failed to generate token",
			slog.Attr{Key: "error", Value: slog.StringValue(err.Error())})
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return token, nil
}

func (a *Auth) RegisterNewUser(
	ctx context.Context,
	email string,
	password string,
) (int64, error) {
	const op = "auth.RegisterNewUser"
	log := a.Logger.With(
		slog.String("op", op),
		slog.String("email", email),
	)

	log.Info("registering user")

	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("failed to generate password hash",
			slog.Attr{Key: "error", Value: slog.StringValue(err.Error())})
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	userID, err := a.userSaver.SaveUser(ctx, email, passHash)
	if err != nil {
		if errors.Is(err, storage.ErrUserExists) {
			log.Error("user already exists",
				slog.Attr{Key: "error", Value: slog.StringValue(err.Error())})
			return 0, fmt.Errorf("%s: %w", op, ErrUserExists)
		}
		log.Error("failed to save user",
			slog.Attr{Key: "error", Value: slog.StringValue(err.Error())})
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return userID, nil
}

func (a *Auth) IsAdmin(ctx context.Context, userID int64) (bool, error) {
	const op = "auth.IsAdmin"

	log := a.Logger.With(
		slog.String("op", op),
		slog.Int64("userID", userID),
	)

	log.Info("checking is admin")

	isAdmin, err := a.userProvider.IsAdmin(ctx, userID)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			log.Warn("user not found",
				slog.Attr{Key: "error", Value: slog.StringValue(err.Error())})
			return false, fmt.Errorf("%s: %w", op, ErrInvalidAppId)
		}
		log.Warn("error check is admin",
			slog.Attr{Key: "error", Value: slog.StringValue(err.Error())})
		return false, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("checked if user is admin", slog.Bool("is_admin", isAdmin))

	return isAdmin, nil
}
