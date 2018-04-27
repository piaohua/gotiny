package roles

import (
	"errors"
	"gotiny/data"
	"gotiny/inter"
	"time"
	"utils"

	"github.com/golang/glog"
)

var Players *PlayerList

//通道关闭信号
type closeFlag int

type PlayerList struct {
	list   map[string]inter.Pid
	gen    *data.UserIDGen  //角色id
	stopCh chan struct{}    // 关闭通道
	msgCh  chan interface{} // 消息通道
}

//初始化
func InitPlayerList() inter.Pid {
	Players = &PlayerList{
		list:   make(map[string]inter.Pid),
		stopCh: make(chan struct{}),
		msgCh:  make(chan interface{}, 100),
		gen:    new(data.UserIDGen),
	}
	Players.initPlayerID()
	go Players.ticker()  //goroutine
	go Players.handler() //goroutine
	return Players
}

//' 发送消息
func (t *PlayerList) Send(msg interface{}) {
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
func (t *PlayerList) Close() {
	t.saveUserId()
	select {
	case <-t.stopCh:
		return
	default:
		for _, p := range t.list {
			p.Send(true) //下线
		}
		t.closed()
	}
}

//关闭
func (t *PlayerList) closed() {
	//关闭消息通道
	t.Send(closeFlag(1))
	//停止发送消息
	close(t.stopCh)
}

//.

//' 消息处理
func (t *PlayerList) handler() {
	glog.Infof("PlayerList start %#v", t.gen)
	for {
		msg, ok := <-t.msgCh
		if !ok {
			glog.Infof("PlayerList handler closed %#v", t.gen)
			//TODO 关闭通道
			return
		}
		//glog.Infof("Player handler msg %v", msg)
		//消息处理
		switch msg.(type) {
		case *data.Broadcast:
			arg := msg.(*data.Broadcast)
			t.broadcast(arg)
		case *data.PlayersDel:
			arg := msg.(*data.PlayersDel)
			t.del(arg.Userid)
		case *data.PlayersSet:
			arg := msg.(*data.PlayersSet)
			t.set(arg.Userid, arg.Pid)
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
func (t *PlayerList) Call(msg interface{}) interface{} {
	arg := new(data.Call)
	arg.Msg = msg
	arg.Pid = make(chan interface{}, 1)
	arg.Timeout = utils.NewTickerMilli(1000)
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
func (t *PlayerList) handler_call(arg *data.Call) {
	switch arg.Msg.(type) {
	case *data.PlayersGet:
		msg := arg.Msg.(*data.PlayersGet)
		arg.Pid <- t.get(msg.Userid)
	case *data.PlayersGenID:
		//msg := arg.(*data.PlayersGenID)
		arg.Pid <- t.genID()
	default:
		arg.Pid <- errors.New("unknown message")
		glog.Errorf("unknown message %v", arg)
	}
}

//.

//' 内部消息时调用

//消息广播
func (t *PlayerList) broadcast(msg interface{}) {
	for _, p := range t.list {
		if p != nil {
			p.Send(msg)
		}
	}
}

//生成ID
func (t *PlayerList) genID() string {
	t.gen.LastUserID = utils.StringAdd(t.gen.LastUserID)
	t.saveUserId()
	return t.gen.LastUserID
}

func (t *PlayerList) saveUserId() {
	if !t.gen.Save() {
		glog.Errorf("UserIDGen save err -> %#v", t.gen)
	}
}

func (t *PlayerList) get(id string) inter.Pid {
	if player, ok := t.list[id]; ok {
		return player
	}
	return nil
}

//删除
func (t *PlayerList) del(id string) {
	if p, ok := t.list[id]; ok {
		p.Close()
	}
}

//设置添加
func (t *PlayerList) set(id string, player inter.Pid) {
	t.list[id] = player
}

//.

//' 初始化玩家唯一ID
func (t *PlayerList) initPlayerID() {
	t.gen.ServerID = data.Conf.ServerId
	t.gen.Get()
	if t.gen.LastUserID == "" {
		t.gen.LastUserID = "10000"
	}
}

//.

//' 计时器
func (t *PlayerList) ticker() {
	tick := time.Tick(3 * time.Minute) //每3分钟
	glog.Infof("onlines ticker started -> %d", 1)
	for {
		select {
		case <-t.stopCh:
			glog.Infof("players closed -> %d", len(t.list))
			return
		default:
		}
		select {
		case <-t.stopCh:
			return
		case <-tick:
			//t.Send()
		}
	}
}

//.

// vim: set foldmethod=marker foldmarker=//',//.:
