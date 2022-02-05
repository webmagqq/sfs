package search

//将目录加载如内存map
import (
	"sfsgo/mysql"
)

var Mlmap map[string]string

func Newmlmap() {
	Mlmap = make(map[string]string)
	rd := mysql.Selects("SELECT tid,jingming from jingmulu")
	tid, jingming := "", ""
	for _, v := range rd {
		tid = v["tid"]
		jingming = v["jingming"]
		Mlmap[tid] = jingming
	}
}
