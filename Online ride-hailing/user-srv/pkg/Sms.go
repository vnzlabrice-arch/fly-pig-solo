package pkg

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

const strUrl = "https://api.ihuyi.com/sms/Submit.json"

func SmsSend(mobile string, code string) T {

	v := url.Values{}
	v.Set("account", "C28419428")                         //APIID（用户中心【文本短信】-【验证码/通知短信】-【产品总览】查看）
	v.Set("password", "510a45300222410494404cd3e9d131cb") //1、APIKEY（用户中心【文本短信】-【验证码/通知短信】-【产品总览】查看）2、动态密码（生成动态密码方式请看该文档末尾的说明）
	v.Set("mobile", mobile)                               //根据发送方式不同：1、完整内容方式提交完整的短信内容，如：您的验证码是：1234。请不要把验证码泄露给其他人。2、模板变量方式模板中的变量内容，多个变量以英文竖线（|）隔开①单变量示例模板内容：您的验证码是：【变量】。请不要把验证码泄露给其他人。参数写法：content=1234最终短信为：您的验证码是：1234。请不要把验证码泄露给其他人。②多变量示例模板内容：订单号：【变量1】，联系人：【变量2】，手机号：【变量3】，金额：【变量4】。参数写法：content=20180515006|张三|136xxxxxxxx|100元最终短信为：订单号：20180515006，联系人：张三，手机号：136xxxxxxxx，金额：100元。支持500字以内的长短信，长短信按多条计费
	v.Set("content", "您的验证码是："+code+"。请不要把验证码泄露给其他人。")    //短信内容，注：模板ID为空时必填

	body := strings.NewReader(v.Encode()) //把form数据编码
	client := &http.Client{}
	req, _ := http.NewRequest("POST", strUrl, body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req) //发送
	if err != nil {
		fmt.Println(err)
	}

	defer resp.Body.Close() //一定要关闭resp.Body
	res, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(res))
	var t T
	json.Unmarshal(res, &t)
	return t
}

type T struct {
	Code  int    `json:"code"`
	Msg   string `json:"msg"`
	Smsid string `json:"smsid"`
}
