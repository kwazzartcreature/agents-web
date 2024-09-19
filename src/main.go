package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"

	"agents-web/src/app/controllers"
	"agents-web/src/app/services"
	"agents-web/src/templates/layouts"
	"agents-web/src/templates/pages"
	"agents-web/src/templates/root"
)

func main() {
	godotenv.Load()

	r := chi.NewRouter()
	r.Handle("/css/*", http.StripPrefix("/css/", http.FileServer(http.Dir("./build/css"))))

	

	auth_service := services.NewAuth(os.Getenv("JWT_KEY"))
	r.Mount("/", controllers.Auth())
	r.Mount("/api", controllers.AuthApi(auth_service))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		token, err := r.Cookie("token")
		if err != nil {
			if err == http.ErrNoCookie {
				http.Redirect(w, r, "/log-in", http.StatusUnauthorized)
				return
			}
		}

		fmt.Println("Token:", token)
    	root.App(layouts.AppLayout(pages.Home())).Render(r.Context(), w)
	})

	r.Post("/clicked", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("WOW"))
	})

	http.ListenAndServe(":8000", r)
}
