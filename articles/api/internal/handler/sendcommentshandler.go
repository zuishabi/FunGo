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

func SendCommentsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.SendCommentsReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewSendCommentsLogic(r.Context(), svcCtx)
		err := l.SendComments(&req)
		response.Response(r, w, nil, err)
	}
}
