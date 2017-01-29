package handlers

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/gocraft/web"
	"github.com/royvandewater/meshchain/record"
)

// CreateRecord handles a POST Record request
func CreateRecord(rw web.ResponseWriter, req *web.Request) {
	data, err := ioutil.ReadAll(req.Body)
	if err != nil {
		respondWithError(rw, 500, err)
		return
	}

	metadata := parseMetadata(req.Header)
	signatureBase64 := parseAuthHeader(req.Header)

	rootRecord, err := record.NewRootRecord(metadata, data, signatureBase64)
	if err != nil {
		respondWithError(rw, 422, err)
		return
	}

	responseBody, err := rootRecord.JSON()
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

func parseMetadata(headers map[string][]string) record.Metadata {
	ID := safeGetFirst(headers, "meshchain-id")
	LocalID := safeGetFirst(headers, "meshchain-local-id")

	return record.Metadata{
		ID:         ID,
		LocalID:    LocalID,
		PublicKeys: headers["meshchain-public-key"],
	}
}

func respondWithError(rw web.ResponseWriter, code int, err error) {
	rw.WriteHeader(code)
	fmt.Fprintf(rw, `{"error": "%v"}`, err.Error())
}

func safeGetFirst(headers map[string][]string, key string) string {
	if len(headers[key]) > 0 {
		return headers[key][0]
	}
	return ""
}
