/**********************************************************
 * Author        : piaohua
 * Email         : 814004090@qq.com
 * Last modified : 2017-01-16 10:00
 * Filename      : recv.go
 * Description   : 机器人
 * *******************************************************/
package robots

import (
	"gotiny/data"
	"gotiny/protocol"

	"github.com/golang/glog"
)

//' 初始化协议注册
func init() {
	p1 := protocol.SLogin{}
	regist(p1.GetCode(), p1, recvLogin)
	p2 := protocol.SEnterRoom{}
	regist(p2.GetCode(), p2, recvComein)
	p3 := protocol.SLaunchVote{}
	regist(p3.GetCode(), p3, recvVote)
	p4 := protocol.SCreateRoom{}
	regist(p4.GetCode(), p4, recvCreate)
	p5 := protocol.SPushDealer{}
	regist(p5.GetCode(), p5, recvDealer)
	p6 := protocol.SDraw{}
	regist(p6.GetCode(), p6, recvDraw)
	p7 := protocol.SGameover{}
	regist(p7.GetCode(), p7, recvGameover)
	p8 := protocol.SLeave{}
	regist(p8.GetCode(), p8, recvLeave)
	p11 := protocol.SUserData{}
	regist(p11.GetCode(), p11, recvdata)
	p12 := protocol.SKick{}
	regist(p12.GetCode(), p12, recvKick)
	p19 := protocol.SRegist{}
	regist(p19.GetCode(), p19, recvRegist)
	p20 := protocol.SEnterFreeRoom{}
	regist(p20.GetCode(), p20, recvEntryFree)
	p21 := protocol.SFreeSit{}
	regist(p21.GetCode(), p21, recvFreeSit)
	p22 := protocol.SFreeBet{}
	regist(p22.GetCode(), p22, recvFreeBet)
	p23 := protocol.SFreeGamestart{}
	regist(p23.GetCode(), p23, recvGamestartFree)
	p24 := protocol.SFreeGameover{}
	regist(p24.GetCode(), p24, recvGameoverFree)
	p25 := protocol.SPushCurrency{}
	regist(p25.GetCode(), p25, recvPushCurrency)
	p26 := protocol.SEnterClassicRoom{}
	regist(p26.GetCode(), p26, recvEntryClassic)
	p27 := protocol.SClassicList{}
	regist(p27.GetCode(), p27, recvClassicList)
	p28 := protocol.SClassicGameover{}
	regist(p28.GetCode(), p28, recvGameover2)
}

//.

//' 接收到服务器登录返回
func recvRegist(stoc *protocol.SRegist, c *Robot) {
	var errcode uint32 = stoc.GetError()
	switch {
	case errcode == 13016:
		glog.Infof("regist err -> %d", errcode)
		//c.SendLogin() //注册成功,登录
	case errcode == 0:
		c.Logined()     //登录成功
		c.regist = true //注册成功
		c.data.Userid = stoc.GetUserid()
		glog.Infof("regist successful -> %s", c.data.Userid)
		c.SendGetClassic()
		c.SendUserData() // 获取玩家数据
		return
	default:
		glog.Infof("regist err -> %d", errcode)
	}
	c.Close()
}

//.

//' 接收到服务器登录返回
func recvLogin(stoc *protocol.SLogin, c *Robot) {
	var errcode uint32 = stoc.GetError()
	switch {
	case errcode == 13009:
		glog.Infof("login passwd err -> %s", c.data.Phone)
	case errcode == 0:
		c.Logined() //登录成功
		c.data.Userid = stoc.GetUserid()
		glog.Infof("login successful -> %s", c.data.Userid)
		c.SendGetClassic()
		c.SendUserData() // 获取玩家数据
		return
	default:
		glog.Infof("login err -> %d", errcode)
	}
	c.Close()
}

//.

//' 接收到玩家数据
func recvdata(stoc *protocol.SUserData, c *Robot) {
	var errcode uint32 = stoc.GetError()
	if errcode != 0 {
		glog.Infof("get data err -> %d", errcode)
	}
	data := stoc.GetData()
	// 设置数据
	c.data.Userid = data.GetUserid()     // 用户id
	c.data.Nickname = data.GetNickname() // 用户昵称
	c.data.Sex = data.GetSex()           // 用户性别,男1 女2 非男非女3
	c.data.Coin = data.GetCoin()         // 金币
	c.data.Diamond = data.GetDiamond()   // 钻石
	//c.rtype = data.GetRoomtype()         // 房间类型
	if c.data.Coin < 200000 {
		go addCoin(c.data.Userid, 10000000)
	} else if c.data.Coin < 500000 {
		go addCoin(c.data.Userid, 5000000)
	} else if c.data.Coin < 1000000 {
		go addCoin(c.data.Userid, 15000000)
	}
	//查找房间-进入房间
	switch c.code {
	case "free":
		c.SendEntryFree() //进入百人场
	case "classic":
		c.SendEntryClassic()
	case "create":
		switch c.rtype {
		case 6, 7: //跑胡子
			c.SendEntryPhz()
		default:
			c.SendEntry()
		}
	default:
		glog.Infof("enter phz room -> %s, %d", c.code, c.rtype)
		if len(c.code) == 6 {
			switch c.rtype {
			case 6, 7: //跑胡子
				c.SendEntryPhz()
			default:
				c.SendEntry()
			}
		} else {
			c.Close()
		}
	}
}

