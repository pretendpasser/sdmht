package entity

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"sdmht/lib/log"
)

const (
	PayloadTypeReq = "req"
	PayloadTypeRsp = "rsp"
)

type Payload struct {
	Version     string `json:"ver"`
	SN          int    `json:"cseq"`
	PayloadType string `json:"type"` // "req" or "rsp"
	MsgType     string `json:"msg"`  // "login" or other
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

func NewReqPayload(sn int, msgType string, content interface{}) Payload {
	return Payload{
		SN:          sn,
		PayloadType: PayloadTypeReq,
		MsgType:     msgType,
		MsgContent:  content,
	}
}

func MsgToPayload(msg []byte) (ret Payload, err error) {
	decoder := json.NewDecoder(bytes.NewReader(msg))
	decoder.UseNumber()
	err = decoder.Decode(&ret)
	if err != nil {
		log.S().Errorw("msgToPayload illegal", "msg", string(msg), "raw", string(msg))
		return
	}

	target, ok := MsgTypes[ret.MsgType]
	if !ok {
		log.S().Errorw("unknown msg type", "type", ret.MsgType, "raw", string(msg))
		return
	}
	t := reflect.TypeOf(target)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	instance := reflect.New(t)
	ptr := instance.Interface()

	data, err := json.Marshal(ret.MsgContent)
	if err != nil {
		log.S().Errorw("json marshal", "err", err)
		return
	}

	decoder = json.NewDecoder(bytes.NewReader(data))
	decoder.UseNumber()
	err = decoder.Decode(&ptr)
	if err != nil {
		log.S().Errorw("json Unmarshal", "err", err)
		return
	}

	ret.MsgContent = ptr

	return
}

func RspMsgToPayload(msg []byte) (ret Payload, err error) {
	decoder := json.NewDecoder(bytes.NewReader(msg))
	decoder.UseNumber()
	err = decoder.Decode(&ret)
	if err != nil {
		log.S().Errorw("msgToPayload illegal", "msg", string(msg), "raw", string(msg))
		return
	}

	msgType := fmt.Sprintf("%s_%s", ret.MsgType, PayloadTypeRsp)
	target, ok := MsgTypes[msgType]
	if !ok {
		log.S().Errorw("unknown msg type", "type", ret.MsgType, "raw", string(msg))
		return
	}
	t := reflect.TypeOf(target)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	instance := reflect.New(t)
	ptr := instance.Interface()

	data, err := json.Marshal(ret.MsgContent)
	if err != nil {
		log.S().Errorw("json marshal", "err", err)
		return
	}

	decoder = json.NewDecoder(bytes.NewReader(data))
	decoder.UseNumber()
	err = decoder.Decode(&ptr)
	if err != nil {
		log.S().Errorw("json Unmarshal", "err", err)
		return
	}

	ret.MsgContent = ptr

	return
}

func PayloadToMsg(payload Payload) []byte {
	data, err := json.Marshal(payload)
	if err != nil {
		log.S().Infow("PayloadToMsg", "err", err)
	}
	return data
}
