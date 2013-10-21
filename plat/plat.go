package plat

import (
	"net/http"
)

type PlatLoginFunc func(w http.ResponseWriter, req *http.Request)

var (
	platMap map[string]PlatLoginFunc
)

func init() {
	platMap = make(map[string]PlatLoginFunc)
	platMap["/bw/juxian/auth"] = OnJuXianAuth
	platMap["/bw/juxian/bill"] = OnJuXianAuth
	platMap["/bw/juxian/check"] = OnJuXianAuth
	platMap["/bw/kw/auth"] = OnKuaiWanAuth
	platMap["/bw/kw/bill"] = OnKuaiWanBill
	platMap["/bw/kw/check"] = OnKuaiWanCheckName
	platMap["/bw/619"] = On619GameAuth
	platMap["/bw/619/bill"] = On619GameAuth
	platMap["/bw/619/check"] = On619GameAuth
}
func InitPlat() {
	for key, val := range platMap {
		http.HandleFunc(key, val)
	}
}
