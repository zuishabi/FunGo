// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package handler

import (
	"fmt"
	"fungo/user/api/internal/svc"
	"fungo/user/api/internal/types"
	"net/http"
	"os"
	"path"
	"strconv"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func getUserAvatarHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetUserAvatarReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		target := path.Join(svcCtx.Config.UserStaticPath, "avatars", strconv.Itoa(int(req.ID)), "avatar.png")
		if _, err := os.Stat(target); err != nil {
			http.ServeFile(w, r, fmt.Sprintf("%s/avatars/default.png", svcCtx.Config.UserStaticPath))
		} else {
			http.ServeFile(w, r, target)
		}
	}
}
