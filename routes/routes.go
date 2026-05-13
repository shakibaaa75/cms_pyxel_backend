package routes

import (
	"net/http"
	"strings"

	"cms-backend/handlers"
	"cms-backend/middleware"
	"cms-backend/utils"
)

func Router() http.Handler {
	mux := http.NewServeMux()

	base := []func(http.HandlerFunc) http.HandlerFunc{
		middleware.SecurityHeadersMiddleware,
		middleware.CORSMiddleware,
		middleware.RateLimitMiddleware,
		middleware.RequestSizeMiddleware(1 << 20),
	}

	adminMiddleware := append(base, middleware.RequireAdminAuth)

	// ── Auth ──────────────────────────────────────────────────────────
	mux.HandleFunc("/api/admin/login", middleware.Chain(handlers.AdminLogin, base...))

	// ── Client Projects ───────────────────────────────────────────────
	mux.HandleFunc("/api/admin/projects", middleware.Chain(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handlers.ListProjects(w, r)
		case http.MethodPost:
			handlers.CreateProject(w, r)
		default:
			utils.JSONErr(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	}, adminMiddleware...))

	mux.HandleFunc("/api/admin/projects/", middleware.Chain(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/updates") && r.Method == http.MethodPost {
			handlers.AddUpdate(w, r)
			return
		}
		switch r.Method {
		case http.MethodPut:
			handlers.UpdateProject(w, r)
		case http.MethodDelete:
			handlers.DeleteProject(w, r)
		default:
			utils.JSONErr(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	}, adminMiddleware...))

	// Public project lookup (no auth, no siteId — token/code is globally unique)
	mux.HandleFunc("/api/project", middleware.Chain(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			utils.JSONErr(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		if r.URL.Query().Has("t") {
			handlers.GetProjectByToken(w, r)
		} else {
			handlers.GetProjectByCode(w, r)
		}
	}, base...))

	// ── Blogs ─────────────────────────────────────────────────────────
	mux.HandleFunc("/api/admin/blogs", middleware.Chain(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handlers.AdminListBlogs(w, r)
		case http.MethodPost:
			handlers.AdminCreateBlog(w, r)
		default:
			utils.JSONErr(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	}, adminMiddleware...))

	mux.HandleFunc("/api/admin/blogs/", middleware.Chain(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPut:
			handlers.AdminUpdateBlog(w, r)
		case http.MethodDelete:
			handlers.AdminDeleteBlog(w, r)
		default:
			utils.JSONErr(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	}, adminMiddleware...))

	mux.HandleFunc("/api/blogs", middleware.Chain(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			utils.JSONErr(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		handlers.PublicListBlogs(w, r)
	}, base...))

	mux.HandleFunc("/api/blogs/", middleware.Chain(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			utils.JSONErr(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		handlers.PublicGetBlog(w, r)
	}, base...))

	return mux
}
