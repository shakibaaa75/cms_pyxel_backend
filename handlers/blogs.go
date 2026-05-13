// package handlers

// import (
// 	"context"
// 	"encoding/json"
// 	"net/http"
// 	"strings"
// 	"time"

// 	"cms-backend/database"
// 	"cms-backend/models"
// 	"cms-backend/utils"

// 	"go.mongodb.org/mongo-driver/v2/bson"
// 	"go.mongodb.org/mongo-driver/v2/mongo/options"
// )

// func AdminListBlogs(w http.ResponseWriter, r *http.Request) {
// 	if !utils.RequireMongo(w) {
// 		return
// 	}
// 	siteID, ok := utils.RequireSiteID(w, r)
// 	if !ok {
// 		return
// 	}
// 	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 	defer cancel()
// 	cursor, err := database.MongoDB.Collection("blogs").Find(ctx, bson.M{"siteId": siteID},
// 		options.Find().SetSort(bson.D{{Key: "createdAt", Value: -1}}))
// 	if err != nil {
// 		utils.JSONErr(w, "failed to fetch", http.StatusInternalServerError)
// 		return
// 	}
// 	var posts []models.CMSBlogPost
// 	cursor.All(ctx, &posts)
// 	if posts == nil {
// 		posts = []models.CMSBlogPost{}
// 	}
// 	utils.JSON200(w, posts)
// }

// func AdminCreateBlog(w http.ResponseWriter, r *http.Request) {
// 	if !utils.RequireMongo(w) {
// 		return
// 	}
// 	var post models.CMSBlogPost
// 	if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
// 		utils.JSONErr(w, "invalid request", http.StatusBadRequest)
// 		return
// 	}
// 	if strings.TrimSpace(post.Title) == "" {
// 		utils.JSONErr(w, "title is required", http.StatusUnprocessableEntity)
// 		return
// 	}
// 	if strings.TrimSpace(post.Slug) == "" {
// 		utils.JSONErr(w, "slug is required", http.StatusUnprocessableEntity)
// 		return
// 	}
// 	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 	defer cancel()
// 	post.ID = bson.NewObjectID()
// 	post.CreatedAt = time.Now()
// 	post.UpdatedAt = time.Now()
// 	if post.Tags == nil {
// 		post.Tags = []string{}
// 	}
// 	if post.Published {
// 		now := time.Now()
// 		post.PublishedAt = &now
// 	}
// 	if _, err := database.MongoDB.Collection("blogs").InsertOne(ctx, post); err != nil {
// 		if strings.Contains(err.Error(), "duplicate key") {
// 			utils.JSONErr(w, "a post with that slug already exists for this site", http.StatusConflict)
// 			return
// 		}
// 		utils.JSONErr(w, "failed to create", http.StatusInternalServerError)
// 		return
// 	}
// 	w.WriteHeader(http.StatusCreated)
// 	utils.JSON200(w, post)
// }

// func AdminUpdateBlog(w http.ResponseWriter, r *http.Request) {
// 	if !utils.RequireMongo(w) {
// 		return
// 	}
// 	rawID := utils.ExtractPathSegment(r.URL.Path, "/api/admin/blogs/", 0)
// 	id, ok := utils.ParseObjectID(w, rawID)
// 	if !ok {
// 		return
// 	}
// 	var fields bson.M
// 	json.NewDecoder(r.Body).Decode(&fields)
// 	delete(fields, "_id")
// 	delete(fields, "createdAt")
// 	fields["updatedAt"] = time.Now()
// 	if pub, ok := fields["published"].(bool); ok && pub {
// 		if _, hasPublishedAt := fields["publishedAt"]; !hasPublishedAt {
// 			fields["publishedAt"] = time.Now()
// 		}
// 	}
// 	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 	defer cancel()
// 	var post models.CMSBlogPost
// 	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
// 	err := database.MongoDB.Collection("blogs").FindOneAndUpdate(ctx, bson.M{"_id": id}, bson.M{"$set": fields}, opts).Decode(&post)
// 	if err != nil {
// 		utils.JSONErr(w, "failed to update", http.StatusInternalServerError)
// 		return
// 	}
// 	utils.JSON200(w, post)
// }

