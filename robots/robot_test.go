/**********************************************************
 * Author        : piaohua
 * Email         : 814004090@qq.com
 * Last modified : 2017-04-26 21:16:43
 * Filename      : robot_test.go
 * Description   : 机器人
 * *******************************************************/
package robots

import (
	"encoding/json"
	"gotiny/data"
	"net/http"
	"net/url"
	"testing"
	"utils"

	"github.com/golang/glog"
	"github.com/gorilla/websocket"
)

func TestRobot(t *testing.T) {
	client()
}

//召唤机器人
func client(Code string, Num uint32) {
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
			Code: Code,
			Num:  Num,
		}
		data, err1 := json.Marshal(reqMsg)
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
