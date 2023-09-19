package entity

import (
	"encoding/json"
	"fmt"
	"reflect"

	"sdmht/lib/log"
)

func MsgToPayload(msg []byte) (ret Payload, err error) {
	err = json.Unmarshal(msg, &ret)
	if err != nil {
		log.S().Errorw("msgToPayload illegal", "msg", string(msg), "raw", string(msg))
		return
	}
	target, ok := MsgTypes[ret.MsgType]
	if !ok {
		log.S().Errorw("unknown msg type", "type", ret.MsgType, "raw", string(msg))
		return
		//return ret, errors.Errorf("%s:%s", entity.ErrMsgCodeTypes[entity.ErrCodeMsgMsgTypeUnknown], ret.MsgType)
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
	err = json.Unmarshal(data, ptr)
	if err != nil {
		log.S().Errorw("json Unmarshal", "err", err)
	}
	ret.MsgContent = ptr

	return
}

func RspMsgToPayload(msg []byte) (ret Payload, err error) {
	err = json.Unmarshal(msg, &ret)
	if err != nil {
		log.S().Errorw("msgToPayload illegal", "msg", string(msg), "raw", string(msg))
		return
	}
	msgType := fmt.Sprintf("%s_rsp", ret.MsgType)
	target, ok := MsgTypes[msgType]
	if !ok {
		log.S().Errorw("unknown msg type", "type", ret.MsgType, "raw", string(msg))
		return
		//return ret, errors.Errorf("%s:%s", entity.ErrMsgCodeTypes[entity.ErrCodeMsgMsgTypeUnknown], ret.MsgType)
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
	err = json.Unmarshal(data, ptr)
	if err != nil {
		log.S().Errorw("json Unmarshal", "err", err)
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
