/**********************************************************
 * Author        : piaohua
 * Email         : 814004090@qq.com
 * Last modified : 2017-10-14 20:54:59
 * Filename      : desk-niu-msg.go
 * Description   :
 * *******************************************************/
package desk

import (
	"gotiny/data"
	"gotiny/errorcode"
	"gotiny/protocol"

	"github.com/gogo/protobuf/proto"
)

// 消息接收响应

func (t *DeskFree) msg_enter(arg *data.P2DEntry) {
	err := t.enter(arg.User, arg.Pid)
	switch err {
	case 0:
		return
	}
	stoc := new(protocol.SEnterFreeRoom)
	switch err {
	case 1:
		stoc.Error = proto.Uint32(errorcode.RoomFull)
		arg.Pid.Send(stoc)
	}
}

func (t *DeskFree) msg_leave(arg *data.P2DLeave) {
	err := t.leave(arg.Userid)
	switch err {
	case 0:
		msg := new(data.D2PSetRoom)
		arg.Pid.Send(msg)
		return
	}
	stoc := new(protocol.SLeave)
	switch err {
	case 1:
		stoc.Error = proto.Uint32(errorcode.GameStartedCannotLeave)
		arg.Pid.Send(stoc)
	}
}

func (t *DeskFree) msg_chat(arg *data.P2DChat) {
	t.broadcasts(arg.Mtype, arg.Userid, arg.Msg)
}

func (t *DeskFree) msg_sitdown(arg *data.P2DSitDown) {
	err := t.sitDown(arg.Userid, arg.Seat, arg.Situp)
	switch err {
	case 0:
		return
	}
	stoc := new(protocol.SFreeSit)
	switch err {
	case 1:
		stoc.Error = proto.Uint32(errorcode.NotInRoom)
		arg.Pid.Send(stoc)
	case 2:
		stoc.Error = proto.Uint32(errorcode.SitDownFailed)
		arg.Pid.Send(stoc)
	case 3:
		stoc.Error = proto.Uint32(errorcode.StandUpFailed)
		arg.Pid.Send(stoc)
	case 4:
		stoc.Error = proto.Uint32(errorcode.DealerSitFailed)
		arg.Pid.Send(stoc)
	}
}

func (t *DeskFree) msg_bedealer(arg *data.P2DBeDealer) {
	err := t.beDealer(arg.Userid, arg.State, arg.Num)
	switch err {
	case 0:
		return
	}
	stoc := new(protocol.SFreeDealer)
	switch err {
	case 1:
		stoc.Error = proto.Uint32(errorcode.NotInRoom)
		arg.Pid.Send(stoc)
	case 2:
		stoc.Error = proto.Uint32(errorcode.BeDealerAlready)
		arg.Pid.Send(stoc)
	case 3:
		stoc.Error = proto.Uint32(errorcode.BeDealerAlreadySit)
		arg.Pid.Send(stoc)
	case 4:
		stoc.Error = proto.Uint32(errorcode.BeDealerNotEnough)
		arg.Pid.Send(stoc)
	case 5:
		stoc.Error = proto.Uint32(errorcode.GameStartedCannotLeave)
		arg.Pid.Send(stoc)
	}
}

func (t *DeskFree) msg_dealerlist(arg *data.P2DBeDealerList) {
	//上庄列表
	//msg := t.res_bedealerlist()
	//arg.Pid.Send(msg)
}

func (t *DeskFree) msg_bet(arg *data.P2DChoiceBet) {
	err := t.choiceBet(arg.Userid, arg.SeatBet, arg.Num)
	switch err {
	case 0:
		return
	}
	stoc := new(protocol.SFreeBet)
	switch err {
	case 1:
		stoc.Error = proto.Uint32(errorcode.GameNotStart)
		arg.Pid.Send(stoc)
	case 2:
		stoc.Error = proto.Uint32(errorcode.BetDealerFailed)
		arg.Pid.Send(stoc)
	case 3:
		stoc.Error = proto.Uint32(errorcode.BetTopLimit)
		arg.Pid.Send(stoc)
	case 4:
		stoc.Error = proto.Uint32(errorcode.NotInRoom)
		arg.Pid.Send(stoc)
	case 5:
		stoc.Error = proto.Uint32(errorcode.BetNotSeat)
		arg.Pid.Send(stoc)
	}
}

//.
