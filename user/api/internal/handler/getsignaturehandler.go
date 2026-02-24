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

func getSignatureHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetSignatureReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewGetSignatureLogic(r.Context(), svcCtx)
		resp, err := l.GetSignature(&req)
		response.Response(r, w, resp, err)
	}
}
