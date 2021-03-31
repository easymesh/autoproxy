package engin

import (
	"github.com/astaxie/beego/logs"
	"io"
	"net"
	"net/url"
	"strings"
	"sync"
	"time"
)

type AuthInfo struct {
	User  string
	Token string
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
	if strings.Index(host,":") == -1 {
		host += ":80"
	}
	return host
}

func WriteFull(w io.Writer, body []byte) error {
	begin := 0
	for  {
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

type connectCopy struct {
	in, out net.Conn
	timeout time.Duration
	flow  uint64
	close chan struct{}
	sync.WaitGroup
}

func (c *connectCopy)iocopy(in net.Conn, out net.Conn, statcall func(uint64))  {
	defer c.Done()
	buff := make([]byte, 8192)
	var err1 error
	var err2 error
	var cnt int
	for  {
		cnt, err1 = in.Read(buff)
		if cnt > 0 {
			statcall(uint64(cnt))
			c.flow += uint64(cnt)
			err2 = WriteFull(out, buff[:cnt])
		}
		if err1 != nil || err2 != nil {
			c.close <- struct{}{}
			break
		}
	}
}

func (c *connectCopy)timer()  {
	ticker := time.NewTicker(c.timeout)

	defer func() {
		c.Done()
		ticker.Stop()

		c.in.Close()
		c.out.Close()
	}()

	for  {
		old := c.flow
		select {
		case <- ticker.C: {
			new := c.flow
			if new == old {
				return
			}
		}
		case <- c.close: {
			return
		}
		}
	}
}

func ConnectCopyWithTimeout(in net.Conn, out net.Conn, tmout int, statcall func(uint64)) {
	c := new(connectCopy)
	c.timeout = time.Duration(tmout) * time.Second
	c.in = in
	c.out = out
	c.close = make(chan struct{}, 2)

	c.Add(3)
	go c.iocopy(in, out, statcall)
	go c.iocopy(out, in, statcall)
	go c.timer()
	c.Wait()

	logs.Info("connect %s <-> %s close", in.RemoteAddr(), out.RemoteAddr())
}
