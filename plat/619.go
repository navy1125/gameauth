package plat

import (
	"../db"
	"../game"
	"crypto/md5"
	"fmt"
	"git.code4.in/logging"
	"github.com/navy1125/config"
	"io"
	"math"
	"net/http"
	"strconv"
	"time"
)

// juxian game auth
func On619GameAuth(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	logging.Debug("619 request login:%s,%s", req.RemoteAddr, req.URL.Path)
	game_plat := game.GetGameNameByUrl(req.URL.Path) + "_" + game.GetPlatNameByUrl(req.URL.Path)
	StationId := req.FormValue("StationId")
	account := req.FormValue("UserId")
	timestr := req.FormValue("Time")
	game_id := req.FormValue("GameId")
	server_id := req.FormValue("ServerId")
	sign := req.FormValue("Sign")
	hash := md5.New()
	mystr := "StationId=" + StationId + "&" + account + "&" + game_id + "&" + server_id + "&" + config.GetConfigStr(game_plat+"_key") + "_" + timestr
	io.WriteString(hash, mystr)
	mysign := fmt.Sprintf("%x", hash.Sum(nil))
	if mysign != sign {
		logging.Debug("md5 check err:%s,%s,%s", mystr, mysign, sign)
		//TODO redirect error page
		http.Redirect(w, req, config.GetConfigStr(game_plat+"_err"), 303)
		return
	}
	if account == "" {
		http.Redirect(w, req, config.GetConfigStr(game_plat+"_err"), 303)
		logging.Debug("account and qid can not be none both:%s", mystr)
		return
	}
	tr, _ := strconv.Atoi(timestr)
	diffsec := math.Abs(float64(time.Now().Unix() - int64(tr)))
	if diffsec > time.Hour.Seconds()*24 {
		http.Redirect(w, req, config.GetConfigStr(game_plat+"_err"), 303)
		logging.Debug("time err:%s,%d,%d", mystr, time.Now().Unix(), diffsec)
		return
	}

	myaccount, myaccid, err := db.GetMyAccount(game.GetPlatNameByUrl(req.URL.Path), account, "")
	if err != nil {
		http.Redirect(w, req, config.GetConfigStr(game_plat+"_err"), 303)
		return
	}
	logging.Debug("request login ok:%s,%d", myaccount, myaccid)
	gameid, _ := strconv.Atoi(game_id)
	serverid, _ := strconv.Atoi(server_id)
	server_id = strconv.Itoa((gameid << 16) + serverid)
	game.AddLoginToken(game.GetGameNameByUrl(req.URL.Path), server_id, myaccount, myaccid, 1, mysign, w, game.GetPlatNameByUrl(req.URL.Path), req, config.GetConfigStr(game_plat+"_err"))
	//para := fmt.Sprintf("?account=%s&accid=%d&server_id=%s", myaccount, myaccid, server_id)
	//http.Redirect(w, req, config.GetConfigStr(game_plat+"_ok")+para, 303)
	//TODO redirect to gamepage
}
