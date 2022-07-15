package utils

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"github.com/ShevchenkoVadim/helperlib/config"
	"io/ioutil"
	"log"
	"os"
)

var tlsContext *tls.Config = nil

func CheckFileIsExist(filePath string) bool {
	_, err := os.Stat(filePath)
	if err == nil {
		return true
	} else if errors.Is(err, os.ErrNotExist) {
		return false
	} else {
		log.Println(err)
		return false
	}

}

func GetTlsContext() (*tls.Config, error) {
	if tlsContext == nil {
		caCert, err := ioutil.ReadFile(config.C.SSLCert.SslCA)
		if err != nil {
			return nil, err
		}

		cert, err := tls.LoadX509KeyPair(config.C.SSLCert.SslPem, config.C.SSLCert.SslKey)
		if err != nil {
			return nil, err
		}

		rootCAs := x509.NewCertPool()
		rootCAs.AppendCertsFromPEM(caCert)

		tlsContext = &tls.Config{
			RootCAs:            rootCAs,
			Certificates:       []tls.Certificate{cert},
			InsecureSkipVerify: true,
			//ServerName:   "localhost", // Optional
		}
	} else {
		return tlsContext, nil
	}
	return tlsContext, nil
}

func LogWrapper(msg ...any) {
	if config.C.Debug {
		log.Println(msg...)
	}
}
