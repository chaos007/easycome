package enum

// 服务器类型，对应的etcd目录
const (
	ServerTypeAgent  = "agent"
	ServerTypeGame   = "game"
	ServerTypeCenter = "center"
	ServerTypeUnique = "unique"
	ServerTypeAll    = "all"
)

// 字符串协议末尾标识
const (
	ToAgentString  = "ToAgent"
	ToGameString   = "ToGame"
	ToAllString    = "ToAll"
	ToCenterString = "ToCenter"
	ToUniqueString = "ToUnique"
)

// session常用量
const (
	SessKeyExchange = 0x1 // 是否已经交换完毕KEY
	SessEncrypt     = 0x2 // 是否可以开始加密
	SessKickedOut   = 0x4 // 踢掉
	SessAuthorized  = 0x8 // 已授权访问
	// SessHasClose    = 0x16 // 已授权访问
)
