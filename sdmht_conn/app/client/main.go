package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	account_entity "sdmht/account/svc/entity"
	sdmht_entity "sdmht/sdmht/svc/entity"
	"sdmht/sdmht_conn/svc/entity"

	"github.com/gorilla/websocket"
)

type HandlerFunc func(c *Client, params interface{})

var cmdList = map[string]HandlerFunc{
	"login":        MakeLoginReq,
	"newLineup":    MakeNewLineupReq,
	"findLineup":   MakeFindLineupReq,
	"updateLineup": MakeUpdateLineupReq,
	"deleteLineup": MakeDeleteLineupReq,
	"newMatch":     MakeNewMatchReq,
	"keepAlive":    MakeKeepAliveReq,
}

var (
	g_account = &account_entity.Account{}
	g_lineup  = []*sdmht_entity.Lineup{}
	g_match   = &sdmht_entity.Match{}
)

func main() {
	testClient()
}

type Client struct {
	NO   int
	conn *websocket.Conn

	SN              int
	closeChan       chan struct{}
	closeWaitChan   chan struct{}
	sendPayloadChan chan entity.Payload
	recvPayloadChan chan entity.Payload

	serverReqPayloadChan chan entity.Payload
}

func NewClient(no int, conn *websocket.Conn) *Client {
	return &Client{
		NO:              no,
		conn:            conn,
		closeChan:       make(chan struct{}, 1),
		closeWaitChan:   make(chan struct{}, 1),
		sendPayloadChan: make(chan entity.Payload, 10),
		recvPayloadChan: make(chan entity.Payload, 10),

		serverReqPayloadChan: make(chan entity.Payload, 10),
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

		payload, err := entity.RspMsgToPayload(msg)
		if err != nil {
			fmt.Println("MsgToPayload err", err)
			payload.PayloadType = entity.PayloadTypeRsp
			payload.Result = entity.NewResult(entity.ErrCodeMsgBadRequest, entity.ErrCodeMsgs[entity.ErrCodeMsgBadRequest])
			c.sendPayloadChan <- payload
			continue
		}
		if payload.PayloadType == entity.PayloadTypeRsp {
			if _, ok := payload.MsgContent.(*sdmht_entity.CommonRes); !ok {
				c.recvPayloadChan <- payload
			}
		} else {
			c.serverReqPayloadChan <- payload
		}
	}
}

// func (c *Client) handleReqMsg(req entity.Payload) (ret entity.Payload) {
// 	ret = entity.NewRespPayload(req, entity.ErrCodeMsgSuccess, "", struct{}{})
// 	return
// }

// func (c *Client) handleRespMsg(payload entity.Payload) {
// 	fmt.Println("handleRespMsg sn:", payload.SN)
// }

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
	errChan := make(chan error)
	fs := flag.NewFlagSet("webinar_conn demo", flag.ExitOnError)
	var (
		addr = fs.String("addr", "localhost:4090", "ws server addr")
	)
	if err := fs.Parse(os.Args[1:]); err != nil {
		panic(err)
	}

	wsServerUrl := "ws://" + *addr + "/sdmht"
	conn, _, err := websocket.DefaultDialer.Dial(wsServerUrl, nil)
	if err != nil {
		panic(err)
	}

	c := NewClient(1, conn)
	go c.Run()

	// go func() {
	// 	for {
	// 		cmdList["keepAlive"](c, nil)
	// 		time.Sleep(3 * time.Second)
	// 	}
	// }()
	fmt.Printf("waiting...")
	go func() {
		notExit := true
		for notExit {
			cmd := ""
			time.Sleep(1 * time.Millisecond)
			fmt.Printf("\nInput cmd > ")
			fmt.Scanln(&cmd)
			switch cmd {
			case "login":
				wechatID, username := "", ""
				fmt.Printf("Input login param [wechatID (username)]> ")
				fmt.Scanln(&wechatID, &username)
				cmdList[cmd](c, &sdmht_entity.LoginReq{
					WeChatID: wechatID,
					UserName: username,
				})
			case "newLineup":
				lineupName, unitsStr := "l1", "11;218;401"
				cardLibrarys, units := []int64{}, []int64{}
				// fmt.Printf("Input newLineup param [lineupName unitStr(eg:1;2;3)]> ")
				// fmt.Scanln(&lineupName, &unitStr)
				cards := rand.Perm(100)
				for i := 0; i < 20; i++ {
					cardLibrarys = append(cardLibrarys, int64(cards[i]))
				}
				for _, unitStr := range strings.Split(unitsStr, ";") {
					unit, _ := strconv.ParseInt(unitStr, 10, 64)
					units = append(units, unit)
				}
				cmdList[cmd](c, &sdmht_entity.NewLineupReq{
					Lineup: sdmht_entity.Lineup{
						AccountID:    g_account.ID,
						Name:         lineupName,
						CardLibrarys: cardLibrarys,
						Units:        units,
					},
				})
			case "findLineup":
				cmdList[cmd](c, &sdmht_entity.FindLineupReq{
					AccountID: g_account.ID,
				})
			case "updateLineup":
				lineupID, lineupName, unitsStr := uint64(0), "", ""
				cardLibrarys, units := []int64{}, []int64{}
				fmt.Printf("Input updateLineup param [lineupID lineupName unitStr(eg:1;2;3)]> ")
				fmt.Scanln(&lineupID, &lineupName, &unitsStr)
				cards := rand.Perm(100)
				for i := 0; i < 20; i++ {
					cardLibrarys = append(cardLibrarys, int64(cards[i]))
				}
				for _, unitStr := range strings.Split(unitsStr, ";") {
					unit, _ := strconv.ParseInt(unitStr, 10, 64)
					units = append(units, unit)
				}
				cmdList[cmd](c, &sdmht_entity.UpdateLineupReq{
					Lineup: sdmht_entity.Lineup{
						ID:           lineupID,
						AccountID:    g_account.ID,
						Name:         lineupName,
						CardLibrarys: cardLibrarys,
						Units:        units,
					},
				})
			case "deleteLineup":
				var lineupID uint64 = 0
				fmt.Printf("Input deleteLineup param [lineupID]> ")
				fmt.Scanln(&lineupID)
				cmdList[cmd](c, &sdmht_entity.DeleteLineupReq{
					ID:        lineupID,
					AccountID: g_account.ID,
				})
			case "newMatch":
			case "help":
				fmt.Println("following cmd can be used\n" +
					"|login|\n" +
					"|newLineup|findLineup|updateLineup|deleteLineup|\n" +
					"|newMatch|xxx|")
			case "exit":
				notExit = false
				errChan <- errors.New("cmd exit")
			case "\n", "\r", "":
				continue
			default:
				fmt.Println("invalid cmd!  ", "|", cmd, "|"+
					"\nfollowing cmd can be used\n"+
					"|login|\n"+
					"|newLineup|findLineup|updateLineup|deleteLineup|\n"+
					"|newMatch|xxx|")
			}
		}
	}()

	log.Println("quit", "err", <-errChan)
	_ = c.conn.Close()
	time.Sleep(1 * time.Second)
	fmt.Println("quit")
}

