package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"miniIM/server/app/api/middleware"
	"miniIM/server/app/task"
	"miniIM/server/app/ws"
	"miniIM/server/conf"
	"miniIM/server/db"
	"miniIM/server/route"
	"time"
)

func main() {
	fmt.Println("应用启动", time.Now().Format("2006-01-02 15:04:05"))
	gin.SetMode(conf.RunMode)
	go ws.Manager.Start() //开启一个新的并发运行
	r := gin.Default()
	r.Use(middleware.CorsMiddleware) //允许跨域中间件
	route.Ws(r)                      //websocket
	route.Web(r)                     //渲染静态资源
	route.Api(r)                     //api路由
	db.DbInit()                      //数据库连接
	db.RedisInit()                   //redis连接
	task.Start()                     //启动定时任务
	r.Run(":" + conf.HttpPort)
}
