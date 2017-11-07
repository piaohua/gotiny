/**********************************************************
 * Author        : piaohua
 * Email         : 814004090@qq.com
 * Last modified : 2017-03-25 17:39:56
 * Filename      : desk.go
 * Description   : 玩牌逻辑
 * *******************************************************/
package desk

import (
	"gotiny/data"
	"gotiny/inter"
)

//进入
func (t *DeskFree) enter(user *data.User, p inter.Pid) int {
	//TODO
	return 0
}

//玩家离开,下注也可以离开
func (t *DeskFree) leave(userid string) int {
	//TODO
	return 0
}

//玩家离开
func (t *DeskFree) sitDown(userid string, seat uint32, st bool) int {
	//TODO
	return 0
}

//0下庄 1上庄 2补庄
func (t *DeskFree) beDealer(userid string, st, num uint32) int {
	//TODO
	return 0
}

//房间消息广播,聊天
func (t *DeskFree) broadcasts(mtype int, userid string, msg []byte) {
	//TODO
}

//下注
func (t *DeskFree) choiceBet(userid string, seatBet, num uint32) int {
	//TODO
	return 0
}
