package router

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httpcoala"
	"github.com/neilli-sable/sideupload/infrastructure/setting"
	"github.com/neilli-sable/sideupload/presenter/handler"
	"github.com/neilli-sable/sideupload/presenter/httpserve"
)

var revision = "undefined"

// Router ルーター
func Router(setting *setting.Setting) *chi.Mux {
	sideupload := handler.SideUploadHandler{}
	r := chi.NewRouter()
	r.Use(httpcoala.Route("HEAD", "GET"))
	r.Use(middleware.StripSlashes)
	r.Use(middleware.RequestID)

	cors := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link", "Content-Disposition"},
		AllowCredentials: true,
		MaxAge:           300,
	})
	r.Use(cors.Handler)

	r.Route("/", func(r chi.Router) {
		r.Get("/info", func(w http.ResponseWriter, r *http.Request) {
			httpserve.JSON(w, struct {
				AppName  string `json:"appName"`
				Revision string `json:"revision"`
			}{
				AppName:  "sideupload",
				Revision: revision,
			}, http.StatusOK)
		})

		r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
			httpserve.JSON(w, struct{}{}, http.StatusOK)
		})

		r.Post("/backup", sideupload.Backup)
		r.Post("/clean", sideupload.Clean)
	})

	return r
}
