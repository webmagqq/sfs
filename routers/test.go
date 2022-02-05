package routers

import (
	"net/http"
	"sfsgo/test"
)

func Test(w http.ResponseWriter, req *http.Request) {
	id := req.URL.Query().Get("id")
	if id == "" {
		w.Write([]byte(""))
		return
	}
	ifs := test.Newifos()
	mdo := ifs.Mdo[id]
	rst := "参数不对"
	if mdo != nil {
		rst = mdo()
	}
	w.Write([]byte(rst))

}
