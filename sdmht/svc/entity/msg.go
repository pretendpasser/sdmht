package entity

import (
	"encoding/json"
	"time"
)

var MsgTypes = make(map[string]interface{})

const (
	MsgTypeLogin     = "login"
	MsgTypeNewMatch  = "new_match_request"
	MsgTypeKeepAlive = "keep_alive_request"
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
	MsgTypes[MsgTypeLogin] = LoginReq{}
	MsgTypes[MsgTypeNewMatch] = NewMatchReq{}
	MsgTypes[MsgTypeKeepAlive] = KeepAliveReq{}
}
