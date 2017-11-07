package socket

import (
	"fmt"
	"io/ioutil"
)

var (
	protoPacket = make(map[string]uint32) //响应协议
	protoUnpack = make(map[string]uint32) //请求协议
)

type protoRoute struct {
	name string
	code uint32
}

var protosUnpack = []protoRoute{
	{name: "CBuy", code: 3000},
	{name: "CWxpayOrder", code: 3002},
	{name: "CWxpayQuery", code: 3004},
	{name: "CApplePay", code: 3006},
	{name: "CShop", code: 3010},
	{name: "CChatText", code: 2003},
	{name: "CChatVoice", code: 2004},
	{name: "CNotice", code: 2008},
	{name: "CLogin", code: 1000},
	{name: "CRegist", code: 1002},
	{name: "CWxLogin", code: 1004},
	{name: "CEnterRoom", code: 4000},
	{name: "CCreateRoom", code: 4001},
	{name: "CLeave", code: 4004},
	{name: "CKick", code: 4005},
	{name: "CReady", code: 4006},
	{name: "CDealer", code: 4008},
	{name: "CBet", code: 4010},
	{name: "CNiu", code: 4011},
	{name: "CLaunchVote", code: 4016},
	{name: "CVote", code: 4017},
	{name: "CGameRecord", code: 4020},
	{name: "CEnterFreeRoom", code: 4040},
	{name: "CFreeDealer", code: 4042},
	{name: "CDealerList", code: 4043},
	{name: "CFreeSit", code: 4044},
	{name: "CFreeBet", code: 4046},
	{name: "CFreeTrend", code: 4051},
	{name: "CEnterClassicRoom", code: 4060},
	{name: "CGetPrize", code: 4081},
	{name: "CPrizeCards", code: 4082},
	{name: "CEnterZiRoom", code: 4100},
	{name: "CCreateZiRoom", code: 4101},
	{name: "CZiGameRecord", code: 4103},
	{name: "CPushDiscard", code: 4106},
	{name: "COperate", code: 4108},
	{name: "CConfig", code: 1020},
	{name: "CUserData", code: 1022},
	{name: "CGetCurrency", code: 1024},
	{name: "CBuildAgent", code: 1026},
	{name: "CBank", code: 1030},
	{name: "CPing", code: 1050},
	{name: "CBankrupts", code: 1052},
	{name: "CPrizeList", code: 1053},
	{name: "CPrizeDraw", code: 1054},
	{name: "CPrizeBox", code: 1055},
	{name: "CClassicList", code: 1060},
	{name: "CVipList", code: 1062},
}

