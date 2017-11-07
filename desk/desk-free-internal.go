/**********************************************************
 * Author        : piaohua
 * Email         : 814004090@qq.com
 * Last modified : 2017-03-30 16:21:11
 * Filename      : desk.go
 * Description   : 玩牌逻辑
 * *******************************************************/
package desk

import (
	"gotiny/data"
	"gotiny/inter"
)

//房间牌桌数据结构
type DeskFree struct {
	id    string         //房间id
	data  *data.DeskData //房间类型基础数据
	cards []uint32       //没摸起的海底牌
	//-------------
	players map[string]*data.User //房间无座玩家
	seats   map[string]uint32     //userid:seat (seat:1~8)
	//-------------
	state   uint32           //房间状态,0准备中,1游戏中
	timer   int              //计时
	closeCh chan bool        //计时器关闭通道
	stopCh  chan struct{}    //关闭通道
	msgCh   chan interface{} //消息通道
	//-------------
	round     uint32                      //局数
	pond      uint32                      //奖池
	dealerNum uint32                      //做庄次数
	dealer    string                      //庄家
	carry     uint32                      //庄家的携带,小于一定值时下庄,字段只做记录,真实数据直接写入玩家数据
	num       uint32                      //当前局下注总数
	dealers   []map[string]uint32         //上庄列表,userid: carry
	bets      map[string]uint32           //userid:num, 玩家下注金额
	seatBets  map[uint32]uint32           //userid:num, 玩家下注金额
	tian      map[string]uint32           //天,seat:value
	di        map[string]uint32           //地
	xuan      map[string]uint32           //玄
	huang     map[string]uint32           //黄
	handCards map[uint32][]uint32         //手牌 seat:cards,seat=(1,2,3,4,5)
	power     map[uint32]uint32           //位置(1-5)对应牌力
	multiple  map[uint32]int32            //结果 seat:num,seat=(1,2,3,4,5),倍数
	score     map[uint32]int32            //位置(1-5)输赢总量
	score2    map[string]int32            //每个闲家输赢总量
	score3    map[uint32]map[string]int32 //位置(1-5)上每个玩家输赢
	score4    map[uint32]uint32           //位置(1-5)分到奖池金额
	score5    map[string]uint32           //玩家分到奖池金额
	//-------------
}

//// external function

//新建一张牌桌
func NewDeskFree(d *data.DeskData) inter.Pid {
	desk := &DeskFree{
		id:      d.Rid,
		data:    d,
		players: make(map[string]*data.User),
		seats:   make(map[string]uint32),
		dealers: make([]map[string]uint32, 0),
		closeCh: make(chan bool, 1),
		stopCh:  make(chan struct{}),
		msgCh:   make(chan interface{}, 100),
	}
	desk.gameInit()   //初始化
	go desk.ticker()  //计时器goroutine
	go desk.handler() //goroutine
	return desk
}

//初始化
func (t *DeskFree) gameInit() {
	t.num = 0                                    //
	t.bets = make(map[string]uint32)             //
	t.seatBets = make(map[uint32]uint32)         //
	t.tian = make(map[string]uint32)             //
	t.di = make(map[string]uint32)               //
	t.xuan = make(map[string]uint32)             //
	t.huang = make(map[string]uint32)            //
	t.handCards = make(map[uint32][]uint32)      //手牌
	t.multiple = make(map[uint32]int32)          //倍数
	t.score = make(map[uint32]int32)             //位置(1-5)输赢总量
	t.score2 = make(map[string]int32)            //个人输赢结果userid: value
	t.score3 = make(map[uint32]map[string]int32) //个人输赢结果userid: value
	t.score4 = make(map[uint32]uint32)           //位置(1-5)分到奖池金额
	t.score5 = make(map[string]uint32)           //分到奖池金额
	t.power = make(map[uint32]uint32)
}

func (t *DeskFree) print() {
}
