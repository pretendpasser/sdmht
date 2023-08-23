package itfs

import (
	"context"

	"sdmht/sdmht/svc/entity"
	sdmht_itfs "sdmht/sdmht/svc/interfaces"
)

type ConnService interface {
	DispatchEventToClient(ctx context.Context, target uint64, event entity.ClientEvent) (entity.DispatchEventToClientReply, error)
	KickClient(ctx context.Context, id uint64) error
	// SwitchToSpeecher(ctx context.Context, req entity.SwitchSpeecherReq) (*entity.SwitchSpeecherRes, error)
}

type ConnManager sdmht_itfs.SignalingService
