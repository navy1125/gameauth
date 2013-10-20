package db

import (
	"github.com/GXTime/logging"
	//"github.com/ziutek/mymysql/mysql"
	mysql "github.com/ziutek/mymysql/autorc"
	_ "github.com/ziutek/mymysql/native" // Native engine
)

var (
	db *mysql.Conn
)

func GetMyAccount(channel, account string) (string, uint32, error) {
	query_string := "select accid from ACCOUNT where account = '" + account + "' and channel = '" + channel + "'"
	row, res, err := db.QueryFirst(query_string)
	row, res, err = db.QueryFirst(query_string)
	if err != nil {
		logging.Error("select err:", err.Error())
		return "", 0, err
	}
	if len(row) == 0 {
		return addMyAccount(channel, account)
	}
	return channel + ":" + account, uint32(row.Int(res.Map("accid"))), nil
}
func addMyAccount(channel, account string) (string, uint32, error) {
	query_string := "insert into ACCOUNT (CHANNEL,ACCOUNT) values( '" + channel + "' , '" + account + "')"
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
