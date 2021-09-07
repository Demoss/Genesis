package pkg

import (
	"encoding/json"
	"fmt"
	"genesis/pkg/logging"
	"genesis/pkg/resources"
	"log"
	"net/http"
)

type Connector struct {
	logger *logging.Logger
}

func NewConnector() *Connector {
	return &Connector{logger: logging.GetLogger()}
}

func (c *Connector) GetBTC() *resources.ResponceBTC {

	res, err := http.Get("https://api.cryptonator.com/api/ticker/btc-uah")
	if err != nil {
		log.Fatal("Failed")

	}

	c.logger.Info("getting BTC")
	var resp resources.ResponceBTC
	err = json.NewDecoder(res.Body).Decode(&resp)
	if err != nil {
		fmt.Println(err)
	}

	err = res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	return &resp

}
