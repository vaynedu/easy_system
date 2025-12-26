package main

import (
	"log"

	"github.com/vaynedu/exam_system/config"
	"github.com/vaynedu/exam_system/router"
)

func main() {
	// 1. 初始化数据库连接
	config.InitDB()

	// 2. 初始化路由
	r := router.InitRouter()

	// 3. 启动服务（端口8080）
	log.Println("服务启动成功：http://127.0.0.1:8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("服务启动失败：", err)
	}
}
