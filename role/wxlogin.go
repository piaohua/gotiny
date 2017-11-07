package role

import (
	"gotiny/inter"
	"gotiny/protocol"
)

func wxLogin(ctos *protocol.CWxLogin, msgCh inter.Pid, ip string) inter.Pid {
	stoc := new(protocol.SWxLogin)
	msgCh.Send(stoc)
	return nil
}