//.

//' 离开房间
func recvLeave(stoc *protocol.SLeave, c *Robot) {
	var seat uint32 = stoc.GetSeat()
	if seat == c.seat {
		c.Close() //下线
	}
	if seat >= 1 && seat <= 4 && seat != c.seat {
		c.SendLeave() //离开
	}
}

//.

//' 创建房间
func recvCreate(stoc *protocol.SCreateRoom, c *Robot) {
	var errcode uint32 = stoc.GetError()
	switch {
	case errcode == 0:
		var code string = stoc.GetRdata().GetInvitecode()
		if code != "" && len(code) == 6 {
			glog.Infof("create room code -> %s", code)
			c.code = code       //设置邀请码
			c.SendEntry()       //进入房间
			Msg2Robots(code, 3) //创建房间成功,邀请3个人进入
		} else {
			glog.Errorf("create room code empty -> %s", code)
		}
	default:
		glog.Infof("create room err -> %d", errcode)
		c.Close() //进入出错,关闭
	}
}

//.

//' 进入房间
func recvComein(stoc *protocol.SEnterRoom, c *Robot) {
	var errcode uint32 = stoc.GetError()
	switch {
	case errcode == 0:
		userinfo := stoc.GetUserinfo()
		for _, v := range userinfo {
			if v.GetUserid() == c.data.Userid {
				c.seat = v.GetSeat()
				glog.Infof("enter room -> %s, %d", c.data.Userid, c.seat)
				c.SendReady() //准备
				break
			}
		}
		data := stoc.GetRoominfo()
		c.rtype = data.GetRtype() // 房间类型
	default:
		glog.Infof("comein err -> %d", errcode)
		c.Close() //进入出错,关闭
	}
}

//.

//' 解散
func recvVote(stoc *protocol.SLaunchVote, c *Robot) {
	var seat uint32 = stoc.GetSeat()
	glog.Infof("vote seat -> %d", seat)
	c.SendVote()
}

// 解散
func recvKick(stoc *protocol.SKick, c *Robot) {
	userid := stoc.GetUserid()
	if userid == c.data.Userid {
		c.Close() //下线
	}
}

//.

//' 结束
func recvGameover(stoc *protocol.SGameover, c *Robot) {
	var round uint32 = stoc.GetRound()
	c.cards = []uint32{} //清除牌
	if round == 0 {
		c.Close() //结束下线
	} else {
		c.SendReady() //准备
	}
}

//.

//' 抓牌
//STATE_DEALER = 1  //抢庄状态
//STATE_BET    = 2  //下注状态
//STATE_NIU    = 3  //选牛状态
func recvDraw(stoc *protocol.SDraw, c *Robot) {
	var state uint32 = stoc.GetState()
	var seat uint32 = stoc.GetSeat()
	var cards []uint32 = stoc.GetCards()
	if seat != c.seat { //自己摸牌
		return
	}
	c.cards = append(c.cards, cards...)
	switch c.rtype {
	case data.ROOM_PRIVATE:
		switch state {
		case 1:
			//抢庄
			c.SendDealer()
		case 2:
			//提交组合
			c.SendNiu()
		}
	case data.ROOM_PRIVATE4:
		switch state {
		case 2:
			//提交组合
			c.SendNiu()
		}
	case data.ROOM_PRIVATE3:
		switch state {
		case 2:
			//提交组合
			c.SendNiu()
		}
	}
}

//.

//' 打庄
func recvDealer(stoc *protocol.SPushDealer, c *Robot) {
	var dealer uint32 = stoc.GetDealer()
	if c.seat == dealer { //做庄不下注
		return
	}
	switch c.rtype {
	case data.ROOM_PRIVATE:
		//下注
		c.SendBet()
	case data.ROOM_PRIVATE3, data.ROOM_PRIVATE4:
		//下注
		c.SendClassicBet()
	}
}

//.

//' free

