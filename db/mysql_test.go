package db

import (
	"fmt"
	"testing"
)

func TestMysql(t *testing.T) {
	InitDatabase("127.0.0.1:3306", "root", "123", "MonitorServer")
	fmt.Println(GetMyAccount("SDO", "wanghaijun"))
}
