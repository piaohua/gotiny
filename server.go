package main

import (
	"flag"
	"gotiny/data"
	_ "gotiny/desk"
	"gotiny/roles"
	"gotiny/rooms"
	"gotiny/socket"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/golang/glog"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	var config string
	flag.StringVar(&config, "conf", "./conf.json", "config path")
	flag.Parse()
	data.LoadConf(config)
	defer glog.Flush()
	glog.Infoln("Config: ", data.Conf)
	//启动服务
	data.InitMgo()                 //数据库连接
	pid1 := roles.InitPlayerList() //在线列表服务
	pid2 := rooms.InitRoomsList()  //房间列表服务
	wsServer := new(socket.WSServer)
	//wsServer
	//wsServer.Addr = data.Conf.ServerHost + ":" + data.Conf.ServerPort
	wsServer.Addr = "0.0.0.0:" + data.Conf.ServerPort
	if wsServer != nil {
		wsServer.Start()
	}
	closeSig := make(chan bool, 1)
	go signalProcess(closeSig) //监听关闭信号
	<-closeSig                 //通道阻塞
	//关闭服务
	pid1.Send(true)
	pid2.Send(true)
	//关闭websocket连接,TODO 先关监听
	if wsServer != nil {
		wsServer.Close()
	}
	<-time.After(10 * time.Second) //延迟关闭
	data.Close()
	//延迟等待
	<-time.After(5 * time.Second) //延迟关闭
}

//监听服务
func signalProcess(closeSig chan bool) {
	ch := make(chan os.Signal, 1)
	//signal.Notify(ch, syscall.SIGUSR1, syscall.SIGUSR2)
	//signal.Notify(ch, syscall.SIGHUP)
	signal.Notify(ch, os.Interrupt, os.Kill) //监听SIGINT和SIGKILL信号
	for {
		msg := <-ch
		switch msg {
		default:
			glog.Infof("close signal : %v ", msg)
			//关闭服务
			closeSig <- true
		case syscall.SIGHUP:
			//TODO
		}
	}
}
