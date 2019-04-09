package configmodel

//HTTPConfig httpmodule的配置
type HTTPConfig struct {
	//HTTPAddr 监听地址
	HTTPAddr string
	Timeout  int32
}
