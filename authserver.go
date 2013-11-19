package main

import (
	"./db"
	"./game"
	"./plat"
	"flag"
	"fmt"
	"github.com/GXTime/logging"
	"github.com/navy1125/config"
	"net"
	"net/http"
	"regexp"
	"strings"
)

func main() {
	flag.Parse()
	config.SetConfig("config", *flag.String("config", "config.xml", "config xml file for start"))
	config.SetConfig("logfilename", *flag.String("logfilename", "/log/logfilename.log", "log file name"))
	config.SetConfig("deamon", *flag.String("deamon", "false", "need run as demo"))
	config.SetConfig("port", *flag.String("port", "8000", "http port "))
	config.SetConfig("log", *flag.String("log", "debug", "logger level "))
	config.SetConfig("loginServerList", *flag.String("loginServerList", "loginServerList.xml", "server list config"))
	config.LoadFromFile(config.GetConfigStr("config"), "global")
	if err := config.LoadFromFile(config.GetConfigStr("config"), "AuthServer"); err != nil {
		fmt.Println(err)
		return
	}
	if err := config.LoadFromFile(config.GetConfigStr("key"), "Key"); err != nil {
		fmt.Println(err)
		return
	}

	logger, err := logging.NewTimeRotationHandler(config.GetConfigStr("logfilename"), "060102-15")
	if err != nil {
		fmt.Println(err)
		return
	}
	logger.SetLevel(logging.DEBUG)
	logging.AddHandler("AU", logger)
	mysqlurl := config.GetConfigStr("mysql")
	if ok, err := regexp.MatchString("^mysql://.*:.*@.*/.*$", mysqlurl); ok == false || err != nil {
		logging.Error("mysql config syntax err:%s", mysqlurl)
		return
	}
	mysqlurl = strings.Replace(mysqlurl, "mysql://", "", 1)
	mysqlurl = strings.Replace(mysqlurl, "@", ":", 1)
	mysqlurl = strings.Replace(mysqlurl, "/", ":", 1)
	mysqlurls := strings.Split(mysqlurl, ":")
	config.SetConfig("dbname", mysqlurls[4])
	db.InitDatabase(mysqlurls[2]+":"+mysqlurls[3], mysqlurls[0], mysqlurls[1], mysqlurls[4])
	mysqlurl = config.GetConfigStr("mysql_checkname")
	if ok, err := regexp.MatchString("^mysql://.*:.*@.*/.*$", mysqlurl); ok == false || err != nil {
		logging.Error("mysql config syntax err:%s", mysqlurl)
		return
	}
	mysqlurl = strings.Replace(mysqlurl, "mysql://", "", 1)
	mysqlurl = strings.Replace(mysqlurl, "@", ":", 1)
	mysqlurl = strings.Replace(mysqlurl, "/", ":", 1)
	mysqlurls = strings.Split(mysqlurl, ":")
	db.InitCheckNameDatabase(mysqlurls[2]+":"+mysqlurls[3], mysqlurls[0], mysqlurls[1], mysqlurls[4])
	plat.InitPlat()
	game.InitLogin()
	game.InitBill()
	//addr, err := net.InterfaceByIndex(0)
	ifname, err := net.InterfaceByName(config.GetConfigStr("ifname"))
	if err != nil {
		fmt.Println("GetLocalIp Err:", err)
		logging.Debug("GetLocalIp Err:%s", err.Error())
	}
	addrs, err := ifname.Addrs()
	if err != nil {
		fmt.Println("GetLocalIp Err:", err)
		logging.Debug("GetLocalIp Err:%s", err.Error())
	}
	err = http.ListenAndServe(strings.Split(addrs[0].String(), "/")[0]+":"+config.GetConfigStr("port"), nil)
	if err != nil {
		fmt.Println(err)
		logging.Debug("ListenAndServe:%s", err.Error())
	}
}
