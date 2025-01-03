package gateway

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
)

// configureTLS sets up a TLS configuration using the given certificate, key, and CA.
func configureTLS(clientCert, clientKey, caCert string, skipCertVerify bool) (*tls.Config, error) {
	certificates := []tls.Certificate{}
	if clientCert != "" && clientKey != "" {
		cert, err := tls.X509KeyPair([]byte(clientCert), []byte(clientKey))
		if err != nil {
			return nil, err
		}
		certificates = append(certificates, cert)
	}

	caCertPool := x509.NewCertPool()
	if caCert != "" {
		if ok := caCertPool.AppendCertsFromPEM([]byte(caCert)); !ok {
			return nil, fmt.Errorf("failed to append CA certificate to pool")
		}
	}

	return &tls.Config{
		Certificates:       certificates,
		RootCAs:            caCertPool,
		InsecureSkipVerify: skipCertVerify,
	}, nil
}
