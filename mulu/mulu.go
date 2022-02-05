package mulu

import (
	"net/http"
	"sfsgo/pubgo"
)

type mulu struct {
	req  *http.Request
	Host string
	Path string
	Uri  string
	Up   string //繁体参数 ?l=0
	Ift  bool   //是否转繁体
}

func Newmulu(req *http.Request) mulu {
	//静态变量返回值用指针*,非静态类则返回copy一份，取决于new
	return mulu{req: req}
}

func (q *mulu) Setto(uif *pubgo.Urlinfo) {
	q.Host = uif.Host
	q.Path = uif.Path
	q.Uri = uif.Uri
	q.Up = uif.Up
	q.Ift = uif.Ift
}
