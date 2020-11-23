package main

import (
	"fmt"
	"web_app/controller"
	"web_app/dao/mysql"
	"web_app/dao/redis"
	"web_app/logger"
	"web_app/pkg/snowflake"
	"web_app/routers"
	"web_app/settings"

	"go.uber.org/zap"
)

// Go Web开发较通用的脚手架模板

func main() {
	// 1.加载配置
	if err := settings.Init(); err != nil {
		fmt.Printf("init settins failed, e:%v\n", err)
		return
	}
	// 2.初始化日志
	if err := logger.Init(settings.Conf.LogConfig, settings.Conf.Mode); err != nil {
		fmt.Printf("init logger failed, e:%v\n", err)
		return
	}
	defer zap.L().Sync()
	// 3.初始化MySql连接
	if err := mysql.Init(settings.Conf.MySQLConfig); err != nil {
		fmt.Printf("init mysql failed, e:%v\n", err)
		return
	}
	defer mysql.Close()
	// 4.初始化Redis连接
	if err := redis.Init(settings.Conf.RedisConfig); err != nil {
		fmt.Printf("init redis failed, e:%v\n", err)
		return
	}
	defer redis.Close()
	// 初始化雪花算法生成UID
	if err := snowflake.Init(settings.Conf.StartTime, settings.Conf.MachineID); err != nil {
		fmt.Printf("init snowflake failed, err:%v\n", err)
		return
	}
	// 初始化gin框架内置的校验器使用的翻译器
	if err := controller.InitTrans("zh"); err != nil {
		fmt.Printf("init initTrans failed, err:%v\n", err)
		return
	}
	// 5.注册路由
	r := routers.SetupRouter(settings.Conf.Mode)
	err := r.Run(fmt.Sprintf(":%d", settings.Conf.Port))
	if err != nil {
		fmt.Printf("run server failed, err:%v\n", err)
		return
	}
}
