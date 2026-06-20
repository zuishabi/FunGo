// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package config

import "github.com/zeromicro/go-zero/rest"

type Config struct {
	rest.RestConf
	DeepSeek struct {
		ApiKey  string
		BaseURL string
		Model   string
	}
	Mysql struct {
		DataSource string
	}
}
