package shortener

import (
	"net/http"
)

func URLHandler(paths map[string]string, handler http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if dest, ok := paths[path]; ok {
			http.Redirect(w, r, dest, http.StatusFound)
			return
		}
		handler.ServeHTTP(w, r)
	})
}
