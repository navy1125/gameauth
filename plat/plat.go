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
	platMap["/zssj/ucjoy/bill"] = OnJuUcjoyBill
	platMap["/zssj/ucjoy/check"] = OnJuXianCheckName
	platMap["/zssj/101yo/auth"] = OnJuXianAuth
	platMap["/zssj/101yo/bill"] = OnJuXianBill
	platMap["/zssj/101yo/check"] = OnJuXianCheckName

	platMap["/zssj/92you/auth"] = OnJuXianAuth
	platMap["/zssj/92you/bill"] = OnJuXianBill
	platMap["/zssj/92you/check"] = OnJuXianCheckName

	platMap["/zssj/007wan/auth"] = OnJuXianAuth
	platMap["/zssj/007wan/bill"] = OnJuXianBill
	platMap["/zssj/007wan/check"] = OnJuXianCheckName

	platMap["/zssj/myxiaoyao/auth"] = OnJuXianAuth
	platMap["/zssj/myxiaoyao/bill"] = OnJuXianBill
	platMap["/zssj/myxiaoyao/check"] = OnJuXianCheckName

	platMap["/zssj/youxiwangguo/auth"] = OnJuXianAuth
	platMap["/zssj/youxiwangguo/bill"] = OnJuXianBill
	platMap["/zssj/youxiwangguo/check"] = OnJuXianCheckName

	platMap["/zssj/717play/auth"] = OnJuXianAuth
	platMap["/zssj/717play/bill"] = OnJuXianBill
	platMap["/zssj/717play/check"] = OnJuXianCheckName

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
