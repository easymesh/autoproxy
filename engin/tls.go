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
	max := new(big.Int).Lsh(big.NewInt(1), 128)   //把 1 左移 128 位，返回给 big.Int
	serialNumber, _ := rand.Int(rand.Reader, max) //返回在 [0, max) 区间均匀随机分布的一个随机值
	subject := pkix.Name{                         //Name代表一个X.509识别名。只包含识别名的公共属性，额外的属性被忽略。
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
		SerialNumber: serialNumber, // SerialNumber 是 CA 颁布的唯一序列号，在此使用一个大随机数来代表它
		Subject:      subject,
		NotBefore:    time.Now(),
		NotAfter:     time.Now().Add(10 * 365 * 24 * time.Hour),
		KeyUsage:     x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature, //KeyUsage 与 ExtKeyUsage 用来表明该证书是用来做服务器认证的
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},               // 密钥扩展用途的序列
		IPAddresses:  ipAddress,
	}
	pk, _ := rsa.GenerateKey(rand.Reader, 2048) //生成一对具有指定字位数的RSA密钥

	//CreateCertificate基于模板创建一个新的证书
	//第二个第三个参数相同，则证书是自签名的
	//返回的切片是DER编码的证书
	derBytes, _ := x509.CreateCertificate(rand.Reader, &template, &template, &pk.PublicKey, pk) //DER 格式

	certOut := bytes.NewBuffer(make([]byte, 0))
	pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})

	keyOut := bytes.NewBuffer(make([]byte, 0))
	pem.Encode(keyOut, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(pk)})

	return tls.X509KeyPair(certOut.Bytes(), keyOut.Bytes())
}
