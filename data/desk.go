package data

import "gotiny/inter"

//玩家
type DeskPlayer struct {
	inter.Pid
	Userid   string `json:"userid"`    // 用户id
	Nickname string `json:"nickname"`  // 用户昵称
	Photo    string `json:"photo"`     // 头像
	Coin     uint32 `json:"coin"`      // 金币
	Diamond  uint32 `json:"diamond"`   // 钻石
	VipLevel int    `json:"vip_level"` // vip等级,充值时改变
	Robot    bool   `json:"robot"`     // 是否是机器人
	Online   bool   `json:"online"`    // 在线
}

func (this *DeskPlayer) GetUserid() string {
	return this.Userid
}

func (this *DeskPlayer) GetName() string {
	return this.Nickname
}

func (this *DeskPlayer) GetPhoto() string {
	return this.Photo
}

func (this *DeskPlayer) GetDiamond() uint32 {
	return this.Diamond
}

func (this *DeskPlayer) GetCoin() uint32 {
	return this.Coin
}

func (this *DeskPlayer) GetVipLevel() int {
	return this.VipLevel
}

//' 跑胡子

//低级别的操作等待执行
type LowLevel struct {
	Seat  uint32   `json:"seat"`  //操作玩家座位号
	Card  uint32   `json:"card"`  //操作牌
	Value uint32   `json:"value"` //掩码
	Cards []uint32 `json:"cards"` //吃
	Bione []uint32 `json:"bione"` //比
	Bitwo []uint32 `json:"bitwo"` //比
}

//操作
type PhzOperate struct {
	Userid string   `json:"userid"` //操作玩家
	Value  uint32   `json:"value"`  //掩码
	Cards  []uint32 `json:"cards"`  //吃
	Bione  []uint32 `json:"bione"`  //比
	Bitwo  []uint32 `json:"bitwo"`  //比
}

//操作出牌
type PhzDiscard struct {
	Userid string `json:"userid"` //操作玩家
	Card   uint32 `json:"card"`   //操作牌
}

//.
