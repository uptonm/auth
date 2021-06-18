package common

import (
	"embed"
	"io/fs"
	"net/http"
	"net/url"
	"strings"

	log "github.com/sirupsen/logrus"
)

// StrictFs is a Custom strict filesystem implementation to
// prevent directory listings for resources.
type StrictFs struct {
	Fs http.FileSystem
}

// Open only allows existing files to be pulled, not directories
func (sfs StrictFs) Open(path string) (http.File, error) {
	// url decode path to support encoded characters
	path, err := url.QueryUnescape(path)
	if err != nil {
		log.Printf("StrictFS error: %s, %s", path, err.Error())
		return nil, err
	}

	// trim trailing slashes to avoid invalid path errors
	// in fiber's filesystem middleware
	if path != "/" {
		path = strings.TrimSuffix(path, "/")
	}

	// open file directly if it exists
	f, err := sfs.Fs.Open(path)
	if err != nil {
		return nil, err
	}

	// prevent directory listings, only show index file if any
	s, err := f.Stat()
	if err == nil && s.IsDir() {
		index := strings.TrimSuffix(path, "/") + "/auth.html"
		if _, err := sfs.Fs.Open(index); err != nil {
			return nil, err
		}
	}
	return f, nil
}

// PickFS returns either an embedded FS or an on-disk FS for the
// given directory path.
func PickFS(useDisk bool, e embed.FS, dir string) http.FileSystem {
	if useDisk {
		return http.Dir(dir)
	}

	efs, err := fs.Sub(e, strings.Trim(dir, "./"))
	if err != nil {
		panic(err)
	}

	return http.FS(efs)
}
