package users

type PrivateToken struct {
	UserID       *string `json:"id"`
	PrivateToken *string `json:"privateToken"`
	ItemId       *string `json:"itemId"`
}
