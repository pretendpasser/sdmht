package sdmht

import (
	"context"

	account "sdmht/account/svc/interfaces"
	"sdmht/lib/kitx"
	"sdmht/lib/log"
	"sdmht/lib/seq"
	"sdmht/sdmht/svc/entity"
	itfs "sdmht/sdmht/svc/interfaces"
)

const (
	ConnServeAddrKey = "sdmgt_conn_addr"
)

// func getConnServeAddr(ctx context.Context) string {
// 	if v := ctx.Value(kitx.MetadataKey(ConnServeAddrKey)); v != nil {
// 		if value, ok := v.([]string); ok {
// 			return value[0]
// 		}
// 	}
// 	return ""
// }

var _ itfs.SignalingService = (*signalingSvc)(nil)

type signalingSvc struct {
	idGenerator seq.IDGenerator
	// repo       itfs.EventParticipantRepo
	eventSvc   itfs.Service
	accountSvc account.Service
	connManger *ConnManager
}

func NewSignalingService(idGenerator seq.IDGenerator,
	eventSvc itfs.Service,
	accountSvc account.Service,
	connManger *ConnManager) itfs.SignalingService {
	return &signalingSvc{
		idGenerator: idGenerator,
		eventSvc:    eventSvc,
		accountSvc:  accountSvc,
		connManger:  connManger,
	}
}

func (s *signalingSvc) NewMatch(ctx context.Context, req *entity.NewMatchReq) (*entity.NewMatchRsp, error) {
	res := &entity.NewMatchRsp{}
	res.Player = &entity.Player{
		ID:    req.Operator,
		Scene: entity.NewScene(),
	}
	return res, nil
}

// func (s *signalingSvc) JoinMatch(ctx context.Context, req entity.JoinMatchReq) (*entity.JoinMatchRes, error) {
// 	res := &entity.JoinMatchRes{}
// 	newMatchID, err := s.idGenerator.NextID()
// 	if err != nil {
// 		return nil, err
// 	}
// 	res.MatchID = newMatchID
// 	return res, nil
// }

func (s *signalingSvc) KeepAlive(ctx context.Context, req *entity.KeepAliveReq) error {
	slog := log.L().With(kitx.TraceIDField(ctx)).Sugar()
	slog.Infow("KeepAlive:Req", "params", req)
	// return s.repo.Update(ctx, 0, req.Operator, "status", entity.UserOnline)
	return nil
}

func (s *signalingSvc) Offline(ctx context.Context, req *entity.LogoutReq) error {
	slog := log.L().With(kitx.TraceIDField(ctx)).Sugar()
	slog.Infow("Offline:Req", "params", req)

	// info, err := s.repo.Get(ctx, 0, req.Operator)
	// if err != nil {
	// 	slog.Errorw("Offline:GetJoinedEvent", "params", req, "err", err)
	// 	return err
	// }

	// return s.LeaveEvent(ctx, entity.LeaveEventReq{Operator: req.Operator, EventID: info.EventID})
	return nil

}