// func AdminDeleteBlog(w http.ResponseWriter, r *http.Request) {
// 	if !utils.RequireMongo(w) {
// 		return
// 	}
// 	rawID := utils.ExtractPathSegment(r.URL.Path, "/api/admin/blogs/", 0)
// 	id, ok := utils.ParseObjectID(w, rawID)
// 	if !ok {
// 		return
// 	}
// 	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 	defer cancel()
// 	if _, err := database.MongoDB.Collection("blogs").DeleteOne(ctx, bson.M{"_id": id}); err != nil {
// 		utils.JSONErr(w, "failed to delete", http.StatusInternalServerError)
// 		return
// 	}
// 	utils.JSON200(w, map[string]string{"message": "deleted"})
// }

// func PublicListBlogs(w http.ResponseWriter, r *http.Request) {
// 	if !utils.RequireMongo(w) {
// 		return
// 	}
// 	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 	defer cancel()
// 	cursor, err := database.MongoDB.Collection("blogs").Find(ctx,
// 		bson.M{"published": true},
// 		options.Find().SetSort(bson.D{{Key: "createdAt", Value: -1}}))
// 	if err != nil {
// 		utils.JSONErr(w, "failed to fetch", http.StatusInternalServerError)
// 		return
// 	}
// 	var posts []models.CMSBlogPost
// 	cursor.All(ctx, &posts)
// 	if posts == nil {
// 		posts = []models.CMSBlogPost{}
// 	}
// 	utils.JSON200(w, posts)
// }

// func PublicGetBlog(w http.ResponseWriter, r *http.Request) {
// 	if !utils.RequireMongo(w) {
// 		return
// 	}
// 	slug := utils.ExtractPathSegment(r.URL.Path, "/api/blogs/", 0)
// 	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 	defer cancel()
// 	var post models.CMSBlogPost
// 	err := database.MongoDB.Collection("blogs").FindOne(ctx, bson.M{"slug": slug, "published": true}).Decode(&post)
// 	if err != nil {
// 		utils.JSONErr(w, "not found", http.StatusNotFound)
// 		return
// 	}
// 	utils.JSON200(w, post)
// }

package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"cms-backend/database"
	"cms-backend/models"
	"cms-backend/utils"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// AdminListBlogs — requires siteId, returns all blogs for that site
