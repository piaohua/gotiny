package socket

import (
	"gotiny/role"

	"github.com/golang/glog"
)

//1000 登录协议
//1002 注册协议
//1004 微信登录协议
func (ws *WSConn) route(code uint32, b []byte) {
	msg, err := unpack(code, b)
	if err != nil {
		glog.Errorf("route err code: %d", code)
		ws.Close()
		return
	}
	switch code {
	case 1000, 1002, 1004:
		if ws.player != nil {
			glog.Errorln("repeat login -> ", code)
			ws.Close()
			return
		}
		ws.player = role.HandlerLogin(msg, ws, ws.GetIPAddr())
	default:
		if ws.player == nil {
			glog.Errorln("player logout -> ", code)
			ws.Close()
			return
		}
		ws.player.Send(msg)
	}
}

//Big Endian
func DecodeUint32(data []byte) uint32 {
	return (uint32(data[0]) << 24) | (uint32(data[1]) << 16) | (uint32(data[2]) << 8) | uint32(data[3])
}

//Big Endian
func EncodeUint32(n uint32) []byte {
	b := make([]byte, 4)
	b[3] = byte(n & 0xFF)
	b[2] = byte((n >> 8) & 0xFF)
	b[1] = byte((n >> 16) & 0xFF)
	b[0] = byte((n >> 24) & 0xFF)
	return b
}

//封包
func Pack(proto uint32, message []byte, count int) []byte {
	buff := make([]byte, 9+len(message))
	msglen := uint32(len(message))
	buff[0] = byte(count)
	copy(buff[1:5], EncodeUint32(proto))
	copy(buff[5:9], EncodeUint32(msglen))
	copy(buff[9:], message)
	return buff
}
