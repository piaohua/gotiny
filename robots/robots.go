/**********************************************************
 * Author        : piaohua
 * Email         : 814004090@qq.com
 * Last modified : 2017-01-16 10:00
 * Filename      : robot.go
 * Description   : 机器人
 * *******************************************************/
package robots

import (
	"encoding/json"
	"gotiny/data"
	"net"
	"net/http"
	"net/url"
	"sync"
	"time"
	"utils"

	"github.com/golang/glog"
	"github.com/gorilla/websocket"
)

type WSServer struct {
	Addr            string        //地址
	MaxConnNum      int           //最大连接数
	PendingWriteNum int           //等待写入消息长度
	MaxMsgLen       uint32        //最大消息长度
	HTTPTimeout     time.Duration //超时时间
	ln              net.Listener  //监听
	handler         *WSHandler    //处理
}

type WSHandler struct {
	maxConnNum      int                //最大连接数
	pendingWriteNum int                //等待写入消息长度
	maxMsgLen       uint32             //最大消息长
	upgrader        websocket.Upgrader //升级http连接
	conns           WebsocketConnSet   //连接集合
	mutexConns      sync.Mutex         //互斥锁
	wg              sync.WaitGroup     //同步机制
}

func (handler *WSHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}

	if !verifyToken(r.Header.Get("Token")) {
		return
	}

	conn, err := handler.upgrader.Upgrade(w, r, nil)
	if err != nil {
		glog.Errorf("upgrade error: %v", err)
		return
	}
	conn.SetReadLimit(int64(handler.maxMsgLen))

	handler.wg.Add(1)
	defer handler.wg.Done()

	handler.mutexConns.Lock()
	if handler.conns == nil {
		handler.mutexConns.Unlock()
		conn.Close()
		return
	}
	if len(handler.conns) >= handler.maxConnNum {
		handler.mutexConns.Unlock()
		conn.Close()
		glog.Errorf("too many connections: %d", len(handler.conns))
		return
	}
	handler.conns[conn] = struct{}{}
	handler.mutexConns.Unlock()

	recvHandler(conn)

	// cleanup
	handler.mutexConns.Lock()
	delete(handler.conns, conn)
	handler.mutexConns.Unlock()
}

func (server *WSServer) Start() {
	ln, err := net.Listen("tcp", server.Addr)
	if err != nil {
		glog.Fatal("%v", err)
	}

	if server.MaxConnNum <= 0 {
		server.MaxConnNum = 30000
		glog.Infof("invalid MaxConnNum, reset to %v", server.MaxConnNum)
	}
	if server.PendingWriteNum <= 0 {
		server.PendingWriteNum = 100
		glog.Infof("invalid PendingWriteNum, reset to %v", server.PendingWriteNum)
	}
	if server.MaxMsgLen <= 0 {
		server.MaxMsgLen = 10240
		glog.Infof("invalid MaxMsgLen, reset to %v", server.MaxMsgLen)
	}
	if server.HTTPTimeout <= 0 {
		server.HTTPTimeout = 10 * time.Second
		glog.Infof("invalid HTTPTimeout, reset to %v", server.HTTPTimeout)
	}

	server.ln = ln
	server.handler = &WSHandler{
		maxConnNum:      server.MaxConnNum,
		pendingWriteNum: server.PendingWriteNum,
		maxMsgLen:       server.MaxMsgLen,
		conns:           make(WebsocketConnSet),
		upgrader: websocket.Upgrader{
			ReadBufferSize:   10240, //default 4096
			WriteBufferSize:  10240, //default 4096
			HandshakeTimeout: server.HTTPTimeout,
			CheckOrigin:      func(_ *http.Request) bool { return true },
		},
	}

	httpServer := &http.Server{
		Addr:           server.Addr,
		Handler:        server.handler,
		ReadTimeout:    server.HTTPTimeout,
		WriteTimeout:   server.HTTPTimeout,
		MaxHeaderBytes: 1024,
	}

	go httpServer.Serve(ln)
}

