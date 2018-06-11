package web

import (
	"context"
	"encoding/json"
	"github.com/go-chi/chi"
	"github.com/ssOleg/go_service/go_server/storage"
	"log"
	"net/http"
	"os"
)

type Router struct {
	Storage storage.DBase // Interface
}

func (router *Router) Ctx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gifID := chi.URLParam(r, "gifID")
		element, err := router.Storage.Get(gifID)
		if err != nil {
			http.Error(w, http.StatusText(404), 404)
			return
		}
		ctx := context.WithValue(r.Context(), "gif", element)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Get items
func (router *Router) getGifs(w http.ResponseWriter, r *http.Request) {
	elements, err := router.Storage.GetAll()
	check(err)
	json.NewEncoder(w).Encode(storage.Results{elements})
}

// Create a new item
func (router *Router) createGifs(w http.ResponseWriter, r *http.Request) {
	var element storage.Element
	err := json.NewDecoder(r.Body).Decode(&element)
	check(err)
	if element == (storage.Element{}) {
		json.NewEncoder(w).Encode("Please use correct format")
		return
	}

	err = router.Storage.Insert(element)
	check(err)

	json.NewEncoder(w).Encode(element)
}

// Get an item
func (router *Router) getGif(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	element, ok := ctx.Value("gif").(storage.Element)
	if !ok {
		http.Error(w, http.StatusText(422), 422)
		return
	}
	json.NewEncoder(w).Encode(element)
}

// Delete an item
func (router *Router) deleteGif(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	element, ok := ctx.Value("gif").(storage.Element)
	if !ok {
		http.Error(w, http.StatusText(422), 422)
		return
	}
	err := router.Storage.Remove(element)
	check(err)
	json.NewEncoder(w).Encode("Element was deleted")
}

func check(e error) {
	if e != nil {
		//TODO: Add better logging
		log.Fatal(e)
		os.Exit(1)
	}
}
