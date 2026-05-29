package init

func init() {
	NoCosInit()
	MysqlInit()
	RedisInit()
	InitServiceData()
	InitCouponData()
}
