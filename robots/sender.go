/**********************************************************
 * Author        : piaohua
 * Email         : 814004090@qq.com
 * Last modified : 2017-01-16 10:00
 * Filename      : sender.go
 * Description   : 机器人
 * *******************************************************/
package robots

import (
	"crypto/md5"
	"encoding/hex"
	"gotiny/protocol"
	"utils"

	"github.com/golang/glog"
	"github.com/golang/protobuf/proto"
)

//' 登录
// 发送注册请求
func (c *Robot) SendRegist() {
	ctos := &protocol.CRegist{}
	ctos.Phone = proto.String(c.data.Phone)
	ctos.Nickname = proto.String(c.data.Nickname)
	h := md5.New()
	h.Write([]byte("piaohua")) // 需要加密的字符串为 123456
	pwd := hex.EncodeToString(h.Sum(nil))
	ctos.Password = proto.String(pwd)
	c.Sender(ctos.GetCode(), ctos)
}

// 发送登录请求
func (c *Robot) SendLogin() {
	ctos := &protocol.CLogin{}
	ctos.Phone = proto.String(c.data.Phone)
	h := md5.New()
	h.Write([]byte("piaohua")) // 需要加密的字符串为 123456
	pwd := hex.EncodeToString(h.Sum(nil))
	ctos.Password = proto.String(pwd)
	c.Sender(ctos.GetCode(), ctos)
}

// 获取玩家数据
func (c *Robot) SendUserData() {
	ctos := &protocol.CUserData{}
	ctos.Userid = proto.String(c.data.Userid)
	c.Sender(ctos.GetCode(), ctos)
}

// 玩家创建房间
func (c *Robot) SendCreate() {
	ctos := &protocol.CCreateRoom{}
	ctos.Round = proto.Uint32(8)
	ctos.Rtype = proto.Uint32(1)
	ctos.Ante = proto.Uint32(1)
	ctos.Count = proto.Uint32(4)
	ctos.Payment = proto.Uint32(0)
	ctos.Rname = proto.String("ddd")
	utils.Sleep(3)
	c.Sender(ctos.GetCode(), ctos)
}

// 玩家进入房间
func (c *Robot) SendEntry() {
	if c.code == "create" { //表示创建房间
		c.SendCreate() //创建一个房间
	} else { //表示进入房间
		ctos := &protocol.CEnterRoom{}
		ctos.Invitecode = proto.String(c.code)
		c.Sender(ctos.GetCode(), ctos)
	}
}

// 玩家准备
func (c *Robot) SendReady() {
	ctos := &protocol.CReady{}
	ctos.Ready = proto.Bool(true)
	utils.Sleep(5)
	c.Sender(ctos.GetCode(), ctos)
}

// 离开
func (c *Robot) SendLeave() {
	ctos := &protocol.CLeave{}
	c.Sender(ctos.GetCode(), ctos)
}

// 解散
func (c *Robot) SendVote() {
	ctos := &protocol.CVote{}
	ctos.Vote = proto.Uint32(0)
	c.Sender(ctos.GetCode(), ctos)
}

//.

//' 牛牛
// 抢庄
func (c *Robot) SendDealer() {
	ctos := &protocol.CDealer{}
	ctos.Dealer = proto.Bool(true)
	num := uint32(utils.RandInt32N(3) + 1) //随机
	ctos.Num = proto.Uint32(num)
	utils.Sleep(3)
	c.Sender(ctos.GetCode(), ctos)
}

// 下注
func (c *Robot) SendBet() {
	ctos := &protocol.CBet{}
	ctos.Seatbet = proto.Uint32(c.seat)
	val := uint32(utils.RandInt32N(3)) + 1
	ctos.Value = proto.Uint32(val)
	utils.Sleep(3)
	c.Sender(ctos.GetCode(), ctos)
}

// 提交组合
func (c *Robot) SendNiu() {
	//TODO
}

// 提交组合
func (c *Robot) SendRecord() {
	ctos := &protocol.CGameRecord{
		Page: proto.Uint32(0),
	}
	utils.Sleep(3)
	c.Sender(ctos.GetCode(), ctos)
}

//.

//' free

// 进入房间
func (c *Robot) SendEntryFree() {
	ctos := &protocol.CEnterFreeRoom{}
	utils.Sleep(2)
	c.Sender(ctos.GetCode(), ctos)
}

// 玩家入坐
func (c *Robot) SendFreeSit() {
	seat := uint32(utils.RandInt32N(7) + 1) //随机
	ctos := &protocol.CFreeSit{
		State: proto.Bool(true),
		Seat:  proto.Uint32(seat),
	}
	//c.sits++ //尝试次数
	utils.Sleep(2)
	c.Sender(ctos.GetCode(), ctos)
}

// 玩家离坐
func (c *Robot) SendFreeStandup() {
	ctos := &protocol.CFreeSit{
		State: proto.Bool(false),
		Seat:  proto.Uint32(c.seat),
	}
	utils.Sleep(2)
	c.Sender(ctos.GetCode(), ctos)
	utils.Sleep(2)
	c.SendLeave()
	utils.Sleep(2)
	c.Close() //下线
}

