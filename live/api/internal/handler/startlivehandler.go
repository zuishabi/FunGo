// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package handler

import (
	"fungo/common/response"
	"net/http"

	"fungo/live/api/internal/logic"
	"fungo/live/api/internal/svc"
)

func StartLiveHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logic.NewStartLiveLogic(r.Context(), svcCtx)
		err := l.StartLive()
		response.Response(r, w, nil, err)
	}
}
