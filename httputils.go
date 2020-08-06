package httputils

import (
	"errors"
	"net/http"
	"os"

	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"sort"
	"strings"
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


// original package CurlCommand = https://github.com/moul/http2curl.git


// CurlCommand contains exec.Command compatible slice + helpers
type CurlCommand []string

// append appends a string to the CurlCommand
func (c *CurlCommand) append(newSlice ...string) {
	*c = append(*c, newSlice...)
}

// String returns a ready to copy/paste command
func (c *CurlCommand) String() string {
	return strings.Join(*c, " ")
}

// nopCloser is used to create a new io.ReadCloser for req.Body
type nopCloser struct {
	io.Reader
}

func bashEscape(str string) string {
	return `'` + strings.Replace(str, `'`, `'\''`, -1) + `'`
}

func (nopCloser) Close() error { return nil }

// GetCurlCommand returns a CurlCommand corresponding to an http.Request
func GetCurlCommand(req *http.Request) (*CurlCommand, error) {
	command := CurlCommand{}

	command.append("curl")

	command.append("-X", bashEscape(req.Method))

	if req.Body != nil {
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			return nil, err
		}
		req.Body = nopCloser{bytes.NewBuffer(body)}
		if len(string(body)) > 0 {
			bodyEscaped := bashEscape(string(body))
			command.append("-d", bodyEscaped)
		}
	}

	var keys []string

	for k := range req.Header {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		command.append("-H", bashEscape(fmt.Sprintf("%s: %s", k, strings.Join(req.Header[k], " "))))
	}

	command.append(bashEscape(req.URL.String()))

	return &command, nil
}