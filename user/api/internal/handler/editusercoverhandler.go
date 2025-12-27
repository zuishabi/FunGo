// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package handler

import (
	"errors"
	"fungo/common/response"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"fungo/user/api/internal/logic"
	"fungo/user/api/internal/svc"
)

func editUserCoverHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logic.NewEditUserCoverLogic(r.Context(), svcCtx)
		err := l.EditUserCover()
		// 获取用户的uid
		uid, ok := r.Context().Value("user_id").(uint64)
		if !ok {
			response.Response(r, w, nil, errors.New("鉴权失败"))
			return
		}

		// 读取传输的文件
		coverFile, coverFileHeader, err := r.FormFile("cover")
		if err != nil {
			response.Response(r, w, nil, errors.New("上传头像失败"))
			return
		}

		dir := filepath.Join("uploads", strconv.Itoa(int(uid)))
		if err := os.MkdirAll(dir, 0o755); err != nil {
			response.Response(r, w, nil, err)
			return
		}

		dstPath := filepath.Join(dir, coverFileHeader.Filename)
		file, err := os.Create(dstPath)
		if err != nil {
			response.Response(r, w, nil, err)
			return
		}
		defer file.Close()
		defer coverFile.Close()

		if _, err := io.Copy(file, coverFile); err != nil {
			response.Response(r, w, nil, errors.New("保存文件失败"))
			return
		}

		matches, err := filepath.Glob(filepath.Join(dir, "*"))
		if err != nil {
			response.Response(r, w, nil, err)
			return
		}

		for _, v := range matches {
			// 跳过目录
			if fi, statErr := os.Stat(v); statErr != nil || fi.IsDir() {
				continue
			}
			// v 是完整路径，需要用 Base() 比较文件名
			if filepath.Base(v) == coverFileHeader.Filename {
				continue
			}
			_ = os.Remove(v) // 也可改为遇错就返回
		}

		response.Response(r, w, nil, nil)
	}
}
