package engin

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"net"
	"time"
)

func TlsConfigClient(servername string, certfile, keyfile string) (*tls.Config, error) {
	var certs tls.Certificate
	var err error

	if certfile != "" && keyfile != "" {
		certs, err = tls.LoadX509KeyPair(certfile, keyfile)
		if err != nil {
			return nil, err
		}
	} else {
		certs, err = newCertificates([]string{})
		if err != nil {
			return nil, err
		}
	}
	return &tls.Config{
		MinVersion:         tls.VersionTLS12,
		MaxVersion:         tls.VersionTLS13,
		ServerName:         servername,
		InsecureSkipVerify: true,
		Certificates:       []tls.Certificate{certs},
	}, nil
}

func TlsConfigServer(certfile, keyfile string) (*tls.Config, error) {
	var certs tls.Certificate
	var err error

	if certfile != "" && keyfile != "" {
		certs, err = tls.LoadX509KeyPair(certfile, keyfile)
		if err != nil {
			return nil, err
		}
	} else {
		certs, err = newCertificates([]string{})
		if err != nil {
			return nil, err
		}
	}
	return &tls.Config{
		MinVersion:   tls.VersionTLS12,
		MaxVersion:   tls.VersionTLS13,
		Certificates: []tls.Certificate{certs},
		ClientAuth:   tls.RequestClientCert,
	}, nil
}

func newCertificates(address []string) (tls.Certificate, error) {
	max := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, _ := rand.Int(rand.Reader, max)
	subject := pkix.Name{
		Organization:       []string{"autoproxy co."},
		OrganizationalUnit: []string{"autoproxy"},
		CommonName:         "Autoproxy Programming",
	}

	ipAddress := make([]net.IP, 0)
	if address != nil {
		for _, v := range address {
			ipAddress = append(ipAddress, net.ParseIP(v))
		}
	}
	ipAddress = append(ipAddress, net.ParseIP("127.0.0.1"))

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject:      subject,
		NotBefore:    time.Now(),
		NotAfter:     time.Now().Add(10 * 365 * 24 * time.Hour),
		KeyUsage:     x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		IPAddresses:  ipAddress,
	}
	pk, _ := rsa.GenerateKey(rand.Reader, 2048)

	derBytes, _ := x509.CreateCertificate(rand.Reader, &template, &template, &pk.PublicKey, pk)

	certOut := bytes.NewBuffer(make([]byte, 0))
	pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})

	keyOut := bytes.NewBuffer(make([]byte, 0))
	pem.Encode(keyOut, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(pk)})

	return tls.X509KeyPair(certOut.Bytes(), keyOut.Bytes())
}
