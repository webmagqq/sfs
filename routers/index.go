package routers

import (
	"fmt"
	"html/template"
	"net/http"
	"sfsgo/gstr"
	"sfsgo/index"
	"sfsgo/pubgo"
)

func Index(w http.ResponseWriter, req *http.Request) {
	uif := pubgo.Newurlinfo(req)
	rd := index.Newindex(req)
	rd.Setto(&uif)
	rd.SetFilecontent()
	rd.SetFilecontent1()

	//统计
	pubgo.Tj.Brows(gstr.Mstr(req.URL.Path, "/", "/"))
	//--组织模板数据
	TemplatesFiles := []string{
		"index/index.html",
		"tplpub/static.html",
		//"tpl/public/header.html",
		//"tpl/public/footer.html", // 多加的文件
	}
	//t, _ := template.ParseFiles(TemplatesFiles...)

	funcMap := template.FuncMap{ //--需要注册的函数
		"jf": pubgo.Jf,
	}
	t, _ := template.New("index.html").Funcs(funcMap).ParseFiles(TemplatesFiles...)
	//--New("index.html") 的 index.html必须是TemplatesFiles第一个文件名

	err := t.ExecuteTemplate(w, "index.html", rd)
	if err != nil {
		fmt.Println(req.URL.Path, err)
	}
}
