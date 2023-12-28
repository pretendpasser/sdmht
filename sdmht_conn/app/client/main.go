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
	"getMatch":     MakeGetMatchReq,
	"joinMatch":    MakeJoinMatchReq,
	"operate":      MakeOperateReq,
	"keepAlive":    MakeKeepAliveReq,
}

var (
	g_account  = &account_entity.Account{}
	g_lineup   = []*sdmht_entity.Lineup{}
	g_match    = &sdmht_entity.Match{}
	g_playerID = 0
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
		c.HandlerReadReqMsg()
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

		// fmt.Println("client[", c.NO, "] recv payload", string(msg))

		payload, err := entity.RspMsgToPayload(msg)
		if err != nil {
			fmt.Println("MsgToPayload err", err)
			payload.PayloadType = entity.PayloadTypeRsp
			payload.Result = entity.NewResult(entity.ErrCodeMsgBadRequest, entity.ErrCodeMsgs[entity.ErrCodeMsgBadRequest])
			c.sendPayloadChan <- payload
			continue
		}
		fmt.Println("Mention Info:", payload.Result)
		if payload.PayloadType == entity.PayloadTypeRsp {
			if _, ok := payload.MsgContent.(*sdmht_entity.CommonRes); !ok {
				c.recvPayloadChan <- payload
			}
		} else {
			payload, err := entity.MsgToPayload(msg)
			if err != nil {
				fmt.Println("MsgToPayload err", err)
				payload.PayloadType = entity.PayloadTypeRsp
				payload.Result = entity.NewResult(entity.ErrCodeMsgBadRequest, entity.ErrCodeMsgs[entity.ErrCodeMsgBadRequest])
				c.sendPayloadChan <- payload
				continue
			}
			c.serverReqPayloadChan <- payload
		}
	}
}

func (c *Client) HandlerReadReqMsg() {
	for {
		payload := <-c.serverReqPayloadChan
		if payload.MsgType == sdmht_entity.MsgTypeSyncMatch {
			res := payload.MsgContent.(*sdmht_entity.Match)
			g_match = res
			ShowScene(g_match)
		}
		respPayload := entity.NewRespPayload(payload, entity.ErrCodeMsgSuccess, "", &sdmht_entity.CommonRes{})
		c.sendPayloadChan <- respPayload
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
			if g_match.ID != 0 && cmd != "operate" && cmd != "op" {
				fmt.Println("You are in a match. Use 'operate [event] [from] [to]'")
				continue
			}
			switch cmd {
			case "login":
				wechatID, username := "", ""
				fmt.Printf("Input login param [wechatID (username)]> ")
				fmt.Scanln(&wechatID, &username)
				if wechatID == "" {
					wechatID = "1"
				}
				if username == "" {
					username = "default"
				}
				cmdList[cmd](c, &sdmht_entity.LoginReq{
					WeChatID: wechatID,
					UserName: username,
				})
				cmdList["findLineup"](c, &sdmht_entity.FindLineupReq{
					AccountID: g_account.ID,
				})
			case "newLineup":
				lineupName, unitsStr := "", ""
				cardLibrarys, units := []int64{}, []int64{}
				fmt.Printf("Input newLineup param [lineupName unitStr(eg:1;2;3)]> ")
				fmt.Scanln(&lineupName, &unitsStr)
				if lineupName == "" {
					lineupName = "l1"
				}
				if unitsStr == "" {
					unitsStr = "11;218;401"
				}
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
			case "newMatch", "nm":
				cmd = "newMatch"
				lineupID := uint64(0)
				fmt.Printf("Input newMatch param [lineupID]> ")
				fmt.Scanln(&lineupID)

				var chooseLineup *sdmht_entity.Lineup
				for _, lineup := range g_lineup {
					if lineup.ID == lineupID && lineup.Enabled {
						chooseLineup = lineup
						break
					}
				}
				if chooseLineup == nil {
					fmt.Printf("invalid lineup id")
					break
				}
				position := make([]int64, 16)
				for i, p := range rand.Perm(16)[:3] {
					position[p] = chooseLineup.Units[i]
				}

				cmdList[cmd](c, &sdmht_entity.NewMatchReq{
					AccountID: g_account.ID,
					LineupID:  lineupID,
					Positions: position,
				})
			case "getMatch", "gm":
				cmd = "getMatch"
				cmdList[cmd](c, &sdmht_entity.GetMatchReq{
					AccountID: g_account.ID,
				})
			case "joinMatch", "jm":
				cmd = "joinMatch"
				lineupID, matchID := uint64(0), uint64(0)
				fmt.Printf("Input newMatch param [lineupID, matchID]> ")
				fmt.Scanln(&lineupID, &matchID)
				var chooseLineup *sdmht_entity.Lineup
				for _, lineup := range g_lineup {
					if lineup.ID == lineupID && lineup.Enabled {
						chooseLineup = lineup
						break
					}
				}
				if chooseLineup == nil {
					fmt.Printf("invalid lineup id")
					break
				}
				position := make([]int64, 16)
				for i, p := range rand.Perm(16)[:3] {
					position[p] = chooseLineup.Units[i]
				}

				cmdList[cmd](c, &sdmht_entity.JoinMatchReq{
					AccountID: g_account.ID,
					MatchID:   matchID,
					Positions: position,
					LineupID:  lineupID,
				})
			case "operate", "op":
				cmd = "operate"
				if g_playerID != int(g_match.WhoseTurn) {
					fmt.Println("not your turn")
					break
				}
				event, from, to := "", int64(0), int64(0)
				fmt.Printf("Event mention: [attack|move|end]\nInput operate param [event, from, to]> ")
				fmt.Scanln(&event, &from, &to)
				cmdList[cmd](c, &sdmht_entity.SyncOperate{
					MatchID:  g_match.ID,
					Operator: int32(g_playerID),
					Event:    event,
					From:     from,
					To:       to,
				})
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
	res, ok := recvPayload.MsgContent.(*sdmht_entity.FindLineupRes)
	if !ok {
		return
	}
	g_lineup = res.Lineups
	fmt.Println("========== Find Lineups ============")
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
}

func MakeNewMatchReq(c *Client, req interface{}) {
	r := req.(*sdmht_entity.NewMatchReq)
	payload := entity.NewReqPayload(c.NewSN(), sdmht_entity.MsgTypeNewMatch, r)
	c.sendPayloadChan <- payload
	recvPayload := <-c.recvPayloadChan
	res, ok := recvPayload.MsgContent.(*sdmht_entity.NewMatchRes)
	if !ok {
		return
	}
	g_match.ID = res.MatchID
	g_playerID = 0
	fmt.Println("========== New Match ============")
	fmt.Println(g_match.ID)
}

func MakeGetMatchReq(c *Client, req interface{}) {
	r := req.(*sdmht_entity.GetMatchReq)
	payload := entity.NewReqPayload(c.NewSN(), sdmht_entity.MsgTypeGetMatch, r)
	c.sendPayloadChan <- payload
	recvPayload := <-c.recvPayloadChan
	res, ok := recvPayload.MsgContent.(*sdmht_entity.GetMatchRes)
	if !ok {
		return
	}
	g_match = &res.Match
	fmt.Println("========== Get Match ============")
	ShowScene(g_match)
}

func MakeJoinMatchReq(c *Client, req interface{}) {
	r := req.(*sdmht_entity.JoinMatchReq)
	payload := entity.NewReqPayload(c.NewSN(), sdmht_entity.MsgTypeJoinMatch, r)
	c.sendPayloadChan <- payload
	recvPayload := <-c.recvPayloadChan
	res, ok := recvPayload.MsgContent.(*sdmht_entity.JoinMatchRes)
	if !ok {
		return
	}
	g_match = &res.Match
	g_playerID = 1
	fmt.Println("========== Join Match ============")
	ShowScene(g_match)
}

func MakeOperateReq(c *Client, req interface{}) {
	r := req.(*sdmht_entity.SyncOperate)
	payload := entity.NewReqPayload(c.NewSN(), sdmht_entity.MsgTypeSyncOperator, r)
	c.sendPayloadChan <- payload
	recvPayload := <-c.recvPayloadChan
	res, ok := recvPayload.MsgContent.(*sdmht_entity.SyncOperateRes)
	if !ok {
		return
	}
	g_match = &res.Match
	fmt.Println("========== After Operate ============")
	ShowScene(g_match)
}

func MakeKeepAliveReq(c *Client, _ interface{}) {
	payload := entity.NewReqPayload(c.NewSN(), sdmht_entity.MsgTypeKeepAlive, &sdmht_entity.KeepAliveReq{})
	c.sendPayloadChan <- payload
}

func ShowScene(m *sdmht_entity.Match) {
	fmt.Println("-------------", m.WhoseTurn, "-------------")
	for who, scene := range m.Scenes {
		fmt.Printf("--- %d ---", who)
		for i, square := range scene.Squares {
			if i%4 == 0 {
				fmt.Println()
			}
			if square <= 0 {
				if scene.UnitsLocation[i] != 0 && who == g_playerID {
					fmt.Printf("■%d\t", scene.UnitsLocation[i])
				} else {
					fmt.Printf("■\t")
				}
			} else {
				if scene.UnitsLocation[i] == 0 {
					fmt.Printf("□\t")
				} else {
					fmt.Printf("□%d\t", scene.UnitsLocation[i])
				}
			}
		}
		fmt.Println()
		if who == g_playerID {
			fmt.Println()
			fmt.Println("library num:", len(scene.CardLibraries),
				"\tcountdown:", scene.DrawCardCountDown,
				"\tcost:", scene.Cost)
			fmt.Println(scene.HandCards)
			for _, unit := range scene.Units {
				fmt.Println("-->\tname:", unit.Name, "\tid:", unit.ID, "\n",
					"Health:", unit.Health, "\tDefend:", unit.Defend, "\n",
					"Attack:", unit.Attack, "\tmove:", unit.Move)
			}
		} else {
			for i, square := range scene.Squares {
				if square > 0 && scene.UnitsLocation[i] != 0 {
					unit := scene.Units[scene.UnitsLocation[i]]
					fmt.Println("-->\tname:", unit.Name, "\tid:", unit.ID, "\n",
						"Health:", unit.Health, "\tDefend:", unit.Defend, "\n",
						"Attack:", unit.Attack, "\tmove:", unit.Move)
				}
			}
		}
	}
}
