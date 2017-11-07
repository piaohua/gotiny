package role

func (t *Player) Handler(msg interface{}) {
	switch msg.(type) {
	/*
		case *protocol.CBuy:
			arg := msg.(*protocol.CBuy)
			t.buy(arg)
		case *protocol.CApplePay:
			arg := msg.(*protocol.CApplePay)
			t.appleOrder(arg)
		case *protocol.CChatText:
			arg := msg.(*protocol.CChatText)
			t.chattext(arg)
		case *protocol.CChatVoice:
			arg := msg.(*protocol.CChatVoice)
			t.chatsound(arg)
		case *protocol.CNotice:
			arg := msg.(*protocol.CNotice)
			t.getNotice(arg)
		case *protocol.CWxpayQuery:
			arg := msg.(*protocol.CWxpayQuery)
			t.wxQuery(arg)
		case *protocol.CWxpayOrder:
			arg := msg.(*protocol.CWxpayOrder)
			t.wxOrder(arg)
		case *protocol.CEnterFreeRoom:
			arg := msg.(*protocol.CEnterFreeRoom)
			t.entryfreeroom(arg)
		case *protocol.CDealerList:
			arg := msg.(*protocol.CDealerList)
			t.freedealerlist(arg)
		case *protocol.CFreeBet:
			arg := msg.(*protocol.CFreeBet)
			t.freebet(arg)
		case *protocol.CFreeSit:
			arg := msg.(*protocol.CFreeSit)
			t.freesit(arg)
		case *protocol.CFreeDealer:
			arg := msg.(*protocol.CFreeDealer)
			t.freedealer(arg)
		case *protocol.CFreeTrend:
			arg := msg.(*protocol.CFreeTrend)
			t.getTrend(arg)
		case *protocol.CEnterClassicRoom:
			arg := msg.(*protocol.CEnterClassicRoom)
			t.entryclassicroom(arg)
		case *protocol.CEnterRoom:
			arg := msg.(*protocol.CEnterRoom)
			t.entryroom(arg)
		case *protocol.CCreateRoom:
			arg := msg.(*protocol.CCreateRoom)
			t.create(arg)
		case *protocol.CBet:
			arg := msg.(*protocol.CBet)
			t.bet(arg)
		case *protocol.CNiu:
			arg := msg.(*protocol.CNiu)
			t.niu(arg)
		case *protocol.CDealer:
			arg := msg.(*protocol.CDealer)
			t.dealer(arg)
		case *protocol.CReady:
			arg := msg.(*protocol.CReady)
			t.ready(arg)
		case *protocol.CLeave:
			arg := msg.(*protocol.CLeave)
			t.leave(arg)
		case *protocol.CKick:
			arg := msg.(*protocol.CKick)
			t.kick(arg)
		case *protocol.CLaunchVote:
			arg := msg.(*protocol.CLaunchVote)
			t.launchVote(arg)
		case *protocol.CVote:
			arg := msg.(*protocol.CVote)
			t.vote(arg)
		case *protocol.CGameRecord:
			arg := msg.(*protocol.CGameRecord)
			t.getRecord(arg)
		case *protocol.CGetPrize:
			arg := msg.(*protocol.CGetPrize)
			t.getPrize(arg)
		case *protocol.CPrizeCards:
			arg := msg.(*protocol.CPrizeCards)
			t.prizeCards(arg)
		case *protocol.CConfig:
			arg := msg.(*protocol.CConfig)
			t.config(arg)
		case *protocol.CUserData:
			arg := msg.(*protocol.CUserData)
			t.getUserDataHdr(arg)
		case *protocol.CBuildAgent:
			arg := msg.(*protocol.CBuildAgent)
			t.buildAgent(arg)
		case *protocol.CGetCurrency:
			arg := msg.(*protocol.CGetCurrency)
			t.getCurrency(arg)
		case *protocol.CBank:
			arg := msg.(*protocol.CBank)
			t.bank(arg)
		case *protocol.CPing:
			arg := msg.(*protocol.CPing)
			t.ping(arg)
		case *protocol.CBankrupts:
			arg := msg.(*protocol.CBankrupts)
			t.bankrupt(arg)
		case *protocol.CPrizeList:
			arg := msg.(*protocol.CPrizeList)
			t.prizeList(arg)
		case *protocol.CPrizeDraw:
			arg := msg.(*protocol.CPrizeDraw)
			t.prizeDraw(arg)
		case *protocol.CPrizeBox:
			arg := msg.(*protocol.CPrizeBox)
			t.prizeBox(arg)
		case *protocol.CClassicList:
			arg := msg.(*protocol.CClassicList)
			t.classicList(arg)
		case *protocol.CVipList:
			arg := msg.(*protocol.CVipList)
			t.vipList(arg)
		case *protocol.CEnterZiRoom:
			arg := msg.(*protocol.CEnterZiRoom)
			t.entryroomzi(arg)
		case *protocol.CCreateZiRoom:
			arg := msg.(*protocol.CCreateZiRoom)
			t.createzi(arg)
		case *protocol.CZiGameRecord:
			arg := msg.(*protocol.CZiGameRecord)
			t.getRecordZi(arg)
		case *protocol.COperate:
			arg := msg.(*protocol.COperate)
			t.operate(arg)
		case *protocol.CPushDiscard:
			arg := msg.(*protocol.CPushDiscard)
			t.discard(arg)
	*/
	default:
		t.Send2Conn(msg)
		//glog.Errorf("unknown message %v", msg)
	}
}
