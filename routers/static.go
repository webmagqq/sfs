package routers

import (
	"net/http"
)

//--静态文件服务
func Static(w http.ResponseWriter, req *http.Request) {
	had := http.StripPrefix("/static/", http.FileServer(http.Dir("./static")))
	had.ServeHTTP(w, req)

}
