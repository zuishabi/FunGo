package handler

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fungo/articles/internal/svc"
	"fungo/articles/model"
	"fungo/common/response"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/redis/go-redis/v9"
)

const maxUploadSize = int64(50 << 20) // 100MB

// saveUploadedFiles 假定 r.MultipartForm 已经准备好，不再调用 ParseMultipartForm
func saveUploadedFiles(r *http.Request, storeDir string) ([]string, error) {
	if r.MultipartForm == nil {
		return nil, nil
	}

	var fhs []*multipart.FileHeader
	if arr := r.MultipartForm.File["files"]; len(arr) > 0 {
		fhs = append(fhs, arr...)
	}
	if arr := r.MultipartForm.File["file"]; len(arr) > 0 {
		fhs = append(fhs, arr...)
	}
	if len(fhs) == 0 {
		return nil, nil
	}

	if err := os.MkdirAll(storeDir, 0755); err != nil {
		return nil, err
	}

	var saved []string
	for _, fh := range fhs {
		src, err := fh.Open()
		if err != nil {
			continue
		}
		data, err := io.ReadAll(src)
		src.Close()
		if err != nil {
			continue
		}

		sum := md5.Sum(data)
		md5name := hex.EncodeToString(sum[:])
		ext := strings.ToLower(filepath.Ext(fh.Filename))
		if ext == "" {
			ext = ".bin"
		}
		outName := md5name + ext
		outPath := filepath.Join(storeDir, outName)

		if _, err := os.Stat(outPath); os.IsNotExist(err) {
			if err := os.WriteFile(outPath, data, 0644); err != nil {
				continue
			}
		}

		saved = append(saved, outName)
	}

	return saved, nil
}

func CreateArticleHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 先限制最大读取大小
		r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)

		// 立即使用相同的 max 调用 ParseMultipartForm，确保后续 FormValue/文件读取一致
		if err := r.ParseMultipartForm(maxUploadSize); err != nil {
			l := strings.ToLower(err.Error())
			if strings.Contains(l, "request body too large") || strings.Contains(l, "multipart: message too large") {
				http.Error(w, "上传文件过大", http.StatusRequestEntityTooLarge)
				return
			}
			// 解析错误返回通用错误
			response.Response(r, w, nil, err)
			return
		}

		// 读取 user_id（确保存在并处理错误）
		uidVal := r.Context().Value("user_id")
		if uidVal == nil {
			response.Response(r, w, nil, errors.New("创建失败"))
			return
		}
		uidNum, ok := uidVal.(json.Number)
		if !ok {
			response.Response(r, w, nil, errors.New("创建失败"))
			return
		}
		uid, err := uidNum.Int64()
		if err != nil {
			response.Response(r, w, nil, errors.New("创建失败"))
			return
		}

		// 现在安全读取表单字段（ParseMultipartForm 已执行）
		content := r.FormValue("content")
		section := r.FormValue("section")
		title := r.FormValue("title")
		sec, err := strconv.Atoi(section)
		if err != nil {
			response.Response(r, w, nil, errors.New("创建失败"))
			return
		}

		tx := svcCtx.Db.Begin()
		temp := model.Article{
			Content: content,
			Section: uint32(sec),
			Author:  uint64(uid),
			Title:   title,
		}
		if err := tx.Create(&temp).Error; err != nil {
			response.Response(r, w, nil, errors.New("创建失败"))
			tx.Rollback()
			return
		}

		// 直接使用已经解析好的 r.MultipartForm 来保存文件
		savedFiles, err := saveUploadedFiles(r, "uploads/"+strconv.Itoa(int(temp.ID))+"/")
		if err != nil {
			l := strings.ToLower(err.Error())
			if strings.Contains(l, "request body too large") || strings.Contains(l, "multipart: message too large") {
				http.Error(w, "上传文件过大", http.StatusRequestEntityTooLarge)
				tx.Rollback()
				return
			}
			response.Response(r, w, nil, err)
			tx.Rollback()
			return
		}

		pics, err := json.Marshal(model.Pictures{Name: savedFiles})
		if err != nil {
			response.Response(r, w, nil, errors.New("创建失败"))
			tx.Rollback()
			return
		}
		if tx.Model(&temp).Update("pictures", string(pics)).Error != nil {
			response.Response(r, w, nil, errors.New("创建失败"))
			tx.Rollback()
			return
		}

		// 再将数据写入redis中
		if err = svcCtx.RedisClient.ZAdd(context.Background(), "section-"+strconv.Itoa(int(temp.Section))+"-list", redis.Z{
			Score:  float64(temp.CreatedAt.UnixNano()) / 1e6,
			Member: temp.ID,
		}).Err(); err != nil {
			response.Response(r, w, nil, errors.New("创建失败"))
			tx.Rollback()
			return
		}

		if err = svcCtx.RedisClient.ZAdd(context.Background(), "section-1-list", redis.Z{
			Score:  float64(temp.CreatedAt.UnixNano()) / 1e6,
			Member: temp.ID,
		}).Err(); err != nil {
			response.Response(r, w, nil, errors.New("创建失败"))
			tx.Rollback()
			return
		}

		tx.Commit()

		response.Response(r, w, nil, nil)
	}
}
