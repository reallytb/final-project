package signin

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"final-project/internal/api/task"
)

var pass = os.Getenv("TODO_PASSWORD")
var pasEnv = os.Getenv("TODO_PASSWORD")

type password struct {
	Password string `json:"password"`
}

type tokenJson struct {
	Token string `json:"token"`
}

func SigninHandler(w http.ResponseWriter, r *http.Request) {
	var pass password
	err := json.NewDecoder(r.Body).Decode(&pass)
	if err != nil {
		log.Println(err.Error())
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		task.SendJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	r.Body.Close()
	if pass.Password != pasEnv {
		log.Println("введён неверный пароль")
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		task.SendJSONError(w, "введён неверный пароль", http.StatusUnauthorized)
		return
	}
	claims := jwt.MapClaims{
		"authorized": true,
		"exp":        time.Now().Add(time.Hour * 8).Unix(),
		"iat":        time.Now().Unix(),
	}
	tokenJWT := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	var tokenJson tokenJson
	tokenJson.Token, err = tokenJWT.SignedString([]byte(pasEnv))
	if err != nil {
		log.Println(err.Error())
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		task.SendJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    tokenJson.Token,
		Expires:  time.Now().Add(8 * time.Hour),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(tokenJson)
}

func Auth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if len(pass) > 0 {
			var jwtString string
			cookie, err := r.Cookie("token")
			if err == nil {
				jwtString = cookie.Value
			}
			token, err := jwt.Parse(jwtString, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("ошибка: неизвестный метод подписи %v", token.Header["alg"])
				}
				return []byte(pass), nil
			})

			if err != nil || !token.Valid {
				http.Error(w, "ошибка: неверный токен", http.StatusUnauthorized)
				return
			}

			if claims, ok := token.Claims.(jwt.MapClaims); ok {
				if exp, ok := claims["exp"].(float64); ok {
					if time.Now().Unix() > int64(exp) {
						http.Error(w, "Token expired", http.StatusUnauthorized)
						return
					}
				}
				if authorized, ok := claims["authorized"].(bool); !ok || !authorized {
					http.Error(w, "Not authorized", http.StatusUnauthorized)
					return
				}
			}
		}
		next(w, r)
	})
}
