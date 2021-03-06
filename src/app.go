package src

import (
	"net"
	"net/http"

	"github.com/teambition/gear"
	"github.com/teambition/gear/logging"
	"github.com/teambition/gear/middleware/favicon"
	"github.com/wangtuanjie/ip17mon"
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
		res = result{IP: "", Status: http.StatusBadRequest, Message: "invalid IP format"}
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

// New return a App instance
func New(dataPath string) *gear.App {
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
	app.UseHandler(logging.Default())

	// add router middleware
	router := gear.NewRouter()
	router.Get("/json/:ip", jsonAPI)
	router.Otherwise(func(ctx *gear.Context) error {
		log := logging.FromCtx(ctx)
		log.Reset() // Reset log, don't logging for non-api request.
		return ctx.HTML(200, indexHTML)
	})
	app.UseHandler(router)

	return app
}
