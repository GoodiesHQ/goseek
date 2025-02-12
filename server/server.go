package server

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/goodieshq/goseek/utils"
	"github.com/rs/zerolog/log"
)

type GoSeek struct {
	browser     *Browser
	authchecker AuthChecker
}

func NewGoSeek(path string, authchecker AuthChecker) (*GoSeek, error) {
	browser, err := NewBrowser(path)
	if err != nil {
		return nil, err
	}
	return &GoSeek{
		browser:     browser,
		authchecker: authchecker,
	}, nil
}

func (gs *GoSeek) Run(port uint16) {
	laddr := fmt.Sprintf(":%d", port)

	router := chi.NewRouter()
	router.Use(AuthMiddleware(gs.authchecker))
	router.HandleFunc("/*", gs.handle)

	http.ListenAndServe(laddr, router)
}

func (gs *GoSeek) handle(w http.ResponseWriter, r *http.Request) {
	apikey := r.URL.Query().Get("apikey")

	path := strings.Trim(r.URL.Path, "/\\")
	if path == "" {
		path = "."
	}

	valid, kind := gs.browser.ValidatePath(path)
	if !valid {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("not found"))
		return
	}

	if kind == KindDir {
		items, err := gs.browser.list(path, true)
		if err != nil {
			log.Error().Err(err).Str("path", path).Msg("error listing directory contents")
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte("server error"))
			return
		}

		var builder strings.Builder
		builder.WriteString("<h1>Contents of:  " + utils.Format(path, true) + "</h1>\n<br>\n")
		log.Info().Str("path", path).Msg("path loaded")

		if path == "." {
			builder.WriteString("[Root Directory]")
		} else {
			builder.WriteString(utils.Href(filepath.Dir(path), "[Parent Directory]", false, apikey))
		}
		builder.WriteString("<br><br>\n")

		for _, item := range items {
			builder.WriteString(utils.Href(filepath.Join(path, item.Name), item.Name, item.Kind == KindDir, apikey))
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(builder.String()))
	} else {
		basename := filepath.Base(path)
		info, _ := gs.browser.root.Stat(path)
		file, err := gs.browser.root.Open(path)
		if err != nil {
			log.Error().Err(err).Str("path", path).Msg("failed to open the file path")
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		defer file.Close()

		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", basename))
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("Content-Length", fmt.Sprintf("%d", info.Size()))

		http.ServeContent(w, r, basename, info.ModTime(), file)
	}
}
