package pkg

import (
	"fmt"
	"strconv"

	"github.com/smartwalle/alipay/v3"
)

func Alipay(OutTradeNo string, TotalAmount float64) string {
	var privateKey = "MIIEowIBAAKCAQEA1XmwRXAvQOEkEShHDPl0TfW0H84ggVjlJErE6lbBtqCLSntT26s5cJ0IwnjOfdC4A8bEARhugbNZVk91wMxxC7kpN9Wlp/EOLo9ejr3DqX4HV+NVcAZnMRepoixOmbSTS8kcfHdzHyWYWy0jP6FhtRrFkvpdzRJqaT9vbCwDmtKRW8G4FZz3sqrA7zemMmwG2fIeYmkb1/RH6ydjZC1IYIZfVpqx/Sbs5p7rcCfiBkT8mf0qovs24MXY7ONy/09wcImWlLmeTsEJhCZSJBpBpNw6v7paEdJZsJUFwNXE2ncbw1ZpU2ugAg3ssN8A1u+OzFPG4GGPUzZBAqwFcFeFMQIDAQABAoIBAQCpvGkQvSOFVb5Euh7MTShRuHRw+YdU8Wa4lY6+uR6rZDl8jlTeC+XPCb1mS4H7/dHihs9gA5oSHJWDEAkTtDqKVL5WO9TMlNlu4eYQXwhNIv0Zv5M2UJEKMsNZRBt3EWIw1waEXMG3WUAr9tbZCCtVQaUouVON8/+SEbM5ckGnfiLSdqc0W+qS9OFdt5YV0bVycJYL3fhg5KLSS629xw8p7fpItNZr47QPOKKF9f+Lv9OpW6JHy92seKsJtsHoikgI4biPdZqyaJaCvsAfUiKGfIIIxduh10jAA4wbJz18GKDh/3JwQrYsbQnbWaaWnuBXQsd/G4sbPSsGVULbklUVAoGBAP0V3BojNMSLo4xh6Y2bPP0ITa1rJ70RVRf8XC7vz5ISj+BdIAH780veL7ZhowLdtZl9DbQ9+KSNECsm2Whne97Cg3Q8OlqhPg+83IIRZ9Ay+7W8JuzZUugKjS8XQahtXzQGC/j4Sh1T9atJSiR/A4aFGm9uVoUOpKOeHEWIk2L/AoGBANfvDSlcl1uf4pvAdQfil0vJd6JkHhcg6el5WzLB5UaQPZAuPez2u1/Q8O9zFL/xoibmrqAaGWnobDhuE7Eig2UEU/SQB8ZlndqZr8c1xKLLqa47Ff3YNTB1nRzbkHt2Uu1uJPAAUUP1OxIqlPlY32+oOyuaucRXJchuY/2WCofPAoGAVHekQzD23pZAWo0fNvzVyRmpWzTFSYvsW5oFZkDxhS4eyOIr2Z2uYObiA7vofP9kKbscBMkeIEVYXR6VNww6wmknGHc1fqQMI5KEAgEvJcSuT8RhNXF/AyTqZAoeBsmiGane+xRbIBiyiJ1oWm/tzErGPOSVickOV/FgPDDOCB8CgYBwOAjkBP2YHWh8PzHS7jxoo4Qr/dHxfSZzMqFlqITJ/i4wXwfJvZQ1QHXmSy1ub1ow01PPWqD4fFS5ouNS/DfC6NPk1nFj9u1pbNOAOP4/CI7fQE+2g4Vo8Ma895KHxz9jqwlBPTj+k1SmpUCUsU41Tf2qNJf0ZMH6/vuyDybb2wKBgCCMht8bajPEFukrsCB2LfnorNC1qSTTh2lZM/oQIRDuu3MzPY0L7erfbop2J5MmE/DoGhSZYyGfiRmI/zMknPXKbn/K6PCBjwLWhpUuMGZDLmjbgGj3XGjdVqjCdvY29FzA1mvnZ/3aO0WkdDkLrZ/VnUj4dZqOZhkRv/5iT18t" // 必须，上一步中使用 RSA签名验签工具 生成的私钥
	client, err := alipay.New("9021000158650697", privateKey, false)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	var p = alipay.TradePagePay{}
	p.NotifyURL = "https://25562ce5.r9.cpolar.cn/collback"
	p.ReturnURL = "https://25562ce5.r9.cpolar.cn/notify"
	p.Subject = "网约车订单支付宝支付"
	p.OutTradeNo = OutTradeNo
	p.TotalAmount = strconv.FormatFloat(TotalAmount, 'f', 2, 64)
	p.ProductCode = "FAST_INSTANT_TRADE_PAY"

	url, err := client.TradePagePay(p)
	if err != nil {
		fmt.Println(err)
	}

	// 这个 payURL 即是用于打开支付宝支付页面的 URL，可将输出的内容复制，到浏览器中访问该 URL 即可打开支付页面。
	var payURL = url.String()
	fmt.Println(payURL)
	return payURL
}
