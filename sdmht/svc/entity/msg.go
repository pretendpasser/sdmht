package entity

import (
	"encoding/json"
	"time"
)

var MsgTypes = make(map[string]interface{})

const (
	MsgTypeLoginRequest     = "login"
	MsgTypeNewMatchRequest  = "new_match_request"
	MsgTypeKeepAliveRequest = "keep_alive_request"
)

type CommonResp struct{}

type ClientEvent struct {
	UserID  uint64
	Type    string
	AtTime  time.Time
	Content json.RawMessage
}

type DispatchEventToClientReply struct {
	UserID    uint64
	OK        bool
	ClientErr string
}

func init() {
	MsgTypes[MsgTypeLoginRequest] = LoginReq{}
	MsgTypes[MsgTypeNewMatchRequest] = NewMatchReq{}
	MsgTypes[MsgTypeKeepAliveRequest] = KeepAliveReq{}
}
