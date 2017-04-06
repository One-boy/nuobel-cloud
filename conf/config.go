package conf

/*

一些配置信息

*/

const (
	//redis连接地址,redis-648392706.tenxcloud.net
	RedisAddr = "localhost:5501"
	// RedisAddr = "redis-648392706.tenxcloud.net:25302"
	// RedisAddr = "182.92.64.152:25302"
	//mariadb连接地址，格式：user:password@tcp(host:port)/DBname?param
	MysqlAddr = "nuobel:aabbccdd123@tcp(localhost:5401)/nuobel?charset=utf8"
	// MysqlAddr = "nuobel:aabbccdd123@tcp(mysql-648392706.tenxcloud.net:48763)/nuobel?charset=utf8"
	// MysqlAddr = "nuobel:aabbccdd123@tcp(182.92.64.152:48763)/nuobel?charset=utf8"
	//md5 hash 盐
	Md5Salt = "flLss"
	//cookie名字
	CookieName = "nuobelId"
	//cookie有效时间
	CookieExpireTime = 7200
	//登陆跳转地址
	LoginUrl = "/nuobel/login"
	//文件存储根文件路径
	FileRootPath = "file"
)
