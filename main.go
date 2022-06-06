package main

import (
	"crypto/tls"
	"fmt"
	. "sslcheck/config"
	"time"
)

const expireDay = 10

func httpsHandshake(domain string) (conn *tls.Conn, err error) {
	conn, err = tls.Dial("tcp", domain+":443", nil)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func hostnameVerify(conn *tls.Conn, domain, hostname string) error {
	err := conn.VerifyHostname(domain)
	if err != nil {
		return err
	}
	return nil
}

func expireDateVerify(conn *tls.Conn) bool {
	// chain list, 0 -> domain ssl
	expire := conn.ConnectionState().PeerCertificates[0].NotAfter
	subDay := expire.Sub(time.Now()) / 24
	if subDay <= expireDay {
		return true
	}
	return false
}

func main() {
	InitAll()
	fmt.Println(Conf.GetString("expireDay"))
	conn, err := tls.Dial("tcp", "www.baidu.com:443", nil)
	if err != nil {
		panic("Server doesn't support SSL certificate err: " + err.Error())
	}
	err = conn.VerifyHostname("www.baidu.com")
	if err != nil {
		panic("hostname 不匹配" + err.Error())
	}
	// 获取到证书链
	expire := conn.ConnectionState().PeerCertificates[0].NotAfter
	fmt.Println(expire)
	fmt.Println(expire.Unix())
	now := time.Now()
	days := expire.Sub(now).Hours() / 24
	fmt.Println(int(days))
}
