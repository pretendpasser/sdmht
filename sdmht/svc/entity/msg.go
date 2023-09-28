package entity

import (
	"encoding/json"
	"time"
)

var MsgTypes = make(map[string]interface{})

const (
	//res
	MsgTypeLogin        = "login"
	MsgTypeNewLineup    = "new_lineup"
	MsgTypeFindLineup   = "find_lineup"
	MsgTypeUpdateLineup = "update_lineup"
	MsgTypeDeleteLineup = "delete_lineup"
	MsgTypeNewMatch     = "new_match"
	MsgTypeGetMatch     = "get_match"
	MsgTypeJoinMatch    = "join_match"
	MsgTypeSyncMatch    = "sync_match"
	MsgTypeSyncOperator = "sync_operator"
	MsgTypeKeepAlive    = "keep_alive"

	// rsp(for decode)
	MsgTypeLoginRsp        = "login_rsp"
	MsgTypeNewLineupRsp    = "new_lineup_rsp"
	MsgTypeFindLineupRsp   = "find_lineup_rsp"
	MsgTypeUpdateLineupRsp = "update_lineup_rsp"
	MsgTypeDeleteLineupRsp = "delete_lineup_rsp"
	MsgTypeNewMatchRsp     = "new_match_rsp"
	MsgTypeGetMatchRsp     = "get_match_rsp"
	MsgTypeJoinMatchRsp    = "join_match_rsp"
	MsgTypeSyncMatchRsp    = "sync_match_rsp"
	MsgTypeSyncOperatorRsp = "sync_operator_rsp"
	MsgTypeKeepAliveRsp    = "keep_alive_rsp"
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
	// req
	MsgTypes[MsgTypeLogin] = LoginReq{}
	MsgTypes[MsgTypeNewLineup] = NewLineupReq{}
	MsgTypes[MsgTypeFindLineup] = FindLineupReq{}
	MsgTypes[MsgTypeUpdateLineup] = UpdateLineupReq{}
	MsgTypes[MsgTypeDeleteLineup] = DeleteLineupReq{}
	MsgTypes[MsgTypeNewMatch] = NewMatchReq{}
	MsgTypes[MsgTypeGetMatch] = GetMatchReq{}
	MsgTypes[MsgTypeJoinMatch] = JoinMatchReq{}
	MsgTypes[MsgTypeSyncMatch] = Match{}
	MsgTypes[MsgTypeSyncOperator] = SyncOperate{}
	MsgTypes[MsgTypeKeepAlive] = KeepAliveReq{}

	// rsp(for decode)
	MsgTypes[MsgTypeLoginRsp] = LoginRes{}
	MsgTypes[MsgTypeNewLineupRsp] = CommonRes{}
	MsgTypes[MsgTypeFindLineupRsp] = FindLineupRes{}
	MsgTypes[MsgTypeUpdateLineupRsp] = CommonRes{}
	MsgTypes[MsgTypeDeleteLineupRsp] = CommonRes{}
	MsgTypes[MsgTypeNewMatchRsp] = NewMatchRes{}
	MsgTypes[MsgTypeGetMatchRsp] = GetMatchRes{}
	MsgTypes[MsgTypeJoinMatchRsp] = JoinMatchRes{}
	MsgTypes[MsgTypeSyncOperatorRsp] = SyncOperateRes{}
	MsgTypes[MsgTypeKeepAliveRsp] = CommonRes{}
}
