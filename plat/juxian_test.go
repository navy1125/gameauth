package plat

import (
	"../db"
	"fmt"
	"github.com/navy1125/config"
	"net/http"
	"testing"
)

func TestMysql(t *testing.T) {
	//config.SetConfig("bw_juxian_key", "value")
	if err := config.LoadFromFile("key.xml", "Key"); err != nil {
		fmt.Println(err)
		return
	}
	db.InitDatabase("127.0.0.1:3306", "root", "123", "MonitorServer")
	http.HandleFunc("/bw/juxian", OnAuthJuxian)
	err := http.ListenAndServe(":12346", nil)
	if err != nil {
		fmt.Println("ListenAndServe: ", err)
	}
}
