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
	MsgTypeGetMatch  = "get_match"
	MsgTypeJoinMatch = "join_match"
	MsgTypeSyncMatch = "sync_match"

	MsgTypeKeepAlive = "keep_alive"
)

type CommonRes struct{}

type ClientEvent struct {
	AccountID uint64
	Type      string
	AtTime    time.Time
	Content   json.RawMessage
}

type DispatchEventToClientReply struct {
	AccountID uint64
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
	MsgTypes[MsgTypeGetMatch] = GetMatchReq{}
	MsgTypes[MsgTypeJoinMatch] = JoinMatchReq{}
	MsgTypes[MsgTypeKeepAlive] = KeepAliveReq{}
}