//进入
func recvEntryFree(stoc *protocol.SEnterFreeRoom, c *Robot) {
	var errcode uint32 = stoc.GetError()
	switch errcode {
	case 0:
		if c.data.Coin >= 300000 {
			c.SendFreeSit()
		}
		data := stoc.GetRoominfo()
		c.rtype = data.GetRtype() // 房间类型
	default:
		c.SendFreeStandup()
	}
}

//坐下
func recvFreeSit(stoc *protocol.SFreeSit, c *Robot) {
	var errcode uint32 = stoc.GetError()
	var seat uint32 = stoc.GetSeat()
	var userid string = stoc.GetUserid()
	if userid != c.data.Userid {
		return
	}
	switch errcode {
	case 0:
		c.seat = seat //坐下位置
	default:
		//if c.sits > 7 { //尝试次数过多
		//	c.SendFreeStandup()
		//} else {
		//	c.SendFreeSit()
		//}
	}
}

//下注
func recvFreeBet(stoc *protocol.SFreeBet, c *Robot) {
	var errcode uint32 = stoc.GetError()
	var userid string = stoc.GetUserid()
	if userid != c.data.Userid {
		return
	}
	switch errcode {
	case 0:
		//if c.bits < 3 {
		//	c.SendFreeBet()
		//}
	default:
		c.SendFreeStandup()
	}
}

const (
	STATE_FREE_READY  = 0 //准备状态
	STATE_FREE_DEALER = 1 //休息中状态
	STATE_FREE_BET    = 2 //下注中状态
)

//开始
func recvGamestartFree(stoc *protocol.SFreeGamestart, c *Robot) {
	var state uint32 = stoc.GetState()
	switch state {
	case STATE_FREE_READY:
		c.SendFreeStandup()
	case STATE_FREE_DEALER:
		if c.data.Coin < 20000 {
			c.SendFreeStandup()
		}
	case STATE_FREE_BET:
		c.SendFreeBet() //下注
	default:
		c.SendFreeStandup()
	}
}

//结束
func recvGameoverFree(stoc *protocol.SFreeGameover, c *Robot) {
	c.round++
	if c.round >= 16 { //打10局下线
		c.SendFreeStandup()
	}
}

//更新金币
func recvPushCurrency(stoc *protocol.SPushCurrency, c *Robot) {
	var coin int32 = stoc.GetCoin()
	if coin > 100000 {
		c.SendChat()
	}
	newcoin := int32(c.data.Coin) + coin
	if newcoin < 20000 { //金币少于一定金额时下线
		c.SendFreeStandup()
		c.SendLeave()
		c.Close() //下线
	} else {
		c.data.Coin = uint32(newcoin)
	}
}

//.

//' classic

//进入
func recvEntryClassic(stoc *protocol.SEnterClassicRoom, c *Robot) {
	var errcode uint32 = stoc.GetError()
	switch errcode {
	case 0:
		userinfo := stoc.GetUserinfo()
		for _, v := range userinfo {
			if v.GetUserid() == c.data.Userid {
				c.seat = v.GetSeat()
				glog.Infof("enter classic room -> %s, %d", c.data.Userid, c.seat)
				c.SendReady() //准备
				break
			}
		}
		data := stoc.GetRoominfo()
		c.rtype = data.GetRtype() // 房间类型
	default:
		c.SendClassicClose()
	}
}

// 结束
func recvClassicList(stoc *protocol.SClassicList, c *Robot) {
	var errcode uint32 = stoc.GetError()
	switch errcode {
	case 0:
		c.classicId = ""
		c.classic = stoc.GetList()
	default:
		c.SendClassicClose()
	}
}

// 结束
func recvGameover2(stoc *protocol.SClassicGameover, c *Robot) {
	data := stoc.GetData()
	c.cards = []uint32{} //清除牌
	for _, v := range data {
		if v.GetSeat() != c.seat {
			continue
		}
		for _, m := range c.classic {
			if m.GetId() != c.classicId {
				continue
			}
			if m.GetMinimum() > v.GetCoin() {
				c.SendClassicClose()
				return
			}
		}
		if v.GetCoin() <= 20000 {
			c.SendClassicClose()
			return
		} else {
			c.SendReady() //准备
			return
		}
	}
	c.SendReady() //准备
}

//.

//' 跑胡子

