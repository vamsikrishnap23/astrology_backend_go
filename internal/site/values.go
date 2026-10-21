package site

import (
	_ "embed"
	"net/http"
)

// valuesJSON lists every value the site endpoints return, with Telugu.
//
//go:embed values.json
var valuesJSON []byte

// ValuesHandler serves GET /site/values.
func ValuesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "use GET")
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Cache-Control", "public, max-age=86400")
	w.Write(valuesJSON)
}
