package db

import (
	"github.com/GXTime/logging"
	//"github.com/ziutek/mymysql/mysql"
	mysql "github.com/ziutek/mymysql/autorc"
	_ "github.com/ziutek/mymysql/native" // Native engine
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
	row, res, err = db.QueryFirst(query_string)
	if err != nil {
		logging.Error("select err:", err.Error())
		return "", 0, err
	}
	if len(row) == 0 {
		return addMyAccount(channel, account, accountid)
	}
	return channel + ":" + accountid + ":" + account, uint32(row.Int(res.Map("accid"))), nil
}
func addMyAccount(channel, account, accountid string) (string, uint32, error) {
	query_string := "insert into ACCOUNT (CHANNEL,ACCOUNT,ACCOUNTID) values( '" + channel + "' , '" + account + "' , " + accountid + ")"
	_, res, err := db.Query(query_string)
	if err != nil {
		logging.Error("insert err:", err.Error())
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

func GetAllCharNameByAccid(myaccid uint32) []CharName {
	query_string := "select ID,NAME from CHARNAME where ACCID = " + strconv.Itoa(int(myaccid))
	rows, res, err := db.Query(query_string)
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
