package engin

import (
	"fmt"
	"io"
	"math/rand"
	"net"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/astaxie/beego/logs"
)

type AuthInfo struct {
	User  string
	Token string
}

var (
	uuid string
)

func init() {
	rand.Seed(time.Now().UnixNano())
	uuid = fmt.Sprintf("%08X-%016X-%08X", rand.Int31(), rand.Int63(), rand.Int31())
}

func GetUUID() string {
	return uuid
}

func IsConnect(address string, timeout int) bool {
	conn, err := net.DialTimeout("tcp", address, time.Duration(timeout)*time.Second)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}

func Address(u *url.URL) string {
	host := u.Host
	if strings.Index(host, ":") == -1 {
		host += ":80"
	}
	return host
}

func WriteFull(w io.Writer, body []byte) error {
	begin := 0
	for {
		cnt, err := w.Write(body[begin:])
		if cnt > 0 {
			begin += cnt
		}
		if begin >= len(body) {
			return err
		}
		if err != nil {
			return err
		}
	}
}

func iocopy(c *sync.WaitGroup, in net.Conn, out net.Conn) {
	defer c.Done()

	size, err := io.Copy(in, out)
	if size == 0 && err != nil {
		logs.Warn("io copy fail", err.Error())
	} else {
		StatUpdate(0, size)
	}

	in.Close()
	out.Close()
}

func Connect(acc *HttpAccess, in net.Conn, out net.Conn) {
	var wg sync.WaitGroup

	StatUpdate(1, 0)

	wg.Add(2)
	go iocopy(&wg, in, out)
	go iocopy(&wg, out, in)
	wg.Wait()

	StatUpdate(-1, 0)

	logs.Info("connect %s <-> %s close", in.RemoteAddr(), out.RemoteAddr())
}
