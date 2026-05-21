package models

import "blog/internal/db"

// SearchPosts performs a full-text search on title and content
func SearchPosts(query string) ([]Post, error) {
	rows, err := db.DB.Query(`
		SELECT id, title, slug, COALESCE(excerpt,''), content, author, created_at, updated_at
		FROM posts
		WHERE title ILIKE '%' || $1 || '%'
		   OR content ILIKE '%' || $1 || '%'
		   OR excerpt ILIKE '%' || $1 || '%'
		ORDER BY created_at DESC
	`, query)
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
