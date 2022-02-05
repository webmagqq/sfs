package main

//搜佛说golang版本
import (
	"log"
	"net/http"
	"os"
	"sfsgo/cache"
	"sfsgo/mysql"
	"sfsgo/pubgo"
	"sfsgo/routers"
	"sfsgo/search"
	"strconv"
)

func main() {
	//--读取初始化设置数据
	gf := pubgo.Newsetfile("./设置.txt")
	ip := gf.Gp("ip")
	port := ""
	port = gf.Gp("port")
	if port == "" {
		port = os.Getenv("ASPNETCORE_PORT") // 默认对IIS asp.net core 支持，获取asp.net core 分配的port
	}
	mysql.NewsqlDb(gf.Gp("root"), gf.Gp("psw"), gf.Gp("dbname")) //数据库连接
	cachect := gf.Gp("cache")
	ct, _ := strconv.Atoi(cachect)              //字符串转int
	search.CacheSearchData = cache.NewCache(ct) //新建搜索缓存
	search.CacheNrData = cache.NewCache(ct)     //新建搜索结构内容缓存
	search.Newmlmap()                           //将目录加载如内存map

	pubgo.Tj = pubgo.Newtongji()

	r := routers.NewEngine()
	r.Addrouter("/static/", routers.Static) //静态文件服务器
	r.Addrouter("/", routers.Index)         //首页
	r.Addrouter("/jingbu/", routers.Search) //经
	r.Addrouter("/lvbu/", routers.Search)   //律
	r.Addrouter("/lunbu/", routers.Search)  //论
	r.Addrouter("/ming/", routers.Ming)     //名
	r.Addrouter("/getonejuzi/", routers.Getonejuzi)
	r.Addrouter("/getjingwen/", routers.Getjingwen)
	r.Addrouter("/cidian/", routers.Cidian)
	r.Addrouter("/quanwen/", routers.Quanwen)
	r.Addrouter("/showdir/", routers.Showdir)
	r.Addrouter("/mulu/", routers.Mulu)
	r.Addrouter("/cipin/", routers.Cipin)
	r.Addrouter("/jing/", routers.Wen)
	r.Addrouter("/lv/", routers.Wen)
	r.Addrouter("/lun/", routers.Wen)
	r.Addrouter("/mapp/", routers.Mapp)
	r.Addrouter("/huchi/", routers.Huchi)
	r.Addrouter("/shuoming/", routers.Shuoming)
	r.Addrouter("/xianyan/", routers.Xianyan)
	r.Addrouter("/llan/", routers.Llan)
	r.Addrouter("/tongji/", routers.Tongji)
	r.Addrouter("/redir/", routers.Redir)
	r.Addrouter("/err/", routers.Err)
	r.Addrouter("/test/", routers.Test)
	r.Addrouter("/huiji/", routers.Huiji)
	r.Addrouter("/webSocket/", routers.Websocket)
	r.Addrouter("/chat/", routers.Chat)
	//r.Addrouter("/", r.test)
	//r.Addrouter("*", r.dir)
	log.Fatal(http.ListenAndServe(ip+":"+port, r))

}
