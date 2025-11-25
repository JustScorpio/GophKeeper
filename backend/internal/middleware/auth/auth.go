// Пакет auth содержит middleware а также вспомогательные функции для аутентификации и авторизации пользователей
package auth

import (
	"errors"
	"net/http"
	"time"

	"github.com/JustScorpio/GophKeeper/backend/internal/customcontext"
	"github.com/golang-jwt/jwt/v4"
)

const (
	// Имя куки с JWT-токеном
	jwtCookieName = "jwt_token"
	//Время жизни токена
	tokenLifeTime = time.Hour * 3
	// Ключ для генерации и расшифровки токена (В РЕАЛЬНОМ ПРИЛОЖЕНИИ ХРАНИТЬ В НАДЁЖНОМ МЕСТЕ)
	secretKey = "supersecretkey"
)

// Claims — структура утверждений, которая включает стандартные утверждения и одно пользовательское UserID
type Claims struct {
	jwt.RegisteredClaims
	Login string `json:"login"`
}

// newJWTString - создаёт токен с логином пользователя
func newJWTString(login string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenLifeTime)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		Login: login, // Сохраняем логин в токене
	})

	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// SetJWTCookie - устанавливает JWT куку
func SetJWTCookie(w http.ResponseWriter, userID string) error {
	newToken, err := newJWTString(userID)
	if err != nil {
		return err
	}

	newCookie := &http.Cookie{
		Name:     jwtCookieName,
		Value:    newToken,
		Path:     "/",
		Expires:  time.Now().Add(tokenLifeTime),
		HttpOnly: true,
	}

	http.SetCookie(w, newCookie)
	return nil
}

// GetLoginFromToken - извлекает логин из JWT токена
func GetLoginFromToken(tokenString string) (string, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err != nil {
		return "", err
	}

	if !token.Valid {
		return "", errors.New("invalid token")
	}

	return claims.Login, nil
}

// AuthMiddleware - middleware для проверки аутентификации (добавляет в контекст login пользователя)
func AuthMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var login string

			cookie, err := r.Cookie(jwtCookieName)
			if err == nil {
				claims := &Claims{}
				token, err := jwt.ParseWithClaims(cookie.Value, claims, func(t *jwt.Token) (interface{}, error) {
					return []byte(secretKey), nil
				})

				if err == nil && token.Valid {
					login = claims.Login // Теперь получаем логин из токена
				}
			}

			// Добавляем логин в контекст (может быть пустым если нет валидного токена)
			ctx := customcontext.WithUserID(r.Context(), login) // Используем логин как UserID
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
