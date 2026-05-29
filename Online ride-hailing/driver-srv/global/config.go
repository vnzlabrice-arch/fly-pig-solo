package global

type configData struct {
	Mysql
	Redis
	RabbitMQ
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
type RabbitMQ struct {
	MQURL string
}
