package role

import (
	"gotiny/data"
	"gotiny/inter"
	"gotiny/roles"
	"gotiny/rooms"

	"github.com/golang/glog"
)

//roles call

//生成用户ID
func genUserid() string {
	arg := new(data.PlayersGenID)
	msg := roles.Players.Call(arg)
	switch msg.(type) {
	case string:
		return msg.(string)
	default:
		glog.Errorf("unknown message %v", msg)
	}
	return ""
}

//获取在线用户
func getPlayer(userid string) inter.Pid {
	arg := new(data.PlayersGet)
	arg.Userid = userid
	msg := roles.Players.Call(arg)
	switch msg.(type) {
	case inter.Pid:
		//已经在线,重复登录检测
		return msg.(inter.Pid)
	default:
		glog.Errorf("unknown message %v", msg)
	}
	return nil
}

//添加新用户
func addPlayer(userid string, newp inter.Pid) {
	msg2 := new(data.PlayersSet)
	msg2.Userid = userid
	msg2.Pid = newp
	roles.Players.Send(msg2) //放入在线列表
}

//rooms call

//获取邀请码
func getCode(atype int) string {
	arg := new(data.GenInvitecode)
	arg.Atype = atype
	msg := rooms.Rooms.Call(arg)
	switch msg.(type) {
	case string:
		return msg.(string)
	default:
		glog.Errorf("unknown message %v", msg)
	}
	return ""
}

//生成房间ID
func genRoomid() string {
	arg := new(data.RoomsGenID)
	msg := rooms.Rooms.Call(arg)
	switch msg.(type) {
	case string:
		return msg.(string)
	default:
		glog.Errorf("unknown message %v", msg)
	}
	return ""
}

//添加新房间
func addRoom(atype int, code string, newr inter.Pid) {
	msg := new(data.RoomsAdd)
	msg.Atype = atype
	msg.Key = code
	msg.Pid = newr
	rooms.Rooms.Send(msg)
}