var protosPacket = []protoRoute{
	{name: "SBuy", code: 3000},
	{name: "SWxpayOrder", code: 3002},
	{name: "SWxpayQuery", code: 3004},
	{name: "SApplePay", code: 3006},
	{name: "SShop", code: 3010},
	{name: "SChatText", code: 2003},
	{name: "SChatVoice", code: 2004},
	{name: "SBroadcast", code: 2006},
	{name: "SNotice", code: 2008},
	{name: "SLogin", code: 1000},
	{name: "SRegist", code: 1002},
	{name: "SWxLogin", code: 1004},
	{name: "SLoginOut", code: 1006},
	{name: "SEnterRoom", code: 4000},
	{name: "SCreateRoom", code: 4001},
	{name: "SCamein", code: 4003},
	{name: "SLeave", code: 4004},
	{name: "SKick", code: 4005},
	{name: "SReady", code: 4006},
	{name: "SDraw", code: 4007},
	{name: "SDealer", code: 4008},
	{name: "SPushDealer", code: 4009},
	{name: "SBet", code: 4010},
	{name: "SNiu", code: 4011},
	{name: "SGameover", code: 4012},
	{name: "SLaunchVote", code: 4016},
	{name: "SVote", code: 4017},
	{name: "SVoteResult", code: 4018},
	{name: "SGameRecord", code: 4020},
	{name: "SEnterFreeRoom", code: 4040},
	{name: "SFreeCamein", code: 4041},
	{name: "SFreeDealer", code: 4042},
	{name: "SDealerList", code: 4043},
	{name: "SFreeSit", code: 4044},
	{name: "SFreeBet", code: 4046},
	{name: "SFreeGamestart", code: 4048},
	{name: "SFreeGameover", code: 4050},
	{name: "SFreeTrend", code: 4051},
	{name: "SEnterClassicRoom", code: 4060},
	{name: "SClassicGameover", code: 4062},
	{name: "SPubDraw", code: 4080},
	{name: "SGetPrize", code: 4081},
	{name: "SPrizeCards", code: 4082},
	{name: "SEnterZiRoom", code: 4100},
	{name: "SCreateZiRoom", code: 4101},
	{name: "SZiCamein", code: 4104},
	{name: "SZiGameover", code: 4102},
	{name: "SZiGameRecord", code: 4103},
	{name: "SPushDeal", code: 4109},
	{name: "SPushDealerDeal", code: 4110},
	{name: "SPushDraw", code: 4105},
	{name: "SPushDiscard", code: 4106},
	{name: "SPushAuto", code: 4107},
	{name: "SOperate", code: 4108},
	{name: "SPushStatus", code: 4111},
	{name: "SConfig", code: 1020},
	{name: "SUserData", code: 1022},
	{name: "SGetCurrency", code: 1024},
	{name: "SBuildAgent", code: 1026},
	{name: "SPushCurrency", code: 1028},
	{name: "SBank", code: 1030},
	{name: "SPing", code: 1050},
	{name: "SBankrupts", code: 1052},
	{name: "SPrizeList", code: 1053},
	{name: "SPrizeDraw", code: 1054},
	{name: "SPrizeBox", code: 1055},
	{name: "SClassicList", code: 1060},
	{name: "SVipList", code: 1062},
	{name: "SPushVip", code: 1063},
}

//初始化
func Init() {
	//request
	for _, v := range protosUnpack {
		registUnpack(v.name, v.code)
	}
	//response
	for _, v := range protosPacket {
		registPacket(v.name, v.code)
	}
	//最后生成MsgID.lua文件
	genMsgID()
}

func registUnpack(key string, code uint32) {
	if _, ok := protoUnpack[key]; ok {
		panic(fmt.Sprintf("%s registered: %d", key, code))
	}
	protoUnpack[key] = code
}

func registPacket(key string, code uint32) {
	if _, ok := protoPacket[key]; ok {
		panic(fmt.Sprintf("%s registered: %d", key, code))
	}
	protoPacket[key] = code
}

//生成文件
func Gen() {
	gen_packet()
	gen_unpack()
	//client
	gen_client_packet()
	gen_client_unpack()
}

//生成打包文件
func gen_packet() {
	var str string
	str += head_packet()
	str += body_packet()
	str += end_packet()
	err := ioutil.WriteFile("./packet.go", []byte(str), 0644)
	if err != nil {
		panic(fmt.Sprintf("write file err -> %v\n", err))
	}
}

//生成解包文件
func gen_unpack() {
	var str string
	str += head_unpack()
	str += body_unpack()
	str += end_unpack()
	err := ioutil.WriteFile("./unpack.go", []byte(str), 0644)
	if err != nil {
		panic(fmt.Sprintf("write file err -> %v\n", err))
	}
}

func body_unpack() string {
	var str string
	for k, v := range protoUnpack {
		str += fmt.Sprintf("case %d:\n\t\tmsg := new(protocol.%s)\n\t\t%s\n\t", v, k, result_unpack())
	}
	return str
}

func body_packet() string {
	var str string
	for k, v := range protoPacket {
		str += fmt.Sprintf("case *protocol.%s:\n\t\tb, err := proto.Marshal(msg.(proto.Message))\n\t\t%s\n\t", k, result_packet(v))
	}
	return str
}

