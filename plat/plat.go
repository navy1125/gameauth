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
	platMap["/bw/juxian"] = OnAuthJuXian
	platMap["/bw/kw"] = OnAuthKuaiWan
	platMap["/bw/619"] = OnAuth619Game
}
func InitPlat() {
	for key, val := range platMap {
		http.HandleFunc(key, val)
	}
}
