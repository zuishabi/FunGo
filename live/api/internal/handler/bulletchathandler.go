// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/zeromicro/go-zero/core/logc"
	"github.com/zeromicro/go-zero/core/threading"
	"github.com/zeromicro/go-zero/rest/httpx"

	"fungo/live/api/internal/logic"
	"fungo/live/api/internal/svc"
	"fungo/live/api/internal/types"
)

func BulletChatHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Buffer size of 16 is chosen as a reasonable default to balance throughput and memory usage.
		// You can change this based on your application's needs.
		// if your go-zero version less than 1.8.1, you need to add 3 lines below.
		// w.Header().Set("Content-Type", "text/event-stream")
		// w.Header().Set("Cache-Control", "no-cache")
		// w.Header().Set("Connection", "keep-alive")
		client := make(chan *types.BulletChatMessageRsp, 16)
		var req types.BulletChatMessageReq
		if err := httpx.Parse(r, &req); err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}

		// 将房间信息写入BulletChatServer中
		messageChan := svcCtx.BulletChatServer.CreateConn(req.RoomID)

		l := logic.NewBulletChatLogic(r.Context(), svcCtx, messageChan)
		threading.GoSafeCtx(r.Context(), func() {
			defer close(client)
			defer func() {
				svcCtx.BulletChatServer.DeleteConn(req.RoomID, messageChan)
			}()

			err := l.BulletChat(client)
			if err != nil {
				logc.Errorw(r.Context(), "BulletChatHandler", logc.Field("error", err))
				return
			}
		})

		for {
			select {
			case data, ok := <-client:
				if !ok {
					return
				}
				output, err := json.Marshal(data)
				if err != nil {
					logc.Errorw(r.Context(), "BulletChatHandler", logc.Field("error", err))
					continue
				}

				if _, err := fmt.Fprintf(w, "data: %s\n\n", string(output)); err != nil {
					logc.Errorw(r.Context(), "BulletChatHandler", logc.Field("error", err))
					return
				}
				if flusher, ok := w.(http.Flusher); ok {
					flusher.Flush()
				}
			case <-r.Context().Done():
				return
			}
		}
	}
}
