package sdmht

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	"sdmht/lib"
	"sdmht/lib/log"
	"sdmht/sdmht_conn/svc/entity"
	itfs "sdmht/sdmht_conn/svc/interfaces"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

var slog = entity.SLog

type Client struct {
	conn *websocket.Conn

	pendingMsgs sync.Map // {SN: pendingMsg}

	addr      string
	sn        int
	accountID uint64
	wechatID  string

	serverStartTime   time.Time
	loginTime         time.Time // 用于客户端重登时，pocmng可能有过去的消息，刚发过来，此时可以判断时间做废弃
	lastHeartBeatTime time.Time
	heartBeatInterval int
	rwLock            sync.RWMutex

	sendPayloadChan chan entity.Payload
	closeChan       chan struct{}
	closeWaitChan   chan struct{}
	quitReason      chan string // 长度必须设2 因为主动t掉 也会触发后续 掉线

	notifyClientEvent chan ClientEvent

	connMng itfs.SignalingService
}

func NewClient(conn *websocket.Conn, addr string, notifyClientEvent chan ClientEvent,
	connMng itfs.SignalingService, serverStartTime time.Time, heartBeatInterval int) *Client {
	return &Client{
		conn:              conn,
		addr:              addr,
		sendPayloadChan:   make(chan entity.Payload, 50),
		closeChan:         make(chan struct{}, 1),
		closeWaitChan:     make(chan struct{}, 1),
		quitReason:        make(chan string, 1),
		notifyClientEvent: notifyClientEvent,
		connMng:           connMng,
		serverStartTime:   serverStartTime,
		lastHeartBeatTime: time.Now(),
		heartBeatInterval: heartBeatInterval,
	}
}

func (c *Client) Run() {
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

	wg.Add(1)
	go func() {
		defer wg.Done()
		c.HandleCheckValid()
	}()

	wg.Wait()

	c.closeWaitChan <- struct{}{}
	quitReason := <-c.quitReason
	log.S().Infow("quit", "client", c, "close reason", quitReason)
	if quitReason == entity.ClientQuitReasonOffline || quitReason == entity.ClientQuitReasonHeartBeatTimeout {
		c.notifyClientEvent <- NewClientEvent(ClientEventTypeRemove, c.AccountID(), c)
		err := c.connMng.Offline(context.TODO(), &entity.LogoutReq{
			Operator: c.AccountID(),
			Reason:   quitReason,
		})
		if err != nil {
			log.S().Errorw("connMng.Offline", "err", err)
		}
	}
}

func (c *Client) HandleReadMsg() {
	ctx := entity.InjectInternalTrace(context.Background(), "")
	for {
		_, msg, err := c.conn.ReadMessage()
		if err != nil {
			slog(ctx).Infow("ReadMessage", "err", err)
			c.close(entity.ClientQuitReasonOffline)
			return
		}
		slog(ctx).Infow("recv msg", "client", c, "payload", json.RawMessage(msg))
		payload, err := entity.MsgToPayload(msg) // parse message
		if err != nil {
			slog(ctx).Warnw("MsgToPayload", "err", err, "msg", string(msg), "client", c)
			payload.PayloadType = entity.PayloadTypeRsp
			payload.Result = entity.NewResult(lib.ErrInvalidArgument, lib.ErrorStrings[lib.ErrInvalidArgument])
			payload.MsgContent = struct{}{}
			c.sendPayloadChan <- payload
		} else {
			if payload.PayloadType == entity.PayloadTypeRsp {
				c.handleRespMsg(payload)
			} else if payload.PayloadType == entity.PayloadTypeReq {
				respPayload := c.handleReqMsg(ctx, payload)
				c.sendPayloadChan <- respPayload
			}
		}
	}
}

func (c *Client) String() string {
	return fmt.Sprintf("client[remote addr:%s, AccountID:%d, loginTime:%v]",
		c.Addr(), c.AccountID(), c.LoginTime())
}

func (c *Client) HandleSendMsg() {
	for {
		select {
		case <-c.closeChan:
			//log.S().Info("HandleSendMsg", "close")
			return
		case payload := <-c.sendPayloadChan:
			msg := entity.PayloadToMsg(payload)
			log.S().Infow("send msg", "msg", json.RawMessage(msg), "payload", payload, "client", c)
			err := c.conn.WriteMessage(1, msg)
			if err != nil {
				log.S().Errorw("write message err", "err", err, "client addr", c.Addr(),
					"msg", json.RawMessage(msg))
			}
		}
	}
}

