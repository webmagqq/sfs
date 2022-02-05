package routers

import (
	"fmt"
	"html/template"
	"net/http"
	"sfsgo/gstr"
	"sfsgo/pubgo"
	"sfsgo/xianyan"
)

func Xianyan(w http.ResponseWriter, req *http.Request) {
	uif := pubgo.Newurlinfo(req)
	rd := xianyan.Newxianyan(req)
	rd.Setto(&uif)
	rd.SetP()
	rd.Setfs()
	rd.SetText()
	rd.SetPg()

	//统计
	pubgo.Tj.Brows(gstr.Mstr(req.URL.Path, "/", "/"))

	//--组织模板数据
	TemplatesFiles := []string{
		"xianyan/xianyan.html",
		"tplpub/static.html",
	}
	//t, _ := template.ParseFiles(TemplatesFiles...)

	funcMap := template.FuncMap{ //--需要注册的函数
		"jf": pubgo.Jf,
	}
	t, _ := template.New("xianyan.html").Funcs(funcMap).ParseFiles(TemplatesFiles...)
	//--New("xianyan.html") 的 xianyan.html必须是TemplatesFiles第一个文件名

	err := t.ExecuteTemplate(w, "xianyan.html", rd)
	if err != nil {
		fmt.Println(req.URL.Path, err)
	}
}
