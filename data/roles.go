package data

import (
	"gotiny/inter"
	"time"
)

//获取消息
type Call struct {
	Msg     interface{}
	Pid     chan interface{}
	Timeout *time.Ticker
}

// 玩家列表

//广播消息
type Broadcast struct {
	Atype uint32      //分包类型
	Msg   interface{} //内容
}

//生成ID
type PlayersGenID struct {
}

//删除
type PlayersDel struct {
	Userid string //玩家ID
}

//设置添加
type PlayersSet struct {
	Userid string    //玩家ID
	Pid    inter.Pid //通道
}

//获取
type PlayersGet struct {
	Userid string
}

//在线
type PlayersOnline struct {
}

//.

// 房间
//p2d = player to desk
//p2s = player to roles
//p2r = player to rooms
//r2d = rooms to desk
//r2p = rooms to player
//d2r = desk to rooms
//d2p = desk to player

//.

// vim: set foldmethod=marker foldmarker=//',//.:
