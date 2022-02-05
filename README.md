# sfs
乾隆大藏经搜索引擎

网站采用 golang+mysql5.7+ Bootstrap v5构架开发。

安装mysql5.7后，将sfs2.5.sql导入数据库。

然后在sfsgo文件夹的“设置.txt”，将下面的mysql配置内容更改为自己的，

root=root
psw=123456
dbname=sfs2.5
cache=5



配置说明


root，数据库用户名
psw，密码
dbname，数据库名称
cache，缓存，1g内存可以填写 10000


注意，ip不用填写。

操作系统对应的执行程序：
Windows -->winsfs.exe
Linux -->linux-webbook
Mac -->linuxsfs


以 Windows 为例子，
打开 winsfs.exe 后，在浏览器上输入：
http://127.0.0.1/ 

局域网用户可以通过服务器的局域网ip访问。
搭建互联网网站请自行搜索相关使用方法。

安装mysql5.7请自行百度搜索，需要一定的技术要求。
由于sfsgo.7z文件包含exe可执行文件，解压时杀毒软件可能错报，进行信任操作即可。


版权声明：该系统是“佛音传播工作室”自主开发的系统。可用于任何非营利性用途，以及二次开发使用。