func head_packet() string {
	return fmt.Sprintf(`// Code generated by protoc-gen-main.
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
	`)
}

func head_unpack() string {
	return fmt.Sprintf(`// Code generated by protoc-gen-main.
// source: gen.go
// DO NOT EDIT!

package socket

import (
	"errors"
	"gotiny/protocol"

	"github.com/golang/protobuf/proto"
)

//解包消息
func unpack(id uint32, b []byte) (interface{}, error) {
	switch id {
	`)
}

func result_packet(code uint32) string {
	return fmt.Sprintf("return %d, b, err", code)
}

func result_unpack() string {
	return fmt.Sprintf(`err := proto.Unmarshal(b, msg)
		return msg, err`)
}

func end_packet() string {
	return fmt.Sprintf(`default:
		return 0, []byte{}, errors.New("unknown msg")
	}
}`)
}

func end_unpack() string {
	return fmt.Sprintf(`default:
		return nil, errors.New("unknown msg")
	}
}`)
}

//生成lua文件
func genMsgID() {
	var str string
	str += fmt.Sprintf("msgID = {")
	for k, v := range protoUnpack {
		str += fmt.Sprintf("\n\t%s = %d,", k, v)
	}
	str += fmt.Sprintf("\n")
	for k, v := range protoPacket {
		str += fmt.Sprintf("\n\t%s = %d,", k, v)
	}
	str += fmt.Sprintf("\n}")
	err := ioutil.WriteFile("./MsgID.lua", []byte(str), 0666)
	if err != nil {
		panic(fmt.Sprintf("write file err -> %v\n", err))
	}
}

/*
msgID = {
	RegisterReq = 0
,	RegisterRes = 1
,	LoginReq = 2
,	LoginRes = 3
,}
func genMsgID() {
	//go_path := os.Getenv("GOPATH")
	//if go_path == "" {
	//	panic(errors.New("GOPATH is not set"))
	//}
	//file, err := os.OpenFile(go_path+"/bin/MsgID.lua", os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0666)
	file, err := os.OpenFile("./MsgID.lua", os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0666)
	defer file.Close()
	if err != nil {
		panic(err)
	}

	_, err = file.WriteString("msgID = {\n")
	if err != nil {
		panic(err)
	}
	ProtobufProcessor.Range(func(id uint16, t reflect.Type) {
		str := fmt.Sprintf("\t%s = %d\n,", t.Elem().Name(), id)
		_, err = file.WriteString(str)
		if err != nil {
			panic(err)
		}
	})
	_, err = file.WriteString("}\n")
	if err != nil {
		panic(err)
	}
}
*/

//生成机器人打包文件
func gen_client_packet() {
	var str string
	str += head_packet()
	str += body_client_packet()
	str += end_packet()
	err := ioutil.WriteFile("../robots/packet.go", []byte(str), 0644)
	if err != nil {
		panic(fmt.Sprintf("write file err -> %v\n", err))
	}
}

//生成机器人解包文件
func gen_client_unpack() {
	var str string
	str += head_unpack()
	str += body_client_unpack()
	str += end_unpack()
	err := ioutil.WriteFile("../robots/unpack.go", []byte(str), 0644)
	if err != nil {
		panic(fmt.Sprintf("write file err -> %v\n", err))
	}
}

func body_client_packet() string {
	var str string
	for k, v := range protoUnpack {
		str += fmt.Sprintf("case *protocol.%s:\n\t\tb, err := proto.Marshal(msg.(proto.Message))\n\t\t%s\n\t", k, result_packet(v))
	}
	return str
}

func body_client_unpack() string {
	var str string
	for k, v := range protoPacket {
		str += fmt.Sprintf("case %d:\n\t\tmsg := new(protocol.%s)\n\t\t%s\n\t", v, k, result_unpack())
	}
	return str
}
