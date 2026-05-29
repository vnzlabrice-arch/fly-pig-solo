package inits

func init() {
	NoCosInit()
	MysqlInit()
	RedisInit()
	InitMQ() // 初始化 RocketMQ 消息队列（用于后台派单）
}
