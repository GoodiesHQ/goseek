package server

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/rs/zerolog/log"
)

type Kind int

const (
	KindNone Kind = iota
	KindFile
	KindDir
)

type Item struct {
	Name string
	Kind Kind
}

type Browser struct {
	root *os.Root
}

func (browser *Browser) list(path string, validated bool) ([]Item, error) {
	if !validated {
		valid, kind := browser.ValidatePath(path)
		if !valid {
			return nil, fmt.Errorf("path is invalid")
		}
		if kind != KindDir {
			return nil, fmt.Errorf("path is not a directory type")
		}
	}

	path = filepath.Join(browser.root.Name(), path)
	entries, err := os.ReadDir(path)
	if err != nil {
		log.Error().Err(err).Msg("failed to read dir")
	}

	items := make([]Item, 0, len(entries))
	for _, entry := range entries {
		var kind Kind

		if entry.IsDir() {
			kind = KindDir
		} else {
			kind = KindFile
		}

		items = append(items, Item{
			Name: entry.Name(),
			Kind: kind,
		})
	}

	return items, nil
}

func (browser *Browser) ValidatePath(path string) (bool, Kind) {
	stat, err := browser.root.Stat(path)
	if err != nil {
		return false, KindNone
	}
	if stat.IsDir() {
		return true, KindDir
	} else {
		return true, KindFile
	}
}

func (browser *Browser) IsValidPath(path string) bool {
	valid, _ := browser.ValidatePath(path)
	return valid
}

func NewBrowser(rootPath string) (*Browser, error) {
	root, err := os.OpenRoot(rootPath)
	if err != nil {
		return nil, err
	}
	return &Browser{root: root}, nil
}
