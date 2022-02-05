package routers

import (
	"net/http"
)

func Chat(w http.ResponseWriter, req *http.Request) {
	http.ServeFile(w, req, "index.html")
}
