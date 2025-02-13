package server

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strings"
	"sync"

	"github.com/go-chi/chi/v5"
	"github.com/goodieshq/goseek/utils"
	"github.com/rs/zerolog/log"
)

type GoSeek struct {
	configPath    string
	config        *GoSeekConfig
	browser       *Browser
	mu            sync.RWMutex
	apikeyChecker ApiKeyChecker
}

// create a new GoSeek server instance
func NewGoSeek(configPath string) (*GoSeek, error) {
	var config GoSeekConfig
	if err := config.LoadConfig(configPath); err != nil {
		return nil, err
	}

	browser, err := NewBrowser(config.Root)
	if err != nil {
		return nil, err
	}

	apikeyChecker := NewApiKeyCheckStatic(config.ApiKeys)

	return &GoSeek{
		configPath:    configPath,
		config:        &config,
		browser:       browser,
		apikeyChecker: apikeyChecker,
	}, nil
}

func (gs *GoSeek) ReloadApiKeys() error {
	var config GoSeekConfig
	if err := config.LoadConfig(gs.configPath); err != nil {
		return err
	}

	gs.mu.Lock()
	defer gs.mu.Unlock()

	gs.config = &config
	gs.apikeyChecker.UpdateApiKeys(config.ApiKeys)

	log.Info().Int("apikey_count", len(gs.config.ApiKeys)).Msg("new API keys loaded")

	return nil
}

// run GoSeek on the provided port
func (gs *GoSeek) Run() {
	// listen on the provided port
	laddr := fmt.Sprintf(":%d", gs.config.Port)

	// create the chi router
	router := chi.NewRouter()

	// check the provided API keys
	router.Use(MiddlewareAPIKeys(gs.apikeyChecker))
	router.HandleFunc("/*", gs.handle)

	http.ListenAndServe(laddr, router)
}

func buildResponseDir(path, apikey string, items []Item) []byte {
	// string builder for crafting the HTML output
	var builder strings.Builder

	// Show the directory title
	builder.WriteString("<h1>Contents of:  " + utils.Format(path, true) + "</h1>\n<br>\n")
	log.Info().Str("path", path).Str("apikey", utils.ApiKeyPrefix(apikey)).Msg("directory viewed")

	// navigate up in a directory if it is not the root
	if path == "." {
		builder.WriteString("[Root Directory]\n<br>\n")
	} else {
		builder.WriteString(utils.Href(filepath.Dir(path), "[Parent Directory]", false, apikey))
	}
	builder.WriteString("<br>\n")

	// create anchor tags for each file/directory entry
	for _, item := range items {
		builder.WriteString(utils.Href(filepath.Join(path, item.Name), item.Name, item.Kind == KindDir, apikey))
	}

	return []byte(builder.String())
}

func (gs *GoSeek) handle(w http.ResponseWriter, r *http.Request) {
	// get the API key from the query
	apikey := r.URL.Query().Get("apikey")

	// take the file path and strip slashes
	path := strings.Trim(r.URL.Path, "/\\")
	if path == "" {
		// root directory
		path = "."
	}

	// make sure the path is valid within the filesystem browser
	valid, kind := gs.browser.ValidatePath(path)
	if !valid {
		// invalid file, 404
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	switch kind {
	case KindDir:
		// `path` is a directory
		items, err := gs.browser.list(path, true)
		if err != nil {
			log.Error().Err(err).Str("path", path).Msg("error listing directory contents")
			http.Error(w, "service unavailable", http.StatusServiceUnavailable)
			return
		}

		response := buildResponseDir(path, apikey, items)

		// Ok, send the HTML
		w.WriteHeader(http.StatusOK)
		w.Write(response)
		return
	case KindFile:
		// get the base filename for the download name
		basename := filepath.Base(path)

		// get stat info of the target file
		info, _ := gs.browser.root.Stat(path)

		// open the target file
		file, err := gs.browser.root.Open(path)
		if err != nil {
			log.Error().Err(err).Str("path", path).Msg("failed to open the file path")
			http.Error(w, "service unavailable", http.StatusServiceUnavailable)
			return
		}
		defer file.Close()

		log.Info().Str("path", path).Str("apikey", utils.ApiKeyPrefix(apikey)).Msg("file downloaded")

		// add headers to facilitate file downloading
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", basename))
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("Content-Length", fmt.Sprintf("%d", info.Size()))

		// send the file for downloading
		http.ServeContent(w, r, basename, info.ModTime(), file)
		return
	}
}
