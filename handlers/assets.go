package handlers

import (
	"fmt"
	"net/http"
)

type AssetsHandler struct {
	Kind string
}

func (a *AssetsHandler) HandleAssets(w http.ResponseWriter, r *http.Request) {
	http.StripPrefix(fmt.Sprint("/", a.Kind), http.FileServer(http.Dir(fmt.Sprint("./public/", a.Kind)))).ServeHTTP(w, r)
}