//关闭连接
func (server *WSServer) Close() {
	server.ln.Close()

	server.handler.mutexConns.Lock()
	for conn := range server.handler.conns {
		conn.Close()
	}
	server.handler.conns = nil
	server.handler.mutexConns.Unlock()

	server.handler.wg.Wait()
}

func verifyToken(Token string) bool {
	// client
	// Key := "XG0e2Ye/KAUJRXaMNnJ5UH1haBvh2FXOoAggE6f2Utw"
	// Now := strconv.FormatInt(utils.Timestamp(), 10)
	// Sign := utils.Md5(Key+Now)
	// Token := Sign+Now+RandNum
	// r.Header.Set("Token")
	// server
	// Key := "XG0e2Ye/KAUJRXaMNnJ5UH1haBvh2FXOoAggE6f2Utw"
	// Token := r.Header.Get("Token")
	// r.Header.Del("Token")
	TokenB := []byte(Token)
	if len(TokenB) >= 42 {
		SignB := TokenB[:32]
		TimeB := TokenB[32:42]
		if utils.Md5(data.Conf.GmKey+string(TimeB)) == string(SignB) {
			return true
		}
	}
	return false
}

//处理
func recvHandler(conn *websocket.Conn) {
	defer conn.Close()
	_, message, err := conn.ReadMessage()
	if err != nil {
		glog.Infof("read err -> %v\n", err)
		return
	}
	req := new(ReqMsg)
	err1 := json.Unmarshal(message, req)
	//glog.Infof("req: %#v, err1: %v", req, err1)
	if err1 != nil {
		glog.Errorf("Unmarshal err -> %v\n", err1)
		return
	}
	Msg2Robots(Message{code: req.Code, rtype: req.Rtype}, req.Num)
}

//机器人服务
type Robots struct {
	host    string           //服务器地址
	port    string           //服务器端口
	phone   string           //注册起始电话号码,10005000101
	online  map[string]bool  //map[phone]状态,true=在线,
	offline map[string]bool  //map[phone]状态,true=离线,false=登录中
	msgCh   chan interface{} //消息通道
	wg      sync.WaitGroup   //同步机制
}

//通道
var RobotsCh chan interface{}

//消息通知
func Msg2Robots(msg interface{}, num uint32) {
	for num > 0 {
		RobotsCh <- msg
		num--
	}
}

func NewRobots() *Robots {
	r := new(Robots)
	r.host = data.Conf.ServerHost
	r.port = data.Conf.ServerPort
	r.phone = data.Conf.RobotPhone
	r.online = make(map[string]bool)
	r.offline = make(map[string]bool)
	r.msgCh = make(chan interface{}, 100)
	return r
}

//启动
func (r *Robots) Start() {
	if r.host == "" {
		panic("host empty")
	}
	if r.port == "" {
		panic("port empty")
	}
	if r.phone == "" {
		panic("phone empty")
	}
	RobotsCh = r.msgCh //通道
	go r.Run()         //启动
	//test
	//go r.runTest()
	//go Msg2Robots(Message{code: "create"}, 30)
	//go Msg2Robots(Message{code: "free"}, 10)
	//go r.runTestFree()
}

//机器人测试
func (r *Robots) runTest() {
	glog.Infof("runTest started -> %d", 1)
	tick := time.Tick(time.Second)
	for {
		select {
		case <-tick:
			glog.Infof("r.online -> %d\n", len(r.online))
			glog.Infof("r.offline -> %d\n", len(r.offline))
			glog.Infof("r.phone -> %s\n", r.phone)
			//TODO:优化
			//运行指定数量机器人(每个创建一个牌局)
			//code = "create" 表示机器人创建房间
			if len(r.online) < 3 {
				go Msg2Robots(Message{code: "create"}, 3)
			}
		}
	}
}

//机器人测试
func (r *Robots) runTestFree() {
	glog.Infof("runTestFree started -> %d", 1)
	tick := time.Tick(5 * time.Minute)
	num := 0
	for {
		select {
		case <-tick:
			now := utils.Timestamp()
			today := utils.TimestampToday()
			if (now - today) < (8 * 3600) {
				continue
			}
			glog.Infof("r.online -> %d\n", len(r.online))
			glog.Infof("r.offline -> %d\n", len(r.offline))
			glog.Infof("r.phone -> %s\n", r.phone)
			//TODO:优化,按时间段运行
			//运行指定数量机器人(每个创建一个牌局)
			//code = "create" 表示机器人创建房间
			go Msg2Robots(Message{code: "free"}, 5)
			if len(r.online) < 60 {
				go Msg2Robots(Message{code: "classic"}, 10)
			}
			if num > 1000 {
				break
			}
			num++
		}
	}
}

