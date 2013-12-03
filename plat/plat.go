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
	platMap["/zssj/juxian/bill"] = OnJuXianBill
	platMap["/zssj/juxian/check"] = OnJuXianCheckName

	platMap["/zssj/619wan/auth"] = OnJuXianAuth
	platMap["/zssj/619wan/bill"] = OnJuXianBill
	platMap["/zssj/619wan/check"] = OnJuXianCheckName
	platMap["/zssj/8558/auth"] = OnJuXianAuth
	platMap["/zssj/8558/bill"] = OnJuXianBill
	platMap["/zssj/8558/check"] = OnJuXianCheckName
	platMap["/zssj/zixia/auth"] = OnJuXianAuth
	platMap["/zssj/zixia/bill"] = OnJuXianBill
	platMap["/zssj/zixia/check"] = OnJuXianCheckName
	platMap["/zssj/tymmo/auth"] = OnJuXianAuth
	platMap["/zssj/tymmo/bill"] = OnJuXianBill
	platMap["/zssj/tymmo/check"] = OnJuXianCheckName
	platMap["/zssj/50pk/auth"] = OnJuXianAuth
	platMap["/zssj/50pk/bill"] = OnJuXianBill
	platMap["/zssj/50pk/check"] = OnJuXianCheckName
	platMap["/zssj/ufojoy/auth"] = OnJuXianAuth
	platMap["/zssj/ufojoy/bill"] = OnJuXianBill
	platMap["/zssj/ufojoy/check"] = OnJuXianCheckName
	platMap["/zssj/kuaiwan/auth"] = OnJuXianAuth
	platMap["/zssj/kuaiwan/bill"] = OnJuXianBill
	platMap["/zssj/kuaiwan/check"] = OnJuXianCheckName
	platMap["/zssj/9cb/auth"] = OnJuXianAuth
	platMap["/zssj/9cb/bill"] = OnJuXianBill
	platMap["/zssj/9cb/check"] = OnJuXianCheckName
	platMap["/zssj/ucjoy/auth"] = OnJuXianAuth
	platMap["/zssj/ucjoy/bill"] = OnJuXianBill
	platMap["/zssj/ucjoy/check"] = OnJuXianCheckName

	platMap["/zssj/kw/auth"] = OnKuaiWanAuth
	platMap["/zssj/kw/bill"] = OnKuaiWanBill
	platMap["/zssj/kw/check"] = OnKuaiWanCheckName
}
func InitPlat() {
	http.Handle("/zssj/js/", http.StripPrefix("/zssj/js/", http.FileServer(http.Dir(config.GetConfigStr("bw_plugin")+"/js"))))
	http.Handle("/zssj/images/", http.StripPrefix("/zssj/images/", http.FileServer(http.Dir(config.GetConfigStr("bw_plugin")+"/images"))))
	for key, val := range platMap {
		http.HandleFunc(key, val)
	}
}