// 玩家下注
func (c *Robot) SendFreeBet() {
	var a1 []uint32 = []uint32{2, 3, 4, 5}
	var c1 []uint32 = []uint32{100, 1000, 10000, 50000, 100000, 200000}
	var coin uint32 = c.data.Coin / 4
	var n1 int32 = utils.RandInt32N(15) //随机
	if coin > 10000000 {
		n1 += 10
	}
	for {
		if n1 <= 0 {
			break
		}
		if coin <= 0 {
			break
		}
		var i2 int
		for i := 5; i >= 0; i-- {
			if coin >= c1[i] {
				i2 = i
				break
			}
		}
		var val int
		switch i2 {
		case 0, 4, 5:
			val = i2
		default:
			val = int(utils.RandInt32N(int32(i2))) + 1 //随机
		}
		var i1 int32 = utils.RandInt32N(4) //随机
		ctos := &protocol.CFreeBet{
			Value: proto.Uint32(c1[val]),
			Seat:  proto.Uint32(a1[i1]),
		}
		//c.bits++
		utils.Sleep(1)
		c.Sender(ctos.GetCode(), ctos)
		//
		n1--
		coin -= c1[val]
	}
}

// 玩家聊天
func (c *Robot) SendChat() {
	if c.rtype == 0 {
		return
	}
	if utils.RandInt32N(10) > 4 {
		return
	}
	content := []byte(chat[utils.RandInt32N(int32(len(chat)))])
	ctos := &protocol.CChatText{
		Content: content,
	}
	utils.Sleep(3)
	c.Sender(ctos.GetCode(), ctos)
}

var chat []string = []string{
	"收金币，200W金币150，微伈：13699755455",
	"卖金币，200W金币160，微伈：13699755455",
	"不要走，决战到天亮",
	"又赢了",
	"哈哈",
	"哈哈哈",
	"nnd",
	"hahahaha",
	"搞什么飞机",
	"这运气",
	"这都行",
	"也是没谁了",
	"怎么搞的，又输，没金币了",
	"输惨了，有没有土豪送点金币",
	"还好，只输了一点",
	"这次一定是要赢",
	"连输2把，下次一定要赢回来",
	"又赢了，手气不错",
	"哈哈，这样都能赢",
	"连胜4把，快超神了",
	"今天要连赢19把",
	"赢得豪爽",
	"这把牌简直逆天了",
	"这游戏赢得豪爽啊",
	"哈哈哈哈",
	"你们继续，我先撤了",
	"晚上继续搞起",
	"我靠，这牌，这么大！",
	"我靠，好大！",
	"牌小也能赢，牛逼",
}

//.

//' classic

// 进入房间
func (c *Robot) SendEntryClassic() {
	ctos := &protocol.CEnterClassicRoom{}
	for _, v := range c.classic {
		min := v.GetMinimum()
		max := v.GetMaximum()
		if c.data.Coin > min && (c.data.Coin < max || max == 0) {
			c.classicId = v.GetId()
			break
		}
	}
	ctos.Id = proto.String(c.classicId)
	utils.Sleep(3)
	c.Sender(ctos.GetCode(), ctos)
}

// 进入房间
func (c *Robot) SendGetClassic() {
	ctos := &protocol.CClassicList{}
	c.Sender(ctos.GetCode(), ctos)
}

// 玩家下注
func (c *Robot) SendClassicBet() {
	val := uint32(utils.RandInt32N(3)) + 1 //随机
	ctos := &protocol.CBet{
		Value:   proto.Uint32(val),
		Seatbet: proto.Uint32(c.seat),
	}
	//c.bits++
	utils.Sleep(1)
	c.Sender(ctos.GetCode(), ctos)
}

func (c *Robot) SendClassicClose() {
	c.classicId = ""
	c.SendLeave()
	c.Close() //下线
}

//.

//' 跑胡子

// 玩家进入房间
func (c *Robot) SendEntryPhz() {
	if c.code == "create" { //表示创建房间
		c.SendCreatePhz() //创建一个房间
	} else { //表示进入房间
		glog.Infof("enter phz room -> %s, %d", c.code, c.rtype)
		ctos := &protocol.CEnterZiRoom{}
		ctos.Invitecode = proto.String(c.code)
		c.Sender(ctos.GetCode(), ctos)
	}
}

// 玩家创建房间
func (c *Robot) SendCreatePhz() {
	ctos := &protocol.CCreateZiRoom{}
	ctos.Round = proto.Uint32(8)
	ctos.Rtype = proto.Uint32(6)
	ctos.Ante = proto.Uint32(1)
	ctos.Count = proto.Uint32(3)
	ctos.Payment = proto.Uint32(0)
	ctos.Rname = proto.String("phz")
	ctos.Xi = proto.Uint32(10)
	utils.Sleep(3)
	c.Sender(ctos.GetCode(), ctos)
}

// 玩家操作(吃碰胡)
func (c *Robot) SendOperatePhz(value uint32, cards, bione, bitwo []uint32) {
	ctos := &protocol.COperate{}
	ctos.Value = proto.Uint32(value) //掩码值
	ctos.Cards = cards               //吃牌
	ctos.Bione = bione               //比牌
	ctos.Bitwo = bitwo               //比牌
	utils.Sleep(3)
	c.Sender(ctos.GetCode(), ctos)
}

// 玩家出牌
func (c *Robot) SendDiscardPhz(card uint32) {
	ctos := &protocol.CPushDiscard{}
	ctos.Card = proto.Uint32(card)
	utils.Sleep(3)
	c.Sender(ctos.GetCode(), ctos)
}

func (c *Robot) SendPhzClose() {
	c.classicId = ""
	c.SendLeave()
	c.Close() //下线
}

//胡碰吃
func (c *Robot) phzOperate(card, value uint32) {
	//TODO
}

//自动操作(提,跑,偎)
func (c *Robot) phzAuto(card, value uint32) {
	//TODO
}

//移除吃碰组合
func (c *Robot) phzChow(card, value uint32, cards, bione, bitwo []uint32) {
	//TODO
}

//出牌
func (c *Robot) phzDiscard() {
	//TODO 选择最优出牌
	c.SendDiscardPhz(c.cards[0])
}

//.

// vim: set foldmethod=marker foldmarker=//',//.:
