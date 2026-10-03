// Package uistatic serves the embedded Admin UI (Svelte/Vite build output).
//
// The UI is embedded from the repository ui/dist directory at compile time,
// so the janus binary is a single deployable artifact. Build it first:
//
//	pnpm --filter ui build
//
// A lone .gitkeep keeps the embed valid on checkouts without a UI build;
// in that case Handler reports ok=false and the Admin server stays API-only.
package uistatic

import (
	"embed"
	"io/fs"
	"mime"
	"net/http"
	"path"
	"strings"
)

//go:embed all:dist
var embedded embed.FS

// Handler returns a file server for the embedded UI. ok is false when
// ui/dist was never built (only .gitkeep embedded).
func Handler() (handler http.Handler, ok bool) {
	sub, err := fs.Sub(embedded, "dist")
	if err != nil {
		return nil, false
	}
	index, err := fs.ReadFile(sub, "index.html")
	if err != nil {
		return nil, false
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet && r.Method != http.MethodHead {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		name := path.Clean("/" + strings.TrimPrefix(r.URL.Path, "/"))
		// Never leak dotfiles (.gitkeep, .well-known internals, ...).
		for seg := range strings.SplitSeq(name, "/") {
			if strings.HasPrefix(seg, ".") && seg != "." && seg != ".." {
				http.NotFound(w, r)
				return
			}
		}
		data := index
		ext := ".html"
		if name != "/" && name != "/index.html" {
			rel := strings.TrimPrefix(name, "/")
			if f, err := fs.ReadFile(sub, rel); err == nil {
				if isFile(sub, rel) {
					data = f
					ext = path.Ext(rel)
				}
			}
		}
		// SPA fallback: unknown paths boot index.html (client-side routing).
		w.Header().Set("Content-Type", mime.TypeByExtension(ext))
		if strings.HasPrefix(name, "/assets/") {
			// Vite content-hashes asset filenames: immutable forever.
			w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
		} else {
			w.Header().Set("Cache-Control", "no-store")
		}
		_, _ = w.Write(data)
	}), true
}

func isFile(fsys fs.FS, name string) bool {
	fi, err := fs.Stat(fsys, name)
	return err == nil && !fi.IsDir()
}
