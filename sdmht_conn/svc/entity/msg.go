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

	// client
	//ClientQuitReasonOffline          = 0 // 客户端退出原因 掉线退出
	//ClientQuitReasonLoginTimeout     = 1 // 登录超时
	//ClientQuitReasonKick             = 2 // 被T掉
	//ClientQuitReasonHeartBeatTimeout = 3 // 心跳超时
	ClientQuitReasonOffline          = "ClientQuitReasonOffline"          // 客户端退出原因 掉线退出
	ClientQuitReasonLoginTimeout     = "ClientQuitReasonLoginTimeout"     // 登录超时
	ClientQuitReasonKick             = "ClientQuitReasonKick"             // 被T掉
	ClientQuitReasonHeartBeatTimeout = "ClientQuitReasonHeartBeatTimeout" // 心跳超时
	ClientQuitReasonUnauthorized     = "ClientQuitReasonUnauthorized"     // token验证失败
	ClientQuitReasonOfflineMsg       = "client close conn"
	ErrClientResTimeout              = "client res time out"
	ErrConnClientNotOnline           = "client not on line"
	ClientDoReqWaitRespTimeout       = 5 // 向客户端发起请求并等待响应超时时间 秒
	//GrpcWaitClientRespTimeout  = ClientDoReqWaitRespTimeout + 2 // grpc服务 要在客户端响应基础+n秒
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
	MsgTypes[MsgTypeSyncOperator] = SyncOperateReq{}
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
	MsgTypes[MsgTypeSyncMatchRsp] = SyncOperateRes{}
	MsgTypes[MsgTypeSyncOperatorRsp] = SyncOperateRes{}
	MsgTypes[MsgTypeKeepAliveRsp] = CommonRes{}
}
