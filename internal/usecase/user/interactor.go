package user

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/immxrtalbeast/flagman-backend/internal/domain"
	"github.com/immxrtalbeast/flagman-backend/internal/lib"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrNotFound           = errors.New("user not found")
)

type UserInteractor struct {
	userRepo  domain.UserRepository
	tokenTTL  time.Duration
	appSecret string
}

func NewUserInteractor(userRepo domain.UserRepository, tokenTTL time.Duration, appSecret string) *UserInteractor {
	return &UserInteractor{
		userRepo:  userRepo,
		tokenTTL:  tokenTTL,
		appSecret: appSecret,
	}
}

func (ui *UserInteractor) CreateUser(ctx context.Context, fullname string, email string, phonenumber string, pass string) (uint, error) {
	const op = "uc.user.create"
	passHash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	user := domain.User{
		FullName:    fullname,
		Email:       email,
		PhoneNumber: phonenumber,
		PassHash:    passHash,
	}
	id, err := ui.userRepo.CreateUser(ctx, &user)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return id, nil
}

func (ui *UserInteractor) Login(ctx context.Context, email string, passhash string) (string, error) {
	const op = "uc.user.login"
	user, err := ui.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}
	if err := bcrypt.CompareHashAndPassword(user.PassHash, []byte(passhash)); err != nil {
		return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}
	return lib.NewToken(user, ui.tokenTTL, ui.appSecret)

}

func (ui *UserInteractor) User(ctx context.Context, id uint) (*domain.User, error) {
	const op = "uc.user.get"
	user, err := ui.userRepo.User(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
	}
	return user, nil
}
