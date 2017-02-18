package pkg

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/DavidCai1993/request"
	"github.com/stretchr/testify/assert"
	"github.com/teambition/gear"
)

func TestGearApp(t *testing.T) {
	app := New("../data/17monipdb.dat")
	app.Start(":8080")
	defer app.Close()

	t.Run("home", func(t *testing.T) {
		assert := assert.New(t)

		res, err := request.Get("http://127.0.0.1:8080").End()
		assert.Nil(err)
		assert.Equal(res.StatusCode, http.StatusOK)
		assert.Equal(res.Header.Get(gear.HeaderContentType), gear.MIMETextHTMLCharsetUTF8)
		body, _ := ioutil.ReadAll(res.Body)
		assert.Equal(string(body), indexHTML)
	})

	t.Run("query IP in JSON", func(t *testing.T) {
		assert := assert.New(t)

		res, err := request.Get("http://127.0.0.1:8080/json/8.8.8.8").End()
		assert.Nil(err)
		assert.Equal(res.StatusCode, http.StatusOK)
		assert.Equal(res.Header.Get(gear.HeaderContentType), gear.MIMEApplicationJSONCharsetUTF8)
		body, _ := ioutil.ReadAll(res.Body)
		assert.True(strings.Contains(string(body), `"IP":"8.8.8.8"`))
		rt := &result{}
		json.Unmarshal(body, rt)
		fmt.Println(rt)
		assert.Equal(rt.IP, "8.8.8.8")
		assert.Equal(rt.Status, 200)
	})

	t.Run("query IP in JSONP", func(t *testing.T) {
		assert := assert.New(t)

		res, err := request.Get("http://127.0.0.1:8080/json/8.8.8.8?callback=test").End()
		assert.Nil(err)
		assert.Equal(res.StatusCode, http.StatusOK)
		assert.Equal(res.Header.Get(gear.HeaderContentType), gear.MIMEApplicationJavaScriptCharsetUTF8)
		body, _ := ioutil.ReadAll(res.Body)
		assert.True(strings.Contains(string(body), "test"))
		assert.True(strings.Contains(string(body), `"IP":"8.8.8.8"`))
	})
}
