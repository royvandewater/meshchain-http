package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gocraft/web"
	"github.com/royvandewater/meshchain/record"
)

func parseAuthHeader(headers map[string][]string) string {
	authHeaders, ok := headers["Authorization"]
	if !ok {
		return ""
	}

	authHeader := ""
	for _, header := range authHeaders {
		if strings.HasPrefix(header, "Bearer ") {
			authHeader = header
			break
		}
	}

	if authHeader == "" {
		return ""
	}

	return strings.Replace(authHeader, "Bearer ", "", 1)
}

// CreateRecord handles a POST Record request
func CreateRecord(rw web.ResponseWriter, req *web.Request) {
	jwt := parseAuthHeader(req.Header)

	if jwt == "" {
		respondWithError(rw, 401, fmt.Errorf("Unauthorized"))
		return
	}

	rec, err := record.NewFromReader(req.Body)
	if err != nil {
		respondWithError(rw, 422, err)
		return
	}

	err = rec.Validate(jwt)
	if err != nil {
		respondWithError(rw, 422, err)
		return
	}

	err = rec.Save()
	if err != nil {
		respondWithError(rw, 500, err)
		return
	}

	responseBody, err := rec.ToJSON()
	if err != nil {
		respondWithError(rw, 500, err)
		return
	}

	rw.WriteHeader(http.StatusCreated)
	fmt.Fprintf(rw, responseBody)
}

// GetRecord handles a GET Record request
func GetRecord(rw web.ResponseWriter, req *web.Request) {
	record.Get("hi")
	rw.WriteHeader(http.StatusNoContent)
}

func respondWithError(rw web.ResponseWriter, code int, err error) {
	rw.WriteHeader(code)
	fmt.Fprintf(rw, `{"error": "%v"}`, err.Error())
}
