package routers

import (
	"fmt"
	"html/template"
	"net/http"
	"sfsgo/gstr"
	"sfsgo/mulu"
	"sfsgo/pubgo"
)

func Mulu(w http.ResponseWriter, req *http.Request) {
	uif := pubgo.Newurlinfo(req)
	rd := mulu.Newmulu(req)
	rd.Setto(&uif)

	//统计
	pubgo.Tj.Brows(gstr.Mstr(req.URL.Path, "/", "/"))
	//--组织模板数据

	TemplatesFiles := []string{
		"mulu/mulu.html",
		"tplpub/static.html",
		"tplpub/header.html",
		"tplpub/footer.html", // 多加的文件
	}
	//t, _ := template.ParseFiles(TemplatesFiles...)

	funcMap := template.FuncMap{ //--需要注册的函数
		"jf": pubgo.Jf,
	}
	t, _ := template.New("mulu.html").Funcs(funcMap).ParseFiles(TemplatesFiles...)
	//--New("index.html") 的 index.html必须是TemplatesFiles第一个文件名

	err := t.ExecuteTemplate(w, "mulu.html", rd)
	if err != nil {
		fmt.Println(req.URL.Path, err)
	}
}