func (c *Client) HandleCheckValid() {
	checkHeartBeatTicker := time.NewTicker(5 * time.Second)
	defer checkHeartBeatTicker.Stop()
	for {
		select {
		case <-c.closeChan:
			return
		case <-checkHeartBeatTicker.C:
			if c.HeartBeatInterval() == 0 {
				continue
			}
			if c.LastHeartBeatTime().Add(3 * time.Duration(c.HeartBeatInterval()) * time.Second).Before(time.Now()) {
				c.close(entity.ClientQuitReasonHeartBeatTimeout)
				return
			}
		}
	}
}

func (c *Client) handleRespMsg(payload entity.Payload) {
	val, loaded := c.pendingMsgs.LoadAndDelete(payload.SN)
	if !loaded { // 已经过期被删除了
		data, _ := json.Marshal(payload)
		log.S().Infow("pending msg not loaded", "payload", json.RawMessage(data))
		return
	}
	pm := val.(pendingMsg)
	pm.rspChan <- payload
}

func (c *Client) handleReqMsg(ctx context.Context, payload entity.Payload) (ret entity.Payload) {
	c.SetLastHeartBeatTime()
	var err error

	switch payload.MsgType {
	case entity.MsgTypeLogin:
		req := payload.MsgContent.(*entity.LoginReq)
		ctx = context.WithValue(ctx, ConnServeAddrKey, c.addr)
		res, err2 := c.connMng.Login(ctx, req)
		if err2 != nil {
			err = err2
			break
		}
		log.S().Infow("client login res", "res", res, "AccountID", res.AccountID)
		now := time.Now()
		c.SetInfo(res.AccountID, req.WeChatID)
		c.SetLoginTime(now)
		c.notifyClientEvent <- NewClientEvent(ClientEventTypeAdd, c.AccountID(), c)
		ret = entity.NewRespPayload(payload, lib.ErrSuccess, "", res)
	case entity.MsgTypeNewLineup:
		content := payload.MsgContent.(*entity.NewLineupReq)
		content.AccountID = c.AccountID()
		err1 := c.connMng.NewLineup(ctx, content)
		if err1 != nil {
			err = err1
			break
		}
		ret = entity.NewRespPayload(payload, lib.ErrSuccess, "", entity.CommonRes{})
	case entity.MsgTypeFindLineup:
		content := payload.MsgContent.(*entity.FindLineupReq)
		content.AccountID = c.AccountID()
		res, err1 := c.connMng.FindLineup(ctx, content)
		if err1 != nil {
			err = err1
			break
		}
		slog(ctx).Infow("find lineup res", "res", res)
		ret = entity.NewRespPayload(payload, lib.ErrSuccess, "", res)
	case entity.MsgTypeUpdateLineup:
		content := payload.MsgContent.(*entity.UpdateLineupReq)
		content.AccountID = c.AccountID()
		err1 := c.connMng.UpdateLineup(ctx, content)
		if err1 != nil {
			err = err1
			break
		}
		ret = entity.NewRespPayload(payload, lib.ErrSuccess, "", entity.CommonRes{})
	case entity.MsgTypeDeleteLineup:
		content := payload.MsgContent.(*entity.DeleteLineupReq)
		content.AccountID = c.AccountID()
		err1 := c.connMng.DeleteLineup(ctx, content)
		if err1 != nil {
			err = err1
			break
		}
		ret = entity.NewRespPayload(payload, lib.ErrSuccess, "", entity.CommonRes{})
	case entity.MsgTypeNewMatch:
		content := payload.MsgContent.(*entity.NewMatchReq)
		rsp, err1 := c.connMng.NewMatch(ctx, content)
		if err1 != nil {
			err = err1
			break
		}
		ret = entity.NewRespPayload(payload, lib.ErrSuccess, "", rsp)
	case entity.MsgTypeGetMatch:
		content := payload.MsgContent.(*entity.GetMatchReq)
		rsp, err1 := c.connMng.GetMatch(ctx, content)
		if err1 != nil {
			err = err1
			break
		}
		ret = entity.NewRespPayload(payload, lib.ErrSuccess, "", rsp)
	case entity.MsgTypeJoinMatch:
		content := payload.MsgContent.(*entity.JoinMatchReq)
		rsp, err1 := c.connMng.JoinMatch(ctx, content)
		if err1 != nil {
			err = err1
			break
		}
		ret = entity.NewRespPayload(payload, lib.ErrSuccess, "", rsp)
	case entity.MsgTypeSyncOperator:
		content := payload.MsgContent.(*entity.SyncOperateReq)
		rsp, err1 := c.connMng.SyncOperate(ctx, content)
		if err1 != nil {
			err = err1
			break
		}
		ret = entity.NewRespPayload(payload, lib.ErrSuccess, "", rsp)
	case entity.MsgTypeKeepAlive:
		req := payload.MsgContent.(*entity.KeepAliveReq)
		req.Operator = c.AccountID()
		err = c.connMng.KeepAlive(ctx, req)
		if err != nil {
			break
		}
		ret = entity.NewRespPayload(payload, lib.ErrSuccess, "", entity.CommonRes{})
	}

	if err != nil {
		log.S().Errorw("client handleReqMsg", "err", err, "payload content", payload.MsgContent)
		var errCode int
		var reason string
		err1, ok := err.(lib.Error)
		if ok {
			errCode = err1.HttpStatusCode()
			reason = err1.Message
		} else {
			errCode = lib.ErrInternal
			reason = err.Error()
		}
		ret = entity.NewRespPayload(payload, errCode, reason, nil)
	}

	return
}

