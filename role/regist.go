package role

import (
	"gotiny/inter"
	"gotiny/protocol"
)

// 玩家goroutine

func regist(ctos *protocol.CRegist, msgCh inter.Pid, ip string) inter.Pid {
	stoc := new(protocol.SRegist)
	msgCh.Send(stoc)
	return nil
}
