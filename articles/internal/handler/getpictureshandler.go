// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package handler

import (
	"fmt"
	"fungo/articles/internal/svc"
	"fungo/articles/internal/types"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetPicturesHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.PictureReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		// 基本校验：id 合法且 md5 不包含路径分隔符（防止目录穿越）
		if req.ID == 0 {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}
		if req.MD5 == "" || strings.ContainsAny(req.MD5, "/\\") {
			http.Error(w, "invalid md5", http.StatusBadRequest)
			return
		}

		// 构造本地文件路径：uploads/<id>/<md5>
		baseDir := "uploads"
		filePath := filepath.Join(baseDir, strconv.Itoa(int(req.ID)), req.MD5)

		// 防护：把构造的路径清理并确保是在 uploads 目录下
		cleanPath := filepath.Clean(filePath)
		absBase, err := filepath.Abs(baseDir)
		if err == nil {
			absFile, err2 := filepath.Abs(cleanPath)
			if err2 == nil {
				if !strings.HasPrefix(absFile, absBase+string(filepath.Separator)) && absFile != absBase {
					http.Error(w, "access denied", http.StatusForbidden)
					return
				}
			}
		}

		f, err := os.Open(cleanPath)
		if err != nil {
			if os.IsNotExist(err) {
				http.NotFound(w, r)
				return
			}
			http.Error(w, "open file error", http.StatusInternalServerError)
			return
		}
		defer f.Close()

		// 获取文件信息用于 Content-Length / Last-Modified
		stat, err := f.Stat()
		if err != nil {
			http.Error(w, "stat file error", http.StatusInternalServerError)
			return
		}
		if stat.IsDir() {
			http.NotFound(w, r)
			return
		}

		// 读取前 512 字节检测 Content-Type
		buf := make([]byte, 512)
		n, _ := f.Read(buf)
		contentType := http.DetectContentType(buf[:n])
		if contentType == "application/octet-stream" {
			// 尝试根据扩展名提供更合适的类型
			ext := strings.ToLower(filepath.Ext(stat.Name()))
			if ext != "" {
				if ct := mimeByExtension(ext); ct != "" {
					contentType = ct
				}
			}
		}

		// 设置响应头
		w.Header().Set("Content-Type", contentType)
		w.Header().Set("Content-Length", strconv.FormatInt(stat.Size(), 10))
		w.Header().Set("Cache-Control", "public, max-age=300")
		w.Header().Set("Last-Modified", stat.ModTime().UTC().Format(http.TimeFormat))

		// 将读指针回到文件开头，然后流式拷贝到响应
		if _, err := f.Seek(0, io.SeekStart); err != nil {
			http.Error(w, "seek error", http.StatusInternalServerError)
			return
		}
		// 返回 200 并写文件内容
		http.ServeContent(w, r, stat.Name(), stat.ModTime(), f)
		fmt.Println("发送成功")
	}
}

func mimeByExtension(ext string) string {
	switch ext {
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".gif":
		return "image/gif"
	case ".webp":
		return "image/webp"
	case ".bmp":
		return "image/bmp"
	case ".svg":
		return "image/svg+xml"
	default:
		return ""
	}
}
