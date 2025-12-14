// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package handler

import (
	"fmt"
	"fungo/game/model"
	"net/http"
	"path/filepath"
	"strconv"

	"fungo/game/internal/svc"
	"fungo/game/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func DownloadGameHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GameFileReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		// 先从数据库中获取文件名
		gameInfo := model.GameInfo{}
		if err := svcCtx.Db.Select("download_file").Where("id = ?", req.ID).First(&gameInfo).Error; err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		if gameInfo.DownloadFile == "" {
			httpx.ErrorCtx(r.Context(), w, http.ErrMissingFile)
			return
		}
		
		// 这会告诉浏览器下载的文件名应该是什么
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", gameInfo.DownloadFile))

		baseDir := "uploads"
		filePath := filepath.Join(baseDir, strconv.Itoa(int(req.ID)), gameInfo.DownloadFile)
		http.ServeFile(w, r, filePath)
	}
}
