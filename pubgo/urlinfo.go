package pubgo

//---公共参数
import (
	"net/http"
)

//--其他包引用必须大写。
type Urlinfo struct {
	Host string
	Path string
	Uri  string
	Ift  bool   //是否转繁体
	Up   string //繁体参数 ?l=0
}

func Newurlinfo(req *http.Request) Urlinfo {
	//静态变量返回值用指针*,非静态类则返回copy一份，取决于new
	u := ""
	if req.URL.Query().Get("l") != "" {
		u = "?l=0"
	}
	return Urlinfo{
		Host: "http://" + req.Host,
		Path: req.URL.Path,
		Uri:  req.RequestURI,
		Up:   u,
		Ift:  u != "",
	}
}
