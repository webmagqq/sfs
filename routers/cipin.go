package routers

import (
	"html/template"
	"net/http"
	"sfsgo/cipin"
	"sfsgo/gstr"
	"sfsgo/pubgo"
)

func Cipin(w http.ResponseWriter, req *http.Request) {
	//数据组织
	uif := pubgo.Newurlinfo(req)
	rd := cipin.Newcipin(req)
	rd.Setto(&uif)
	rd.SetP()
	rd.SetKwhtml()
	rd.SetPgHtml()
	//--组织模板数据
	//统计
	pubgo.Tj.Brows(gstr.Mstr(req.URL.Path, "/", "/"))

	TemplatesFiles := []string{
		"cipin/cipin.html",
		"tplpub/static.html",
		"tplpub/header.html",
		"tplpub/footer.html", // 多加的文件
	}
	funcMap := template.FuncMap{ //--需要注册的函数
		"jf": pubgo.Jf, //简繁转换
	}
	t, _ := template.New("cipin.html").Funcs(funcMap).ParseFiles(TemplatesFiles...)
	//--  quanwen.html必须是TemplatesFiles第一个文件名
	t.ExecuteTemplate(w, "cipin.html", rd)

}
