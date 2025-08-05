package client

import (
	"ayode.org/visor/config"
	"net/http"
)

func (c *Client) SetCookies(cookies []config.Cookie) {
	var httpCookies []*http.Cookie
	for _, cookie := range cookies {
		httpCookie := http.Cookie{
			Name:  cookie.Name,
			Value: cookie.Value,
		}
		httpCookies = append(httpCookies, &httpCookie)
	}
	c.jar.SetCookies(c.u, httpCookies)
}

func (c *Client) SetReqCookies(req *http.Request, cookies []config.Cookie) {
	for _, cookie := range cookies {
		httpCookie := http.Cookie{
			Name:  cookie.Name,
			Value: cookie.Value,
		}
		req.AddCookie(&httpCookie)
	}
}
