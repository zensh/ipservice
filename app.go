package main

import (
	"encoding/json"
	"flag"
	"net"
	"net/http"
	"os"
	"regexp"
	"time"

	"github.com/teambition/gear"
	"github.com/teambition/gear/logging"
	"github.com/teambition/gear/middleware/favicon"
	"github.com/wangtuanjie/ip17mon"
)

var (
	portReg  = regexp.MustCompile(`^\d+$`)
	port     = flag.String("port", "8080", `Server port.`)
	dataPath = flag.String("data", "", "IP data file path.")
)

type result struct {
	IP      string
	Status  int
	Message string
	Data    interface{}
}

func jsonAPI(ctx *gear.Context) error {
	var ip net.IP
	var res result

	callback := ctx.Query("callback")
	ipStr := ctx.Param("ip")
	if ipStr == "" {
		ip = ctx.IP()
	} else {
		ip = net.ParseIP(ipStr)
	}

	if ip == nil {
		res = result{IP: "", Status: http.StatusBadRequest, Message: "Invalid IP format"}
	} else {
		loc, err := ip17mon.Find(ip.String())
		if err != nil {
			res = result{IP: ip.String(), Status: http.StatusNotFound, Message: err.Error()}
		} else {
			res = result{IP: ip.String(), Status: http.StatusOK, Data: loc}
		}
	}

	if callback == "" {
		return ctx.JSON(res.Status, res)
	}
	return ctx.JSONP(res.Status, callback, res)
}

func app(port, dataPath string) *gear.ServerListener {
	// init IP db
	err := ip17mon.Init(dataPath)
	if err != nil {
		panic(err)
	}

	// create app
	app := gear.New()

	// add favicon middleware
	app.Use(favicon.NewWithIco(faviconData))

	// add logger middleware
	logger := logging.New(os.Stdout)
	logger.SetLogConsume(func(log logging.Log, _ *gear.Context) {
		now := time.Now()
		delete(log, "Start")
		delete(log, "Type")
		switch res, err := json.Marshal(log); err == nil {
		case true:
			logger.Output(now, logging.InfoLevel, string(res))
		default:
			logger.Output(now, logging.WarningLevel, err.Error())
		}
	})
	app.UseHandler(logger)

	// add router middleware
	router := gear.NewRouter()
	router.Get("/json/:ip", jsonAPI)
	router.Otherwise(func(ctx *gear.Context) error {
		log := logging.FromCtx(ctx)
		log.Reset() // Reset log, don't logging for non-api request.
		return ctx.HTML(200, indexHTML)
	})
	app.UseHandler(router)

	// start app
	logging.Info("IP Service start " + port)
	return app.Start(port)
}

func main() {
	flag.Parse()
	if portReg.MatchString(*port) {
		*port = ":" + *port
	}
	if *port == "" || *dataPath == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}
	srv := app(*port, *dataPath)
	srv.Wait()
}
