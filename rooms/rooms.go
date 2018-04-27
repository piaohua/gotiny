package rooms

import (
	"errors"
	"gotiny/data"
	"gotiny/inter"
	"time"
	"utils"

	"github.com/golang/glog"
)

//全局
var Rooms *RoomsList

//通道关闭信号
type closeFlag int

//牌桌列表结构
type RoomsList struct {
	list        map[string]inter.Pid //牌桌列表
	listFree    map[string]inter.Pid //牌桌列表
	listClassic map[string]inter.Pid //牌桌列表
	listPaohuzi map[string]inter.Pid //牌桌列表
	listCoin    map[string]inter.Pid //牌桌列表
	gen         *data.RoomIDGen      //房间id
	stopCh      chan struct{}        //关闭通道
	msgCh       chan interface{}     //消息通道
}

//初始化列表
func InitRoomsList() inter.Pid {
	Rooms = &RoomsList{
		list:        make(map[string]inter.Pid),
		listFree:    make(map[string]inter.Pid),
		listClassic: make(map[string]inter.Pid),
		listPaohuzi: make(map[string]inter.Pid),
		listCoin:    make(map[string]inter.Pid),
		msgCh:       make(chan interface{}, 100),
		stopCh:      make(chan struct{}),
		gen:         new(data.RoomIDGen),
	}
	Rooms.InitRoomID()
	go Rooms.ticker()  //goroutine
	go Rooms.handler() //goroutine
	return Rooms
}

//' 计时器
func (t *RoomsList) ticker() {
	tick := time.Tick(10 * time.Minute) //每局限制了最高10分钟
	glog.Infof("rooms ticker started -> %#v", t.gen)
	for {
		select {
		case <-t.stopCh:
			glog.Infof("rooms closed list len -> %d", len(t.list))
			glog.Infof("rooms closed gen -> %#v", t.gen)
			return
		default: //防止阻塞
		}
		select {
		case <-t.stopCh:
			return
		case <-tick:
			//逻辑处理
			//TODO 过期清理
			//t.Send(new(data.PlayersOnline))
		}
	}
}

//.

//' 发送消息
func (t *RoomsList) Send(msg interface{}) {
	if t.msgCh == nil {
		glog.Errorf("msg channel closed %x", msg)
		return
	}
	if len(t.msgCh) == cap(t.msgCh) {
		glog.Errorf("send msg channel full -> %d", len(t.msgCh))
		return
	}
	select {
	case <-t.stopCh:
		return
	case t.msgCh <- msg:
	}
}

//.

//' 关闭列表
func (t *RoomsList) Close() {
	t.saveRoomId()
	select {
	case <-t.stopCh:
		return
	default:
		for _, p := range t.list {
			p.Send(true) //下线
		}
		for _, p := range t.listFree {
			p.Send(true) //下线
		}
		for _, p := range t.listClassic {
			p.Send(true) //下线
		}
		for _, p := range t.listPaohuzi {
			p.Send(true) //下线
		}
		for _, p := range t.listCoin {
			p.Send(true) //下线
		}
		t.closed()
	}
}

//关闭
func (t *RoomsList) closed() {
	//关闭消息通道
	t.Send(closeFlag(1))
	//停止发送消息
	close(t.stopCh)
}

//.

//' 消息处理
func (t *RoomsList) handler() {
	glog.Infof("RoomsList start %#v", t.gen)
	for {
		msg, ok := <-t.msgCh
		if !ok {
			glog.Infof("RoomsList handler closed %#v", t.gen)
			//TODO 关闭通道
			return
		}
		//glog.Infof("DeskPhz handler msg %v", msg)
		//消息处理
		switch msg.(type) {
		case *data.RoomsMatch:
			t.entry(msg.(*data.RoomsMatch)) //route
		case *data.RoomsCreate:
			t.create(msg.(*data.RoomsCreate)) //route
		case *data.RoomsAdd:
			t.add(msg.(*data.RoomsAdd)) //route
		case *data.RoomsDel:
			t.del(msg.(*data.RoomsDel)) //route
		case *data.Call:
			t.handler_call(msg.(*data.Call)) //route
		case bool:
			t.Close()
		case closeFlag:
			return //msg channel closed
		default:
			glog.Errorf("unknown message %v", msg)
		}
	}
}

//.