func MakeLoginReq(c *Client, req interface{}) {
	r := req.(*sdmht_entity.LoginReq)
	payload := entity.NewReqPayload(c.NewSN(), sdmht_entity.MsgTypeLogin, r)
	c.sendPayloadChan <- payload
	recvPayload := <-c.recvPayloadChan
	res := recvPayload.MsgContent.(*sdmht_entity.LoginRes)
	g_account.ID = res.AccountID
}

func MakeNewLineupReq(c *Client, req interface{}) {
	r := req.(*sdmht_entity.NewLineupReq)
	payload := entity.NewReqPayload(c.NewSN(), sdmht_entity.MsgTypeNewLineup, r)
	c.sendPayloadChan <- payload
	MakeFindLineupReq(c, &sdmht_entity.FindLineupReq{
		AccountID: g_account.ID,
	})
}

func MakeFindLineupReq(c *Client, req interface{}) {
	r := req.(*sdmht_entity.FindLineupReq)
	payload := entity.NewReqPayload(c.NewSN(), sdmht_entity.MsgTypeFindLineup, r)
	c.sendPayloadChan <- payload
	recvPayload := <-c.recvPayloadChan
	res := recvPayload.MsgContent.(*sdmht_entity.FindLineupRes)
	g_lineup = res.Lineups
	fmt.Println("========== Lineups ============")
	for _, lineup := range g_lineup {
		fmt.Println(lineup.ID, lineup)
	}
}

func MakeUpdateLineupReq(c *Client, req interface{}) {
	r := req.(*sdmht_entity.UpdateLineupReq)
	payload := entity.NewReqPayload(c.NewSN(), sdmht_entity.MsgTypeUpdateLineup, r)
	c.sendPayloadChan <- payload
	MakeFindLineupReq(c, &sdmht_entity.FindLineupReq{
		AccountID: g_account.ID,
	})
}

func MakeDeleteLineupReq(c *Client, req interface{}) {
	r := req.(*sdmht_entity.DeleteLineupReq)
	payload := entity.NewReqPayload(c.NewSN(), sdmht_entity.MsgTypeDeleteLineup, r)
	c.sendPayloadChan <- payload
	// recvPayload := <-c.recvPayloadChan
	// _ = recvPayload.MsgContent.(*sdmht_entity.CommonRes)
}

func MakeNewMatchReq(c *Client, _ interface{}) {
	payload := entity.NewReqPayload(c.NewSN(), sdmht_entity.MsgTypeNewMatch, &sdmht_entity.NewMatchReq{})
	c.sendPayloadChan <- payload
}

func MakeKeepAliveReq(c *Client, _ interface{}) {
	payload := entity.NewReqPayload(c.NewSN(), sdmht_entity.MsgTypeKeepAlive, &sdmht_entity.KeepAliveReq{})
	c.sendPayloadChan <- payload
}
