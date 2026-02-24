// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"
	"fungo/user/model"

	"fungo/user/api/internal/svc"
	"fungo/user/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetSignatureLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetSignatureLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSignatureLogic {
	return &GetSignatureLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetSignatureLogic) GetSignature(req *types.GetSignatureReq) (resp *types.GetSignatureRsp, err error) {
	target := model.User{}
	l.svcCtx.Db.Select("signature").Where("id = ?", req.ID).First(&target)

	return &types.GetSignatureRsp{Content: target.Signature}, nil
}
