package common

import (
	"net"
	"strings"
	"unsafe"
)

func IsIPv4(ipAddr string) bool {
	ip := net.ParseIP(ipAddr)
	return ip != nil && strings.Count(ipAddr, ".") >= 2
}

func IsIPv6(ipAddr string) bool {
	ip := net.ParseIP(ipAddr)
	return ip != nil && strings.Count(ipAddr, ":") >= 2
}

func IsIPEqual(ipAddr1 string, ipAddr2 string) bool {
	ip1 := net.ParseIP(ipAddr1)
	ip2 := net.ParseIP(ipAddr1)
	return ip1 != nil && ip2 != nil && ip1.Equal(ip2)
}

func Str2bytes(s string) []byte {
	x := (*[2]uintptr)(unsafe.Pointer(&s))
	b := [3]uintptr{x[0], x[1], x[1]}
	return *(*[]byte)(unsafe.Pointer(&b))
}

func Bytes2str(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
