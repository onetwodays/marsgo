package trace

import (
	"fmt"
	"net"
	"os"
	"sync"
	"time"
)

const (
	// MaxPackageSize .
	_maxPackageSize = 1024 * 32
	// safe udp package size // MaxPackageSize = 508 _dataChSize = 4096)
	// max memory usage 1024 * 32 * 4096 -> 128MB
	_dataChSize                 = 4096
	_defaultWriteChannalTimeout = 50 * time.Millisecond
	_defaultWriteTimeout        = 200 * time.Millisecond
)

// reporter trace reporter.
type reporter interface {
	WriteSpan(sp *Span) error
	Close() error
}

// newReport with network address
func newReport(network, address string, timeout time.Duration, protocolVersion int32) reporter {
	if timeout == 0 {
		timeout = _defaultWriteTimeout
	}
	report := &connReport{
		network: network,
		address: address,
		dataCh:  make(chan []byte, _dataChSize), //4096个字节数组
		done:    make(chan struct{}), //同步chan
		timeout: timeout,
		version: protocolVersion,
	}
	go report.daemon() //启动一个协程发数据
	return report
}

type connReport struct {
	version int32
	rmx     sync.RWMutex
	closed  bool

	network, address string

	dataCh chan []byte //字节数组的通道

	conn net.Conn

	done chan struct{}

	timeout time.Duration
}

func (c *connReport) daemon() {
	for b := range c.dataCh {
		c.send(b)  //每次发送一个字节数组
	}
	c.done <- struct{}{} //发送完了,给一个通知.
}

func (c *connReport) WriteSpan(sp *Span) error {
	data, err := marshalSpan(sp, c.version) //probuf格式发出
	if err != nil {
		return err
	}
	return c.writePackage(data)
}

func (c *connReport) writePackage(data []byte) error {
	c.rmx.RLock()
	defer c.rmx.RUnlock()
	if c.closed {
		return fmt.Errorf("report already closed")
	}
	if len(data) > _maxPackageSize {
		return fmt.Errorf("package too large length %d > %d", len(data), _maxPackageSize)
	}
	select {
	case c.dataCh <- data:
		return nil
	case <-time.After(_defaultWriteChannalTimeout):
		return fmt.Errorf("write to data channel timeout")
	}
}

func (c *connReport) Close() error {
	c.rmx.Lock()
	c.closed = true
	c.rmx.Unlock()

	t := time.NewTimer(time.Second)
	close(c.dataCh)
	select {
	case <-t.C: //有数据来了
		c.closeConn() //定时器关闭.
		return fmt.Errorf("close report timeout force close")
	case <-c.done:
		return c.closeConn() //数据发完就关闭
	}
}

func (c *connReport) send(data []byte) {
	if c.conn == nil {
		if err := c.reconnect(); err != nil {
			c.Errorf("connect error: %s retry after second", err)
			time.Sleep(time.Second)
			return
		}
	}
	c.conn.SetWriteDeadline(time.Now().Add(100 * time.Microsecond)) //100ms发完
	if _, err := c.conn.Write(data); err != nil {
		c.Errorf("write to conn error: %s, close connect", err)
		c.conn.Close()
		c.conn = nil
	}
}

func (c *connReport) reconnect() (err error) {
	c.conn, err = net.DialTimeout(c.network, c.address, c.timeout)
	return
}

func (c *connReport) closeConn() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

func (c *connReport) Errorf(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
}