// DoRequest 向客户端发送请求 并等待回包
func (c *Client) DoRequest(req entity.Payload) (resp entity.Payload, err error) {
	pendingMsg := newPendingMsg(req)
	c.pendingMsgs.Store(pendingMsg.req.SN, pendingMsg)
	c.sendPayloadChan <- req
	timeout := time.NewTimer(entity.ClientDoReqWaitRespTimeout * time.Second)
	select {
	case <-timeout.C:
		c.pendingMsgs.Delete(pendingMsg.req.SN)
		err = errors.New(entity.ErrClientResTimeout)
	case resp = <-pendingMsg.rspChan:
	}
	log.S().Infow("client DoRequest", "req", req, "resp", resp, "err", err)
	return
}

func (c *Client) CloseWait(quitReason string) {
	log.S().Infow("closewait", zap.String("closeReason", quitReason), zap.Uint64("clientid", c.AccountID()))
	c.close(quitReason)
	<-c.closeWaitChan
}

// close 关闭
// 只有自动下线直接关closeChan，不然先关conn，通过conn触发closeChan
// 因为read conn没有chan的实现
func (c *Client) close(closeReason string) {
	log.S().Infow("close", zap.String("closeReason", closeReason), zap.Uint64("clientid", c.AccountID()))
	select {
	case c.quitReason <- closeReason:
	default:
		log.S().Info("quit reason full when close", zap.Uint64("clientID", c.AccountID()))
	}
	_ = c.conn.Close()
	if closeReason == entity.ClientQuitReasonOffline || closeReason == entity.ClientQuitReasonUnauthorized {
		close(c.closeChan)
	}
}

func (c *Client) SetLoginTime(t time.Time) {
	c.rwLock.Lock()
	defer c.rwLock.Unlock()
	c.loginTime = t
}

func (c *Client) LoggedIn() bool {
	c.rwLock.RLock()
	defer c.rwLock.RUnlock()
	return c.loginTime != time.Time{}
}

func (c *Client) SetInfo(accountID uint64, wechatID string) {
	c.rwLock.Lock()
	defer c.rwLock.Unlock()
	c.accountID = accountID
	c.wechatID = wechatID
}

func (c *Client) AccountID() uint64 {
	c.rwLock.RLock()
	defer c.rwLock.RUnlock()
	return c.accountID
}

func (c *Client) LoginTime() time.Time {
	c.rwLock.RLock()
	defer c.rwLock.RUnlock()
	return c.loginTime
}

func (c *Client) HeartBeatInterval() int {
	c.rwLock.RLock()
	defer c.rwLock.RUnlock()
	return c.heartBeatInterval
}

func (c *Client) SetHeartBeatInterval(interval int) {
	c.rwLock.Lock()
	defer c.rwLock.Unlock()
	c.heartBeatInterval = interval
}

func (c *Client) LastHeartBeatTime() time.Time {
	c.rwLock.RLock()
	defer c.rwLock.RUnlock()
	return c.lastHeartBeatTime
}

func (c *Client) SetLastHeartBeatTime() {
	c.rwLock.Lock()
	defer c.rwLock.Unlock()
	c.lastHeartBeatTime = time.Now()
}

func (c *Client) Addr() string {
	c.rwLock.RLock()
	defer c.rwLock.RUnlock()
	return c.addr
}

func (c *Client) NewSN() int {
	c.rwLock.Lock()
	defer c.rwLock.Unlock()
	if c.sn == 65535 {
		c.sn = 0
	} else {
		c.sn++
	}
	return c.sn
}

func (c *Client) SN() int {
	c.rwLock.RLock()
	defer c.rwLock.RUnlock()
	return c.sn
}

type pendingMsg struct {
	req     entity.Payload
	rspChan chan entity.Payload
}

func newPendingMsg(req entity.Payload) pendingMsg {
	return pendingMsg{
		req:     req,
		rspChan: make(chan entity.Payload, 1),
	}
}

//func (c *Client) HandleRequest(payload entity.Payload) (resp entity.Payload, error error) {
//	// TODO
//	return entity.Payload{}, nil
//}
