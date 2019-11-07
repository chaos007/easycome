package enum

import "errors"

// ErrSessionCloseByLogic 逻辑层向下通知的会话关闭
var ErrSessionCloseByLogic = errors.New("Session Close By Logic")

// ErrPlayerNotLogin 玩家还没登陆
var ErrPlayerNotLogin = errors.New("Player Have Not Login Yet")
