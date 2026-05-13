package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type CMSBlogPost struct {
	ID          bson.ObjectID `bson:"_id,omitempty"         json:"id"`
	SiteID      string        `bson:"siteId"                json:"siteId"`
	Title       string        `bson:"title"                 json:"title"`
	Slug        string        `bson:"slug"                  json:"slug"`
	Excerpt     string        `bson:"excerpt"               json:"excerpt"`
	Content     string        `bson:"content"               json:"content"`
	Image       string        `bson:"image"                 json:"image"`
	Category    string        `bson:"category"              json:"category"`
	Tags        []string      `bson:"tags"                  json:"tags"`
	ReadTime    int           `bson:"readTime"              json:"readTime"`
	Published   bool          `bson:"published"             json:"published"`
	PublishedAt *time.Time    `bson:"publishedAt,omitempty" json:"publishedAt,omitempty"`
	CreatedAt   time.Time     `bson:"createdAt"             json:"createdAt"`
	UpdatedAt   time.Time     `bson:"updatedAt"             json:"updatedAt"`
}