//关闭
func (r *Robots) Close() {
	close(RobotsCh)
}

//处理
func (r *Robots) Run() {
	defer func() {
		glog.Infof("Robots closed online -> %d\n", len(r.online))
		glog.Infof("Robots closed offline -> %d\n", len(r.offline))
		glog.Infof("Robots closed phone -> %s\n", r.phone)
	}()
	glog.Infof("Robots started -> %s", r.phone)
	tick := time.Tick(time.Minute)
	for {
		select {
		case m, ok := <-r.msgCh:
			if !ok {
				return
			}
			switch m.(type) {
			case Message:
				msg := m.(Message)
				var code string = msg.code
				var rtype uint32 = msg.rtype
				var phone string
				for k, v := range r.offline {
					if v {
						phone = k
						r.offline[k] = false
						go RunRobot(r.host, r.port, phone, code, rtype, false)
						break
					}
				}
				if len(phone) == 0 {
					phone = r.phone
					r.phone = utils.StringAdd(r.phone)
					go RunRobot(r.host, r.port, phone, code, rtype, true)
				}
				//phone = r.phone
				//r.phone = utils.StringAdd(r.phone)
				//go RunRobot(r.host, r.port, phone, code, rypte, r.msgCh)
				glog.Infof("phone -> %s", phone)
			case ReLogin:
				msg := m.(ReLogin)
				glog.Infof("ReLogin -> %#v", msg)
				go RunRobot(r.host, r.port, msg.phone, msg.code, msg.rtype, false)
			case Login:
				msg := m.(Login)
				glog.Infof("login -> %v", msg)
				delete(r.offline, msg.phone)
				r.online[msg.phone] = true
			case Logout:
				msg := m.(Logout)
				glog.Infof("logout -> %v", msg)
				if _, ok := r.online[msg.phone]; ok {
					delete(r.online, msg.phone)
					r.offline[msg.phone] = true
				}
			}
		case <-tick:
			//逻辑处理
		}
	}
}

//启动一个机器人
func RunRobot(host, port, phone, code string, rtype uint32, regist bool) {
	glog.Infof("run robot -> %s, rtype -> %d", phone, rtype)
	u := url.URL{Scheme: "ws", Host: host + ":" + port, Path: "/"}
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		glog.Errorf("robot run dial -> %v", err)
		return
	}
	//new robot
	robot := newRobot(conn, 100, 100)
	robot.code = code //设置邀请码
	robot.data.Phone = phone
	robot.data.Nickname = RandNickName()
	robot.rtype = rtype
	go robot.writePump()
	if regist {
		go robot.SendRegist() //发起请求,注册-登录-进入房间
		go robot.Relogin()    //注册失败时尝试登录
	} else {
		go robot.SendLogin() //登录
	}
	go robot.online()
	robot.readPump()
	// cleanup
	robot.Close()
}

//rand nickname
func RandNickName() string {
	list := GetRobotList()
	var num int32 = int32(len(list))
	if num == 0 {
		return utils.RandStr(6)
	}
	var key int32 = utils.RandInt32N(num)
	robot := list[key]
	return robot.Nickname
}

func addCoin(userid string, coin int32) {
	//reqMsg := &images.ReqMsg{
	//	Userid: userid,
	//	Rtype:  data.LogType9,
	//	Itemid: data.COIN,
	//	Amount: coin,
	//}
	//data1, err1 := json.Marshal(reqMsg)
	//if err1 != nil {
	//	glog.Errorf("reqMsg Marshal err %v", err1)
	//	return
	//}
	//_, err2 := images.Gm("ReqMsg", string(data1))
	//if err2 != nil {
	//	glog.Errorf("reqMsg err %v", err2)
	//}
}
