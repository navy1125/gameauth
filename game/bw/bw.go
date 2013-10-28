package bw

import (
	"github.com/GXTime/logging"
	"github.com/navy1125/config"
	"github.com/navy1125/gotcp/bw/base"
	//"github.com/navy1125/gotcp/bw/common"
	"github.com/navy1125/gotcp/gotcp"
	//"math/rand"
	Cmd "./common"
	"net"
	"time"
)

type GameLoginBw struct {
	task *gotcp.Task
}

func NewGameLoginBw() *GameLoginBw {
	return &GameLoginBw{}
}

func (self *GameLoginBw) Init(name string) bool {
	conn := self.Connect()
	tick := time.Tick(time.Second)
	for {
		select {
		case <-tick:
			if conn == nil {
				conn = self.Connect()
			}
			if conn != nil {
				self.task = gotcp.NewTask(conn, "BW")
				self.task.SetHandleReadFun(base.HandleReadFunBw)
				self.task.SetHandleWriteFun(base.HandleWriteFunBw)
				self.task.SetHandleParseFun(base.HandleParseBw)
				self.task.SetHandleHeartBteaFun(base.HandleHeartBeatRequestBw, time.Second*10)
				self.task.SetHandleMessage(&handleMessageMap)
				cmd := Cmd.NewStRequestLoginLoginCmd()
				self.task.SendCmd(cmd)
				self.task.Id = 1 //rand.Int63()
				self.task.Name = name
				self.task.Start()
				<-self.task.StopChan
				conn = nil
			}
		}
	}
	return true
}
func (self *GameLoginBw) Connect() *net.TCPConn {
	cfg := config.NewConfig()
	if err := cfg.LoadFromFile(config.GetConfigStr("loginServerList"), "LoginServerList"); err != nil {
		logging.Error("init game err,%s", err.Error())
		return nil
	}
	addr := cfg.GetConfigStr("ip") + ":" + cfg.GetConfigStr("port")
	raddr, _ := net.ResolveTCPAddr("tcp", addr)
	conn, err := net.DialTCP("tcp", nil, raddr)
	if err != nil {
		logging.Error("conn err:%s,1%s", addr, err.Error())
		return nil
	}
	logging.Debug("new connection:%s", conn.RemoteAddr())
	return conn
}

func (self *GameLoginBw) Login(zoneid uint32, myaccount string, myaccid uint32, isAdult uint32, token string) error {
	cmd := Cmd.NewStWebLoginUserTokenWebGateUserCmd()
	cmd.Zoneid = zoneid
	cmd.Accid = myaccid
	for i, v := range myaccount {
		cmd.Account[i] = byte(v)
	}
	cmd.Lifetime = 3600 //token过期时间,0表示只登陆一次
	cmd.UserType = 1    ///ChannelType
	cmd.Type = isAdult
	for i, v := range token {
		cmd.Token[i] = byte(v)
	}
	if self.task != nil {
		self.task.SendCmd(cmd)
	}

	return nil
}

type GameBillingBw struct {
	task *gotcp.Task
}

func NewGameBillingBw() *GameBillingBw {
	return &GameBillingBw{}
}

func (self *GameBillingBw) Init(name string) bool {
	conn := self.Connect()
	tick := time.Tick(time.Second)
	for {
		select {
		case <-tick:
			if conn == nil {
				conn = self.Connect()
			}
			if conn != nil {
				self.task = gotcp.NewTask(conn, "BW")
				self.task.SetHandleReadFun(base.HandleReadFunBw)
				self.task.SetHandleWriteFun(base.HandleWriteFunBw)
				self.task.SetHandleParseFun(base.HandleParseBw)
				self.task.SetHandleHeartBteaFun(base.HandleHeartBeatRequestBw, time.Second*10)
				self.task.SetHandleMessage(&handleMessageMap)
				cmd := Cmd.NewStRequestLoginBillUserCmd()
				cmd.Version = 20211111
				self.task.SendCmd(cmd)
				self.task.Id = 1 //rand.Int63()
				self.task.Name = name
				self.task.Start()
				<-self.task.StopChan
				conn = nil
			}
		}
	}
	return true
}
func (self *GameBillingBw) Connect() *net.TCPConn {
	cfg := config.NewConfig()
	if err := cfg.LoadFromFile(config.GetConfigStr("loginServerList"), "BillServerList"); err != nil {
		logging.Error("init game err,%s", err.Error())
		return nil
	}
	addr := cfg.GetConfigStr("ip") + ":" + cfg.GetConfigStr("port")
	raddr, _ := net.ResolveTCPAddr("tcp", addr)
	conn, err := net.DialTCP("tcp", nil, raddr)
	if err != nil {
		logging.Error("conn err:%s,1%s", addr, err.Error())
		return nil
	}
	logging.Debug("new connection:%s", conn.RemoteAddr())
	return conn
}

func (self *GameBillingBw) Billing(zoneid uint32, myaccid uint32, charid uint32, moneynum uint32) error {
	cmd := Cmd.NewStReturnPlatformToGoldZoneBillUserCmd()
	cmd.Zoneid = zoneid
	cmd.Accid = myaccid
	cmd.Charid = charid
	cmd.DwNum = moneynum
	if self.task != nil {
		self.task.SendCmd(cmd)
	}
	logging.Info("bw billing :%d,%d,%d,%d", zoneid, myaccid, charid, moneynum)

	return nil
}
