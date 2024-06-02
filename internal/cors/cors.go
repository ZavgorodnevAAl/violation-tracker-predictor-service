package cors

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type Cors struct {
	targetUrl string
	port      string
}

func NewCors(targetUrl string, port string) Cors {
	return Cors{
		targetUrl: targetUrl,
		port:      port,
	}
}

func (c *Cors) Run() {
	target, err := url.Parse(c.targetUrl)
	if err != nil {
		log.Fatalf("Ошибка при разборе URL: %s", err)
	}

	proxy := httputil.NewSingleHostReverseProxy(target)

	modifyResponse := func(resp *http.Response) error {
		resp.Header.Set("Access-Control-Allow-Headers", "*")
		resp.Header.Set("Access-Control-Allow-Methods", "*")
		resp.Header.Set("Access-Control-Allow-Origin", "*")
		return nil
	}

	proxy.Transport = &transport{modifyResponse}

	go func() {
		err = http.ListenAndServe(c.port, proxy)
		if err != nil {
			log.Fatalf("Ошибка при запуске сервера: %s", err)
		}
	}()

}
