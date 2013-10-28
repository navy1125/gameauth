package game

import (
	"./bw"
	"github.com/GXTime/logging"
	"regexp"
	"strconv"
	"strings"
)

//login func
type GameLoginFunc func(zoneid uint32, myaccount string, myaccid uint32, isAdult uint32, token string) error
type GameLogin interface {
	Init(name string) bool
	Login(zoneid uint32, myaccount string, myaccid uint32, isAdult uint32, token string) error
}
type GameBill interface {
	Init(name string) bool
	Billing(zoneid uint32, myaccid uint32, charid uint32, moneynum uint32) error
}

var (
	loginMap map[string]GameLogin
	billMap  map[string]GameBill
)

func init() {
	loginMap = make(map[string]GameLogin)
	AddLogin("bw", bw.NewGameLoginBw())
	billMap = make(map[string]GameBill)
	AddBill("bw", bw.NewGameBillingBw())

}
func InitLogin() bool {
	for key, val := range loginMap {
		go val.Init(key)
	}
	return true
}
func AddLogin(name string, game GameLogin) {
	loginMap[name] = game
}
func InitBill() bool {
	for key, val := range billMap {
		go val.Init(key)
	}
	return true
}
func AddBill(name string, game GameBill) {
	billMap[name] = game
}

func AddLoginToken(game string, server_id string, myaccount string, myaccid uint32, isAdult uint32, token string) error {
	if gamefun, ok := loginMap[game]; ok == true {
		zoneid, _ := strconv.Atoi(server_id)
		return gamefun.Login(uint32(zoneid), myaccount, myaccid, isAdult, token)
	} else {
		logging.Error("not game func for game:%s", game)
	}
	return nil
}
func Billing(game string, server_id string, myaccid uint32, charid uint32, moneynum uint32) error {
	if gamefun, ok := billMap[game]; ok == true {
		zoneid, _ := strconv.Atoi(server_id)
		return gamefun.Billing(uint32(zoneid), myaccid, charid, moneynum)
	} else {
		logging.Error("not game func for game:%s", game)
	}
	return nil
}

//get gamename by urlstring usl master as /game/plat
func GetGameNameByUrl(url string) string {
	if ok, err := regexp.MatchString("^/.*/.*/.*$", url); ok == false || err != nil {
		logging.Error("url string err:%s", url)
		return ""
	}
	return strings.Split(url, "/")[1]
}

//get platname by urlstring usl master as /game/plat
func GetPlatNameByUrl(url string) string {
	if ok, err := regexp.MatchString("^/.*/.*/.*$", url); ok == false || err != nil {
		logging.Error("url string err:%s", url)
		return ""
	}
	return strings.Split(url, "/")[2]
}
