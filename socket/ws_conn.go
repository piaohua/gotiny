package socket

import (
	"bytes"
	"errors"
	"gotiny/data"
	"gotiny/inter"
	"strings"
	"time"
	"utils"

	"github.com/golang/glog"
	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second // Time allowed to write a message to the peer.
	pongWait       = 2 * time.Minute  // Time allowed to read the next pong message from the peer.
	pingPeriod     = 30 * time.Second // Send pings to peer with this period. Must be less than pongWait.
	maxMessageSize = 1024             // Maximum message size allowed from peer.
	waitForLogin   = 20 * time.Second // 连接建立后5秒内没有收到登陆请求,断开socket
)

type WebsocketConnSet map[*websocket.Conn]struct{}

type WSConn struct {
	conn      *websocket.Conn  // websocket连接
	maxMsgLen uint32           // 最大消息长度
	index     int              // 包序
	player    inter.Pid        // 玩家消息通道
	stopCh    chan struct{}    // 关闭通道
	msgCh     chan interface{} // 消息通道
}

//通道关闭信号
type closeFlag int

//' 创建连接
func newWSConn(conn *websocket.Conn, pendingWriteNum int, maxMsgLen uint32) *WSConn {
	return &WSConn{
		conn:      conn,
		maxMsgLen: maxMsgLen,
		msgCh:     make(chan interface{}, pendingWriteNum),
		stopCh:    make(chan struct{}),
	}
}

//.

//' 获取连接地址
func (ws *WSConn) GetIPAddr() string {
	return strings.Split(ws.remoteAddr(), ":")[0]
}

//连接地址
func (ws *WSConn) localAddr() string {
	return ws.conn.LocalAddr().String()
}

func (ws *WSConn) remoteAddr() string {
	return ws.conn.RemoteAddr().String()
}

//.

//' 断开连接,先断开旧连接,再连接新的
func (ws *WSConn) Close() {
	select {
	case <-ws.stopCh:
		return
	default:
		if ws.player != nil {
			//ws.player.Close() //下线
			ws.player.Send(1) //下线
		}
		ws.player = nil
		ws.closed()
		ws.conn.Close()
	}
}

//关闭
func (ws *WSConn) closed() {
	//关闭消息通道
	t.Send(closeFlag(1))
	//停止发送消息
	close(ws.stopCh)
}

//.

//' 发送消息
func (ws *WSConn) Send(msg interface{}) {
	if ws.msgCh == nil {
		glog.Errorf("WSConn msg channel closed %x", msg)
		return
	}
	if len(ws.msgCh) == cap(ws.msgCh) {
		glog.Errorf("send msg channel full -> %d", len(ws.msgCh))
		return
	}
	select {
	case <-ws.stopCh:
		return
	case ws.msgCh <- msg:
	}
}

//.

//' 获取信息
//只做查询
func (ws *WSConn) Call(msg interface{}) interface{} {
	arg := new(data.Call)
	arg.Msg = msg
	arg.Pid = make(chan interface{}, 1)
	arg.Timeout = utils.NewTickerMilli(1000)
	//arg.Timeout := utils.NewTicker(1)
	ws.Send(arg)
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
func (ws *WSConn) handler_call(arg *data.Call) {
	switch arg.Msg.(type) {
	case int:
	default:
		arg.Pid <- errors.New("unknown message")
		glog.Errorf("unknown message %v", arg)
	}
}

//.

//' 登录超时断开连接
func (ws *WSConn) ticker() {
	tick := time.Tick(waitForLogin)
	//glog.Infoln("ws ticker started ... ")
	select {
	case <-tick:
		if ws.player == nil {
			glog.Errorln("ws login time is too long ... ")
			ws.Close()
		}
	}
}

//.

//index(1byte) + proto(4byte) + msgLen(4byte) + msg
//' read goroutine
func (ws *WSConn) readPump() {
	defer ws.Close()
	ws.conn.SetReadLimit(maxMessageSize)
	ws.conn.SetReadDeadline(time.Now().Add(pongWait))
	ws.conn.SetPongHandler(func(string) error { ws.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	// 消息缓冲
	msgbuf := bytes.NewBuffer(make([]byte, 0, 1024))
	// 消息长度
	var length int = 0
	// 包序长度
	var index int = 0
	// 协议编号
	var proto uint32 = 0
	for {
		n, message, err := ws.conn.ReadMessage()
		if err != nil {
			glog.Errorf("Read error: %s, %d\n", err, n)
			break
		}
		// 数据添加到消息缓冲
		m, err := msgbuf.Write(message)
		if err != nil {
			glog.Errorf("Buffer write error: %s, %d\n", err, m)
			return
		}
		// 消息分割循环
		for {
			// 消息头
			if length == 0 && msgbuf.Len() >= 9 {
				index = int(msgbuf.Next(1)[0])             //包序
				proto = DecodeUint32(msgbuf.Next(4))       //协议号
				length = int(DecodeUint32(msgbuf.Next(4))) //消息长度
				// 检查超长消息
				if length > 1024 {
					glog.Errorf("Message too length: %d\n", length)
					return
				}
			} else if length == 0 {
				//不足一条消息
				break
			}
			//fmt.Printf("index: %d, proto: %d, length: %d, len: %d\n", index, proto, length, msgbuf.Len())
			//glog.Infof("index: %d, proto: %d, length: %d, len: %d\n", index, proto, length, msgbuf.Len())
			// 消息体
			if length >= 0 && msgbuf.Len() >= length {
				//fmt.Printf("Client messge: %s\n", string(msgbuf.Next(length)))
				//glog.Infof("index: %d, proto: %d, length: %d, len: %d\n", index, proto, length, msgbuf.Len())
				//msg := msgbuf.Next(length)
				//glog.Infof("Client messge: %s\n", string(msg))
				//包序验证
				ws.index++
				ws.index = ws.index % 256
				if ws.index != index {
					//glog.Errorf("Message index error: %d, %d\n", index, ws.index)
					//return
				}
				//路由
				ws.route(proto, msgbuf.Next(length))
				length = 0
			} else {
				break
			}
		}
	}
}

//.

//' write goroutine TODO write Buff
func (ws *WSConn) writePump() {
	tick := time.Tick(pingPeriod)
	for {
		select {
		case msg, ok := <-ws.msgCh:
			if !ok {
				ws.write(websocket.CloseMessage, []byte{})
				return
			}
			switch msg.(type) {
			case proto.Message:
				err := ws.write(websocket.TextMessage, msg)
				if err != nil {
					return
				}
			case bool:
				ws.Close()
			case closeFlag:
				return //msg channel closed
			default:
				glog.Errorf("unknown message %v", msg)
			}
		case <-tick:
			err := ws.write(websocket.PingMessage, []byte{})
			if err != nil {
				return
			}
		}
	}
}

//写入
func (ws *WSConn) write(mt int, msg interface{}) error {
	var message []byte
	switch msg.(type) {
	case []byte:
		message = msg.([]byte)
	case proto.Message:
		id, b, err := packet(msg)
		if err != nil {
			glog.Errorf("write msg err %v", msg)
			return err
		}
		//msg = Pack(id, b, ws.index)
		message = Pack(id, b, 0)
	}
	if uint32(len(message)) > ws.maxMsgLen {
		glog.Errorf("write msg too long -> %d", len(message))
		return errors.New("write msg too long")
	}
	ws.conn.SetWriteDeadline(time.Now().Add(writeWait))
	return ws.conn.WriteMessage(mt, message)
}

//.

// vim: set foldmethod=marker foldmarker=//',//.:
