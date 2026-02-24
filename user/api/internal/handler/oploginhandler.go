// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package handler

import (
	"fungo/common/response"
	"net/http"

	"fungo/user/api/internal/logic"
	"fungo/user/api/internal/svc"
	"fungo/user/api/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func opLoginHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.OpLoginReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewOpLoginLogic(r.Context(), svcCtx)
		resp, err := l.OpLogin(&req)
		response.Response(r, w, resp, err)
	}
}
