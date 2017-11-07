// Code generated by protoc-gen-main.
// source: gen.go
// DO NOT EDIT!

package socket

import (
	"errors"
	"gotiny/protocol"

	"github.com/golang/protobuf/proto"
)

//打包消息
func packet(msg interface{}) (uint32, []byte, error) {
	switch msg.(type) {
	case *protocol.SPrizeBox:
		b, err := proto.Marshal(msg.(proto.Message))
		return 1055, b, err
	case *protocol.SClassicGameover:
		b, err := proto.Marshal(msg.(proto.Message))
		return 4062, b, err
	case *protocol.SGetPrize:
		b, err := proto.Marshal(msg.(proto.Message))
		return 4081, b, err
	case *protocol.SPushDiscard:
		b, err := proto.Marshal(msg.(proto.Message))
		return 4106, b, err
	case *protocol.SVoteResult:
		b, err := proto.Marshal(msg.(proto.Message))
		return 4018, b, err
	case *protocol.SFreeTrend:
		b, err := proto.Marshal(msg.(proto.Message))
		return 4051, b, err
	case *protocol.SChatVoice:
		b, err := proto.Marshal(msg.(proto.Message))
		return 2004, b, err
	case *protocol.SNotice:
		b, err := proto.Marshal(msg.(proto.Message))
		return 2008, b, err
	case *protocol.SCamein:
		b, err := proto.Marshal(msg.(proto.Message))
		return 4003, b, err
	case *protocol.SPushVip:
		b, err := proto.Marshal(msg.(proto.Message))
		return 1063, b, err
	case *protocol.SPushDealer:
		b, err := proto.Marshal(msg.(proto.Message))
		return 4009, b, err
	case *protocol.SVote:
		b, err := proto.Marshal(msg.(proto.Message))
		return 4017, b, err
	case *protocol.SFreeDealer:
		b, err := proto.Marshal(msg.(proto.Message))
		return 4042, b, err
	case *protocol.SPushAuto:
		b, err := proto.Marshal(msg.(proto.Message))
		return 4107, b, err
	case *protocol.SBankrupts:
		b, err := proto.Marshal(msg.(proto.Message))
		return 1052, b, err
	case *protocol.SEnterRoom:
		b, err := proto.Marshal(msg.(proto.Message))
		return 4000, b, err
	case *protocol.SBet:
		b, err := proto.Marshal(msg.(proto.Message))
		return 4010, b, err
	case *protocol.SZiGameRecord:
		b, err := proto.Marshal(msg.(proto.Message))
		return 4103, b, err
	case *protocol.SApplePay:
		b, err := proto.Marshal(msg.(proto.Message))
		return 3006, b, err
	case *protocol.SCreateRoom:
		b, err := proto.Marshal(msg.(proto.Message))
		return 4001, b, err
	case *protocol.SFreeGameover:
		b, err := proto.Marshal(msg.(proto.Message))
		return 4050, b, err
	case *protocol.SBuildAgent:
		b, err := proto.Marshal(msg.(proto.Message))
		return 1026, b, err
	case *protocol.SWxpayOrder:
		b, err := proto.Marshal(msg.(proto.Message))
		return 3002, b, err
	case *protocol.SReady:
		b, err := proto.Marshal(msg.(proto.Message))
		return 4006, b, err
	case *protocol.SDraw:
		b, err := proto.Marshal(msg.(proto.Message))
		return 4007, b, err
	case *protocol.SGameRecord:
		b, err := proto.Marshal(msg.(proto.Message))
		return 4020, b, err
	case *protocol.SPushCurrency:
		b, err := proto.Marshal(msg.(proto.Message))
		return 1028, b, err
	case *protocol.SBank:
		b, err := proto.Marshal(msg.(proto.Message))
		return 1030, b, err
	case *protocol.SShop:
		b, err := proto.Marshal(msg.(proto.Message))
		return 3010, b, err
	case *protocol.SLogin:
		b, err := proto.Marshal(msg.(proto.Message))
		return 1000, b, err
	case *protocol.SDealer:
		b, err := proto.Marshal(msg.(proto.Message))
		return 4008, b, err
	case *protocol.SCreateZiRoom:
		b, err := proto.Marshal(msg.(proto.Message))
		return 4101, b, err
	case *protocol.SGetCurrency:
		b, err := proto.Marshal(msg.(proto.Message))
		return 1024, b, err
	case *protocol.SClassicList:
		b, err := proto.Marshal(msg.(proto.Message))
		return 1060, b, err
	case *protocol.SBuy:
		b, err := proto.Marshal(msg.(proto.Message))
		return 3000, b, err
	case *protocol.SRegist:
		b, err := proto.Marshal(msg.(proto.Message))
		return 1002, b, err
	case *protocol.SFreeBet:
		b, err := proto.Marshal(msg.(proto.Message))
		return 4046, b, err
	case *protocol.SOperate:
		b, err := proto.Marshal(msg.(proto.Message))
		return 4108, b, err
	case *protocol.SPing:
		b, err := proto.Marshal(msg.(proto.Message))
		return 1050, b, err
	case *protocol.SPrizeDraw:
		b, err := proto.Marshal(msg.(proto.Message))
		return 1054, b, err
	case *protocol.SLoginOut:
		b, err := proto.Marshal(msg.(proto.Message))
		return 1006, b, err
	case *protocol.SLaunchVote:
		b, err := proto.Marshal(msg.(proto.Message))
		return 4016, b, err
	case *protocol.SPushDeal:
		b, err := proto.Marshal(msg.(proto.Message))
		return 4109, b, err
	case *protocol.SChatText:
		b, err := proto.Marshal(msg.(proto.Message))
		return 2003, b, err
	case *protocol.SGameover:
		b, err := proto.Marshal(msg.(proto.Message))
		return 4012, b, err
	case *protocol.SFreeCamein:
		b, err := proto.Marshal(msg.(proto.Message))
		return 4041, b, err
	case *protocol.SDealerList:
		b, err := proto.Marshal(msg.(proto.Message))
		return 4043, b, err
	case *protocol.SPushStatus:
		b, err := proto.Marshal(msg.(proto.Message))
		return 4111, b, err
	case *protocol.SEnterClassicRoom:
		b, err := proto.Marshal(msg.(proto.Message))
		return 4060, b, err
	case *protocol.SPrizeCards:
		b, err := proto.Marshal(msg.(proto.Message))
		return 4082, b, err
	case *protocol.SZiGameover:
		b, err := proto.Marshal(msg.(proto.Message))
		return 4102, b, err
	case *protocol.SPushDraw:
		b, err := proto.Marshal(msg.(proto.Message))
		return 4105, b, err
	case *protocol.SConfig:
		b, err := proto.Marshal(msg.(proto.Message))
		return 1020, b, err
	case *protocol.SPushDealerDeal:
		b, err := proto.Marshal(msg.(proto.Message))
		return 4110, b, err
	case *protocol.SUserData:
		b, err := proto.Marshal(msg.(proto.Message))
		return 1022, b, err
	case *protocol.SKick:
		b, err := proto.Marshal(msg.(proto.Message))
		return 4005, b, err
	case *protocol.SFreeSit:
		b, err := proto.Marshal(msg.(proto.Message))
		return 4044, b, err
	case *protocol.SZiCamein:
		b, err := proto.Marshal(msg.(proto.Message))
		return 4104, b, err
	case *protocol.SPrizeList:
		b, err := proto.Marshal(msg.(proto.Message))
		return 1053, b, err
	case *protocol.SLeave:
		b, err := proto.Marshal(msg.(proto.Message))
		return 4004, b, err
	case *protocol.SFreeGamestart:
		b, err := proto.Marshal(msg.(proto.Message))
		return 4048, b, err
	case *protocol.SEnterZiRoom:
		b, err := proto.Marshal(msg.(proto.Message))
		return 4100, b, err
	case *protocol.SNiu:
		b, err := proto.Marshal(msg.(proto.Message))
		return 4011, b, err
	case *protocol.SEnterFreeRoom:
		b, err := proto.Marshal(msg.(proto.Message))
		return 4040, b, err
	case *protocol.SPubDraw:
		b, err := proto.Marshal(msg.(proto.Message))
		return 4080, b, err
	case *protocol.SVipList:
		b, err := proto.Marshal(msg.(proto.Message))
		return 1062, b, err
	case *protocol.SWxpayQuery:
		b, err := proto.Marshal(msg.(proto.Message))
		return 3004, b, err
	case *protocol.SBroadcast:
		b, err := proto.Marshal(msg.(proto.Message))
		return 2006, b, err
	case *protocol.SWxLogin:
		b, err := proto.Marshal(msg.(proto.Message))
		return 1004, b, err
	default:
		return 0, []byte{}, errors.New("unknown msg")
	}
}