//' 获取信息
//只做查询
func (t *RoomsList) Call(msg interface{}) interface{} {
	arg := new(data.Call)
	arg.Msg = msg
	arg.Pid = make(chan interface{}, 1)
	arg.Timeout = utils.NewTickerMilli(2000)
	//arg.Timeout := utils.NewTicker(1)
	t.Send(arg)
	select {
	case res := <-arg.Pid:
		arg.Timeout.Stop()
		return res
	case <-arg.Timeout.C:
		arg.Timeout.Stop()
		return nil
	}
	return nil
}

//获取消息处理
func (t *RoomsList) handler_call(arg *data.Call) {
	switch arg.Msg.(type) {
	case *data.RoomsExist:
		msg := arg.Msg.(*data.RoomsExist)
		arg.Pid <- t.exist(msg)
	case *data.RoomsGenID:
		//msg := arg.(*data.RoomsGenID)
		arg.Pid <- t.genID()
	case *data.GenInvitecode:
		msg := arg.Msg.(*data.GenInvitecode)
		arg.Pid <- t.genCode(msg.Atype)
	default:
		arg.Pid <- errors.New("unknown message")
		glog.Errorf("unknown message %v", arg)
	}
}

//.

//' 生成ID
func (t *RoomsList) InitRoomID() {
	t.gen.ServerID = data.Conf.ServerId
	t.gen.Get()
	if t.gen.LastRoomID == "" {
		t.gen.LastRoomID = "1"
	}
}

func (t *RoomsList) genID() string {
	t.gen.LastRoomID = utils.StringAdd(t.gen.LastRoomID)
	//失败也更新ID
	t.saveRoomId()
	return t.gen.LastRoomID
}

func (t *RoomsList) saveRoomId() {
	if !t.gen.Save() {
		glog.Errorf("RoomIDGen save err -> %#v", t.gen)
	}
}

//.

//' 增删查
//添加
func (t *RoomsList) add(arg *data.RoomsAdd) {
	switch arg.Atype {
	case data.ROOMS_PRIVATE:
		t.list[arg.Key] = arg.Pid
	case data.ROOMS_FREE:
		t.listFree[arg.Key] = arg.Pid
	case data.ROOMS_CLASSIC:
		t.listClassic[arg.Key] = arg.Pid
	case data.ROOMS_PAOHUZI:
		t.listPaohuzi[arg.Key] = arg.Pid
	case data.ROOMS_COIN:
		t.listCoin[arg.Key] = arg.Pid
	default:
		glog.Errorf("unknown message %v", arg)
	}
}

//移除
func (t *RoomsList) del(arg *data.RoomsDel) {
	switch arg.Atype {
	case data.ROOMS_PRIVATE:
		delete(t.list, arg.Key)
	case data.ROOMS_FREE:
		delete(t.listFree, arg.Key)
	case data.ROOMS_CLASSIC:
		delete(t.listClassic, arg.Key)
	case data.ROOMS_PAOHUZI:
		delete(t.listPaohuzi, arg.Key)
	case data.ROOMS_COIN:
		delete(t.listCoin, arg.Key)
	default:
		glog.Errorf("unknown message %v", arg)
	}
}

//存在
func (t *RoomsList) exist(arg *data.RoomsExist) (ok bool) {
	switch arg.Atype {
	case data.ROOMS_PRIVATE:
		_, ok = t.list[arg.Key]
	case data.ROOMS_FREE:
		_, ok = t.listFree[arg.Key]
	case data.ROOMS_CLASSIC:
		_, ok = t.listClassic[arg.Key]
	case data.ROOMS_PAOHUZI:
		_, ok = t.listPaohuzi[arg.Key]
	case data.ROOMS_COIN:
		_, ok = t.listCoin[arg.Key]
	default:
		glog.Errorf("unknown message %v", arg)
	}
	return
}

//进入,TODO 为房间维护一个状态,方便匹配查找
//暂时使用Call timeout 匹配查找
func (t *RoomsList) entry(arg *data.RoomsMatch) {
	msg := new(data.R2PEntry)
	msg.Atype = arg.Atype
	msg.Id = arg.Id
	switch arg.Atype {
	case data.ROOMS_PRIVATE:
		msg.Pid = t.list[arg.Key]
	case data.ROOMS_PAOHUZI:
		msg.Pid = t.listPaohuzi[arg.Key]
	case data.ROOMS_FREE:
		msg.Pid = t.matchFree(arg)
	case data.ROOMS_CLASSIC:
		msg.Pid = t.matchClassic(arg)
	case data.ROOMS_COIN:
		msg.Pid = t.matchCoin(arg)
	default:
		glog.Errorf("unknown message %v", arg)
	}
	arg.Pid.Send(msg)
}

