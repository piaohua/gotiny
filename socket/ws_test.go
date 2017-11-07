package socket

import (
	"bytes"
	"fmt"
	"gotiny/protocol"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
)

func TestProto(t *testing.T) {
	packet := &protocol.CConfig{}
	message, _ := proto.Marshal((proto.Message)(packet))
	t.Logf("message -> %+v", message)
	v := &protocol.CConfig{}
	b := []byte{}
	err := proto.Unmarshal(b, proto.Message(v))
	t.Log(err)
	t.Logf("v -> %+v", v)
}

// 启动一个服务
func TestRunServer(t *testing.T) {
	closeSig := make(chan bool, 1)
	Run(closeSig)
}

func Run(closeSig chan bool) {
	wsServer := new(WSServer)
	//wsServer
	wsServer.Addr = "127.0.0.1:7001"
	if wsServer != nil {
		wsServer.Start()
	}
	<-closeSig
	if wsServer != nil {
		wsServer.Close()
	}
}

const (
	HeaderLen uint32 = 1 //包头长度
	PROTOLen  uint32 = 4 //协议头长度
	DataLen   uint32 = 4 //数据长度
	HANDDLen  uint32 = 9 //消息头总长度
)

// 发送一条消息
func TestClient(t *testing.T) {
	var HeaderLen uint32 = 1 //包头长度
	var PROTOLen uint32 = 4
	var DataLen uint32 = 4 //包信息数据长度占位长度
	var HANDDLen uint32 = 9
	var count uint32 = 1
	var p uint32 = 1022
	packet := &protocol.CRegist{
		Nickname: proto.String("wwww"),
		Phone:    proto.String("1111"),
		Password: proto.String("xxxx"),
	}
	message, _ := proto.Marshal((proto.Message)(packet))
	t.Logf("message -> %+v", message)
	msglen := uint32(len(message))
	buff := make([]byte, int(HANDDLen)+len(message))
	t.Logf("buff -> %+v", buff)
	buff[0] = byte(count)
	t.Logf("count -> %d, buff -> %+v", count, buff)
	t.Logf("p -> %d, %+v", p, EncodeUint32(p))
	t.Logf("msglen -> %d, %+v", msglen, EncodeUint32(msglen))
	copy(buff[HeaderLen:HeaderLen+PROTOLen], EncodeUint32(p))
	copy(buff[HeaderLen+PROTOLen:HeaderLen+PROTOLen+DataLen], EncodeUint32(msglen))
	copy(buff[HANDDLen:HANDDLen+msglen], message)
	t.Logf("buff -> %+v", buff)
	client(buff)
}

// 发送一条消息
func client(buff []byte) {
	var addr string = "127.0.0.1:7001"
	u := url.URL{Scheme: "ws", Host: addr, Path: "/"}
	c, _, err := websocket.DefaultDialer.Dial(u.String(),
		http.Header{"Token": {""}})
	//fmt.Printf("c -> %+v\n", c)
	fmt.Printf("err -> %+v\n", err)
	if err != nil {
		fmt.Printf("err -> %+v\n", err)
	}
	if c != nil {
		c.WriteMessage(websocket.TextMessage, buff)
		c.Close()
	}
}

// 发送多条消息,粘包测试
func TestRunClient(t *testing.T) {
	// 注册协议请求消息
	packet := &protocol.CRegist{
		Nickname: proto.String("wwww"),
		Phone:    proto.String("1111"),
		Password: proto.String("xxxx"),
	}
	// 打包protobuf协议消息
	message, _ := proto.Marshal((proto.Message)(packet))
	t.Logf("message -> %+v", message)
	// 打包完整协议消息
	buff := Pack(packet.GetCode(), message, 0)
	// 消息长度
	blen := len(buff)
	t.Logf("buff -> %+v, len(buff) -> %d", buff, blen)
	// 链接服务器
	conn, err := client_conn()
	if err != nil {
		t.Logf("conn error: %s\n", err)
		return
	}
	// 消息缓冲
	msgbuf := bytes.NewBuffer(make([]byte, 0, 1024))
	// 写入5条消息
	for i := 0; i < 5; i++ {
		msgbuf.Write(buff)
	}
	t.Logf("msgbuf len -> %d", msgbuf.Len())
	t.Logf("msgbuf -> %v", msgbuf)
	// 发送一条完整消息
	client_write(conn, msgbuf.Next(blen))
	time.Sleep(time.Second)
	// 发送一条不完整的消息头
	client_write(conn, msgbuf.Next(2))
	time.Sleep(time.Second)
	// 发送消息剩下部分
	client_write(conn, msgbuf.Next(blen-2))
	time.Sleep(time.Second)
	// 三条消息分多段发送
	client_write(conn, msgbuf.Next(blen+2))
	time.Sleep(time.Second)
	client_write(conn, msgbuf.Next(-2+blen-6))
	time.Sleep(time.Second)
	client_write(conn, msgbuf.Next(6+3))
	time.Sleep(time.Second)
	client_write(conn, msgbuf.Next(-3+blen))
	time.Sleep(time.Second)
	// 关闭连接
	conn.Close()
}

//召唤机器人
func client_write(c *websocket.Conn, buff []byte) {
	c.WriteMessage(websocket.TextMessage, buff)
}

//召唤机器人
func client_conn() (*websocket.Conn, error) {
	var addr string = "127.0.0.1:7001"
	u := url.URL{Scheme: "ws", Host: addr, Path: "/"}
	c, _, err := websocket.DefaultDialer.Dial(u.String(),
		http.Header{"Token": {""}})
	return c, err
}
