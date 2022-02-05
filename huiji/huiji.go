package huiji

import (
	"net/http"
	"sfsgo/pubgo"
)

type huiji struct {
	req  *http.Request
	Host string
	Path string
	Uri  string
	Up   string //繁体参数 ?l=0
	Ift  bool   //是否转繁体
}

func Newhuiji(req *http.Request) huiji {
	//静态变量返回值用指针*,非静态类则返回copy一份，取决于new
	return huiji{req: req}
}

func (m *huiji) Setto(uif *pubgo.Urlinfo) {
	m.Host = uif.Host
	m.Path = uif.Path
	m.Uri = uif.Uri
	m.Up = uif.Up
	m.Ift = uif.Ift
}
