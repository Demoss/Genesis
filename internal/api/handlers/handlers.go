package handlers

import (
	"fmt"
	"genesis/global"
	"genesis/pkg"
	"net/http"
)

var con = pkg.NewConnector()

func GetBTC(w http.ResponseWriter, r *http.Request) {
	if global.Logged == 1 {
		resp := con.GetBTC()
		w.Write([]byte(fmt.Sprintf("1 %v = %v %v", resp.Ticker.Base, resp.Ticker.Price, resp.Ticker.Target)))
	} else {
		w.Write([]byte("you are not logged in"))
	}
}
