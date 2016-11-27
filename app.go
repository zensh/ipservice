package main

import (
	"net"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/teambition/gear"
	"github.com/wangtuanjie/ip17mon"
)

var helpInfo = `
IP-Service 0.1.0
Qing Yan <admin@zensh.com>

OPTIONS:
	--port=<port>            Server port (default: ":8080").
	--data=<path>            IP data file path.
`

var portReg = regexp.MustCompile(`^\d+$`)

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

func home(ctx *gear.Context) error {
	html := `
<h1>IP Service</h1>
<p>Source Code: <a href="https://github.com/zensh/ipservice">github.com/zensh/ipservice</a></p>
<p>IP Database: <a href="http://www.ipip.net/about.html">IPIP.net</a></p>`
	return ctx.HTML(200, html)
}

func main() {
	port := "8080"
	dataPath := ""
	for _, arg := range os.Args {
		switch {
		case strings.HasPrefix(arg, "--port="):
			port = arg[7:]
		case strings.HasPrefix(arg, "--data="):
			dataPath = arg[7:]
		}
	}
	if portReg.MatchString(port) {
		port = ":" + port
	}
	if port == "" || dataPath == "" {
		os.Stdout.Write([]byte(helpInfo + "\n"))
		os.Exit(1)
	}

	// init IP db
	err := ip17mon.Init(dataPath)
	if err != nil {
		panic(err)
	}

	// start app
	app := gear.New()
	router := gear.NewRouter()
	router.Get("/json/:ip", jsonAPI)
	router.Otherwise(home)
	app.UseHandler(router)
	os.Stdout.Write([]byte("IP Service start at: " + port))
	app.Error(app.Listen(port))
}
