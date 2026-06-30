package cozylife

import (
	"encoding/json"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

const tcpPort = 5555

const (
	cmdQuery = 2
	cmdSet   = 3
)

type tcpClient struct {
	ip      string
	timeout time.Duration
	conn    net.Conn
}

func newClient(ip string, timeout time.Duration) *tcpClient {
	return &tcpClient{ip: ip, timeout: timeout}
}

func (c *tcpClient) connect() error {
	addr := fmt.Sprintf("%s:%d", c.ip, tcpPort)
	conn, err := net.DialTimeout("tcp", addr, c.timeout)
	if err != nil {
		return err
	}
	c.conn = conn
	return nil
}

func (c *tcpClient) disconnect() {
	if c.conn != nil {
		c.conn.Close()
		c.conn = nil
	}
}

func makeSN() string {
	return strconv.FormatInt(time.Now().UnixMilli(), 10)
}

func makePacket(cmd int, payload map[string]any) ([]byte, string) {
	sn := makeSN()
	var msg any
	switch cmd {
	case cmdQuery:
		msg = map[string]any{"attr": []int{0}}
	case cmdSet:
		keys := make([]int, 0, len(payload))
		for k := range payload {
			n, _ := strconv.Atoi(k)
			keys = append(keys, n)
		}
		msg = map[string]any{"attr": keys, "data": payload}
	default:
		msg = map[string]any{}
	}
	envelope := map[string]any{"pv": 0, "cmd": cmd, "sn": sn, "msg": msg}
	b, _ := json.Marshal(envelope)
	return append(b, '\r', '\n'), sn
}

// onlySend sends a SET command without waiting for a response; reconnects once on error.
func (c *tcpClient) onlySend(cmd int, payload map[string]any) error {
	pkt, _ := makePacket(cmd, payload)
	c.conn.SetWriteDeadline(time.Now().Add(c.timeout))
	if _, err := c.conn.Write(pkt); err != nil {
		c.disconnect()
		if err2 := c.connect(); err2 != nil {
			return err2
		}
		c.conn.SetWriteDeadline(time.Now().Add(c.timeout))
		_, err = c.conn.Write(pkt)
		return err
	}
	return nil
}

func (c *tcpClient) query() (map[string]any, error) {
	pkt, sn := makePacket(cmdQuery, nil)
	c.conn.SetWriteDeadline(time.Now().Add(c.timeout))
	if _, err := c.conn.Write(pkt); err != nil {
		c.disconnect()
		if err2 := c.connect(); err2 != nil {
			return nil, err2
		}
		c.conn.SetWriteDeadline(time.Now().Add(c.timeout))
		if _, err2 := c.conn.Write(pkt); err2 != nil {
			return nil, err2
		}
	}
	buf := make([]byte, 4096)
	for range 10 {
		c.conn.SetReadDeadline(time.Now().Add(c.timeout))
		n, err := c.conn.Read(buf)
		if err != nil {
			return nil, err
		}
		text := strings.TrimSpace(string(buf[:n]))
		if !strings.Contains(text, sn) {
			continue
		}
		var resp map[string]any
		if err := json.Unmarshal([]byte(text), &resp); err != nil {
			return nil, err
		}
		msg, _ := resp["msg"].(map[string]any)
		if msg == nil {
			continue
		}
		data, _ := msg["data"].(map[string]any)
		if data == nil {
			continue
		}
		return data, nil
	}
	return nil, fmt.Errorf("no matching response (sn=%s)", sn)
}

func (c *tcpClient) control(payload map[string]any) error {
	return c.onlySend(cmdSet, payload)
}
