package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"fmt"

	"github.com/teambition/gear"
	"github.com/wangtuanjie/ip17mon"
)

type config struct {
	Port       string
	IPDataFile string
}

type result struct {
	Status  int
	Message string
	Data    interface{}
}

func readConfig() *config {
	cfg := &config{}
	data, err := ioutil.ReadFile("./config.json")
	if err == nil {
		err = json.Unmarshal(data, cfg)
	}
	if err != nil {
		panic(err)
	}
	return cfg
}

func ipservice(ctx *gear.Context) error {
	ip := ctx.Param("ip")
	if ip == "" {
		return nil
	}

	callback := ctx.Query("callback")
	loc, err := ip17mon.Find(ip)
	res := result{Status: http.StatusOK, Data: loc}

	if err != nil {
		res.Status = http.StatusNotFound
		res.Message = fmt.Sprintf(`%s: %s`, err.Error(), ip)
	}
	if callback == "" {
		return ctx.JSON(res.Status, res)
	}
	return ctx.JSONP(res.Status, callback, res)
}

func home(ctx *gear.Context) error {
	html := `<h1>IP Service</h1>`
	return ctx.HTML(200, html)
}

func main() {
	// init config
	cfg := readConfig()
	if cfg.Port == "" {
		cfg.Port = ":3000"
	}

	// init IP db
	err := ip17mon.Init(cfg.IPDataFile)
	if err != nil {
		panic(err)
	}

	// start app
	app := gear.New()
	router := gear.NewRouter()
	router.Get("/ip/:ip", ipservice)
	router.Otherwise(home)
	app.UseHandler(router)
	app.Error(app.Listen(cfg.Port))
}
