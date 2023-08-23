package seq

import (
	"encoding/hex"
	"time"

	"github.com/google/uuid"
	"github.com/sony/sonyflake"
)

type IDGenerator interface {
	NextID() (uint64, error)
}

type options struct {
	machineIDFunc      func() (uint16, error)
	machineIDValidator func(uint16) bool
}

type SonyflakeOption func(o *options)

func WithMachineIDFunc(fun func() (uint16, error)) SonyflakeOption {
	return func(o *options) {
		o.machineIDFunc = fun
	}
}

func WithMachineIDValidator(validateFunc func(uint16) bool) SonyflakeOption {
	return func(o *options) {
		o.machineIDValidator = validateFunc
	}
}

var (
	// the staring date of c3_platform project
	startTime      = time.Date(2021, 7, 21, 0, 0, 0, 0, time.UTC)
	globalNode     *sonyflake.Sonyflake
	defaultOptions = options{}
)

func init() {
	globalNode = sonyflake.NewSonyflake(sonyflake.Settings{StartTime: startTime})
}

func NewNode(opts ...SonyflakeOption) *sonyflake.Sonyflake {
	o := defaultOptions
	for _, opt := range opts {
		opt(&o)
	}
	var st sonyflake.Settings
	st.StartTime = startTime
	st.MachineID = o.machineIDFunc
	st.CheckMachineID = o.machineIDValidator

	return sonyflake.NewSonyflake(st)
}

func New(opts ...SonyflakeOption) *sonyflake.Sonyflake {
	return NewNode(opts...)
}

func NextID() (uint64, error) {
	return globalNode.NextID()
}

func NextUUID() (string, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(id[:]), nil
}
