package data

import "gotiny/inter"

//房间列表类型
const (
	ROOMS_PRIVATE = 1 //私人房间列表
	ROOMS_FREE    = 2 //百人自由场
	ROOMS_CLASSIC = 3 //经典场
	ROOMS_PAOHUZI = 4 //跑胡子私人场
	ROOMS_COIN    = 5 //跑胡子金币场
)

//' 房间列表

//添加
type RoomsAdd struct {
	Atype int       //房间列表类型
	Key   string    //roomid或invitecode
	Pid   inter.Pid //通道
}

//存在
type RoomsExist struct {
	Atype int    //房间列表类型
	Key   string //roomid或invitecode
}

//移除
type RoomsDel struct {
	Atype int    //房间列表类型
	Key   string //roomid或invitecode
}

//获取
type RoomsMatch struct {
	Atype   int       //房间列表类型
	Rtype   uint32    //玩法类型
	Ltype   uint32    //房间级别类型
	Id      string    //匹配
	Key     string    //roomid或invitecode
	Diamond uint32    //进入时钻石额度
	Pid     inter.Pid //玩家通道
}

//生成随机码
type RoomsGenID struct {
	Atype int //房间列表类型
}

//生成随机码
type GenInvitecode struct {
	Atype int //房间列表类型
}

//创建或匹配
type RoomsCreate struct {
	Atype int       //房间列表类型
	Id    string    //匹配
	Data  *DeskData //牌桌数据
	Pid   inter.Pid //玩家通道
}

//.

//' 房间
//p2d = player to desk
//p2s = player to roles
//p2r = player to rooms
//r2d = rooms to desk
//r2p = rooms to player
//d2r = desk to rooms
//d2p = desk to player

//获取房间数据
type P2DGetData struct {
}

//进入房间
type P2DEntry struct {
	User *User     //玩家数据
	Pid  inter.Pid //玩家通道
}

//离开房间
type P2DLeave struct {
	Userid string    //玩家ID
	Pid    inter.Pid //玩家通道
}

//打印输出
type P2DPrint struct {
}

//召唤机器人
type P2DRobot struct {
}

//踢除玩家
type P2DKick struct {
	Userid string    //玩家ID
	Seat   uint32    //座位位置
	Vip    uint32    //vip level
	Pid    inter.Pid //玩家通道
}

//广播聊天
type P2DChat struct {
	Mtype  int    //类型
	Userid string //玩家ID
	Msg    []byte //内容
}

//玩家准备
type P2DReady struct {
	Userid string    //玩家ID
	Ready  bool      //准备
	Pid    inter.Pid //玩家通道
}

//投票解散
type P2DVote struct {
	Userid string    //玩家ID
	Vtype  bool      //状态true发起 false投票
	Vote   uint32    //投票
	Pid    inter.Pid //玩家通道
}

//phz
//出牌
type P2DDesert struct {
	Userid string    //玩家ID
	Card   uint32    //牌
	Pid    inter.Pid //玩家通道
}

//操作
type P2DOperate struct {
	Userid string    //玩家ID
	Value  uint32    //操作值
	Cards  []uint32  //
	Bione  []uint32  //
	Bitwo  []uint32  //
	Pid    inter.Pid //玩家通道
}

//刮奖
type P2DPrize struct {
	Userid string    //玩家ID
	Card   uint32    //牌
	Pid    inter.Pid //玩家通道
}

//设置手牌
type P2DSetHands struct {
	Userid string   //玩家ID
	Round  uint32   //局数
	Cards  []uint32 //手牌
}

//获取手牌
type P2DGetHands struct {
	Userid string //玩家ID
}

//抢庄 (看牌抢庄)
type P2DChoiceDealer struct {
	Userid   string    //玩家ID
	IsDealer bool      //是否抢庄
	Num      uint32    //倍数
	Pid      inter.Pid //玩家通道
}

//下注
type P2DChoiceBet struct {
	Userid  string    //玩家ID
	SeatBet uint32    //下注位置
	Num     uint32    //下注数量
	Pid     inter.Pid //玩家通道
}

//提交组合
type P2DChoiceNiu struct {
	Userid string    //玩家ID
	Num    uint32    //牛掩码
	Cards  []uint32  //组合
	Pid    inter.Pid //玩家通道
}

//玩家坐下
type P2DSitDown struct {
	Userid string    //玩家ID
	Seat   uint32    //位置
	Situp  bool      //坐下或站起
	Pid    inter.Pid //玩家通道
}

//抢庄,没人上庄时都可以选择上庄,可以多次上庄，已经上庄的人可以补庄
type P2DBeDealer struct {
	Userid string    //玩家ID
	State  uint32    //状态0下庄 1上庄 2补庄
	Num    uint32    //数量
	Pid    inter.Pid //玩家通道
}

//上庄列表
type P2DBeDealerList struct {
	Pid inter.Pid //玩家通道
}

//房间趋势
type P2DTrend struct {
	Pid inter.Pid //玩家通道
}

//进入房间
type D2PSetRoom struct {
	Pid inter.Pid //玩家通道
}

//超时解散
type D2DDismiss struct {
}

//.
