package model

//Set的NX XX的参数
type RedisSetParam int

const (
	Set_No_NX_XX RedisSetParam = iota //不使用模式
	Set_NX                            //只在键不存在时，才对键进行设置操作
	Set_XX                            //只在键已经存在时，才对键进行设置操作
)

type RedisBitOperation int

const (
	Bit_Defualt RedisBitOperation = iota //无效
	Bit_And                              //逻辑与
	Bit_Or                               //逻辑或
	Bit_Xor                              //逻辑异或
	Bit_Not                              //逻辑非
)
