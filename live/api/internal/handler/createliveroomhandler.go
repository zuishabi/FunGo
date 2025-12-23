// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package handler

import (
	"fungo/common/response"
	"net/http"

	"fungo/live/api/internal/logic"
	"fungo/live/api/internal/svc"
)

func CreateLiveRoomHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logic.NewCreateLiveRoomLogic(r.Context(), svcCtx)
		err := l.CreateLiveRoom()
		response.Response(r, w, nil, err)
	}
}
