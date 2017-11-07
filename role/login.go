package role

import (
	"gotiny/inter"
	"gotiny/protocol"
)

func login(ctos *protocol.CLogin, msgCh inter.Pid, ip string) inter.Pid {
	stoc := new(protocol.SLogin)
	msgCh.Send(stoc)
	return nil
}

// vim: set foldmethod=marker foldmarker=//',//.:
