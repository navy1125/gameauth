package plat

import (
	"../db"
	"../game"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"git.code4.in/logging"
	"github.com/navy1125/config"
	"github.com/xuyu/iconv"
	"io"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// juxian game auth
func OnJuXianAuth(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	game_plat := game.GetGameNameByUrl(req.URL.Path) + "_" + game.GetPlatNameByUrl(req.URL.Path)
	qid := req.FormValue("qid")
	timestr := req.FormValue("time")
	server_id := req.FormValue("server_id")
	sign := req.FormValue("sign")
	isAdult := req.FormValue("isAdult")
	account := req.FormValue("account")
	logging.Debug("juxian request login:%s,%s,%s", qid, req.RemoteAddr, req.URL.Path)
	hash := md5.New()
	mystr := "qid=" + qid + "&time=" + timestr + "&server_id=" + server_id + config.GetConfigStr(game_plat+"_key")
	io.WriteString(hash, mystr)
	logging.Debug("%s", mystr)
	mysign := fmt.Sprintf("%x", hash.Sum(nil))
	if strings.ToLower(mysign) != strings.ToLower(sign) {
		hash.Reset()
		mystr1 := "qid=" + qid + "&time=" + timestr + "&server_id=" + server_id + "&sign=" + config.GetConfigStr(game_plat+"_key")
		io.WriteString(hash, mystr1)
		logging.Debug("%s", mystr)
		mysign := fmt.Sprintf("%x", hash.Sum(nil))
		if strings.ToLower(mysign) != strings.ToLower(sign) {
			hash.Reset()
			mystr2 := qid + timestr + server_id + config.GetConfigStr(game_plat+"_key")
			io.WriteString(hash, mystr1)
			logging.Debug("%s", mystr2)
			mysign := fmt.Sprintf("%x", hash.Sum(nil))
			if strings.ToLower(mysign) != strings.ToLower(sign) {
				logging.Debug("md5 check err:%s,%s,%s,%s,%s", mystr, mystr1, mystr2, mysign, sign)
				//TODO redirect error page
				http.Redirect(w, req, config.GetConfigStr(game_plat+"_err"), 303)
				return
			}
		}
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

	myaccount, myaccid, err := db.GetMyAccount(game.GetPlatNameByUrl(req.URL.Path), account, qid)
	if err != nil {
		http.Redirect(w, req, config.GetConfigStr(game_plat+"_err"), 303)
		return
	}
	logging.Debug("request login ok:%s,%d", myaccount, myaccid)
	adult, _ := strconv.Atoi(isAdult)
	game.AddLoginToken(game.GetGameNameByUrl(req.URL.Path), server_id, myaccount, myaccid, uint32(adult), mysign, w, game.GetPlatNameByUrl(req.URL.Path), req, config.GetConfigStr(game_plat+"_err"))
	//para := fmt.Sprintf("?account=%s&accid=%d&server_id=%s", myaccount, myaccid, server_id)

	//http.Redirect(w, req, config.GetConfigStr(game_plat+"_ok")+para, 303)
	//TODO redirect to gamepage
}
func OnJuXianBill(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	game_plat := game.GetGameNameByUrl(req.URL.Path) + "_" + game.GetPlatNameByUrl(req.URL.Path)
	qid := req.FormValue("qid")
	order_amount := req.FormValue("order_amount")
	order_id := req.FormValue("order_id")
	server_id := req.FormValue("server_id")
	sign := req.FormValue("sign")
	logging.Debug("juxian request bill:%s,%s,%s", qid, req.RemoteAddr, req.URL.Path)
	hash := md5.New()
	mystr := qid + order_amount + order_id + server_id + config.GetConfigStr(game_plat+"_key")
	io.WriteString(hash, mystr)
	mysign := fmt.Sprintf("%x", hash.Sum(nil))
	if strings.ToLower(mysign) != strings.ToLower(sign) {
		hash.Reset()
		mystr1 := "qid=" + qid + "&order_amount=" + order_amount + "&order_id=" + order_id + "&server_id=" + server_id + "&sign=" + config.GetConfigStr(game_plat+"_key")
		io.WriteString(hash, mystr1)
		mysign := fmt.Sprintf("%x", hash.Sum(nil))
		if strings.ToLower(mysign) != strings.ToLower(sign) {
			hash.Reset()
			mystr2 := "qid=" + qid + "&order_amount=" + order_amount + "&order_id=" + order_id + "&server_id=" + server_id + "&sign=" + config.GetConfigStr(game_plat+"_key")
			io.WriteString(hash, mystr2)
			mysign := fmt.Sprintf("%x", hash.Sum(nil))
			if strings.ToLower(mysign) != strings.ToLower(sign) {
				logging.Debug("md5 check err:%s,%s,%s,%s,%s", mystr, mystr1, mystr2, mysign, sign)
				w.Write([]byte("0"))
				return
			}
		}
	}

	myaccid, err := db.GetMyAccountByAccountId(game.GetPlatNameByUrl(req.URL.Path), qid)
	if err != nil {
		logging.Debug("bill account err:%s", qid)
		w.Write([]byte("3"))
		return
	}
	zoneid, _ := strconv.Atoi(server_id)
	gameid, _ := db.GetZonenameByZoneid(uint32(zoneid))
	logging.Debug("request bill ok:%s,%d,%d,%d", qid, myaccid, gameid, zoneid)
	serverid, _ := strconv.Atoi(server_id)
	server_id = strconv.Itoa((gameid << 16) + serverid)
	names := db.GetAllZoneCharNameByAccid(server_id, myaccid)
	if names == nil || len(names) == 0 {
		logging.Debug("bill checkname err:%d,%s", myaccid, qid)
		w.Write([]byte("3"))
		return
	}
	moneynum, _ := strconv.ParseFloat(order_amount, 32)
	game.Billing(game.GetGameNameByUrl(req.URL.Path), server_id, myaccid, 0, uint32(moneynum*100))
	w.Write([]byte("1"))
}

// juxian game auth
func OnJuXianCheckName(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	game_plat := game.GetGameNameByUrl(req.URL.Path) + "_" + game.GetPlatNameByUrl(req.URL.Path)
	qid := req.FormValue("qid")
	server_id := req.FormValue("server_id")
	sign := req.FormValue("sign")
	logging.Debug("juxian request checkname:%s,%s,%s", qid, req.RemoteAddr, req.URL.Path)
	hash := md5.New()
	mystr := "qid=" + qid + "&server_id=" + server_id + config.GetConfigStr(game_plat+"_key")
	io.WriteString(hash, mystr)
	mysign := fmt.Sprintf("%x", hash.Sum(nil))
	if strings.ToLower(mysign) != strings.ToLower(sign) {
		hash.Reset()
		mystr1 := "qid=" + qid + "&server_id=" + server_id + "&sign=" + config.GetConfigStr(game_plat+"_key")
		io.WriteString(hash, mystr1)
		mysign := fmt.Sprintf("%x", hash.Sum(nil))
		if strings.ToLower(mysign) != strings.ToLower(sign) {
			hash.Reset()
			mystr2 := qid + server_id + config.GetConfigStr(game_plat+"_key")
			io.WriteString(hash, mystr2)
			mysign := fmt.Sprintf("%x", hash.Sum(nil))
			if strings.ToLower(mysign) != strings.ToLower(sign) {
				logging.Debug("juxian md5 check err:%s,%s,%s,%s,%s", mystr, mystr1, mystr2, mysign, sign)
				ret := ErrorState{Result: "-1"}
				b, _ := json.Marshal(ret)
				w.Write(b)
				return
			}
			return
		}
	}

	//account := qid

	myaccid, err := db.GetMyAccountByAccountId(game.GetPlatNameByUrl(req.URL.Path), qid)
	if err != nil {
		http.Redirect(w, req, config.GetConfigStr(game_plat+"_err"), 303)
		return
	}
	zoneid, _ := strconv.Atoi(server_id)
	gameid, _ := db.GetZonenameByZoneid(uint32(zoneid))
	logging.Debug("request check ok:%s,%d,%d,%d", qid, myaccid, gameid, zoneid)
	serverid, _ := strconv.Atoi(server_id)
	server_id = strconv.Itoa((gameid << 16) + serverid)
	names := db.GetAllZoneCharNameByAccid(server_id, myaccid)
	if names == nil || len(names) == 0 {
		logging.Debug("names err:%d", myaccid)
		ret := ErrorState{Result: "-3"}
		b, _ := json.Marshal(ret)
		w.Write(b)
		return
	}
	nicks := NickNameList{
		Result: "1",
		Data:   make([]NickName, len(names)),
	}
	convert, _ := iconv.Open("GB2312", "UTF-8")
	for i, name := range names {
		charname, _ := convert.ConvString(name.CharName)
		fmt.Println(name.CharId, charname)
		nicks.Data[i].Nickid = name.CharId
		nicks.Data[i].Nickname = charname
	}
	b, _ := json.Marshal(nicks)
	fmt.Println(string(b))
	w.Write(b)
}

func OnJuUcjoyBill(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	game_plat := game.GetGameNameByUrl(req.URL.Path) + "_" + game.GetPlatNameByUrl(req.URL.Path)
	qid := req.FormValue("accountid")
	gameid := req.FormValue("gameid")
	server_id := req.FormValue("serverid")
	orderid := req.FormValue("orderid")
	order_amount := req.FormValue("point")
	giftcoin := req.FormValue("giftcoin")
	timestamp := req.FormValue("timestamp")
	remark := req.FormValue("remark")
	sign := req.FormValue("sign")
	logging.Debug("ucjoy request bill:%s,%s,%s", qid, req.RemoteAddr, req.URL.Path)
	hash := md5.New()
	mystr := qid + gameid + server_id + orderid + order_amount + giftcoin + timestamp + remark + config.GetConfigStr(game_plat+"_key")
	io.WriteString(hash, mystr)
	mysign := fmt.Sprintf("%x", hash.Sum(nil))
	if strings.ToLower(mysign) != strings.ToLower(sign) {
		logging.Debug("md5 check err:%s,%s,%s", mystr, mysign, sign)
		w.Write([]byte("4"))
		return
	}

	myaccid, err := db.GetMyAccountByAccountId(game.GetPlatNameByUrl(req.URL.Path), qid)
	if err != nil {
		logging.Debug("bill account err:%s", qid)
		w.Write([]byte("3"))
		return
	}
	zoneid, _ := strconv.Atoi(server_id)
	if zoneid == 11211 {
		zoneid = 11012
	}
	if zoneid == 11238 {
		zoneid = 11013
	}
	if zoneid == 11351 {
		zoneid = 11014
	}
	if zoneid == 11397 {
		zoneid = 11015
	}
	if zoneid == 11396 {
		zoneid = 11999
	}
	if zoneid == 11491 {
		zoneid = 11016
	}
	if zoneid == 11551 {
		zoneid = 11017
	}
	if zoneid == 11592 {
		zoneid = 11018
	}
	if zoneid == 11651 {
		zoneid = 11019
	}
	if zoneid == 11694 {
		zoneid = 11998
	}
	mygameid, _ := db.GetZonenameByZoneid(uint32(zoneid))
	logging.Debug("request bill ok:%s,%d,%d,%d", qid, myaccid, mygameid, zoneid)
	serverid := zoneid
	server_id = strconv.Itoa((mygameid << 16) + serverid)
	names := db.GetAllZoneCharNameByAccid(server_id, myaccid)
	if names == nil || len(names) == 0 {
		logging.Debug("bill checkname err:%d,%s", myaccid, qid)
		w.Write([]byte("3"))
		return
	}
	moneynum, _ := strconv.ParseFloat(order_amount, 32)
	game.Billing(game.GetGameNameByUrl(req.URL.Path), server_id, myaccid, 0, uint32(moneynum))
	w.Write([]byte("1"))
}
