package routers

import (
	"fmt"
	"net/http"
)

func Llan(w http.ResponseWriter, req *http.Request) {
	fmt.Println("1")
	/*
		l := req.URL.Query().Get("l")
		i := req.URL.Query().Get("i")
		p := req.URL.Query().Get("p")
		k := req.URL.Query().Get("k")
		ip := ClientIP(req)
		if glip(ip) {
			return
		}
		//fmt.Println(time.Now())
		rd := mysql.Exesql("INSERT INTO llantji(lxing,idkey,llan,purl,ip,rq) VALUES(?,?,?,?,?,NOW())", l, i, k, p, ip)
		if rd > -1 {
			w.Write([]byte("1"))
		} else {
			w.Write([]byte("0"))
		}
	*/
}

/*
func glip(ip string) bool {
	glip := "116.179.37|111.206.221|111.206.198"
	tip := gstr.Do(ip, "", ".", false, true)
	return strings.Contains(glip, tip)
}
func ClientIP(r *http.Request) string {
	xForwardedFor := r.Header.Get("X-Forwarded-For")
	ip := strings.TrimSpace(strings.Split(xForwardedFor, ",")[0])
	if ip != "" {
		return ip
	}
	ip = strings.TrimSpace(r.Header.Get("X-Real-Ip"))
	if ip != "" {
		return ip
	}
	if ip, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr)); err == nil {
		return ip
	}
	return ""
}
*/
