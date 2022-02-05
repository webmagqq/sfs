package pubgo

//--读取配置文件
import (
	"io/ioutil"
	"sfsgo/gstr"
	"strings"
)

//---公共参数

type setfile struct {
	filetext string
}

func Newsetfile(n string) setfile {
	set, _ := ioutil.ReadFile(n)
	return setfile{
		filetext: string(set[:]),
	}
}
func (s *setfile) Gp(name string) string {
	rs := gstr.Mstr(s.filetext, name+"=", "\r\n")
	return strings.TrimSpace(rs)
}
