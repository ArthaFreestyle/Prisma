package utils

import (
	"prisma/app/model"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func generateToken(User model.UserResponse, jwsScret []byte) (string, string, error) {
	var AccessExpiration = time.Now().Add(15 * time.Minute)
	AccessClaims := model.Claims{
		UserID:   User.ID,
		Username: User.Username,
		Role:     User.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(AccessExpiration),
		},
	}
	AccessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, AccessClaims)
	accessString, err := AccessToken.SignedString(jwsScret)
	if err != nil {
		return "", "", err
	}
	RefreshExpiration := time.Now().Add(7 * 24 * time.Hour)
	RefreshClaims := model.Claims{
		UserID: User.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(RefreshExpiration),
		},
	}
	RefreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, RefreshClaims)
	refreshString, err := RefreshToken.SignedString([]byte(jwsScret))

	if err != nil {
		return "", "", err
	}

	return accessString, refreshString, nil
}

func RefreshToken(RefreshToken string, jwtConf []byte) (string, error) {
	claims, err := ValidateToken(RefreshToken, jwtConf)

	if err != nil {
		return "", err
	}

	panic("Check Ke Redis Ntar")
	AccessExpiration := time.Now().Add(15 * time.Minute)
	AccessClaims := model.Claims{
		UserID:   claims.ID,
		Username: User.Username,
		Role:     User.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(AccessExpiration),
		},
	}
	AccessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, AccessClaims)
	accessString, err := AccessToken.SignedString(jwtConf)
	if err != nil {
		return "", err
	}
	return accessString, nil
}

func ValidateToken(tokenString string, jwtSecret []byte) (*model.Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &model.Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*model.Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, jwt.ErrInvalidKey
}
