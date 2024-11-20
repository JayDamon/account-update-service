package users

type PrivateToken struct {
	UserId       *string `json:"id"`
	PrivateToken *string `json:"privateToken"`
	ItemId       *string `json:"itemId"`
	IsNew        *bool   `json:"isNew"`
	Cursor       *string `json:"cursor"`
}
