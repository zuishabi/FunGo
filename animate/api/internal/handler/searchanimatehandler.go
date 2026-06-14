// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package handler

import (
	"fungo/animate/api/internal/types"
	"fungo/common/response"
	"net/http"

	"fungo/animate/api/internal/logic"
	"fungo/animate/api/internal/svc"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func searchAnimateHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.SearchAnimateReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		l := logic.NewSearchAnimateLogic(r.Context(), svcCtx)
		resp, err := l.SearchAnimate(&req)
		response.Response(r, w, resp, err)
	}
}
