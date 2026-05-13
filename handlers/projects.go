package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"cms-backend/config"
	"cms-backend/database"
	"cms-backend/models"
	"cms-backend/utils"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func requireMongo(w http.ResponseWriter) bool {
	if database.MongoDB == nil {
		utils.JSONErr(w, "database not configured", http.StatusServiceUnavailable)
		return false
	}
	return true
}

func ListProjects(w http.ResponseWriter, r *http.Request) {
	if !requireMongo(w) {
		return
	}
	siteID, ok := utils.RequireSiteID(w, r)
	if !ok {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cursor, err := database.MongoDB.Collection("projects").Find(ctx,
		bson.M{"siteId": siteID},
		options.Find().SetSort(bson.D{{Key: "createdAt", Value: -1}}))
	if err != nil {
		utils.JSONErr(w, "failed to fetch", http.StatusInternalServerError)
		return
	}
	var list []models.Project
	cursor.All(ctx, &list)
	if list == nil {
		list = []models.Project{}
	}
	utils.JSON200(w, list)
}

func CreateProject(w http.ResponseWriter, r *http.Request) {
	if !requireMongo(w) {
		return
	}
	siteID, ok := utils.RequireSiteID(w, r)
	if !ok {
		return
	}
	var p models.Project
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		utils.JSONErr(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if err := utils.ValidateProject(&p); err != nil {
		utils.JSONErr(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	p.ID = bson.NewObjectID()
	p.SiteID = siteID
	p.AccessCode = utils.GenerateAccessCode()
	p.MagicToken = utils.GenerateMagicToken()
	p.MagicExpiry = time.Now().Add(90 * 24 * time.Hour)
	p.Updates = []models.Update{}
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()
	if p.Status == "" {
		p.Status = "planning"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if _, err := database.MongoDB.Collection("projects").InsertOne(ctx, p); err != nil {
		utils.JSONErr(w, "failed to create project", http.StatusInternalServerError)
		return
	}

	magicLink := fmt.Sprintf("%s/track?t=%s", config.FRONTEND_URL, p.MagicToken)
	emailSent := false
	if p.ClientEmail != "" {
		utils.SendProjectWelcomeEmail(p.ClientEmail, p.ClientName, p.Title, magicLink)
		emailSent = true
	}

	slog.Info("project created", "id", p.ID, "title", p.Title, "siteId", siteID)
	w.WriteHeader(http.StatusCreated)
	utils.JSON200(w, map[string]any{
		"project":    p,
		"magicLink":  magicLink,
		"emailSent":  emailSent,
		"accessCode": p.AccessCode,
	})
}

func UpdateProject(w http.ResponseWriter, r *http.Request) {
	if !requireMongo(w) {
		return
	}
	siteID, ok := utils.RequireSiteID(w, r)
	if !ok {
		return
	}
	rawID := utils.ExtractPathSegment(r.URL.Path, "/api/admin/projects/", 0)
	id, ok := utils.ParseObjectID(w, rawID)
	if !ok {
		return
	}

	var fields bson.M
	json.NewDecoder(r.Body).Decode(&fields)
	delete(fields, "_id")
	delete(fields, "siteId")
	delete(fields, "createdAt")
	delete(fields, "magicToken")
	delete(fields, "magicExpiry")
	fields["updatedAt"] = time.Now()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var p models.Project
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	err := database.MongoDB.Collection("projects").FindOneAndUpdate(ctx,
		bson.M{"_id": id, "siteId": siteID},
		bson.M{"$set": fields}, opts).Decode(&p)
	if err != nil {
		utils.JSONErr(w, "project not found", http.StatusNotFound)
		return
	}
	utils.JSON200(w, p)
}

func DeleteProject(w http.ResponseWriter, r *http.Request) {
	if !requireMongo(w) {
		return
	}
	siteID, ok := utils.RequireSiteID(w, r)
	if !ok {
		return
	}
	rawID := utils.ExtractPathSegment(r.URL.Path, "/api/admin/projects/", 0)
	id, ok := utils.ParseObjectID(w, rawID)
	if !ok {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if _, err := database.MongoDB.Collection("projects").DeleteOne(ctx,
		bson.M{"_id": id, "siteId": siteID}); err != nil {
		utils.JSONErr(w, "failed to delete", http.StatusInternalServerError)
		return
	}
	utils.JSON200(w, map[string]string{"message": "deleted"})
}

func AddUpdate(w http.ResponseWriter, r *http.Request) {
	if !requireMongo(w) {
		return
	}
	siteID, ok := utils.RequireSiteID(w, r)
	if !ok {
		return
	}
	rawID := utils.ExtractPathSegment(r.URL.Path, "/api/admin/projects/", 0)
	id, ok := utils.ParseObjectID(w, rawID)
	if !ok {
		return
	}

	var u models.Update
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		utils.JSONErr(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if err := utils.ValidateUpdate(&u); err != nil {
		utils.JSONErr(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	u.ID = utils.GenerateSecureID()
	u.CreatedAt = time.Now()
	if u.Images == nil {
		u.Images = []string{}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	setFields := bson.M{"updatedAt": time.Now()}
	if u.Progress > 0 {
		setFields["progress"] = u.Progress
	}
	update := bson.M{
		"$push": bson.M{"updates": u},
		"$set":  setFields,
	}

	var p models.Project
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	err := database.MongoDB.Collection("projects").FindOneAndUpdate(ctx,
		bson.M{"_id": id, "siteId": siteID},
		update, opts).Decode(&p)
	if err != nil {
		utils.JSONErr(w, "project not found", http.StatusNotFound)
		return
	}
	utils.JSON200(w, p)
}

// GetProjectByCode and GetProjectByToken are public routes — no siteId needed
// since access codes and magic tokens are globally unique per project.

func GetProjectByCode(w http.ResponseWriter, r *http.Request) {
	if !requireMongo(w) {
		return
	}
	code := strings.TrimSpace(r.URL.Query().Get("code"))
	if code == "" {
		utils.JSONErr(w, "access code required", http.StatusBadRequest)
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var p models.Project
	err := database.MongoDB.Collection("projects").FindOne(ctx,
		bson.M{"accessCode": code}).Decode(&p)
	if err != nil {
		time.Sleep(100 * time.Millisecond)
		utils.JSONErr(w, "project not found", http.StatusNotFound)
		return
	}
	utils.JSON200(w, p)
}

func GetProjectByToken(w http.ResponseWriter, r *http.Request) {
	if !requireMongo(w) {
		return
	}
	token := strings.TrimSpace(r.URL.Query().Get("t"))
	if token == "" {
		utils.JSONErr(w, "token required", http.StatusBadRequest)
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var p models.Project
	err := database.MongoDB.Collection("projects").FindOne(ctx,
		bson.M{"magicToken": token}).Decode(&p)
	if err != nil {
		utils.JSONErr(w, "invalid or expired link", http.StatusNotFound)
		return
	}
	if time.Now().After(p.MagicExpiry) {
		utils.JSONErr(w, "link has expired", http.StatusGone)
		return
	}
	utils.JSON200(w, p)
}