//创建
func (t *RoomsList) create(arg *data.RoomsCreate) {
	arg.Data.Code = t.genCode(arg.Atype)
	arg.Data.Rid = t.genID()
	msg := new(data.R2PCreate)
	msg.Data = arg.Data
	msg.Atype = arg.Atype
	arg.Pid.Send(msg)
}

//.

//' 匹配
func (t *RoomsList) match(arg *data.RoomsMatch) inter.Pid {
	switch arg.Atype {
	case data.ROOMS_PRIVATE:
		if r, ok := t.list[arg.Key]; ok {
			return r
		}
	case data.ROOMS_FREE:
		return t.matchFree(arg)
	case data.ROOMS_CLASSIC:
		return t.matchClassic(arg)
	case data.ROOMS_PAOHUZI:
		if r, ok := t.listPaohuzi[arg.Key]; ok {
			return r
		}
	case data.ROOMS_COIN:
		return t.matchCoin(arg)
	default:
		glog.Errorf("unknown message %v", arg)
	}
	return nil
}

//匹配一个自由场 TODO 匹配规则
func (t *RoomsList) matchFree(arg *data.RoomsMatch) inter.Pid {
	for _, r := range t.listFree {
		if t.matched(r, arg) {
			return r
		}
	}
	return nil
}

//匹配一个经典场 TODO 匹配规则
func (t *RoomsList) matchClassic(arg *data.RoomsMatch) inter.Pid {
	for _, r := range t.listClassic {
		if t.matched(r, arg) {
			return r
		}
	}
	return nil
}

//匹配一个自由场 TODO 匹配规则
func (t *RoomsList) matchCoin(arg *data.RoomsMatch) inter.Pid {
	for _, r := range t.listCoin {
		if t.matched(r, arg) {
			return r
		}
	}
	return nil
}

//匹配
func (t *RoomsList) matched(r inter.Pid, arg *data.RoomsMatch) bool {
	msg := r.Call(arg)
	switch msg.(type) {
	case bool:
		return msg.(bool)
	default:
		glog.Errorf("unknown message %v", msg)
	}
	return false
}

//.

//' 生成邀请码
//存在
func (t *RoomsList) genCode(atype int) string {
	switch atype {
	case data.ROOMS_PRIVATE:
		return t.genInvitecode()
	case data.ROOMS_FREE:
		return t.genInvitecodeFree()
	case data.ROOMS_CLASSIC:
		//return t.genID()
		return t.genInvitecodeClassic()
	case data.ROOMS_PAOHUZI:
		return t.genInvitecodePaohuzi()
	case data.ROOMS_COIN:
		//return t.genID()
		return t.genInvitecodeCoin()
	default:
		glog.Errorf("unknown message %d", atype)
	}
	return ""
}

//生成一个牌桌邀请码,全列表中唯一
func (t *RoomsList) genInvitecode() (s string) {
	s = utils.RandStr(6)
	if _, ok := t.list[s]; ok { //是否已经存在
		return t.genInvitecode() //重复尝试,TODO:一定次数后放弃尝试
	}
	return
}

//生成一个牌桌邀请码,全列表中唯一
func (t *RoomsList) genInvitecodeFree() (s string) {
	s = utils.RandStr(7)
	if _, ok := t.listFree[s]; ok { //是否已经存在
		return t.genInvitecodeFree() //重复尝试,TODO:一定次数后放弃尝试
	}
	return
}

//生成一个牌桌邀请码,全列表中唯一
func (t *RoomsList) genInvitecodePaohuzi() (s string) {
	s = utils.RandStr(6)
	if _, ok := t.listPaohuzi[s]; ok { //是否已经存在
		return t.genInvitecodePaohuzi() //重复尝试,TODO:一定次数后放弃尝试
	}
	return
}

//生成一个牌桌邀请码,全列表中唯一
func (t *RoomsList) genInvitecodeClassic() (s string) {
	s = utils.RandStr(6)
	if _, ok := t.listClassic[s]; ok { //是否已经存在
		return t.genInvitecodeClassic() //重复尝试,TODO:一定次数后放弃尝试
	}
	return
}

//生成一个牌桌邀请码,全列表中唯一
func (t *RoomsList) genInvitecodeCoin() (s string) {
	s = utils.RandStr(6)
	if _, ok := t.listCoin[s]; ok { //是否已经存在
		return t.genInvitecodeCoin() //重复尝试,TODO:一定次数后放弃尝试
	}
	return
}

//.

// vim: set foldmethod=marker foldmarker=//',//.:
