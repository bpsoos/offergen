package endpoint

import (
	"embed"
	"io/fs"
)

//go:embed static/media/*
var media embed.FS

//go:embed static/styles/*
var styles embed.FS

func (e *Handler) Media() fs.FS {
	media, err := fs.Sub(media, "static/media")
	if err != nil {
		panic("missing subdir")
	}
	return media
}
func (e *Handler) Styles() fs.FS {
	styles, err := fs.Sub(styles, "static/styles")
	if err != nil {
		panic("missing subdir")
	}
	return styles
}
