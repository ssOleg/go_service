package web

import (
	"net/http"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)


func GetRouter(handler Router) *chi.Mux {
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello page"))
	})

	router.Route("/gifs", func(r chi.Router) {
		r.Get("/", handler.getGifs)
		r.Post("/", handler.createGifs)

		r.Route("/{gifID}", func(r chi.Router) {
			r.Use(handler.Ctx)
			r.Get("/", handler.getGif)
			r.Delete("/", handler.deleteGif)
		})
	})

	return router
}

