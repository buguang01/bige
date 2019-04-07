package config

import "time"

//HTTPConfig httpmodule的配置
type HTTPConfig struct {
	//HTTPAddr 监听地址
	HTTPAddr string
	Timeout  time.Duration
}
