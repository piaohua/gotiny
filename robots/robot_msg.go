/**********************************************************
 * Author        : piaohua
 * Email         : 814004090@qq.com
 * Last modified : 2017-04-26 21:16:43
 * Filename      : robot_test.go
 * Description   : 机器人
 * *******************************************************/
package robots

type ReqMsg struct {
	Code  string `json:code`  //邀请码
	Num   uint32 `json:num`   //数量
	Rtype uint32 `json:rtype` //房间类型
}

//通知消息体
type Message struct {
	code  string
	rtype uint32
}

type Login struct {
	phone string
}

type ReLogin struct {
	phone string
	code  string
	rtype uint32
}

type Logout struct {
	phone string
	code  string
}
