package routers

import (
	"html/template"
	"net/http"
	"sfsgo/err"
	"sfsgo/gstr"
	"sfsgo/pubgo"
)

func Err(w http.ResponseWriter, req *http.Request) {
	//数据组织
	uif := pubgo.Newurlinfo(req)
	rd := err.Newerr(req)
	rd.Setto(&uif)
	w.WriteHeader(404)
	//--组织模板数据
	//统计
	pubgo.Tj.Brows(gstr.Mstr(req.URL.Path, "/", "/"))

	TemplatesFiles := []string{
		"err/err.html",
		"tplpub/static.html",
		"tplpub/header.html",
		"tplpub/footer.html", // 多加的文件
	}
	funcMap := template.FuncMap{ //--需要注册的函数
		"jf": pubgo.Jf, //简繁转换
	}
	t, _ := template.New("err.html").Funcs(funcMap).ParseFiles(TemplatesFiles...)
	//--  quanwen.html必须是TemplatesFiles第一个文件名
	t.ExecuteTemplate(w, "err.html", rd)

}
