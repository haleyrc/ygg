package share

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/haleyrc/ygg/cloudflare/r2"
	"github.com/haleyrc/ygg/server"
)

func NewSite() []server.Route {
	return []server.Route{
		{Methods: []string{"GET"}, Path: "", Func: index},
		{Methods: []string{"POST"}, Path: "/upload", Func: upload()},
	}
}

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, `<!doctype html>
	<html>
		<head>
			<title>Ygg - Share</title>
		</head>
		<body>
			<h1>Share a file</h1>
			<form method="POST" action="/public/share/upload" enctype="multipart/form-data">
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

func upload() http.HandlerFunc {
	client, err := r2.New(context.Background(), r2.Credentials{
		AccountID:       os.Getenv("CLOUDFLARE_ACCOUNT_ID"),
		AccessKeyID:     os.Getenv("CLOUDFLARE_ACCESS_KEY_ID"),
		AccessKeySecret: os.Getenv("CLOUDFLARE_ACCESS_KEY_SECRET"),
	})
	if err != nil {
		panic(err)
	}
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		start := time.Now()

		log.Println("-->", r.URL.Path)

		f, header, err := r.FormFile("file")
		if err != nil {
			log.Println("\t", err)
			http.Redirect(w, r, "/public/share", http.StatusMovedPermanently)
			return
		}

		filename := fmt.Sprintf("%s%s", uuid.New(), filepath.Ext(header.Filename))
		if err := client.PutFile(ctx, "share", filename, f); err != nil {
			log.Println("\t", err)
		}

		url := fmt.Sprintf("https://share.ryanchaley.com/%s", filename)
		log.Println("<--", time.Since(start), url)

		http.Redirect(w, r, "/public/share", http.StatusMovedPermanently)
	}
}
