package plat

import (
	"../db"
	"../game"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/GXTime/logging"
	"github.com/navy1125/config"
	"io"
	"math"
	"net/http"
	"strconv"
	"time"
)

type ErrorState struct {
	Result string `json:"result"`
}

type NickName struct {
	Nickid   string `json:"nickid,omitempty"`
	Nickname string `json:"nickname,omitempty"`
}
type NickNameList struct {
	Result string     `json:"result"`
	Data   []NickName `json:"data,omitempty"`
}

// juxian game auth
func OnKuaiWanAuth(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	logging.Debug("kw request login:%s,%s", req.RemoteAddr, req.URL.Path)
	game_plat := game.GetGameNameByUrl(req.URL.Path) + "_" + game.GetPlatNameByUrl(req.URL.Path)
	account := req.FormValue("login_name")
	timestr := req.FormValue("time")
	game_id := req.FormValue("game_id")
	server_id := req.FormValue("server_id")
	from_id := req.FormValue("from_id")
	sign := req.FormValue("token")
	hash := md5.New()
	mystr := "from_id=" + from_id + "&game_id=" + game_id + "&login_name=" + account + "&server_id=" + server_id + "&time=" + timestr + config.GetConfigStr(game_plat+"_key")
	io.WriteString(hash, mystr)
	mysign := fmt.Sprintf("%x", hash.Sum(nil))
	if mysign != sign {
		logging.Debug("kuaiwan md5 check err:%s,%s,%s", mystr, mysign, sign)
		ret := ErrorState{Result: "-1"}
		b, _ := json.Marshal(ret)
		w.Write(b)
		return
	}
	if account == "" || len(account) > 20 {
		logging.Debug("account err:%s", mystr)
		ret := ErrorState{Result: "-2"}
		b, _ := json.Marshal(ret)
		w.Write(b)
		return
	}
	if game_id == "" || server_id == "" {
		logging.Debug("account err:%s", mystr)
		ret := ErrorState{Result: "-3"}
		b, _ := json.Marshal(ret)
		w.Write(b)
		return
	}
	if from_id == "" {
		logging.Debug("account err:%s", mystr)
		ret := ErrorState{Result: "-4"}
		b, _ := json.Marshal(ret)
		w.Write(b)
		return
	}
	tr, _ := strconv.Atoi(timestr)
	diffsec := math.Abs(float64(time.Now().Unix() - int64(tr)))
	if diffsec > time.Hour.Seconds()*24 {
		ret := ErrorState{Result: "-5"}
		b, _ := json.Marshal(ret)
		w.Write(b)
		return
	}

	account = from_id + ":" + account
	myaccount, myaccid, err := db.GetMyAccount(game.GetPlatNameByUrl(req.URL.Path), account, "")
	if err != nil {
		ret := ErrorState{Result: "-5"}
		b, _ := json.Marshal(ret)
		w.Write(b)
		return
	}
	logging.Debug("request login ok:%s,%d", myaccount, myaccid)
	gameid, _ := strconv.Atoi(game_id)
	serverid, _ := strconv.Atoi(server_id)
	server_id = strconv.Itoa(gameid<<16) + strconv.Itoa(serverid)
	game.AddLoginToken(game.GetGameNameByUrl(req.URL.Path), server_id, myaccount, myaccid, 1, mysign)
	para := fmt.Sprintf("?account=%s&accid=%d&server_id=%s", myaccount, myaccid, server_id)
	http.Redirect(w, req, config.GetConfigStr(game_plat+"_ok")+para, 303)
	//TODO redirect to gamepage
}

