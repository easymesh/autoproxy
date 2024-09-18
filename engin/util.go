package engin

import (
	"fmt"
	"io"
	"math/rand"
	"net"
	"net/url"
	"os"
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

/* covert HOST to ip address with port
[2409:8a28:c4d:a181::42d]:8070 -> [2409:8a28:c4d:a181::42d]:8070
[2409:8a28:c4d:a181::42d]:8070 -> [2409:8a28:c4d:a181::42d]:8070
2409:8a28:c4d:a181::42d:8070 -> [2409:8a28:c4d:a181::42d]:8070
2409:8a28:c4d:a181::42d:8070 -> [2409:8a28:c4d:a181::42d]:8070
2409:8a28:c4d:a181::42d -> [2409:8a28:c4d:a181::42d]:80
2409:8a28:c4d:a181::42d -> [2409:8a28:c4d:a181::42d]:443
192.168.3.1:111 -> 192.168.3.1:111
192.168.3.1:111 -> 192.168.3.1:111
192.168.1.1 -> 192.168.1.1:80
192.168.1.1 -> 192.168.1.1:443
demo.abc.a:111 -> demo.abc.a:111
demo.abc.a:111 -> demo.abc.a:111
demo.abc -> demo.abc:80
demo.abc -> demo.abc:443
*/

func parseAddress(schema string, host string) string {
	defaut_port := "80"
	if strings.ToLower(schema) == "https" {
		defaut_port = "443"
	}

	count := strings.Count(host, ":")
	if count == 0 {
		return host + ":" + defaut_port
	}

	if count > 1 {
		index := strings.LastIndex(host, ":")
		addr := host[:index]
		port := host[index+1:]

		ip := net.ParseIP(addr)
		if ip != nil {
			if len(ip.To4()) == net.IPv4len {
				return ip.String() + ":" + port
			} else {
				return "[" + ip.String() + "]:" + port
			}
		}

		ip = net.ParseIP(host)
		if ip != nil {
			if len(ip.To4()) == net.IPv4len {
				return ip.String() + ":" + defaut_port
			} else {
				return "[" + ip.String() + "]:" + defaut_port
			}
		}
	}
	return host
}

func Address(u *url.URL) string {
	addr := parseAddress(u.Scheme, u.Host)
	logs.Info("Parse URL %s to %s", u.String(), addr)
	return addr
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
		if err != io.EOF {
			logs.Info("io copy fail", err.Error())
		}
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

func GetFileTimestamp(file string) time.Time {
	info, err := os.Stat(file)
	if err != nil {
		logs.Warning(err.Error())
		return time.Time{}
	}
	return info.ModTime()
}
