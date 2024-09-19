package controllers

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"

	"agents-web/src/app/services"
	"agents-web/src/templates/layouts"
	"agents-web/src/templates/pages"
	"agents-web/src/templates/root"
)

func Auth() *chi.Mux {
	r := chi.NewRouter()

	r.Get("/log-in", func(w http.ResponseWriter, r *http.Request) {
		root.App(layouts.CenterLayout(pages.Login())).Render(r.Context(), w)
	})

	return r
}


type SignUpRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}
func AuthApi(auth_service *services.Auth) *chi.Mux {
	r := chi.NewRouter()

	r.Post("/v1/sign-in", func(w http.ResponseWriter, r *http.Request) {
		superlogin := os.Getenv("SUPERADMIN_LOGIN")
		superpassword := os.Getenv("SUPERADMIN_PASSWORD")

		var request SignUpRequest
		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if request.Login == superlogin && request.Password == superpassword {
			token, err := auth_service.Generate(services.Payload{
				Id:    request.Login,
				Email: request.Login,
				Super: true,
			})
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			
			http.SetCookie(w, &http.Cookie{
				Name:  "token",
				Value: token,
				Path: "/",
			})
		}

		w.Write([]byte("WOW"))
	})

	return r
}

func AuthMiddleware(auth_service *services.Auth) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token, err := r.Cookie("Authorization")

			if err != nil {
				if err == http.ErrNoCookie {
					http.Redirect(w, r, "/log-in", http.StatusUnauthorized)
					return
				}
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			valid := auth_service.Verify(token.Value)
			if !valid {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
