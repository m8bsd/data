# ✒️ Inkwell — Go + HTMX + Bootstrap Blog

A full-featured blog with CRUD operations built with:
- **Backend**: Go (stdlib `net/http`, `html/template`, `database/sql`)
- **Frontend**: HTMX 1.9 + Bootstrap 5.3
- **Database**: Neon (serverless PostgreSQL)
- **Driver**: `github.com/lib/pq`

---

## 🚀 Quick Start

### 1. Prerequisites

- Go 1.22+
- A [Neon](https://neon.tech) account (free tier works great)

### 2. Clone & configure

```bash
git clone <your-repo>
cd blog

# Copy and fill in your Neon connection string
cp .env.example .env
```

Edit `.env`:
```
DATABASE_URL="postgres://user:pass@ep-xxx.region.aws.neon.tech/dbname?sslmode=require"
```

Get your connection string from: **Neon Console → Your Project → Connection Details → Connection string**

### 3. Install dependencies

```bash
go mod tidy
```

### 4. Run

```bash
make run
# or
go run ./cmd/main.go
```

Open http://localhost:8080

---

## 🗂️ Project Structure

```
blog/
├── cmd/
│   └── main.go                  # Entry point, router
├── internal/
│   ├── db/
│   │   └── db.go                # Neon connection + schema migration
│   ├── handlers/
│   │   └── handlers.go          # HTTP request handlers
│   ├── models/
│   │   ├── post.go              # Post struct + CRUD queries
│   │   └── search.go            # Search query
│   └── templates/
│       ├── renderer.go          # Template engine wrapper
│       ├── layout.html          # Base layout (navbar, footer)
│       ├── home.html            # Post list + live search
│       ├── post.html            # Single post view
│       ├── new_post.html        # Create form
│       ├── edit_post.html       # Edit form
│       └── partials/
│           ├── post_list.html   # Card grid (HTMX target)
│           ├── post_form.html   # Shared form fields
│           └── form_errors.html # Validation errors fragment
├── static/
│   ├── css/app.css              # Custom styles
│   └── js/app.js                # HTMX enhancements
├── .env.example
├── .air.toml                    # Hot-reload config
├── Makefile
└── go.mod
```

---

## ✨ Features

| Feature | Implementation |
|---|---|
| List all posts | `GET /` |
| View a post | `GET /posts/{slug}` |
| Create post | `GET /posts/new` → `POST /posts` |
| Edit post | `GET /posts/{id}/edit` → `POST /posts/{id}/update` |
| Delete (redirect) | `POST /posts/{id}/delete` |
| Delete (inline) | `POST /posts/{id}/delete-row` (HTMX swap) |
| Live search | `GET /search?q=...` (HTMX partial swap) |
| Auto slug | Generated from title, deduplicated |
| DB schema | Auto-created on first run |

---

## 🔥 Hot Reload (optional)

```bash
make install-air
make dev
```

---

## 🌐 Deploy

### Render / Railway / Fly.io
Set the `DATABASE_URL` environment variable in your platform's dashboard, then:

```bash
make build
# Upload/push the binary
```

### Docker (optional)
```dockerfile
FROM golang:1.22-alpine AS build
WORKDIR /app
COPY . .
RUN go mod tidy && go build -o inkwell ./cmd/main.go

FROM alpine
WORKDIR /app
COPY --from=build /app/inkwell .
COPY --from=build /app/internal/templates ./internal/templates
COPY --from=build /app/static ./static
EXPOSE 8080
CMD ["./inkwell"]
```
