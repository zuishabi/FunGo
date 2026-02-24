// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package handler

import (
	"fungo/common/response"
	"net/http"

	"fungo/animate/api/internal/logic"
	"fungo/animate/api/internal/svc"
	"fungo/animate/api/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func getAllAnimatesHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetAllAnimtesReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewGetAllAnimatesLogic(r.Context(), svcCtx)
		resp, err := l.GetAllAnimates(&req)
		response.Response(r, w, resp, err)
	}
}
