package model

type CommentDto struct {
	Id          int64       `json:"id"`
	Commentator string      `json:"commentator"`
	Description string      `json:"description"`
	CommentType CommentType `json:"commentType"`
	CreatedAt   string      `json:"createdAt"`
}

type CommentType string

const (
	Internal CommentType = "INTERNAL"
	External             = "EXTERNAL"
)
