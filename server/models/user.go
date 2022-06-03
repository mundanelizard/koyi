package models

import (
	"context"
	"errors"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/mundanelizard/koyi/server/config"
	"github.com/mundanelizard/koyi/server/helpers"
	"time"
)

const (
	userCollectionName        = "users"
	tokenClaimsCollectionName = "user-tokens"
)

type TokenClaims struct {
	ID           *string    `json:"id"`
	DeviceId     *string    `json:"deviceId"`
	AccessToken  *string    `json:"accessToken"`
	RefreshToken *string    `json:"refreshToken"`
	RefreshClaim *UserClaim `json:"refreshClaim"`
	AccessClaim  *UserClaim `json:"accessClaim"`
	CreatedAt    time.Time  `json:"createdAt"`
}

type User struct {
	Email       *string      `json:"email"`
	PhoneNumber *string      `json:"phoneNumber"`
	Password    *string      `json:"-" bson:"password"`
	Metadata    *interface{} `json:"metadata"`
	ID          *string      `json:"id"`
	IsDeleted   bool         `json:"deleted"`
	IsVerified  bool         `json:"isVerified"`
	CreatedAt   time.Time    `json:"createdAt"`
	UpdatedAt   time.Time    `json:"updatedAt"`
}

type UserClaim struct {
	Email       *string `json:"email"`
	PhoneNumber *string `json:"phoneNumber"`
	ID          *string `json:"id"`
	jwt.StandardClaims
}

func Count(ctx *context.Context, filter interface{}) (int64, error) {
	collection := helpers.GetCollection(config.UserDatabaseName, userCollectionName)
	count, err := collection.CountDocuments(*ctx, filter)
	return count, err
}

func (user *User) Exists(ctx *context.Context) (bool, error) {
	var count int64
	var err error

	if user.Email != nil {
		count, err = Count(ctx, map[string]string{"email": *user.Email})
	} else if user.PhoneNumber != nil {
		count, err = Count(ctx, map[string]string{"email": *user.Email})
	} else {
		return false, errors.New("empty user object")
	}

	if err != nil {
		return false, err
	}

	if count > 0 {
		return true, nil
	}

	return false, nil
}

func (user *User) FillDefaults() {
	id := uuid.New().String()
	user.ID = &id
	user.IsDeleted = false
	user.IsVerified = false
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
}

func (user *User) createHistory() {
	ctx, cancel := context.WithTimeout(context.Background(), config.BackgroundTaskTimeout)
	defer cancel()

	eh := EmailHistory{
		UserId: user.ID,
		Email:  user.Email,
	}
	eh.Create(&ctx)

	ph := PasswordHistory{
		UserId:   user.ID,
		Password: user.Email,
	}
	ph.Create(&ctx)

	pnh := PhoneNumberHistory{
		UserId:      user.ID,
		PhoneNumber: user.PhoneNumber,
	}
	pnh.Create(&ctx)
}

func (user *User) Create(ctx *context.Context) error {
	exists, err := user.Exists(ctx)

	if err != nil {
		return err
	}

	if exists {
		return errors.New("user already exits")
	}

	user.FillDefaults()
	collection := helpers.GetCollection(config.UserDatabaseName, userCollectionName)

	_, err = collection.InsertOne(*ctx, user)

	if err != nil {
		return err
	}

	go user.createHistory()

	return nil
}

func (user *User) GenerateTokensAndPersistClaims(ctx *context.Context, device *Device) (*string, *string, error) {
	accessClaim, accessToken, err := user.CreateClaim(config.AccessTokenDuration, config.AccessTokenSecretKey)
	if err != nil {
		return nil, nil, err
	}

	refreshClaim, refreshToken, err := user.CreateClaim(config.RefreshTokenDuration, config.RefreshTokenSecretKey)
	if err != nil {
		return nil, nil, err
	}

	claims := &TokenClaims{
		CreatedAt:    time.Now(),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		RefreshClaim: refreshClaim,
		AccessClaim:  accessClaim,
		DeviceId:     device.ID,
	}

	err = claims.Persist(ctx)

	if err != nil {
		return nil, nil, err
	}

	return accessToken, refreshToken, nil
}

func (user *User) CreateClaim(duration time.Duration, secret string) (*UserClaim, *string, error) {
	claims := &UserClaim{
		ID:          user.ID,
		Email:       user.Email,
		PhoneNumber: user.PhoneNumber,
		StandardClaims: jwt.StandardClaims{
			Issuer:   config.JWTIssuerName,
			IssuedAt: time.Now().Unix(),
			Audience: config.JWTAudience,
			Id:       uuid.New().String(),
			// todo => NotBefore "nbf"
			// todo => Subject   "sub"
			// Stays valid for 10 hours
			ExpiresAt: time.Now().Local().
				Add(duration).Unix(),
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).
		SignedString([]byte(secret))

	return claims, &token, err
}

func (tc *TokenClaims) Persist(ctx *context.Context) error {
	id := uuid.New().String()
	tc.ID = &id

	collection := helpers.GetCollection(config.UserDatabaseName, tokenClaimsCollectionName)
	_, err := collection.InsertOne(*ctx, tc)

	return err
}
