package mapp

import (
	"net/http"
	"sfsgo/pubgo"
)

type mapp struct {
	req  *http.Request
	Host string
	Path string
	Uri  string
	Up   string //繁体参数 ?l=0
	Ift  bool   //是否转繁体
}

func Newmapp(req *http.Request) mapp {
	//静态变量返回值用指针*,非静态类则返回copy一份，取决于new
	return mapp{req: req}
}

func (m *mapp) Setto(uif *pubgo.Urlinfo) {
	m.Host = uif.Host
	m.Path = uif.Path
	m.Uri = uif.Uri
	m.Up = uif.Up
	m.Ift = uif.Ift
}
