package configdata

// Level 道具表
type Level struct {
	ID         int64  `json:"id"`
	Level      int64  `json:"level"`
	Type       int32  `json:"type"`
	Name       string `json:"name"`
	Exp        int64  `json:"exp"`
	JSONString string `json:"-"`
}
