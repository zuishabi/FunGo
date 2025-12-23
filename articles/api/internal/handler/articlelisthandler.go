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

func ArticleListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ArticleListReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewArticleListLogic(r.Context(), svcCtx)
		resp, err := l.ArticleList(&req)
		response.Response(r, w, resp, err)
	}
}
