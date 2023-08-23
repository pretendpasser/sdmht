package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	sdmht_entity "sdmht/sdmht/svc/entity"
	"sdmht/sdmht_conn/svc/entity"

	"github.com/gorilla/websocket"
)

type HandlerFunc func(c *Client, n uint64, n1 uint64)

var cmdList = map[string]HandlerFunc{
	// "getEvent":       MakeGetEventReq,
	// "joinEvent":      MakeJoinEventReq,
	// "leaveEvent":     MakeLeaveEventReq,
	// "switchspeecher": MakeSwitchSpeecherReq,
	"newMatch":  MakeNewMatchReq,
	"keepAlive": MakeKeepAliveReq,
}

func main() {
	testClient()
}

type Client struct {
	NO   int
	conn *websocket.Conn

	SN              int
	token           string
	closeChan       chan struct{}
	closeWaitChan   chan struct{}
	sendPayloadChan chan entity.Payload
}

func NewClient(no int, conn *websocket.Conn, token string) *Client {
	return &Client{
		NO:              no,
		conn:            conn,
		token:           token,
		closeChan:       make(chan struct{}, 1),
		closeWaitChan:   make(chan struct{}, 1),
		sendPayloadChan: make(chan entity.Payload, 10),
	}
}

//func (c *Client) String() string {
//	return fmt.Sprintf("client[%d]", c.NO)
//}

func (c *Client) Run() {
	fmt.Printf("client[%d] run...\n", c.NO)
	defer func() {
		fmt.Printf("client[%d] quit\n", c.NO)
	}()
	wg := sync.WaitGroup{}

	wg.Add(1)

	go func() {
		defer wg.Done()
		c.HandleReadMsg()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		c.HandleSendMsg()
	}()

	wg.Wait()

	c.closeWaitChan <- struct{}{}
}

func (c *Client) HandleReadMsg() {
	for {
		_, msg, err := c.conn.ReadMessage()
		if err != nil {
			fmt.Printf("client[%d] ReadMessage err: %s \n", c.NO, err.Error())
			close(c.closeChan)
			break
		}

		fmt.Println("client[", c.NO, "] recv payload", string(msg))

		payload, err := entity.MsgToPayload(msg)
		if err != nil {
			fmt.Println("MsgToPayload err", err)
			payload.PayloadType = entity.PayloadTypeRsp
			payload.Result = entity.NewResult(entity.ErrCodeMsgBadRequest, entity.ErrCodeMsgs[entity.ErrCodeMsgBadRequest])
			c.sendPayloadChan <- payload
			continue
		}
	}
}

func (c *Client) handleReqMsg(req entity.Payload) (ret entity.Payload) {
	ret = entity.NewRespPayload(req, entity.ErrCodeMsgSuccess, "", struct{}{})
	return
}

func (c *Client) handleRespMsg(payload entity.Payload) {
	fmt.Println("handleRespMsg sn:", payload.SN)
}

func (c *Client) HandleSendMsg() {
	for {
		select {
		case <-c.closeChan:
			return
		case payload := <-c.sendPayloadChan:
			fmt.Println("client[", c.NO, "] send payload:", payload)
			msg := entity.PayloadToMsg(payload)
			err := c.conn.WriteMessage(1, msg)
			if err != nil {
				fmt.Println("WriteMessage err", err)
			}
		}
	}
}

func (c *Client) NewSN() int {
	c.SN++
	return c.SN
}

func testClient() {
	fs := flag.NewFlagSet("webinar_conn demo", flag.ExitOnError)
	var (
		token         = fs.String("token", "", "token")
		willDo        = fs.String("cmd", "", "cmd")
		addr          = fs.String("addr", "localhost:4090", "ws server addr")
		eventID       = fs.Uint64("eventid", 0, "eventid")
		newSpeecherID = fs.Uint64("newspeecherid", 0, "newspeecherid")
	)
	if err := fs.Parse(os.Args[1:]); err != nil {
		panic(err)
	}

	wsServerUrl := "ws://" + *addr + "/sdmht"
	conn, _, err := websocket.DefaultDialer.Dial(wsServerUrl, nil)
	if err != nil {
		panic(err)
	}

	c := NewClient(1, conn, *token)
	go c.Run()

	cmdList[*willDo](c, *eventID, *newSpeecherID)
	go func() {
		for {
			cmdList["keepAlive"](c, 0, 0)
			time.Sleep(3 * time.Second)
		}
	}()

	fmt.Println("recv quit signal")
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer stop()
	<-ctx.Done()
	_ = c.conn.Close()
	time.Sleep(1 * time.Second)
	fmt.Println("quit")
}

// func MakeGetEventReq(c *Client, _ uint64, _ uint64) {
// 	payload := entity.NewReqPayload(c.NewSN(), webinar_entity.MsgTypeGetEvents, c.token, webinar_entity.FindRelatedEventReq{
// 		Begin: 0,
// 		End:   4828089600,
// 	})
// 	c.sendPayloadChan <- payload
// }

// func MakeJoinEventReq(c *Client, eventID uint64, _ uint64) {
// 	payload := entity.NewReqPayload(c.NewSN(), webinar_entity.MsgTypeJoinEvent, c.token, webinar_entity.JoinEventReq{
// 		EventID: eventID,
// 	})
// 	c.sendPayloadChan <- payload
// }

// func MakeLeaveEventReq(c *Client, eventID uint64, _ uint64) {
// 	payload := entity.NewReqPayload(c.NewSN(), webinar_entity.MsgTypeLeaveEvent, c.token, webinar_entity.LeaveEventReq{
// 		EventID: eventID,
// 	})
// 	c.sendPayloadChan <- payload
// }

func MakeNewMatchReq(c *Client, _ uint64, _ uint64) {
	payload := entity.NewReqPayload(c.NewSN(), sdmht_entity.MsgTypeNewMatchRequest, c.token, &sdmht_entity.NewMatchReq{
		Operator:   11,
		CardConfig: 22,
	})
	c.sendPayloadChan <- payload
}

func MakeKeepAliveReq(c *Client, _ uint64, _ uint64) {
	payload := entity.NewReqPayload(c.NewSN(), sdmht_entity.MsgTypeKeepAliveRequest, c.token, &sdmht_entity.KeepAliveReq{})
	c.sendPayloadChan <- payload
}
