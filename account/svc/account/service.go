package account

import (
	"context"
	"strings"
	"time"

	"sdmht/account/svc/entity"
	itfs "sdmht/account/svc/interfaces"
	"sdmht/lib"
	"sdmht/lib/log"
	"sdmht/lib/utils"

	"github.com/golang-jwt/jwt"
)

var _ itfs.Service = (*service)(nil)

var (
	TokenTTL = 24 * time.Hour
)

type service struct {
	accountRepo itfs.AccountRepo
	tokenRepo   itfs.TokenRepo
}

func NewService(accountRepo itfs.AccountRepo,
	tokenRepo itfs.TokenRepo) *service {
	return &service{
		accountRepo: accountRepo,
		tokenRepo:   tokenRepo,
	}

}

func (s *service) Register(ctx context.Context, req *entity.RegisterReq) error {
	if req.WechatID == "" {
		return lib.NewError(lib.ErrInvalidArgument, "invalid wechatid")
	}
	_, err := s.accountRepo.GetByWechatID(ctx, req.WechatID)
	if err != nil {
		if strings.Compare(err.Error(), "sql: no rows in result set") != 0 {
			log.S().Errorw("Register:get account by wechatid fail", "err", err)
			return err
		}
	} else {
		log.S().Errorw("account exist", "req", req)
		return lib.NewError(lib.ErrInvalidArgument, "account exist")
	}

	account := &entity.Account{
		WeChatID: req.WechatID,
		UserName: req.UserName,
	}
	err = s.accountRepo.Add(ctx, account)
	if err != nil {
		log.S().Errorw("Register:add account fail", "err", err)
		return err
	}
	return nil
}

func (s *service) Login(ctx context.Context, req *entity.LoginReq) (res *entity.LoginRes, err error) {
	log.S().Infow("Login:req", "wechatID", req.WechatID)

	account, err := s.accountRepo.GetByWechatID(ctx, req.WechatID)
	if err != nil {
		log.S().Errorw("Login:get account by wechatid fail", "err", err)
		return nil, err
	}

	claim := entity.Claims{
		StandardClaims: jwt.StandardClaims{
			IssuedAt: time.Now().Unix(),
		},
		AccountID: account.ID,
		WechatID:  req.WechatID,
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claim).SignedString([]byte(utils.GetJWTSecret()))
	if err != nil {
		log.S().Errorw("Login:jwt fail", "err", err)
		return nil, err
	}
	if err := s.tokenRepo.Add(ctx, token, account.ID, TokenTTL); err != nil {
		log.S().Errorw("Login:token add fail", "err", err)
		return nil, err
	}

	account.LastLoginAt = time.Now()
	err = s.accountRepo.Update(ctx, account)
	if err != nil {
		log.S().Errorw("Login:update account fail", "err", err)
		return nil, err
	}

	res = &entity.LoginRes{Token: token, Account: account}
	return res, nil
}

func (s *service) Logout(ctx context.Context, token string) error {
	if err := s.tokenRepo.Delete(ctx, token); err != nil {
		log.S().Warnf("delete token failed")
	}
	return nil
}

func (s *service) Authenticate(ctx context.Context, token string) (*entity.Account, error) {
	claim, err := s.verifyToken(ctx, token)
	if err != nil {
		log.S().Errorf("verify token err: %v", err)
		return nil, err
	}

	account, err := s.accountRepo.Get(ctx, claim.AccountID)
	if err != nil {
		log.S().Errorw("Authenticate:FindByID err", "err", err)
		return nil, err
	}

	return account, nil
}

func (s *service) verifyToken(ctx context.Context, tokenString string) (*entity.Claims, error) {
	var claims entity.Claims
	token, err := jwt.ParseWithClaims(tokenString, &claims, func(t *jwt.Token) (interface{}, error) { return []byte(utils.GetJWTSecret()), nil })
	if err != nil || !token.Valid {
		return nil, lib.NewError(lib.ErrUnauthorized, "token is invalid")
	}

	accountID, err := s.tokenRepo.Get(ctx, tokenString, TokenTTL)
	if err != nil {
		return nil, lib.NewError(lib.ErrUnauthorized, "token not existed")
	}
	if accountID != claims.AccountID {
		return nil, lib.NewError(lib.ErrUnauthorized, "accountid mismatch")
	}

	return &claims, nil
}

func (s *service) GetAccount(ctx context.Context, id uint64) (*entity.Account, error) {
	account, err := s.accountRepo.Get(ctx, id)
	if err != nil {
		log.S().Errorw("get account fail", "err", err)
		return nil, err
	}
	return account, nil
}
