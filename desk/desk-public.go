package desk

import (
	"encoding/json"
	"gotiny/data"
	"net/http"
	"net/url"
	"utils"

	"github.com/golang/glog"
	"github.com/gorilla/websocket"
)

//' 机器人
// test
type ReqMsg struct {
	Code  string `json:code`  //邀请码
	Num   uint32 `json:num`   //数量
	Rtype uint32 `json:rtype` //房间类型
}

//召唤机器人
func client(Code string, Num, Rtype uint32) {
	glog.Infof("code: %s, num: %d, Rtype: %d", Code, Num, Rtype)
	var robot_host string = data.Conf.RobotHost
	var robot_port string = data.Conf.RobotPort
	var robot_addr string = robot_host + ":" + robot_port
	u := url.URL{Scheme: "ws", Host: robot_addr, Path: "/"}
	TimeStr := GmTime()
	Token := GmToke(TimeStr)
	c, _, err := websocket.DefaultDialer.Dial(u.String(),
		http.Header{"Token": {Token}})
	if err != nil {
		glog.Errorf("dial err -> %v\n", err)
	}
	if c != nil {
		reqMsg := &ReqMsg{
			Code:  Code,
			Num:   Num,
			Rtype: Rtype,
		}
		data, err1 := json.Marshal(reqMsg)
		//glog.Infof("req: %s, err1: %v", string(data), err1)
		if err1 == nil {
			c.WriteMessage(websocket.TextMessage, data)
		}
		c.Close()
	}
}

//字符串时间
func GmTime() string {
	Time := utils.Timestamp()
	TimeStr := utils.String(Time)
	return TimeStr
}

// Sign := utils.Md5(Key+Now)
// Token := Sign+Now+RandNum
func GmToke(TimeStr string) string {
	Sign := utils.Md5(data.Conf.GmKey + TimeStr)
	Token := Sign + TimeStr + utils.RandStr(6)
	return Token
}

//.

// vim: set foldmethod=marker foldmarker=//',//.:
