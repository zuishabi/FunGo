// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package handler

import (
	"fmt"
	"net/http"

	"fungo/animate/api/internal/svc"
	"fungo/animate/api/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func getAnimateHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetAnimateReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		if req.Name == "main" {
			http.ServeFile(w, r, fmt.Sprintf("%s/animates/%d/%d/index.m3u8", svcCtx.Config.AnimatePath, req.ID, req.NID))
		} else {
			http.ServeFile(w, r, fmt.Sprintf("%s/animates/%d/%d/%s", svcCtx.Config.AnimatePath, req.ID, req.NID, req.Name))
		}
	}
}
