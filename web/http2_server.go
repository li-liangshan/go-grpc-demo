package web

import (
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo"
)

func request(c echo.Context) error {
	req := c.Request()
	format := "<pre><strong>Request Information</strong>\n\n<code>Protocol: %s\nHost: %s\nRemote Address: %s\nMethod: %s\nPath: %s\n</code></pre>"
	return c.HTML(http.StatusOK, fmt.Sprintf(format, req.Proto, req.Host, req.RemoteAddr, req.Method, req.URL.Path))
}

func stream(c echo.Context) error {
	res := c.Response()
	gone := res.CloseNotify()

	res.Header().Set(echo.HeaderContentType, echo.MIMETextHTMLCharsetUTF8)
	res.WriteHeader(http.StatusOK)
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	fmt.Fprint(res, "<pre><strong>Clock Stream</strong>\n\n<code>")
	for {
		fmt.Fprintf(res, "%v\n", time.Now())
		res.Flush()
		select {
		case <-ticker.C:
			fmt.Println("time:", time.Now())
		case <-gone:
			fmt.Println("request gone!")
			break
		}
	}
}

// go run $GOROOT/src/crypto/tls/generate_cert.go --host localhost 命令会生一个cert.pem 和key.pem 文件

func Http2ServerRun() {
	e := echo.New()
	e.GET("/request", request)
	e.GET("/stream", stream)
	e.Logger.Fatal(e.StartTLS(":1323", "cert.pem", "key.pem"))
}
