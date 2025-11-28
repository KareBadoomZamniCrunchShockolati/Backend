package model

import (
	"time"
)

type Comment struct {
	ID         uint
	EntityType CommentType
	EntityID   uint
	UserID     uint
	Content    string
	ParentID   *uint
	CreatedAt  time.Time
	UpdatedAt  time.Time
	Replies    []*Comment // Add this field for nested structure
}

// Helper method to add replies
func (c *Comment) AddReply(reply *Comment) {
	if c.Replies == nil {
		c.Replies = make([]*Comment, 0)
	}
	c.Replies = append(c.Replies, reply)
}
