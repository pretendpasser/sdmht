package sdmht

import (
	"context"
	"errors"

	account_entity "sdmht/account/svc/entity"
	account "sdmht/account/svc/interfaces"
	"sdmht/lib"
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

func (s *signalingSvc) Login(ctx context.Context, req *entity.LoginReq) (*entity.LoginRes, error) {
	err := s.accountSvc.Register(ctx, &account_entity.RegisterReq{
		WechatID: req.WeChatID,
		UserName: req.UserName,
	})
	if err != nil {
		if !errors.Is(err, lib.NewError(lib.ErrInvalidArgument, "account exist")) {
			log.S().Errorw("Login: account register fail", "err", err)
			return nil, err
		}
	}

	res, err := s.accountSvc.Login(ctx, &account_entity.LoginReq{
		WechatID: req.WeChatID,
	})
	if err != nil {
		log.S().Errorw("Login: account login fail", "err", err)
		return nil, err
	}

	if err := s.connManger.User2ConnRepo().Add(ctx, res.Account.ID, res.Account.WeChatID); err != nil {
		log.S().Errorw("Login: add user client fail", "err", err)
		return nil, err
	}

	return &entity.LoginRes{
		AccountID: res.Account.ID,
	}, nil
}

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

func (s *signalingSvc) NewMatch(ctx context.Context, req *entity.NewMatchReq) (*entity.NewMatchRes, error) {
	res := &entity.NewMatchRes{}
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
