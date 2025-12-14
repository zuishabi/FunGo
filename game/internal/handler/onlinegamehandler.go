// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package handler

import (
	"fmt"
	"fungo/common/response"
	"fungo/game/internal/logic"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"fungo/game/internal/svc"
	"fungo/game/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func OnlineGameHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.OnlineGameReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		// 拼接到本地 uploads 目录
		if req.FileName == "play" {
			req.FileName = "index.html"
		}
		fullPath := filepath.Join("uploads/"+strconv.Itoa(int(req.ID))+"/", req.FileName)
		fmt.Println(req.FileName)

		// 额外安全检查：确保最终文件在 uploads 目录下
		absBase, err := filepath.Abs("uploads")
		if err != nil {
			http.Error(w, "server error", http.StatusInternalServerError)
			return
		}
		absFile, err := filepath.Abs(fullPath)
		if err != nil {
			http.NotFound(w, r)
			return
		}
		if !strings.HasPrefix(absFile, absBase) {
			http.NotFound(w, r)
			return
		}

		// 执行游玩在线游戏的逻辑
		if err := logic.NewOnlineGameLogic(r.Context(), svcCtx).OnlineGame(&req); err != nil {
			response.Response(r, w, nil, err)
		}

		// 最终返回文件（http.ServeFile 会自动设置 Content-Type）
		http.ServeFile(w, r, absFile)
	}
}