func init() {
	p1 := protocol.SEnterZiRoom{}
	regist(p1.GetCode(), p1, recvEntryPhz)
	p2 := protocol.SZiCamein{}
	regist(p2.GetCode(), p2, recvComeinPhz)
	p3 := protocol.SZiGameover{}
	regist(p3.GetCode(), p3, recvPhzOver)
	p4 := protocol.SPushDeal{}
	regist(p4.GetCode(), p4, recvPhzDeal)
	p5 := protocol.SPushDealerDeal{}
	regist(p5.GetCode(), p5, recvPhzDealerDeal)
	p6 := protocol.SPushDraw{}
	regist(p6.GetCode(), p6, recvPhzDraw)
	p7 := protocol.SPushDiscard{}
	regist(p7.GetCode(), p7, recvPhzDiscard)
	p8 := protocol.SPushAuto{}
	regist(p8.GetCode(), p8, recvPhzAuto)
	p9 := protocol.SOperate{}
	regist(p9.GetCode(), p9, recvPhzOperate)
	p10 := protocol.SPushStatus{}
	regist(p10.GetCode(), p10, recvPhzStatus)
}

//进入
func recvEntryPhz(stoc *protocol.SEnterZiRoom, c *Robot) {
	var errcode uint32 = stoc.GetError()
	switch errcode {
	case 0:
		userinfo := stoc.GetUserinfo()
		for _, v := range userinfo {
			if v.GetUserid() == c.data.Userid {
				c.seat = v.GetSeat()
				glog.Infof("enter phz room -> %s, %d", c.data.Userid, c.seat)
				c.SendReady() //准备
				break
			}
		}
		data := stoc.GetRoominfo()
		c.rtype = data.GetRtype() // 房间类型
	default:
		c.SendPhzClose()
	}
}

//进入
func recvComeinPhz(stoc *protocol.SZiCamein, c *Robot) {
	v := stoc.GetUserinfo()
	if v.GetUserid() == c.data.Userid {
		c.seat = v.GetSeat()
		c.SendReady() //准备
	}
}

// 结算广播接口，游戏结束
func recvPhzOver(stoc *protocol.SZiGameover, c *Robot) {
	var round uint32 = stoc.GetRound()
	c.cards = []uint32{} //清除牌
	if round == 0 {
		c.SendPhzClose() //结束下线
	} else {
		c.SendReady() //准备
	}
}

// 发牌
func recvPhzDeal(stoc *protocol.SPushDeal, c *Robot) {
	if stoc.GetSeat() == c.seat {
		c.cards = stoc.GetCards()
	}
}

// 庄家发牌
func recvPhzDealerDeal(stoc *protocol.SPushDealerDeal, c *Robot) {
	if stoc.GetSeat() == c.seat {
		card := stoc.GetCard()
		c.cards = append(c.cards, card)
		value := stoc.GetValue()
		//胡
		c.phzOperate(card, value)
	}
}

// 摸牌
func recvPhzDraw(stoc *protocol.SPushDraw, c *Robot) {
	card := stoc.GetCard()
	c.cards = stoc.GetCards() //同步手牌
	value := stoc.GetValue()
	//胡碰吃
	c.phzOperate(card, value)
}

// 出牌
func recvPhzDiscard(stoc *protocol.SPushDiscard, c *Robot) {
	var errcode uint32 = stoc.GetError()
	switch errcode {
	case 0:
	default:
		glog.Errorf("discard error %d", errcode)
		return
	}
	if c.seat == stoc.GetSeat() {
		return
	}
	card := stoc.GetCard()
	value := stoc.GetValue()
	//吃碰
	c.phzOperate(card, value)
}

//自动操作(提,跑,偎)
func recvPhzAuto(stoc *protocol.SPushAuto, c *Robot) {
	if c.seat == stoc.GetSeat() {
		return
	}
	card := stoc.GetCard()
	value := stoc.GetValue()
	//从手牌中去掉
	c.phzAuto(card, value)
}

// 玩家操作(吃碰胡)
func recvPhzOperate(stoc *protocol.SOperate, c *Robot) {
	if c.seat != stoc.GetSeat() {
		return
	}
	card := stoc.GetCard()
	value := stoc.GetValue()
	cards := stoc.GetCards()
	bione := stoc.GetBione()
	bitwo := stoc.GetBitwo()
	var errcode uint32 = stoc.GetError()
	switch errcode {
	case 0:
	default:
		glog.Errorf("operate error %d", errcode)
		glog.Errorf("operate hs %v", c.cards)
		glog.Errorf("operate value: %d, cards: %v, bione: %v, bitwo: %v",
			value, cards, bione, bitwo)
		return
	}
	//移除牌
	c.phzChow(card, value, cards, bione, bitwo)
}

//房间状态
func recvPhzStatus(stoc *protocol.SPushStatus, c *Robot) {
	if c.seat != stoc.GetSeat() {
		return
	}
	if stoc.GetStatus() == 2 { //出牌
		c.phzDiscard()
	}
}

//.

// vim: set foldmethod=marker foldmarker=//',//.:
