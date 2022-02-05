package huchi

import (
	"net/http"
	"sfsgo/pubgo"
)

type huchi struct {
	req  *http.Request
	Host string
	Path string
	Uri  string
	Up   string //繁体参数 ?l=0
	Ift  bool   //是否转繁体
}

func Newhuchi(req *http.Request) huchi {
	//静态变量返回值用指针*,非静态类则返回copy一份，取决于new
	return huchi{req: req}
}

func (h *huchi) Setto(uif *pubgo.Urlinfo) {
	h.Host = uif.Host
	h.Path = uif.Path
	h.Uri = uif.Uri
	h.Up = uif.Up
	h.Ift = uif.Ift
}
