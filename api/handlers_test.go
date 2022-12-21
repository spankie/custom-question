package api

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

var mux = http.NewServeMux()

func TestMain(m *testing.M) {
	client = &http.Client{
		Timeout: time.Second * 10,
		Transport: localRoundTripper{
			handler: mux,
		},
	}

	os.Exit(m.Run())
}

type localRoundTripper struct {
	handler http.Handler
}

func (l localRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	w := httptest.NewRecorder()
	l.handler.ServeHTTP(w, req)
	return w.Result(), nil
}

func mockRoute(url string, status int, body []byte) {
	mux.HandleFunc(url, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(status)
		_, err := io.WriteString(w, string(body))
		if err != nil {
			panic(err)
		}
	})
}

func TestSearchHandler(t *testing.T) {
	mockRoute("/search", http.StatusOK, []byte(`
	{
		"tracks": {
			"hits": [
				{
					"track": {
						"key": "1",
						"title": "sungba",
						"subtitle": "asake"
					}
				}
			]
		}
	}`))

	mockRoute("/songs/get-count", http.StatusOK, []byte(`
	{
		"id":"40333609",
		"total":3457848,
		"type":"tag"
	}`))
	res := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(res)
	c.Request = httptest.NewRequest(http.MethodGet, "/search?name=sungba", nil)
	SearchHandler(c)

	if res.Code != http.StatusOK {
		t.Errorf("status code should be %v intead got %v", http.StatusOK, res.Code)
	}

	searchResult := struct {
		Result []Track
	}{}
	bytes := res.Body.Bytes()
	log.Printf("hits: %v", string(bytes))
	err := json.Unmarshal(bytes, &searchResult)
	if err != nil {
		t.Errorf("could not marshal response: %v", err)
	}

	if len(searchResult.Result) == 0 {
		t.Errorf("we should get at least one hit")
	}
}
