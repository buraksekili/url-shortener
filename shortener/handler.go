package shortener

import (
	"net/http"
)

func URLHandler(paths map[string]string, fallback http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		path := request.URL.Path
		if dest, ok := paths[path]; ok {
			http.Redirect(writer, request, dest, http.StatusFound)
			return
		}
		fallback.ServeHTTP(writer, request)
	})
}
