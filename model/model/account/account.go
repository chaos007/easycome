package account

// Account 账号
type Account struct {
	ID            int64  `xorm:"autoincr id"`
	UID           string `xorm:"index uid"`
	UserName      string
	Password      string `xorm:"index"`
	LastServerKey string
	Cookie        string
	LastLoginTime int64
	SignUpTime    int64
	MachineCode   string `xorm:"index"`
}
