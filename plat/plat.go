package plat

import (
	"github.com/navy1125/config"
	"net/http"
)

type PlatLoginFunc func(w http.ResponseWriter, req *http.Request)

var (
	platMap map[string]PlatLoginFunc
)

func init() {
	platMap = make(map[string]PlatLoginFunc)
	platMap["/zssj/juxian/auth"] = OnJuXianAuth
	platMap["/zssj/juxian/bill"] = OnJuXianAuth
	platMap["/zssj/juxian/check"] = OnJuXianAuth
	platMap["/zssj/kw/auth"] = OnKuaiWanAuth
	platMap["/zssj/kw/bill"] = OnKuaiWanBill
	platMap["/zssj/kw/check"] = OnKuaiWanCheckName
	platMap["/zssj/619"] = On619GameAuth
	platMap["/zssj/619/bill"] = On619GameAuth
	platMap["/zssj/619/check"] = On619GameAuth
}
func InitPlat() {
	http.Handle("/zssj/js/", http.StripPrefix("/bw/js/", http.FileServer(http.Dir(config.GetConfigStr("bw_plugin")+"/js"))))
	http.Handle("/zssj/images/", http.StripPrefix("/bw/images/", http.FileServer(http.Dir(config.GetConfigStr("bw_plugin")+"/images"))))
	for key, val := range platMap {
		http.HandleFunc(key, val)
	}
}
