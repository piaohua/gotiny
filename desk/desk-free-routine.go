package desk

import (
	"errors"
	"gotiny/data"
	"time"
	"utils"

	"github.com/golang/glog"
)

//' 计时器
func (t *DeskFree) ticker() {
	glog.Infof("DeskFree ticker start id %s", t.id)
	tick := time.Tick(time.Second)
	for {
		select {
		case <-t.closeCh:
			glog.Infof("player close -> %s", t.id)
			return
		default:
		}
		select {
		case <-tick:
			//超时判断
			t.Send(1)
		}
	}
}

//.

//' 发送消息
func (t *DeskFree) Send(msg interface{}) {
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

//.

//' 获取信息
//只做查询
func (t *DeskFree) Call(msg interface{}) interface{} {
	arg := new(data.Call)
	arg.Msg = msg
	arg.Pid = make(chan interface{}, 1)
	arg.Timeout = utils.NewTickerMilli(200)
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
func (t *DeskFree) handler_call(arg *data.Call) {
	switch arg.Msg.(type) {
	default:
		arg.Pid <- errors.New("unknown message")
		glog.Errorf("unknown message %v", arg)
	}
}

//.

//' 关闭列表
func (t *DeskFree) Close() {
	select {
	case <-t.stopCh:
		return
	default:
		//通知玩家房间关闭
		t.closed()
	}
}

//关闭
func (t *DeskFree) closed() {
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
func (t *DeskFree) handler() {
	glog.Infof("DeskFree start id %s", t.id)
	for {
		msg, ok := <-t.msgCh
		if !ok {
			glog.Infof("DeskFree handler closed id %s", t.id)
			//TODO 关闭通道
			return
		}
		//glog.Infof("DeskFree handler msg %v", msg)
		//消息处理
		switch msg.(type) {
		case *data.P2DChoiceBet:
			t.msg_bet(msg.(*data.P2DChoiceBet))
		case *data.P2DEntry:
			t.msg_enter(msg.(*data.P2DEntry))
		case *data.P2DLeave:
			t.msg_leave(msg.(*data.P2DLeave))
		case *data.P2DChat:
			t.msg_chat(msg.(*data.P2DChat))
		case *data.P2DPrint:
			t.print()
		case *data.P2DRobot:
			if t.data.Count > 1 {
				go client(t.data.Code,
					(t.data.Count - 1), t.data.Rtype)
			}
		case int:
			//t.timeout() //timeout
		case bool:
			t.Close()
		case *data.Call:
			t.handler_call(msg.(*data.Call)) //route
		default:
			glog.Errorf("unknown message %v", msg)
		}
	}
}

//.
