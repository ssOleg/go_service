package web

import (
	"testing"
	"net/http"
	"net/http/httptest"
	"github.com/ssOleg/go_service/go_server/storage"
	"github.com/go-chi/chi"
)



//func baseEndpoint(w http.ResponseWriter, req *http.Request) {
//	w.WriteHeader(200)
//	w.Write([]byte("Hello page"))
//}

func TRouter() *chi.Mux {
	dbStructure := storage.DataBase{ConnectionPoint: "localhost"}
	webRouter := Router{Storage: dbStructure}
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

