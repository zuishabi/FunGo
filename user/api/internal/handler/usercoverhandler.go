// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package handler

import (
	"mime"
	"net/http"
	"path/filepath"
	"strconv"

	"fungo/user/api/internal/svc"
	"fungo/user/api/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func userCoverHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UserCoverReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		// 从upload目录读取头像文件并进行传输
		// uploads/<room_id>/<name>
		filePath := filepath.Join("uploads", strconv.FormatUint(req.UserID, 10))
		matches, err := filepath.Glob(filepath.Join(filePath, "*"))
		if err != nil || len(matches) == 0 {
			filePath = "uploads/default.png"
		} else {
			filePath = matches[0]
		}

		// 设置 Content-Type，避免被当成 text/html 等
		if ctype := mime.TypeByExtension(filepath.Ext(filePath)); ctype != "" {
			w.Header().Set("Content-Type", ctype)
		}
		w.Header().Set("Cache-Control", "public, max-age=3600")

		http.ServeFile(w, r, filePath)
	}
}
