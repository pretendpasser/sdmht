package entity

import (
	"time"

	sdmht_entity "sdmht/sdmht/svc/entity"
)

const (
	PayloadTypeReq = "req"
	PayloadTypeRsp = "rsp"

	//ClientQuitReasonOffline          = 0 // 客户端退出原因 掉线退出
	//ClientQuitReasonLoginTimeout     = 1 // 登录超时
	//ClientQuitReasonKick             = 2 // 被T掉
	//ClientQuitReasonHeartBeatTimeout = 3 // 心跳超时

	ClientQuitReasonOffline          = "ClientQuitReasonOffline"          // 客户端退出原因 掉线退出
	ClientQuitReasonLoginTimeout     = "ClientQuitReasonLoginTimeout"     // 登录超时
	ClientQuitReasonKick             = "ClientQuitReasonKick"             // 被T掉
	ClientQuitReasonHeartBeatTimeout = "ClientQuitReasonHeartBeatTimeout" // 心跳超时
	ClientQuitReasonUnauthorized     = "ClientQuitReasonUnauthorized"     // token验证失败

	ClientQuitReasonOfflineMsg = "client close conn"

	ErrCodeMsgSuccess          = 200
	ErrCodeMsgBadRequest       = 400
	ErrCodeMsgUnauthorized     = 401
	ErrCodeMsgForbidden        = 403
	ErrCodeMsgMethodNotAllowed = 405
	ErrCodeMsgGone             = 410
	ErrCodeMsgInternal         = 500
	ErrCodeMsgNotImplemented   = 501
	ErrCodeMsgNotFount         = 701

	ErrClientRespTimeout = "client resp time out"

	ErrPocConnClientNotOnline = "client not on line"

	ClientDoReqWaitRespTimeout = 5 // 向客户端发起请求并等待响应超时时间 秒
	//GrpcWaitClientRespTimeout  = ClientDoReqWaitRespTimeout + 2 // grpc服务 要在客户端响应基础+n秒

	MaxSN = 65535
)

// ErrCodeMsgs {errCode: errMsg}
var ErrCodeMsgs = map[int]string{
	ErrCodeMsgSuccess:          "ErrCodeMsgSuccess",
	ErrCodeMsgBadRequest:       "ErrCodeMsgBadRequest",
	ErrCodeMsgUnauthorized:     "ErrCodeMsgUnauthorized",
	ErrCodeMsgForbidden:        "ErrCodeMsgForbidden",
	ErrCodeMsgMethodNotAllowed: "ErrCodeMsgMethodNotAllowed",
	ErrCodeMsgGone:             "ErrCodeMsgGone",
	ErrCodeMsgInternal:         "ErrCodeMsgInternal",
	ErrCodeMsgNotImplemented:   "ErrCodeMsgNotImplemented",
	ErrCodeMsgNotFount:         "ErrCodeMsgNotFount",
}

var MsgTypes = sdmht_entity.MsgTypes

type Payload struct {
	Version     string `json:"ver"`
	SN          int    `json:"cseq"`
	PayloadType string `json:"type"` // "req" or "rsp"
	MsgType     string `json:"msg"`  // "login" or other
	Token       string `json:"token"`
	Result
	MsgContent interface{} `json:"body"` // req时 反射解出来 是个实际对应结构的对象指针
}

type Result struct {
	Code   int    `json:"code"`
	Reason string `json:"code_desc"`
}

type LogoutEvent struct {
	PocTermName string
	Time        time.Time
}

func NewResult(code int, reason string) Result {
	return Result{Code: code, Reason: reason}
}

func NewRespPayload(req Payload, code int, reason string, resp interface{}) Payload {
	return Payload{
		SN:          req.SN,
		PayloadType: PayloadTypeRsp,
		MsgType:     req.MsgType,
		Result:      NewResult(code, reason),
		MsgContent:  resp,
	}
}

func NewReqPayload(sn int, msgType string, token string, content interface{}) Payload {
	return Payload{
		SN:          sn,
		PayloadType: PayloadTypeReq,
		MsgType:     msgType,
		Token:       token,
		MsgContent:  content,
	}
}
