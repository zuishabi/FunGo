// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package handler

import (
	"fungo/common/response"
	"net/http"

	"fungo/animate/api/internal/logic"
	"fungo/animate/api/internal/svc"
)

func recommendAnimateListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logic.NewRecommendAnimateListLogic(r.Context(), svcCtx)
		resp, err := l.RecommendAnimateList()
		response.Response(r, w, resp, err)
	}
}
