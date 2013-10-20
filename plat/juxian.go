package plat

import (
	"../db"
	"../game"
	"crypto/md5"
	"fmt"
	"github.com/GXTime/logging"
	"github.com/navy1125/config"
	"io"
	"math"
	"net/http"
	"strconv"
	"time"
)

// juxian game auth
func OnAuthJuxian(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	logging.Debug("request login:%s,%s", req.RemoteAddr, req.URL.Path)
	game_plat := game.GetGameNameByUrl(req.URL.Path) + "_" + game.GetPlatNameByUrl(req.URL.Path)
	http.Redirect(w, req, config.GetConfigStr(game_plat+"_err"), 303)
	qid := req.FormValue("qid")
	timestr := req.FormValue("time")
	server_id := req.FormValue("server_id")
	sign := req.FormValue("sign")
	isAdult := req.FormValue("isAdult")
	account := req.FormValue("account")
	hash := md5.New()
	mystr := "qid=" + qid + "&time=" + timestr + "&server_id=" + server_id + config.GetConfigStr(game_plat+"_key")
	io.WriteString(hash, mystr)
	mysign := fmt.Sprintf("%x", hash.Sum(nil))
	if mysign != sign {
		logging.Debug("md5 check err:%s,%s,%s", mystr, mysign, sign)
		//TODO redirect error page
		http.Redirect(w, req, config.GetConfigStr(game_plat+"_err"), 303)
		return
	}
	if account == "" {
		account = qid
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

	myaccount, myaccid, err := db.GetMyAccount(game.GetPlatNameByUrl(req.URL.Path), account)
	if err != nil {
		http.Redirect(w, req, config.GetConfigStr(game_plat+"_err"), 303)
		return
	}
	logging.Debug("request login ok:%s,%d", myaccount, myaccid)
	adult, _ := strconv.Atoi(isAdult)
	game.AddLoginToken(game.GetGameNameByUrl(req.URL.Path), server_id, myaccount, myaccid, uint32(adult), mysign)
	para := fmt.Sprintf("?account=%s&accid=%d&server_id=%s", myaccount, myaccid, server_id)
	http.Redirect(w, req, config.GetConfigStr(game_plat+"_ok")+para, 303)
	//TODO redirect to gamepage
}
