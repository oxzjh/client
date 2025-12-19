package client

import (
	"io"
	"net/http"
	"net/url"
)

type IClient interface {
	Read() ([]byte, error)
	Write([]byte) error
	WriteJson(v any) error
	Close() error
}

type HTTP interface {
	SetHeaders(map[string]string)
	Get(string) (*http.Response, error)
	Download(string, io.Writer) error
	DownloadProgress(string, io.WriteCloser, func(int64, int64)) error
	Post(string, string, io.Reader) (*http.Response, error)
	PostJson(string, any) (*http.Response, error)
	PostForm(string, url.Values) (*http.Response, error)
	Upload(string, io.Reader) (*http.Response, error)
	UploadProgress(string, io.Reader, int64, func(int64)) (*http.Response, error)
	UploadForm(string, string, string, io.Reader, map[string]string) (*http.Response, error)
}
