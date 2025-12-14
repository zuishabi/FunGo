// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package handler

import (
	"encoding/json"
	"errors"
	"fungo/common/response"
	"fungo/game/internal/svc"
	"fungo/game/model"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func UploadGameHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uid, ok := r.Context().Value("user_id").(uint64)
		if !ok {
			response.Response(r, w, nil, errors.New("请先登录"))
			return
		}
		r.Body = http.MaxBytesReader(w, r.Body, 500<<20)

		if err := r.ParseMultipartForm(500 << 20); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		// 将上传的信息保存到数据库中
		downloadFile, downloadFileHeader, err := r.FormFile("download_file")
		fileName := ""
		if err == nil {
			fileName = downloadFileHeader.Filename
			defer downloadFile.Close()
		}

		onlineFiles := r.MultipartForm.File["online_files"]

		links := r.Form.Get("links")

		// 图片内容信息
		pictures := r.MultipartForm.File["pictures"]
		pictureInfos := model.Pictures{Name: make([]string, len(pictures))}
		for i, v := range pictures {
			pictureInfos.Name[i] = v.Filename
		}
		pictureInfosJson, _ := json.Marshal(pictureInfos)

		// 封面信息
		cover, coverFileHeader, err := r.FormFile("cover")
		if err != nil {
			response.Response(r, w, nil, errors.New("保存封面失败"))
			return
		}
		defer cover.Close()

		info := &model.GameInfo{
			Title:            r.FormValue("title"),
			Author:           uid,
			Description:      r.FormValue("description"),
			Links:            links,
			Pictures:         string(pictureInfosJson),
			Cover:            coverFileHeader.Filename,
			DownloadFile:     fileName,
			CanOnlinePlaying: !(len(onlineFiles) == 0),
		}
		if err := svcCtx.Db.Create(info).Error; err != nil {
			response.Response(r, w, nil, err)
			return
		}

		// 保存目录
		uploadDir := filepath.Join("uploads", strconv.FormatUint(info.ID, 10))
		if err := os.MkdirAll(uploadDir, 0755); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		// 保存上传的封面
		coverPath := filepath.Join(uploadDir, coverFileHeader.Filename)
		coverFile, err := os.Create(coverPath)
		if err == nil {
			if _, err = io.Copy(coverFile, cover); err != nil {
				response.Response(r, w, nil, errors.New("保存文件失败"))
				coverFile.Close()
				return
			}
			coverFile.Close()
		}

		// 保存上传的图片
		for i := range pictures {
			fh := pictures[i]
			src, err := fh.Open()
			if err != nil {
				response.Response(r, w, nil, errors.New("图片上传失败"))
				return
			}

			dstPath := filepath.Join(uploadDir, fh.Filename)
			dst, err := os.Create(dstPath)
			if err != nil {
				src.Close()
				response.Response(r, w, nil, errors.New("图片上传失败"))
				return
			}

			if _, err := io.Copy(dst, src); err != nil {
				response.Response(r, w, nil, errors.New("图片上传失败"))
				return
			}
			dst.Close()
			src.Close()
		}

		// 保存上传的游戏文件
		if downloadFile != nil {
			dstPath := filepath.Join(uploadDir, downloadFileHeader.Filename)
			dst, err := os.Create(dstPath)
			if err == nil {
				if _, err = io.Copy(dst, downloadFile); err != nil {
					dst.Close()
					response.Response(r, w, nil, errors.New("保存文件失败"))
					return
				}
				dst.Close()
			}
		}

		// 保存在线游玩的游戏文件
		for i := range onlineFiles {
			fh := onlineFiles[i]
			src, err := fh.Open()
			if err != nil {
				response.Response(r, w, nil, errors.New("文件上传失败"))
				return
			}

			dstPath := filepath.Join(uploadDir, fh.Filename)
			dst, err := os.Create(dstPath)
			if err != nil {
				src.Close()
				response.Response(r, w, nil, errors.New("文件上传失败"))
				return
			}

			if _, err := io.Copy(dst, src); err != nil {
				response.Response(r, w, nil, errors.New("文件上传失败"))
				return
			}
			dst.Close()
			src.Close()
		}
	}
}
