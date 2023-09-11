package entity

import (
	"encoding/json"
	"time"
)

var MsgTypes = make(map[string]interface{})

const (
	MsgTypeLogin = "login"

	MsgTypeNewLineup    = "new_lineup"
	MsgTypeFindLineup   = "find_lineup"
	MsgTypeUpdateLineup = "update_lineup"
	MsgTypeDeleteLineup = "delete_lineup"

	MsgTypeNewMatch  = "new_match"
	MsgTypeJoinMatch = "join_match"
	MsgTypeEndMatch  = "end_match"

	MsgTypeKeepAlive = "keep_alive"
)

type CommonRes struct{}

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
	MsgTypes[MsgTypeNewLineup] = NewLineupReq{}
	MsgTypes[MsgTypeFindLineup] = FindLineupReq{}
	MsgTypes[MsgTypeUpdateLineup] = UpdateLineupReq{}
	MsgTypes[MsgTypeDeleteLineup] = DeleteLineupReq{}
	MsgTypes[MsgTypeNewMatch] = NewMatchReq{}
	MsgTypes[MsgTypeJoinMatch] = JoinMatchReq{}
	MsgTypes[MsgTypeEndMatch] = EndMatchReq{}
	MsgTypes[MsgTypeKeepAlive] = KeepAliveReq{}
}
