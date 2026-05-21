package handlers

import (
	"blog/internal/models"
	"blog/internal/templates"
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

// Home — GET /
func Home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	posts, err := models.GetAllPosts()
	if err != nil {
		http.Error(w, "Failed to load posts", http.StatusInternalServerError)
		return
	}
	templates.Render(w, "home.html", map[string]any{
		"Title": "My Blog",
		"Posts": posts,
	})
}

// ShowPost — GET /posts/{slug}
func ShowPost(w http.ResponseWriter, r *http.Request) {
	slug := strings.TrimPrefix(r.URL.Path, "/posts/")
	if slug == "" {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	post, err := models.GetPostBySlug(slug)
	if err == sql.ErrNoRows {
		http.NotFound(w, r)
		return
	}
	if err != nil {
		http.Error(w, "Failed to load post", http.StatusInternalServerError)
		return
	}

	templates.Render(w, "post.html", map[string]any{
		"Title": post.Title,
		"Post":  post,
	})
}

// NewPostForm — GET /posts/new
func NewPostForm(w http.ResponseWriter, r *http.Request) {
	templates.Render(w, "new_post.html", map[string]any{
		"Title": "New Post",
	})
}

// CreatePost — POST /posts
func CreatePost(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	title := strings.TrimSpace(r.FormValue("title"))
	excerpt := strings.TrimSpace(r.FormValue("excerpt"))
	content := strings.TrimSpace(r.FormValue("content"))
	author := strings.TrimSpace(r.FormValue("author"))

	// Validation
	var errors []string
	if title == "" {
		errors = append(errors, "Title is required")
	}
	if content == "" {
		errors = append(errors, "Content is required")
	}
	if author == "" {
		author = "Admin"
	}

	if len(errors) > 0 {
		// HTMX partial: return error fragment
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusUnprocessableEntity)
		templates.RenderPartial(w, "partials/form_errors.html", map[string]any{
			"Errors": errors,
		})
		return
	}

	post, err := models.CreatePost(title, excerpt, content, author)
	if err != nil {
		http.Error(w, "Failed to create post", http.StatusInternalServerError)
		return
	}

	// HTMX redirect after success
	w.Header().Set("HX-Redirect", fmt.Sprintf("/posts/%s", post.Slug))
	w.WriteHeader(http.StatusOK)
}

// EditPostForm — GET /posts/{id}/edit
func EditPostForm(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimSuffix(strings.TrimPrefix(r.URL.Path, "/posts/"), "/edit")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	post, err := models.GetPostByID(id)
	if err == sql.ErrNoRows {
		http.NotFound(w, r)
		return
	}
	if err != nil {
		http.Error(w, "Failed to load post", http.StatusInternalServerError)
		return
	}

	templates.Render(w, "edit_post.html", map[string]any{
		"Title": "Edit: " + post.Title,
		"Post":  post,
	})
}

// UpdatePost — POST /posts/{id}/update  (HTMX-friendly, avoids PUT parsing issues)
func UpdatePost(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimSuffix(strings.TrimPrefix(r.URL.Path, "/posts/"), "/update")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	title := strings.TrimSpace(r.FormValue("title"))
	excerpt := strings.TrimSpace(r.FormValue("excerpt"))
	content := strings.TrimSpace(r.FormValue("content"))
	author := strings.TrimSpace(r.FormValue("author"))

	var errors []string
	if title == "" {
		errors = append(errors, "Title is required")
	}
	if content == "" {
		errors = append(errors, "Content is required")
	}

	if len(errors) > 0 {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusUnprocessableEntity)
		templates.RenderPartial(w, "partials/form_errors.html", map[string]any{
			"Errors": errors,
		})
		return
	}

	post, err := models.UpdatePost(id, title, excerpt, content, author)
	if err != nil {
		http.Error(w, "Failed to update post", http.StatusInternalServerError)
		return
	}

	w.Header().Set("HX-Redirect", fmt.Sprintf("/posts/%s", post.Slug))
	w.WriteHeader(http.StatusOK)
}

// DeletePost — POST /posts/{id}/delete
func DeletePost(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimSuffix(strings.TrimPrefix(r.URL.Path, "/posts/"), "/delete")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	if err := models.DeletePost(id); err != nil {
		http.Error(w, "Failed to delete post", http.StatusInternalServerError)
		return
	}

	// HTMX redirect to home
	w.Header().Set("HX-Redirect", "/")
	w.WriteHeader(http.StatusOK)
}

// DeletePostRow — POST /posts/{id}/delete-row (for list HTMX swap)
func DeletePostRow(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimSuffix(strings.TrimPrefix(r.URL.Path, "/posts/"), "/delete-row")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	if err := models.DeletePost(id); err != nil {
		http.Error(w, "Failed to delete post", http.StatusInternalServerError)
		return
	}

	// Return empty — HTMX will remove the element
	w.WriteHeader(http.StatusOK)
}

// SearchPosts — GET /search?q=...
func SearchPosts(w http.ResponseWriter, r *http.Request) {
	q := strings.TrimSpace(r.URL.Query().Get("q"))
	var posts []models.Post
	var err error

	if q != "" {
		posts, err = models.SearchPosts(q)
	} else {
		posts, err = models.GetAllPosts()
	}

	if err != nil {
		http.Error(w, "Search failed", http.StatusInternalServerError)
		return
	}

	templates.RenderPartial(w, "partials/post_list.html", map[string]any{
		"Posts": posts,
		"Query": q,
	})
}
