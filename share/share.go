package share

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"path/filepath"
	"time"

	"github.com/google/uuid"
)

type FileSaver interface {
	PutFile(ctx context.Context, dir, path string, r io.Reader) error
}

func NewController(fs FileSaver) (*Controller, error) {
	return &Controller{fs: fs}, nil
}

type Controller struct {
	fs FileSaver
}

func (c *Controller) Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, `<!doctype html>
	<html>
		<head>
			<title>Ygg - Share</title>
		</head>
		<body>
			<h1>Share a file</h1>
			<form method="POST" enctype="multipart/form-data">
				<label>
					File:
					<input type="file" id="file" name="file" />
				</label>
				<button type="submit">Share</button>
			</form>
		</body>
	</html>
	`)
}

func (c *Controller) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	start := time.Now()

	referrer := r.Referer()

	log.Println("-->", r.URL.Path)

	f, header, err := r.FormFile("file")
	if err != nil {
		log.Println("\t", err)
		http.Redirect(w, r, referrer, http.StatusMovedPermanently)
		return
	}

	filename := fmt.Sprintf("%s%s", uuid.New(), filepath.Ext(header.Filename))
	if err := c.fs.PutFile(ctx, "share", filename, f); err != nil {
		log.Println("\t", err)
	}

	url := fmt.Sprintf("https://share.ryanchaley.com/%s", filename)
	log.Println("<--", time.Since(start), url)

	http.Redirect(w, r, referrer, http.StatusMovedPermanently)
}
