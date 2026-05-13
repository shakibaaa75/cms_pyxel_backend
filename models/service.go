package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type CMSService struct {
	ID              bson.ObjectID    `bson:"_id,omitempty" json:"id"`
	SiteID          string           `bson:"siteId"        json:"siteId"`
	Title           string           `bson:"title"         json:"title"`
	Slug            string           `bson:"slug"          json:"slug"`
	Image           string           `bson:"image"         json:"image"`
	Description     string           `bson:"description"   json:"description"`
	FullDescription string           `bson:"fullDescription" json:"fullDescription"`
	Features        []string         `bson:"features"      json:"features"`
	Benefits        []string         `bson:"benefits"      json:"benefits"`
	Process         []CMSProcessStep `bson:"process"       json:"process"`
	Gallery         []string         `bson:"gallery"       json:"gallery"`
	FAQs            []CMSFAQ         `bson:"faqs"          json:"faqs"`
	PriceRange      string           `bson:"priceRange"    json:"priceRange"`
	Timeline        string           `bson:"timeline"      json:"timeline"`
	Published       bool             `bson:"published"     json:"published"`
	SortOrder       int              `bson:"sortOrder"     json:"sortOrder"`
	CreatedAt       time.Time        `bson:"createdAt"     json:"createdAt"`
	UpdatedAt       time.Time        `bson:"updatedAt"     json:"updatedAt"`
}

type CMSProcessStep struct {
	Step        int    `bson:"step"        json:"step"`
	Title       string `bson:"title"       json:"title"`
	Description string `bson:"description" json:"description"`
}

type CMSFAQ struct {
	Question string `bson:"question" json:"question"`
	Answer   string `bson:"answer"   json:"answer"`
}
