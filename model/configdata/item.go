package configdata

// Item 道具表
type Item struct {
	ID          int64  `json:"id"`
	Type        int32  `json:"type"`
	ItemType    int32  `json:"item_type"`
	Name        string `json:"name"`
	Value       int32  `json:"value"`
	Sell        int32  `json:"sell"`
	ItemSkillID int32  `json:"item_skill_id"`
}
