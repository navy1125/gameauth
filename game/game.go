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

var (
	gameMap map[string]GameLogin
)

func init() {
	gameMap = make(map[string]GameLogin)
	AddGame("bw", bw.NewGameLoginBw())

}
func InitGame() bool {
	for key, val := range gameMap {
		go val.Init(key)
	}
	return true
}
func AddGame(name string, game GameLogin) {
	gameMap[name] = game
}

func AddLoginToken(game string, server_id string, myaccount string, myaccid uint32, isAdult uint32, token string) error {
	if gamefun, ok := gameMap[game]; ok == true {
		zoneid, _ := strconv.Atoi(server_id)
		return gamefun.Login(uint32(zoneid), myaccount, myaccid, isAdult, token)
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
