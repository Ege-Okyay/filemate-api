package utils

import (
	"bytes"
	"net/http"
	"net/http/httptest"
)

func PerformRequest(r http.Handler, method, path string, payload []byte) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, bytes.NewBuffer(payload))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	return w
}
