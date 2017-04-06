package mysql

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	. "nuobelcloud/conf"
	. "nuobelcloud/tools"
)

/*

mariadb连接

*/

/*
//建库和建表语句
CREATE DATABASE `nuobel` charset=utf8;
//用户表
create table `user`(
`uid` int(11) PRIMARY KEY auto_increment,
`username` varchar(32) not null,
`nickname` varchar(64) not null,
`password` varchar(32) not null,
`email` varchar(50) not null,
`phone` varchar(16) not null,
`sex` enum("other","man","woman") not null,
`volume(M)` int(11) not null default 200,
`dir1` varchar(32) not null default "",
`createtime` datetime not null,
`changetime` timestamp
)auto_increment=1000000 charset=utf8;
//增加用户
grant select,update,delete,insert on nuobel.* to nuobel@'%' identified by 'aabbccdd123';
*/
var DB *sql.DB

func mysqlConnectInit() (*sql.DB, error) {
	//mysql open后实际是没有建立连接，是一个连接池，执行query等命令时才连接.

	db, err := sql.Open("mysql", MysqlAddr)
	//最大连接数，默认是0不限制
	db.SetMaxOpenConns(2000)
	//闲时连接数
	db.SetMaxIdleConns(1000)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}

func init() {
	var err error
	DB, err = mysqlConnectInit()
	if err != nil {
		panic(err)
	}
	Log(INFO, "mysql连接成功.")
}
