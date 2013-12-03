package bw

import (
	"git.code4.in/logging"
	"github.com/navy1125/config"
	"github.com/navy1125/gotcp/bw/base"
	//"github.com/navy1125/gotcp/bw/common"
	"github.com/navy1125/gotcp/gotcp"
	//"math/rand"
	"../../db"
	Cmd "./common"
	"fmt"
	"html/template"
	"net"
	"net/http"
	"strconv"
	"time"
)

type BwGameTemplate struct {
	GameName   string
	Zonename   string
	Token      string
	Zoneid     string
	Zonetype   string
	Nettype    string
	Starttype  string
	Patchurl   string
	Allurl     string
	Setuppath  string
	Loginurl   string
	Loginaddr  string
	Loginport  string
	Configpath string
	Logintype  string
}

var (
	tmpl_data *BwGameTemplate
	tmpl      *template.Template
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
	tmpl_data = &BwGameTemplate{}
	tmpl_data.GameName = config.GetConfigStr("bw_gamename")
	tmpl_data.Starttype = config.GetConfigStr("bw_starttype")
	tmpl_data.Patchurl = config.GetConfigStr("bw_patchurl")
	tmpl_data.Allurl = config.GetConfigStr("bw_allurl")
	tmpl_data.Setuppath = config.GetConfigStr("bw_setuppath") + "#" + strconv.Itoa(int(time.Now().Unix()))
	tmpl_data.Loginaddr = config.GetConfigStr("bw_loginaddr")
	tmpl_data.Loginport = config.GetConfigStr("bw_loginport")
	tmpl_data.Configpath = config.GetConfigStr("bw_configpath")
	var err error
	tmpl, err = template.ParseFiles(config.GetConfigStr("bw_plugin") + "/templates/game.html")
	if err != nil {
		logging.Debug("open game page error:%s", err.Error())
		return false
	}
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

func (self *GameLoginBw) Login(zoneid uint32, myaccount string, myaccid uint32, isAdult uint32, token string, w http.ResponseWriter, platname string, req *http.Request, err_url string) error {
	game := 0
	game, tmpl_data.Zonename = db.GetZonenameByZoneid(zoneid)
	tmpl_data.Token = token
	tmpl_data.Zoneid = strconv.Itoa(int(zoneid))
	tmpl_data.Zonetype = "0"
	tmpl_data.Nettype = "0"
	starttype := req.FormValue("starttype")
	if starttype != "" {
		tmpl_data.Starttype = starttype
	}
	logintype := req.FormValue("logintype")
	tmpl_data.Configpath = config.GetConfigStr("bw_configpath") + platname + ".xml"
	cmd := Cmd.NewStWebLoginUserTokenWebGateUserCmd()
	cmd.Zoneid = uint32((game << 16) + int(zoneid))
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
	if logintype == "1" {
		page := `<!DOCTYPE html>
	<html>
		<head>
			<title></title>
		</head>
		<body>
			<script type="text/javascript">
				window.location.href = "%s";
			</script>
		</body>
	</html>
		`
		url := fmt.Sprintf("/api/entergame?server=%d&servername=%s&token=%s", zoneid, tmpl_data.Zonename, token)
		w.Write([]byte(fmt.Sprintf(page, url)))
		return nil
	}
	tmpl_data.Loginurl = err_url
	err := tmpl.Execute(w, tmpl_data)
	if err != nil {
		logging.Debug("excute game page error:%s", err.Error())
		return err
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
