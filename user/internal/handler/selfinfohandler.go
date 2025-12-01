// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package handler

import (
	"fungo/common/response"
	"net/http"

	"fungo/user/internal/logic"
	"fungo/user/internal/svc"
)

func selfInfoHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logic.NewSelfInfoLogic(r.Context(), svcCtx)
		resp, err := l.SelfInfo()
		response.Response(r, w, resp, err)
	}
}
