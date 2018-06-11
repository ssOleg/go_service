package web

import (
	"github.com/go-chi/chi"
	"github.com/ssOleg/go_service/go_server/storage"
	"net/http"
	"net/http/httptest"
	"testing"
)

//func baseEndpoint(w http.ResponseWriter, req *http.Request) {
//	w.WriteHeader(200)
//	w.Write([]byte("Hello page"))
//}

func TRouter() *chi.Mux {
	dbStructure := storage.DataBase{ConnectionPoint: "localhost"}
	webRouter := Router{Storage: &dbStructure}
	router := GetRouter(webRouter)
	//router.HandleFunc("/", baseEndpoint)

	return router
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

func TestGet(t *testing.T) {
	request, _ := http.NewRequest("GET", "/", nil)
	response := httptest.NewRecorder()
	TRouter().ServeHTTP(response, request)
	checkResponseCode(t, http.StatusOK, response.Code)

	if body := response.Body.String(); body != "Hello page" {
		t.Errorf("Expected an empty array. Got %s", body)
	}

}

func BenchmarkInitialGet(b *testing.B) {
	request, _ := http.NewRequest("GET", "/", nil)
	response := httptest.NewRecorder()

	for i := 0; i < b.N; i++ {
		TRouter().ServeHTTP(response, request)
	}
}
