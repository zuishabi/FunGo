// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package handler

import (
	"fungo/common/response"
	"net/http"

	"fungo/game/internal/logic"
	"fungo/game/internal/svc"
	"fungo/game/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func GameDetailHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GameDetailReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewGameDetailLogic(r.Context(), svcCtx)
		resp, err := l.GameDetail(&req)
		response.Response(r, w, resp, err)
	}
}
