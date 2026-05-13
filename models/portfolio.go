package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type CMSProject struct {
	ID              bson.ObjectID `bson:"_id,omitempty"         json:"id"`
	SiteID          string        `bson:"siteId"                json:"siteId"`
	Title           string        `bson:"title"                 json:"title"`
	Slug            string        `bson:"slug"                  json:"slug"`
	Image           string        `bson:"image"                 json:"image"`
	Gallery         []string      `bson:"gallery"               json:"gallery"`
	Cost            string        `bson:"cost"                  json:"cost"`
	Client          string        `bson:"client"                json:"client"`
	Year            string        `bson:"year"                  json:"year"`
	Location        string        `bson:"location"              json:"location"`
	Category        string        `bson:"category"              json:"category"`
	Description     string        `bson:"description"           json:"description"`
	FullDescription string        `bson:"fullDescription"       json:"fullDescription"`
	Features        []string      `bson:"features"              json:"features"`
	Testimonial     *CMSLeave     `bson:"testimonial,omitempty" json:"testimonial,omitempty"`
	Featured        bool          `bson:"featured"              json:"featured"`
	Published       bool          `bson:"published"             json:"published"`
	SortOrder       int           `bson:"sortOrder"             json:"sortOrder"`
	CreatedAt       time.Time     `bson:"createdAt"             json:"createdAt"`
	UpdatedAt       time.Time     `bson:"updatedAt"             json:"updatedAt"`
}

type CMSLeave struct {
	Text   string `bson:"text"   json:"text"`
	Author string `bson:"author" json:"author"`
	Role   string `bson:"role"   json:"role"`
}