func AdminListBlogs(w http.ResponseWriter, r *http.Request) {
	if !utils.RequireMongo(w) {
		return
	}
	siteID, ok := utils.RequireSiteID(w, r)
	if !ok {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cursor, err := database.MongoDB.Collection("blogs").Find(ctx, bson.M{"siteId": siteID},
		options.Find().SetSort(bson.D{{Key: "createdAt", Value: -1}}))
	if err != nil {
		utils.JSONErr(w, "failed to fetch", http.StatusInternalServerError)
		return
	}
	var posts []models.CMSBlogPost
	cursor.All(ctx, &posts)
	if posts == nil {
		posts = []models.CMSBlogPost{}
	}
	utils.JSON200(w, posts)
}

// AdminCreateBlog — creates blog for specific site
func AdminCreateBlog(w http.ResponseWriter, r *http.Request) {
	if !utils.RequireMongo(w) {
		return
	}
	siteID, ok := utils.RequireSiteID(w, r)
	if !ok {
		return
	}
	var post models.CMSBlogPost
	if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
		utils.JSONErr(w, "invalid request", http.StatusBadRequest)
		return
	}
	if strings.TrimSpace(post.Title) == "" {
		utils.JSONErr(w, "title is required", http.StatusUnprocessableEntity)
		return
	}
	if strings.TrimSpace(post.Slug) == "" {
		utils.JSONErr(w, "slug is required", http.StatusUnprocessableEntity)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	post.ID = bson.NewObjectID()
	post.SiteID = siteID
	post.CreatedAt = time.Now()
	post.UpdatedAt = time.Now()
	if post.Tags == nil {
		post.Tags = []string{}
	}
	if post.Published {
		now := time.Now()
		post.PublishedAt = &now
	}
	if _, err := database.MongoDB.Collection("blogs").InsertOne(ctx, post); err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			utils.JSONErr(w, "a post with that slug already exists for this site", http.StatusConflict)
			return
		}
		utils.JSONErr(w, "failed to create", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	utils.JSON200(w, post)
}

// AdminUpdateBlog — updates blog, ensures site isolation
func AdminUpdateBlog(w http.ResponseWriter, r *http.Request) {
	if !utils.RequireMongo(w) {
		return
	}
	siteID, ok := utils.RequireSiteID(w, r)
	if !ok {
		return
	}
	rawID := utils.ExtractPathSegment(r.URL.Path, "/api/admin/blogs/", 0)
	id, ok := utils.ParseObjectID(w, rawID)
	if !ok {
		return
	}
	var fields bson.M
	json.NewDecoder(r.Body).Decode(&fields)
	delete(fields, "_id")
	delete(fields, "createdAt")
	delete(fields, "siteId")
	fields["updatedAt"] = time.Now()
	if pub, ok := fields["published"].(bool); ok && pub {
		if _, hasPublishedAt := fields["publishedAt"]; !hasPublishedAt {
			fields["publishedAt"] = time.Now()
		}
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var post models.CMSBlogPost
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	err := database.MongoDB.Collection("blogs").FindOneAndUpdate(ctx,
		bson.M{"_id": id, "siteId": siteID},
		bson.M{"$set": fields}, opts).Decode(&post)
	if err != nil {
		utils.JSONErr(w, "failed to update", http.StatusInternalServerError)
		return
	}
	utils.JSON200(w, post)
}

// AdminDeleteBlog — deletes blog for specific site
func AdminDeleteBlog(w http.ResponseWriter, r *http.Request) {
	if !utils.RequireMongo(w) {
		return
	}
	siteID, ok := utils.RequireSiteID(w, r)
	if !ok {
		return
	}
	rawID := utils.ExtractPathSegment(r.URL.Path, "/api/admin/blogs/", 0)
	id, ok := utils.ParseObjectID(w, rawID)
	if !ok {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if _, err := database.MongoDB.Collection("blogs").DeleteOne(ctx, bson.M{"_id": id, "siteId": siteID}); err != nil {
		utils.JSONErr(w, "failed to delete", http.StatusInternalServerError)
		return
	}
	utils.JSON200(w, map[string]string{"message": "deleted"})
}

// PublicListBlogs — lists published blogs, optionally filtered by siteId
func PublicListBlogs(w http.ResponseWriter, r *http.Request) {
	if !utils.RequireMongo(w) {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"published": true}
	if siteID := strings.TrimSpace(r.URL.Query().Get("siteId")); siteID != "" {
		filter["siteId"] = siteID
	}

	cursor, err := database.MongoDB.Collection("blogs").Find(ctx, filter,
		options.Find().SetSort(bson.D{{Key: "createdAt", Value: -1}}))
	if err != nil {
		utils.JSONErr(w, "failed to fetch", http.StatusInternalServerError)
		return
	}
	var posts []models.CMSBlogPost
	cursor.All(ctx, &posts)
	if posts == nil {
		posts = []models.CMSBlogPost{}
	}
	utils.JSON200(w, posts)
}

// PublicGetBlog — gets single published blog by slug, optionally filtered by siteId
func PublicGetBlog(w http.ResponseWriter, r *http.Request) {
	if !utils.RequireMongo(w) {
		return
	}
	slug := utils.ExtractPathSegment(r.URL.Path, "/api/blogs/", 0)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"slug": slug, "published": true}
	if siteID := strings.TrimSpace(r.URL.Query().Get("siteId")); siteID != "" {
		filter["siteId"] = siteID
	}

	var post models.CMSBlogPost
	err := database.MongoDB.Collection("blogs").FindOne(ctx, filter).Decode(&post)
	if err != nil {
		utils.JSONErr(w, "not found", http.StatusNotFound)
		return
	}
	utils.JSON200(w, post)
}
