package utils

import (
	"net"
	"net/http"
	"strings"
)

func ClientIP(r *http.Request) string {
	IPAddress := r.Header.Get("X-Forwarded-For")
	if IPAddress == "" {
		IPAddress = r.Header.Get("X-Real-Ip")
	}
	if IPAddress == "" {
		// as http.Request.RemoteAddr contains ip:port combination
		pairs := strings.Split(r.RemoteAddr, ":")
		if len(pairs) == 2 {
			IPAddress = pairs[0]
		} else {
			IPAddress = r.RemoteAddr
		}
	}

	strs := strings.Split(IPAddress, ",")
	if len(strs) > 1 {
		IPAddress = strs[0] // take the first one
	}

	return strings.TrimSpace(IPAddress)
}

func ValidIP(ip string) bool {
	if ip == "" {
		return false
	}
	return net.ParseIP(ip) != nil
}

func IPToUInt32(ip string) uint32 {
	ipnr := net.ParseIP(ip)
	var sum uint32

	sum += uint32(ipnr[12]) << 24
	sum += uint32(ipnr[13]) << 16
	sum += uint32(ipnr[14]) << 8
	sum += uint32(ipnr[15])

	return sum
}