// juxian game auth
func OnKuaiWanBill(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	logging.Debug("kw request login:%s,%s", req.RemoteAddr, req.URL.Path)
	game_plat := game.GetGameNameByUrl(req.URL.Path) + "_" + game.GetPlatNameByUrl(req.URL.Path)
	account := req.FormValue("login_name")
	timestr := req.FormValue("time")
	game_id := req.FormValue("game_id")
	server_id := req.FormValue("server_id")
	from_id := req.FormValue("from_id")
	sign := req.FormValue("token")
	hash := md5.New()
	mystr := "from_id=" + from_id + "&game_id=" + game_id + "&login_name=" + account + "&server_id=" + server_id + "&time=" + timestr + config.GetConfigStr(game_plat+"_key")
	io.WriteString(hash, mystr)
	mysign := fmt.Sprintf("%x", hash.Sum(nil))
	if mysign != sign {
		logging.Debug("kuaiwan md5 check err:%s,%s,%s", mystr, mysign, sign)
		ret := ErrorState{Result: "-1"}
		b, _ := json.Marshal(ret)
		w.Write(b)
		return
	}
	if account == "" || len(account) > 20 {
		logging.Debug("account err:%s", mystr)
		ret := ErrorState{Result: "-2"}
		b, _ := json.Marshal(ret)
		w.Write(b)
		return
	}
	if game_id == "" || server_id == "" {
		logging.Debug("account err:%s", mystr)
		ret := ErrorState{Result: "-3"}
		b, _ := json.Marshal(ret)
		w.Write(b)
		return
	}
	if from_id == "" {
		logging.Debug("account err:%s", mystr)
		ret := ErrorState{Result: "-4"}
		b, _ := json.Marshal(ret)
		w.Write(b)
		return
	}
	tr, _ := strconv.Atoi(timestr)
	diffsec := math.Abs(float64(time.Now().Unix() - int64(tr)))
	if diffsec > time.Hour.Seconds()*24 {
		ret := ErrorState{Result: "-5"}
		b, _ := json.Marshal(ret)
		w.Write(b)
		return
	}

	account = from_id + ":" + account
	myaccount, myaccid, err := db.GetMyAccount(game.GetPlatNameByUrl(req.URL.Path), account, "")
	if err != nil {
		ret := ErrorState{Result: "-5"}
		b, _ := json.Marshal(ret)
		w.Write(b)
		return
	}
	logging.Debug("request login ok:%s,%d", myaccount, myaccid)
	gameid, _ := strconv.Atoi(game_id)
	serverid, _ := strconv.Atoi(server_id)
	server_id = strconv.Itoa(gameid<<16) + strconv.Itoa(serverid)
	game.AddLoginToken(game.GetGameNameByUrl(req.URL.Path), server_id, myaccount, myaccid, 1, mysign)
	para := fmt.Sprintf("?account=%s&accid=%d&server_id=%s", myaccount, myaccid, server_id)
	http.Redirect(w, req, config.GetConfigStr(game_plat+"_ok")+para, 303)
	//TODO redirect to gamepage
}

// juxian game auth
func OnKuaiWanCheckName(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	logging.Debug("kw request login:%s,%s", req.RemoteAddr, req.URL.Path)
	game_plat := game.GetGameNameByUrl(req.URL.Path) + "_" + game.GetPlatNameByUrl(req.URL.Path)
	account := req.FormValue("login_name")
	game_id := req.FormValue("game_id")
	server_id := req.FormValue("server_id")
	from_id := req.FormValue("from_id")
	sign := req.FormValue("token")
	hash := md5.New()
	mystr := "from_id=" + from_id + "&game_id=" + game_id + "&login_name=" + account + "&server_id=" + server_id + config.GetConfigStr(game_plat+"_key")
	io.WriteString(hash, mystr)
	mysign := fmt.Sprintf("%x", hash.Sum(nil))
	if mysign != sign {
		logging.Debug("kuaiwan md5 check err:%s,%s,%s", mystr, mysign, sign)
		ret := ErrorState{Result: "-1"}
		b, _ := json.Marshal(ret)
		w.Write(b)
		return
	}
	if from_id == "" {
		logging.Debug("account err:%s", mystr)
		ret := ErrorState{Result: "-2"}
		b, _ := json.Marshal(ret)
		w.Write(b)
		return
	}
	if account == "" || len(account) > 20 {
		logging.Debug("account err:%s", mystr)
		ret := ErrorState{Result: "-3"}
		b, _ := json.Marshal(ret)
		w.Write(b)
		return
	}

	account = from_id + ":" + account
	_, myaccid, err := db.GetMyAccount(game.GetPlatNameByUrl(req.URL.Path), account, "")
	if err != nil {
		ret := ErrorState{Result: "-3"}
		b, _ := json.Marshal(ret)
		w.Write(b)
		return
	}
	names := db.GetAllCharNameByAccid(myaccid)
	if names == nil || len(names) == 0 {
		ret := ErrorState{Result: "-3"}
		b, _ := json.Marshal(ret)
		w.Write(b)
		return
	}
	nicks := NickNameList{
		Result: "1",
		Data:   make([]NickName, len(names)),
	}
	for i, name := range names {
		nicks.Data[i].Nickid = name.CharId
		nicks.Data[i].Nickname = name.CharName
	}
	b, _ := json.Marshal(nicks)
	w.Write(b)
}
