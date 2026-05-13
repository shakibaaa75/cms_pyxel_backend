package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Project struct {
	ID          bson.ObjectID `bson:"_id,omitempty"  json:"id"`
	SiteID      string        `bson:"siteId"         json:"siteId"`
	AccessCode  string        `bson:"accessCode"     json:"accessCode"`
	MagicToken  string        `bson:"magicToken"     json:"-"`
	MagicExpiry time.Time     `bson:"magicExpiry"    json:"-"`
	ClientName  string        `bson:"clientName"     json:"clientName"`
	ClientEmail string        `bson:"clientEmail"    json:"clientEmail"`
	Title       string        `bson:"title"          json:"title"`
	Description string        `bson:"description"    json:"description"`
	Status      string        `bson:"status"         json:"status"`
	Progress    int           `bson:"progress"       json:"progress"`
	StartDate   string        `bson:"startDate"      json:"startDate"`
	EndDate     string        `bson:"endDate"        json:"endDate"`
	Address     string        `bson:"address"        json:"address"`
	Budget      string        `bson:"budget"         json:"budget"`
	Updates     []Update      `bson:"updates"        json:"updates"`
	CreatedAt   time.Time     `bson:"createdAt"      json:"createdAt"`
	UpdatedAt   time.Time     `bson:"updatedAt"      json:"updatedAt"`
}

type Update struct {
	ID          string    `bson:"id"          json:"id"`
	Title       string    `bson:"title"       json:"title"`
	Description string    `bson:"description" json:"description"`
	Phase       string    `bson:"phase"       json:"phase"`
	Progress    int       `bson:"progress"    json:"progress"`
	Images      []string  `bson:"images"      json:"images"`
	CreatedAt   time.Time `bson:"createdAt"   json:"createdAt"`
}
