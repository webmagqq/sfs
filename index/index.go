package index

import (
	"html/template"
	"net/http"
	"sfsgo/pubgo"
	"strings"
)

type index struct {
	req          *http.Request
	Host         string
	Path         string
	Uri          string
	Up           string //繁体参数 ?l=0
	Ift          bool   //是否转繁体
	Filecontent  template.HTML
	Filecontent1 template.HTML
}

func Newindex(req *http.Request) index {
	return index{req: req}
}
func (i *index) Setto(uif *pubgo.Urlinfo) {
	i.Host = uif.Host
	i.Path = uif.Path
	i.Uri = uif.Uri
	i.Up = uif.Up
	i.Ift = uif.Ift
}
func (i *index) SetFilecontent() {
	fc := pubgo.Of("./hc/佛教，践行绝对真理的教育.txt")
	fc = pubgo.Jf(fc, i.Up != "")
	i.Filecontent = template.HTML(fc)
}
func (i *index) SetFilecontent1() {
	fc := pubgo.Of("./hc/s.txt")
	fc = pubgo.Jf(fc, i.Up != "")
	if i.Ift {
		fc = strings.Replace(fc, "\">", i.Up+"\">", -1)
	}
	i.Filecontent1 = template.HTML(fc)
}
