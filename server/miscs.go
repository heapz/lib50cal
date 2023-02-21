package server

import "net/http"

func httpError(w http.ResponseWriter, code int) { http.Error(w, http.StatusText(code), code) }
