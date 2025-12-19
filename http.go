package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
)

type progressReader struct {
	reader   io.Reader
	progress func(int64)
	loaded   int64
}

func (pr *progressReader) Read(p []byte) (n int, err error) {
	if n, err = pr.reader.Read(p); n > 0 {
		pr.loaded += int64(n)
		pr.progress(pr.loaded)
	}
	return
}

type httpClient struct {
	baseURL string
	headers map[string]string
	client  *http.Client
}

func (hc *httpClient) SetHeaders(headers map[string]string) {
	hc.headers = headers
}

func (hc *httpClient) Get(route string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, hc.baseURL+route, nil)
	if err != nil {
		return nil, err
	}
	return hc.send(req)
}

func (hc *httpClient) Download(route string, w io.Writer) error {
	res, err := hc.Get(route)
	if err != nil {
		return err
	}
	_, err = io.Copy(w, res.Body)
	res.Body.Close()
	return err
}

func (hc *httpClient) DownloadProgress(route string, w io.WriteCloser, onProgress func(int64, int64)) error {
	res, err := hc.Get(route)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	var loaded int64
	buf := make([]byte, 32<<10)
	for {
		n, err := res.Body.Read(buf)
		if n > 0 {
			w.Write(buf[:n])
			loaded += int64(n)
			onProgress(loaded, res.ContentLength)
		}
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
	}
}

func (hc *httpClient) Post(route, contentType string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodPost, hc.baseURL+route, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", contentType)
	return hc.send(req)
}

func (hc *httpClient) PostJson(route string, data any) (*http.Response, error) {
	b, _ := json.Marshal(data)
	return hc.Post(route, "application/json", bytes.NewReader(b))
}

func (hc *httpClient) PostForm(route string, data url.Values) (*http.Response, error) {
	return hc.Post(route, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
}

func (hc *httpClient) Upload(route string, body io.Reader) (*http.Response, error) {
	return hc.Post(route, "application/octet-stream", body)
}

func (hc *httpClient) UploadProgress(route string, r io.Reader, total int64, onProgress func(int64)) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodPost, hc.baseURL+route, &progressReader{
		reader:   r,
		progress: onProgress,
	})
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/octet-stream")
	req.ContentLength = total
	return hc.send(req)
}

func (hc *httpClient) UploadForm(route, fieldname, filename string, file io.Reader, fields map[string]string) (*http.Response, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(fieldname, filename)
	if err != nil {
		return nil, err
	}
	if _, err = io.Copy(part, file); err != nil {
		return nil, err
	}
	for k, v := range fields {
		writer.WriteField(k, v)
	}
	writer.Close()
	return hc.Post(route, writer.FormDataContentType(), body)
}

func (hc *httpClient) send(req *http.Request) (*http.Response, error) {
	if hc.headers != nil {
		for k, v := range hc.headers {
			req.Header.Set(k, v)
		}
	}
	res, err := hc.client.Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(res.Body)
		res.Body.Close()
		return nil, fmt.Errorf("status code: %d, body: %s", res.StatusCode, b)
	}
	return res, nil
}

func NewHTTP(baseURL string, headers map[string]string, proxy string) (HTTP, error) {
	client := &http.Client{}
	if proxy != "" {
		u, err := url.Parse(proxy)
		if err != nil {
			return nil, err
		}
		client.Transport = &http.Transport{Proxy: http.ProxyURL(u)}
	}
	return &httpClient{baseURL, headers, client}, nil
}

func ParseResponse(res *http.Response, v any) error {
	err := json.NewDecoder(res.Body).Decode(v)
	res.Body.Close()
	return err
}
