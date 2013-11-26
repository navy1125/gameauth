package db

import (
	"git.code4.in/logging"
	"github.com/xuyu/iconv"
	//"github.com/ziutek/mymysql/mysql"
	mysql "github.com/ziutek/mymysql/autorc"
	_ "github.com/ziutek/mymysql/thrsafe" // You may also use the native engine
	"strconv"
)

var (
	db           *mysql.Conn
	db_checkname *mysql.Conn
)

type CharName struct {
	CharId   string
	CharName string
}

func GetMyAccount(channel, account, accountid string) (string, uint32, error) {
	if accountid == "" {
		accountid = "0"
	}
	query_string := "select accid from ACCOUNT where account = '" + account + "' and channel = '" + channel + "'"
	row, res, err := db.QueryFirst(query_string)
	if err != nil {
		logging.Error("select err:%s", err.Error())
		return "", 0, err
	}
	if len(row) == 0 {
		return addMyAccount(channel, account, accountid)
	}
	return channel + ":" + account + ":" + accountid, uint32(row.Int(res.Map("accid"))), nil
}
func GetMyAccountByAccountId(channel, accountid string) (uint32, error) {
	if accountid == "" {
		accountid = "0"
	}
	query_string := "select accid from ACCOUNT where accountid = '" + accountid + "' and channel = '" + channel + "'"
	row, res, err := db.QueryFirst(query_string)
	if err != nil {
		logging.Error("select err:%s", err.Error())
		return 0, err
	}
	if len(row) == 0 {
		return 0, nil
	}
	return uint32(row.Int(res.Map("accid"))), nil
}
func addMyAccount(channel, account, accountid string) (string, uint32, error) {
	query_string := "insert into ACCOUNT (CHANNEL,ACCOUNT,ACCOUNTID) values( '" + channel + "' , '" + account + "' , " + accountid + ")"
	_, res, err := db.Query(query_string)
	if err != nil {
		logging.Error("insert err:%s", err.Error())
		return "", 0, err
	}
	return channel + ":" + account, uint32(res.InsertId()), nil
}
func InitDatabase(addr, name, password, dbname string) {
	db = mysql.New("tcp", "", addr, name, password, dbname)
}
func InitCheckNameDatabase(addr, name, password, dbname string) {
	db_checkname = mysql.New("tcp", "", addr, name, password, dbname)
}

func GetAllZoneCharNameByAccid(server_id string, myaccid uint32) []CharName {
	query_string := "select ID,NAME from CHARNAME where ACCID = " + strconv.Itoa(int(myaccid)) + " and  ZONE = " + server_id
	rows, res, err := db_checkname.Query(query_string)
	if err != nil {
		logging.Error("select err:%s", err.Error())
		return nil
	}
	map_id := res.Map("ID")
	map_name := res.Map("NAME")
	names := make([]CharName, len(rows))
	for i, row := range rows {
		names[i].CharId = row.Str(map_id)
		names[i].CharName = row.Str(map_name)
	}
	return names
}

func GetZonenameByZoneid(zoneid uint32) (int, string) {
	query_string := "select GAME,NAME from zoneInfo where zone = " + strconv.Itoa(int(zoneid))
	rows, res, err := db.Query(query_string)
	if err != nil {
		logging.Error("select err:%s,%s", err.Error(), zoneid)
		return 0, ""
	}
	if len(rows) == 0 {
		return 3055, "none"
	}
	game := res.Map("GAME")
	name := res.Map("NAME")
	convert, err1 := iconv.Open("GB2312", "UTF-8")
	for _, row := range rows {
		gamename := row.Str(name)
		if err1 == nil {
			var err2 error
			gamename, err2 = convert.ConvString(row.Str(name))
			if err2 != nil {
				gamename = row.Str(name)
			}
		}
		return row.Int(game), gamename
	}
	return 0, ""
}
