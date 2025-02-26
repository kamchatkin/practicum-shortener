package auth

import (
	"context"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/kamchatkin/practicum-shortener/internal/data"
	"github.com/kamchatkin/practicum-shortener/internal/storage"
	"net/http"
	"time"
)

const CookineName = "goadv.auth"

type Claims struct {
	jwt.RegisteredClaims
	UserID int64
}

const TokenExp = time.Hour * 24 * 360
const SecretKey = "di4kooj*o"

const DefaultUserID = -1

// BuildJWTString создаёт токен и возвращает его в виде строки.
func BuildJWTString() (string, error) {

	// сделал и подумал что можно бы и
	db, _ := storage.NewStorage()
	userID, err := data.RegisterUser(context.TODO(), db)
	if err != nil {
		return "", err
	}

	// создаём новый токен с алгоритмом подписи HS256 и утверждениями — Claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			// когда создан токен
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenExp)),
		},
		// собственное утверждение
		UserID: userID,
	})

	// создаём строку токена
	tokenString, err := token.SignedString([]byte(SecretKey))
	if err != nil {
		return "", err
	}

	// возвращаем строку токена
	return tokenString, nil
}

func GetUserID(tokenString string) int64 {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims,
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(SecretKey), nil
		})
	if err != nil {
		return DefaultUserID
	}

	if !token.Valid {
		return DefaultUserID
	}

	return claims.UserID
}

// GetUserIDFromCookie возвращает ID пользователя из Cookie
func GetUserIDFromCookie(r *http.Request) int64 {
	cookie, err := r.Cookie(CookineName)
	if err != nil {
		return DefaultUserID
	}

	return GetUserID(cookie.Value)
}
