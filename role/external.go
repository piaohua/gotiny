package role

import (
	"gotiny/inter"
	"gotiny/protocol"

	"github.com/golang/glog"
)

//登录处理
func HandlerLogin(msg interface{}, msgCh inter.Pid, ip string) inter.Pid {
	switch msg.(type) {
	case *protocol.CLogin:
		arg := msg.(*protocol.CLogin)
		return login(arg, msgCh, ip)
	case *protocol.CRegist:
		arg := msg.(*protocol.CRegist)
		return regist(arg, msgCh, ip)
	case *protocol.CWxLogin:
		arg := msg.(*protocol.CWxLogin)
		return wxLogin(arg, msgCh, ip)
	default:
		glog.Errorf("unknown message %v", msg)
	}
	return nil
}

// vim: set foldmethod=marker foldmarker=//',//.:
