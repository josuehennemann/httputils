package httputils

import (
	"errors"
	"net/http"
	"os"
)

var domainHttps string

func SetDomainHttps(d string) {
	domainHttps = d
}
func FileSystem(fs http.FileSystem) ListenOnlyFilesFilesystem {
	return ListenOnlyFilesFilesystem{fs: fs}
}
func RedirectTLS(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "https://"+domainHttps+r.RequestURI, http.StatusMovedPermanently)
}

type ListenOnlyFilesFilesystem struct {
	fs http.FileSystem
}

func (fs ListenOnlyFilesFilesystem) Open(name string) (http.File, error) {
	f, err := fs.fs.Open(name)
	if err != nil {
		return nil, err
	}

	if isdir, _ := f.Stat(); isdir == nil || isdir.IsDir() {
		return nil, errors.New("Is dir: " + name)
	}

	return myReaddirFile{f}, nil
}

type myReaddirFile struct {
	http.File
}

func (f myReaddirFile) Readdir(count int) ([]os.FileInfo, error) {
	return nil, nil
}
