package data

import "gotiny/inter"

// 玩家

//存储
type PlayerSave struct {
}

//添加奖励
type AddPrize struct {
	Diamond int32 //变动钻石额度
	Coin    int32 //变动金币额度
	Ltype   int   //变动日志类型
}

//vip赠送
type VipGive struct {
	Diamond uint32 //变动钻石额度
	Money   uint32 //变动币额度
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

//进入私人房间的响应
type R2PEntry struct {
	Atype int       //房间列表类型
	Id    string    //匹配
	Pid   inter.Pid //房间通道,nil时表示进入或匹配失败
}

//创建或匹配的响应
type R2PCreate struct {
	Atype int       //房间列表类型
	Data  *DeskData //返回时更新了roomid和invitecode
}

//获取是否在房间
type R2PGetRoom struct {
}

// vim: set foldmethod=marker foldmarker=//',//.:
