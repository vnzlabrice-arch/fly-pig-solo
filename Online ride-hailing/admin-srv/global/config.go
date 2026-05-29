package global

type AppConfig struct {
	Mysql
	Redis
	RocketMQ // 使用 RocketMQ 替代 RabbitMQ（用于后台派单消息队列）
}

type Mysql struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
}

type Redis struct {
	Host     string
	Port     int
	Password string
	Database int
}

type NaCos struct {
	Host      string
	Port      int
	PublicId  string
	NaCosName string
	Password  string
	Group     string
	DataId    string
}

// RocketMQ 配置（用于后台派单消息队列）
type RocketMQ struct {
	NameServers []string `mapstructure:"name_servers"` // NameServer地址列表，如 ["127.0.0.1:9876"]
	GroupName   string   `mapstructure:"group_name"`   // 生产者/消费者组名
	Topic       string   `mapstructure:"topic"`        // 主题：派单消息
}

type RabbitMQ struct {
	MQURL string
}
