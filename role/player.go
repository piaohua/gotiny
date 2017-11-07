package role

import (
	"errors"
	"gotiny/data"
	"gotiny/inter"
	"time"
	"utils"

	"github.com/golang/glog"
)

//用户
type Player struct {
	*data.User                  // 基础数据
	conn       inter.Pid        // 连接
	room       inter.Pid        // 房间
	state      bool             // 数据回存状态
	closeCh    chan bool        // 计时器关闭通道
	stopCh     chan struct{}    // 关闭通道
	msgCh      chan interface{} // 消息通道
}

//创建
func NewPlayer(user *data.User) *Player {
	this := &Player{
		User:    user,
		closeCh: make(chan bool, 1),
		stopCh:  make(chan struct{}),
		msgCh:   make(chan interface{}, 100),
	}
	go this.ticker()  //goroutine
	go this.handler() //goroutine
	return this
}

//' 计时器
func (this *Player) ticker() {
	glog.Infof("Player ticker start id %s", this.GetUserid())
	tick := time.Tick(5 * time.Minute) //每5分钟
	for {
		select {
		case <-this.closeCh:
			glog.Infof("player close -> %s", this.GetUserid())
			return
		default:
		}
		select {
		case <-tick:
			this.Send(new(data.PlayerSave))
		}
	}
}

//.

//' 发送消息
func (t *Player) Send(msg interface{}) {
	if t.msgCh == nil {
		glog.Errorf("msg channel closed %v", msg)
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

//返回消息
func (t *Player) Send2Conn(msg interface{}) {
	if t.conn == nil {
		glog.Errorf("conn channel closed %v", msg)
		return
	}
	t.conn.Send(msg)
}

//.

//' 获取信息
//只做查询
func (t *Player) Call(msg interface{}) interface{} {
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
func (t *Player) handler_call(arg *data.Call) {
	switch arg.Msg.(type) {
	case int:
	default:
		arg.Pid <- errors.New("unknown message")
		glog.Errorf("unknown message %v", arg)
	}
}

//.

//' 关闭列表
func (t *Player) Close() {
	select {
	case <-t.stopCh:
		return
	default:
		t.conn.Send(true) //下线消息
		t.closed()
	}
}

//关闭
func (t *Player) closed() {
	//停止发送消息
	close(t.stopCh)
	//关闭计时器
	if t.closeCh != nil {
		t.closeCh <- true
		close(t.closeCh) //关闭计时器
		t.closeCh = nil  //消除计时器
	}
	//关闭消息通道
	close(t.msgCh)
}

//.

//' 消息处理
func (t *Player) handler() {
	glog.Infof("Player start id %s", t.GetUserid())
	for {
		msg, ok := <-t.msgCh
		if !ok {
			glog.Infof("Player handler closed %s", t.GetUserid())
			//TODO 关闭通道
			return
		}
		//glog.Infof("Player handler msg %v", msg)
		//消息处理
		switch msg.(type) {
		case *data.PlayerSave:
			//TODO
		case bool:
			t.Close()
		case int:
			t.conn = nil //断开连接
		case *data.Call:
			t.handler_call(msg.(*data.Call)) //route
		default:
			glog.Errorf("unknown message %v", msg)
		}
	}
}

//.

// vim: set foldmethod=marker foldmarker=//',//.:
