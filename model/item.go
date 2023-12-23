package model

type Item struct {
	ItemId      uint64 `json:"item_id"`
	UserId      uint64 `json:"user_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      int    `json:"status"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

type ItemInput struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      int    `json:"status"`
}
