package model

type CommentDto struct {
	tableName struct{} `sql:"comment" pg:",discard_unknown_columns"`

	Id            int64       `sql:"id" json:"id"`
	Commentator   string      `sql:"commentator" json:"commentator"`
	Description   string      `sql:"description" json:"description"`
	CommentType   CommentType `sql:"comment_type" json:"commentType"`
	CreatedAt     string      `sql:"created_at" json:"createdAt"`
	ApplicationId int64       `sql:"application_id" json:"applicationId"`
}

type CommentType string

const (
	Internal CommentType = "INTERNAL"
	External             = "EXTERNAL"
)
