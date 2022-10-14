package server

import (
	"fmt"
	"net"
	"strings"
	"time"

	"sdmht/lib/log"
)

const (
	CONN_LOGIN_TIMEOUT        = 10 * time.Second
	CONN_KA_TIMEOUT_MAX_CNT   = 3
	CONN_DEFAULT_READBUF_SIZE = 2048
	CONN_SEND_QUEUE_SIZE      = 20
)

var _ Conn = (*conn)(nil)

type conn struct {
	conn net.Conn
	// neid        uint64
	srv Server
	// sn          uint32
	// pendingMsgs sync.Map // sn -> msg
	sendC  chan *Message
	closeC chan struct{}
	// logined     bool
	// isProxy     bool
}

func NewConn(c net.Conn, srv Server) *conn {
	conn := &conn{
		conn: c,
		srv:  srv,
		// sn:     0,
		closeC: make(chan struct{}),
		// sendC:  make(chan *Message, CONN_SEND_QUEUE_SIZE),
	}

	return conn
}

func (c *conn) Serve() error {
	// go c.timeoutCheck()
	go c.sendRoutine()
	defer func() {
		c.conn.Close()
		close(c.closeC)
	}()
	return c.readRoutine()
}

func (c *conn) sendRoutine() {
	defer func() {
		log.S().Infow("quit send routine", "conn", c)
	}()

	for {
		select {
		case <-c.closeC:
			return
		case msg := <-c.sendC:
			if err := c.sendMsg(msg); err != nil {
				log.S().Errorw("send msg failed", "conn", c, "err", err)
			}
		}
	}
}

func (c *conn) sendMsg(msg *Message) error {
	b, err := msg.Marshal()
	if err != nil {
		log.S().Errorw("ne msg marshal failed", "conn", c, "err", err, "ne_msg", msg)
		return err
	}

	log.S().Debugw("send msg", "ne_msg", msg, "ne_msg_len", len(b), "conn", c)

	_, err = c.conn.Write(b)
	return err
}

func (c *conn) readRoutine() error {
	buf := make([]byte, CONN_DEFAULT_READBUF_SIZE)
	// var parser Parser

	defer func() {
		log.S().Infow("quit read routine", "conn", c)
	}()

	for {
		n, err := c.conn.Read(buf)
		if err != nil {
			if !strings.Contains(err.Error(), "use of closed network connection") {
				log.S().Errorw("conn read", "conn", c, "err", err)
			}
			return err
		}
		fmt.Printf("recv: %d, data: %x, str:%s\n", n, buf[:n], string(buf))
		// if _, err := parser.Write(buf[:n]); err != nil {
		// 	log.S().Errorw("parser write failed", "conn", c, "err", err)
		// 	return err
		// }
		// msgs, err := parser.Parse()
		// for _, msg := range msgs {
		// 	c.handleMsg(msg)
		// }
		// if err != nil {
		// 	log.S().Errorw("parser parse failed", "conn", c, "err", err)
		// 	return err
		// }
	}
}

func (c *conn) Close() error {
	return c.conn.Close()
}
