package models

import (
	"blog/internal/db"
	"fmt"
	"regexp"
	"strings"
	"time"
)

type Post struct {
	ID        int
	Title     string
	Slug      string
	Excerpt   string
	Content   string
	Author    string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (p *Post) FormattedDate() string {
	return p.CreatedAt.Format("January 2, 2006")
}

func (p *Post) ShortExcerpt() string {
	if p.Excerpt != "" {
		return p.Excerpt
	}
	// Auto-generate from content
	text := strings.ReplaceAll(p.Content, "\n", " ")
	if len(text) > 160 {
		return text[:160] + "..."
	}
	return text
}

// GenerateSlug creates a URL-safe slug from a title
func GenerateSlug(title string) string {
	slug := strings.ToLower(title)
	re := regexp.MustCompile(`[^a-z0-9\s-]`)
	slug = re.ReplaceAllString(slug, "")
	re2 := regexp.MustCompile(`[\s-]+`)
	slug = re2.ReplaceAllString(slug, "-")
	slug = strings.Trim(slug, "-")
	return slug
}

// EnsureUniqueSlug appends a number if the slug already exists
func EnsureUniqueSlug(slug string, excludeID int) string {
	original := slug
	counter := 1
	for {
		var count int
		var err error
		if excludeID > 0 {
			err = db.DB.QueryRow(`SELECT COUNT(*) FROM posts WHERE slug=$1 AND id != $2`, slug, excludeID).Scan(&count)
		} else {
			err = db.DB.QueryRow(`SELECT COUNT(*) FROM posts WHERE slug=$1`, slug).Scan(&count)
		}
		if err != nil || count == 0 {
			break
		}
		slug = fmt.Sprintf("%s-%d", original, counter)
		counter++
	}
	return slug
}

// GetAllPosts returns all posts ordered by newest first
func GetAllPosts() ([]Post, error) {
	rows, err := db.DB.Query(`
		SELECT id, title, slug, COALESCE(excerpt,''), content, author, created_at, updated_at
		FROM posts ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var p Post
		if err := rows.Scan(&p.ID, &p.Title, &p.Slug, &p.Excerpt, &p.Content, &p.Author, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}
	return posts, nil
}

// GetPostBySlug fetches a single post by its slug
func GetPostBySlug(slug string) (*Post, error) {
	var p Post
	err := db.DB.QueryRow(`
		SELECT id, title, slug, COALESCE(excerpt,''), content, author, created_at, updated_at
		FROM posts WHERE slug=$1
	`, slug).Scan(&p.ID, &p.Title, &p.Slug, &p.Excerpt, &p.Content, &p.Author, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

// GetPostByID fetches a single post by its ID
func GetPostByID(id int) (*Post, error) {
	var p Post
	err := db.DB.QueryRow(`
		SELECT id, title, slug, COALESCE(excerpt,''), content, author, created_at, updated_at
		FROM posts WHERE id=$1
	`, id).Scan(&p.ID, &p.Title, &p.Slug, &p.Excerpt, &p.Content, &p.Author, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

// CreatePost inserts a new post into the database
func CreatePost(title, excerpt, content, author string) (*Post, error) {
	slug := GenerateSlug(title)
	slug = EnsureUniqueSlug(slug, 0)

	var p Post
	err := db.DB.QueryRow(`
		INSERT INTO posts (title, slug, excerpt, content, author)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, title, slug, COALESCE(excerpt,''), content, author, created_at, updated_at
	`, title, slug, excerpt, content, author).Scan(
		&p.ID, &p.Title, &p.Slug, &p.Excerpt, &p.Content, &p.Author, &p.CreatedAt, &p.UpdatedAt,
	)
	return &p, err
}

// UpdatePost updates an existing post
func UpdatePost(id int, title, excerpt, content, author string) (*Post, error) {
	slug := GenerateSlug(title)
	slug = EnsureUniqueSlug(slug, id)

	var p Post
	err := db.DB.QueryRow(`
		UPDATE posts
		SET title=$1, slug=$2, excerpt=$3, content=$4, author=$5, updated_at=NOW()
		WHERE id=$6
		RETURNING id, title, slug, COALESCE(excerpt,''), content, author, created_at, updated_at
	`, title, slug, excerpt, content, author, id).Scan(
		&p.ID, &p.Title, &p.Slug, &p.Excerpt, &p.Content, &p.Author, &p.CreatedAt, &p.UpdatedAt,
	)
	return &p, err
}

// DeletePost removes a post by ID
func DeletePost(id int) error {
	_, err := db.DB.Exec(`DELETE FROM posts WHERE id=$1`, id)
	return err
}
