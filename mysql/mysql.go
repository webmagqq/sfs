package mysql

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

var MysqlDb *sql.DB
var MysqlDbErr error

// 初始化mysql链接
func NewsqlDb(root, psw, dbname string) {

	// 打开连接
	MysqlDb, MysqlDbErr = sql.Open("mysql", root+":"+psw+"@(localhost:3306)/"+dbname+"?charset=utf8mb4") // "root:123456@(localhost:3306)/sfs2.5?charset=utf8mb4"
	if MysqlDbErr != nil {
		//log.Println("dbDSN: " + dbDSN)
		panic("数据源配置不正确: " + MysqlDbErr.Error())
	}
	MysqlDb.SetMaxOpenConns(15) // 最大连接数
	MysqlDb.SetMaxIdleConns(5)  // 闲置连接数
	//MysqlDb.SetConnMaxLifetime(100 * time.Second) // 最大连接周期
	if MysqlDbErr = MysqlDb.Ping(); nil != MysqlDbErr {
		panic("数据库链接失败: " + MysqlDbErr.Error())
	}
}
func Selects(query string, args ...interface{}) []map[string]string {

	//sqlStr := "SELECT * FROM jingmulu limit 30" //可以换成其它的查询语句,可以得到相应的查询结果,不用每次都去构建存放的结构体
	rows, err := MysqlDb.Query(query, args...)
	//rows, err := MysqlDb.Query(query)
	res := make([]map[string]string, 0)
	if err != nil {
		fmt.Println(err)
		return res
	}
	defer rows.Close()

	//列出所有查询结果的字段名
	cols, _ := rows.Columns()

	//values是每个列的值，这里获取到byte里
	values := make([][]byte, len(cols))     //建立接口  --[]byte方便转换为字符串string(v)
	scans := make([]interface{}, len(cols)) //建立接口指针的接口
	for i := range values {
		scans[i] = &values[i] //将接口转换为指针类型的接口,这样才能接收rows.Scan返回值
	}

	//遍历rows
	for rows.Next() {
		_ = rows.Scan(scans...) //让每一行数据都填充到[][]byte里面
		row := make(map[string]string)
		for k, v := range values { //每行数据是放在values里面，现在把它挪到row里
			key := cols[k]
			row[key] = string(v)
		}
		res = append(res, row)
	}
	//fmt.Println(res)
	return res
}

func Exesql(query string, args ...interface{}) int64 {
	ret, err := MysqlDb.Exec(query, args...)
	if err != nil {
		fmt.Println(err.Error())
		return -1
	}
	rowsaffected, _ := ret.RowsAffected() //影响行数
	return rowsaffected
}

// 插入数据
func Inserts(query string, args ...interface{}) (rowsed int64, lastID int64) {
	ret, err := MysqlDb.Exec(query, args...)
	if err != nil {
		fmt.Println(err.Error())
		return -1, -1
	}
	rowsaffected, _ := ret.RowsAffected() //影响行数
	lastInsertID, _ := ret.LastInsertId() //插入数据的最后主键id   insert 操作有效
	return rowsaffected, lastInsertID
}
