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

func UserArticleListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UserArticleListReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewUserArticleListLogic(r.Context(), svcCtx)
		resp, err := l.UserArticleList(&req)
		response.Response(r, w, resp, err)
	}
}
