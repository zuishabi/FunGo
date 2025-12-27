// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package handler

import (
	"context"
	"errors"
	"fmt"
	"fungo/common/distributedLock"
	"fungo/common/response"
	"fungo/live/model"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"fungo/live/api/internal/logic"
	"fungo/live/api/internal/svc"
)

func EditLiveRoomReqHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logic.NewEditLiveRoomReqLogic(r.Context(), svcCtx)
		uid := r.Context().Value("user_id").(uint64)

		// 获取用户直播间
		room := model.UserRoom{}
		if svcCtx.Db.Where("uid = ?", uid).First(&room).Error != nil {
			response.Response(r, w, nil, errors.New("未创建直播间"))
			return
		}
		roomInfo := model.RoomInfo{
			RoomID: room.RoomID,
		}
		svcCtx.Db.Where("room_id = ?", room.RoomID).First(&roomInfo)

		coverFile, coverFileHeader, err := r.FormFile("cover")
		if err == nil {
			roomInfo.Cover = coverFileHeader.Filename

			// 1\) 确保目录 `uploads/<room_id>` 存在
			dir := filepath.Join("uploads", fmt.Sprintf("%d", room.RoomID))
			if err := os.MkdirAll(dir, 0o755); err != nil {
				response.Response(r, w, nil, err)
				return
			}

			// 2\) 再创建文件 `uploads/<room_id>/<filename>`
			dstPath := filepath.Join(dir, coverFileHeader.Filename)
			file, err := os.Create(dstPath)
			if err != nil {
				response.Response(r, w, nil, err)
				return
			}
			defer file.Close()
			defer coverFile.Close()

			if _, err := io.Copy(file, coverFile); err != nil {
				response.Response(r, w, nil, err)
				return
			}
		}
		title := r.FormValue("title")
		if title != "" {
			// 保存标题信息
			roomInfo.Title = title
		}
		description := r.FormValue("description")
		if description != "" {
			// 保存描述信息
			roomInfo.Description = description
		}

		// 获取一个分布式锁，只有获得锁的才能修改内
		lock := distributedLock.NewDistributedLock(svcCtx.RedisClient, fmt.Sprintf("live-room-%d-lock", room.RoomID))
		if err := lock.Lock(context.Background(), 0); err == nil {
			svcCtx.Db.Save(&roomInfo)

			// 检查是否开播，开播再同步更新到redis中
			exists, _ := svcCtx.RedisClient.Exists(context.Background(), fmt.Sprintf("live-room-%d", room.RoomID)).Result()
			if exists == 1 {
				svcCtx.RedisClient.HMSet(context.Background(), fmt.Sprintf("live-room-%d", room.RoomID),
					"title", roomInfo.Title,
					"description", roomInfo.Description,
					"cover", roomInfo.Cover,
				)
			}

			lock.Unlock(context.Background())
		}
		err = l.EditLiveRoomReq()
		response.Response(r, w, nil, err)
	}
}
