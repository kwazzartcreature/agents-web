package services

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

type Payload struct {
	Id    string `json:"id"`
	Email string `json:"email"`
	Super bool   `json:"super"`
}

type Auth struct {
	key []byte
}

func NewAuth (key string) *Auth {
	return &Auth{
		key: []byte(key),
	}
}

func (a *Auth) Generate(payload Payload) (string, error) {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    payload.Id,
		"email": payload.Email,
		"super": payload.Super,
	})
	return t.SignedString((a.key))
}


func (a *Auth) Verify(token_string string) bool {
	token, err := jwt.Parse(token_string, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return a.key, nil
	})

	if err != nil {
		fmt.Println("Error parsing token:", err)
		return false
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		fmt.Println("Token is valid. Claims:", claims)
		return true
	} else {
		fmt.Println("Token is invalid or claims are not of type jwt.MapClaims.")
		return false
	}
}

func (a *Auth) Regenerate (token_string string) {
	valid := a.Verify(token_string)
	if !valid {
		fmt.Printf("Invalid token: %s", token_string)
		return
	}

	// TODO fetch actual data
}