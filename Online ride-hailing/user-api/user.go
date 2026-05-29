package main

import (
	"flag"
	"fmt"
	"time"

	"user-api/internal/config"
	"user-api/internal/handler"
	"user-api/internal/job"
	"user-api/internal/middleware"
	"user-api/internal/svc"
	"user-api/internal/ws"
	_ "user-srv/init"

	"github.com/robfig/cron/v3"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/user-api.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	hub := ws.NewHub()
	go hub.Run()

	server := rest.MustNewServer(c.RestConf, rest.WithCors("*"))
	server.Use(middleware.CorsHandler())
	defer server.Stop()

	ctx := svc.NewServiceContext(c, hub)
	handler.RegisterHandlers(server, ctx)

	server.AddRoute(rest.Route{
		Method:  "GET",
		Path:    "/ws",
		Handler: ws.ServeWebSocket(hub),
	})

	go startCronJob(ctx, hub)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	fmt.Println("WebSocket endpoint: /ws?user_id={user_id}")
	server.Start()
}

// startCronJob 启动定时任务
func startCronJob(ctx *svc.ServiceContext, hub *ws.Hub) {
	time.Sleep(3 * time.Second)

	paymentReminderJob := job.NewPaymentReminderJob(ctx, hub)

	c := cron.New()
	_, err := c.AddFunc("*/1 * * * *", func() {
		logx.Info("定时任务触发：检查未支付订单")
		paymentReminderJob.Run()
	})
	_, err = c.AddFunc("*/1 * * * *", func() {
		logx.Info("定时任务触发：超时取消订单")
		paymentReminderJob.TimeoutCancelOrders()
	})
	if err != nil {
		logx.Errorf("添加定时任务失败: %v", err)
		return
	}

	logx.Info("定时任务已启动：每分钟检查未支付订单和超时订单")
	c.Start()

	select {}
}
