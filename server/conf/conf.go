package conf

const (
	AppName  string = "im"
	HttpPort string = "2580"
	RunMode  string = "release" //debug test release
	//mysql
	DbHost     string = "localhost"
	DbPort     string = "3306"
	DbName     string = "im"
	DbUser     string = "root"
	DbPassword string = "root"
	//redis
	RedisAddr     string = "localhost:6379"
	RedisPassword string = ""
	RedisDb       int    = 1 //默认数据库
	RedisPrefix   string = AppName + "."
	//JWT
	JwtSecret string = "hjJlbE8LW7NAO15v"
	//token超时时间(秒)
	TokenExp int64 = 60 * 60 * 24 * 7
	//im
	ImUserPrefix  string = AppName + "_"
	ImMsgQueueKey string = AppName + ".msg_list"
	ImLineUserKey string = AppName + ".user_line"
	ImOfflineMsg  string = AppName + ".offline_"
)
