package tool

import (
	"net"
	"net/http"
	"strings"

	"github.com/zeromicro/go-zero/rest/httpx"
)

// GetIP returns request real ip.
func GetClientIP(r *http.Request) string {

	ip := httpx.GetRemoteAddr(r)
	for _, v := range strings.Split(ip, ",") {
		if net.ParseIP(strings.TrimSpace(v)) != nil {
			return v
		}
	}

	ip = strings.TrimSpace(r.Header.Get("X-Real-IP"))
	if net.ParseIP(ip) != nil {
		return ip
	}

	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return ""
	}

	if net.ParseIP(ip) != nil {
		return ip
	}

	return ""
}
