// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package handler

import (
	"fungo/common/response"
	"net/http"

	"fungo/ai/api/internal/logic"
	"fungo/ai/api/internal/svc"
	"fungo/ai/api/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func searchAnimateByAIHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.SearchAnimateByAIReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewSearchAnimateByAILogic(r.Context(), svcCtx)
		resp, err := l.SearchAnimateByAI(&req)
		response.Response(r, w, resp, err)
	}
}
