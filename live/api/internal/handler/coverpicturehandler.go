// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package handler

import (
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"fungo/live/api/internal/svc"
	"fungo/live/api/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func CoverPictureHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CoverPictureReq
		if err := httpx.Parse(r, &req); err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}

		// uploads/<room_id>/<name>
		filePath := filepath.Join("uploads", strconv.FormatUint(req.RoomID, 10), req.Name)

		fi, err := os.Stat(filePath)
		if err != nil || fi.IsDir() {
			http.NotFound(w, r)
			return
		}

		// 设置 Content-Type，避免被当成 text/html 等
		if ctype := mime.TypeByExtension(filepath.Ext(filePath)); ctype != "" {
			w.Header().Set("Content-Type", ctype)
		}
		w.Header().Set("Cache-Control", "public, max-age=3600")

		http.ServeFile(w, r, filePath)
	}
}
