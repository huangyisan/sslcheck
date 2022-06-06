package main

import (
	"crypto/tls"
	"fmt"
	. "sslcheck/config"
	"sync"
	"time"
)

const expireDay = 10

var (
	wg sync.WaitGroup
)

func httpsHandshake(domain string) (conn *tls.Conn, err error) {
	conn, err = tls.Dial("tcp", domain+":443", nil)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func hostnameVerify(conn *tls.Conn, domain string) error {
	err := conn.VerifyHostname(domain)
	if err != nil {
		return err
	}
	return nil
}

func expireDateVerify(conn *tls.Conn) bool {
	// chain list, 0 -> domain ssl
	expire := conn.ConnectionState().PeerCertificates[0].NotAfter
	subDay := expire.Sub(time.Now()).Hours() / 24
	if subDay <= expireDay {
		return true
	}
	return false
}

func main() {
	InitAll()

	domainList := Conf.GetStringSlice("domains")
	fmt.Println(domainList)
	successChan := make(chan string, len(domainList))
	failedChan := make(chan string, len(domainList))
	expireChan := make(chan string, len(domainList))
	for _, domain := range domainList {
		wg.Add(1)
		testDomain := domain
		go func(testDomain string) {
			defer wg.Done()
			conn, err := httpsHandshake(testDomain)
			if err != nil {
				failedChan <- testDomain
				return
			}
			err = hostnameVerify(conn, testDomain)
			if err != nil {
				failedChan <- testDomain
				return
			}
			if expireDateVerify(conn) {
				expireChan <- testDomain
				return
			} else {
				successChan <- testDomain
				return
			}
		}(testDomain)
	}

	go func() {
		wg.Wait()
		close(successChan)
		close(failedChan)
		close(expireChan)
	}()

	for v := range expireChan {
		fmt.Printf("%v\n", v)
	}

}
