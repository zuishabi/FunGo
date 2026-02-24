// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package handler

import (
	"fungo/animate/api/internal/types"
	"fungo/common/response"
	"net/http"

	"fungo/animate/api/internal/logic"
	"fungo/animate/api/internal/svc"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func getFavoritesHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetFavoritesReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		l := logic.NewGetFavoritesLogic(r.Context(), svcCtx)
		resp, err := l.GetFavorites(&req)
		response.Response(r, w, resp, err)
	}
}
