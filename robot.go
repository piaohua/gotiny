/**********************************************************
 * Author        : piaohua
 * Email         : 814004090@qq.com
 * Last modified : 2017-04-28 13:26:03
 * Filename      : robot.go
 * Description   : 机器人
 * *******************************************************/
package main

import (
	"flag"
	"gotiny/data"
	"gotiny/robots"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/golang/glog"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	//加载配置
	var config string
	flag.StringVar(&config, "conf", "./conf.json", "config path")
	flag.Parse()
	data.LoadConf(config)
	glog.Infoln("Config: ", data.Conf)
	defer glog.Flush()
	//启动远程连接服务
	wsServer := new(robots.WSServer)
	//wsServer
	wsServer.Addr = data.Conf.RobotHost + ":" + data.Conf.RobotPort
	if wsServer != nil {
		wsServer.Start()
	}
	//rbServer
	rbServer := robots.NewRobots()
	if rbServer != nil {
		rbServer.Start() //启动服务
	}
	//启动服务
	//images.InitGm()
	//
	closeSig := make(chan bool, 1)
	go signalProcess(closeSig) //监听关闭信号
	<-closeSig                 //通道阻塞
	//关闭服务
	//关闭websocket连接
	if wsServer != nil {
		wsServer.Close()
	}
	//关闭rbServer连接
	if rbServer != nil {
		rbServer.Close() //关闭服务
	}
	//关闭退出消息
	//
	//延迟等待
	<-time.After(10 * time.Second) //延迟关闭
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
