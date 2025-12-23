// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package handler

import (
	"fungo/articles/api/internal/logic"
	"fungo/articles/api/internal/svc"
	"fungo/articles/api/internal/types"
	"fungo/common/response"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func LikeHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.LikeReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewLikeLogic(r.Context(), svcCtx)
		err := l.Like(&req)
		response.Response(r, w, nil, err)
	}
}
