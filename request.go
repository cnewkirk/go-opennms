package opennms

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// request describes one HTTP call. contentType "" means the default
// application/json (mirroring the Python wrapper's session header,
// which is live-validated even on bodyless requests); set
// omitContentType for endpoints that reject it (e.g. GraphML).
type request struct {
	method          string
	path            string
	v2              bool
	params          url.Values
	body            []byte
	contentType     string
	omitContentType bool
	accept          string
}

// endpoint builds a full URL, using the v2 base if v2 is true.
func (c *Client) endpoint(path string, v2 bool) string {
	base := c.v1URL
	if v2 {
		base = c.v2URL
	}
	return base + "/" + strings.TrimLeft(path, "/")
}

// retryable reports whether the attempt should be retried, mirroring
// urllib3's Retry(status_forcelist=(500, 502, 503, 504)).
func retryable(resp *http.Response, err error) bool {
	if err != nil {
		return true
	}
	switch resp.StatusCode {
	case 500, 502, 503, 504:
		return true
	}
	return false
}

// do sends the request, retrying on connection errors and HTTP
// 500/502/503/504 with a 0.5s exponential backoff factor.
func (c *Client) do(ctx context.Context, r request) (*http.Response, error) {
	u := c.endpoint(r.path, r.v2)
	if len(r.params) > 0 {
		u += "?" + r.params.Encode()
	}
	accept := r.accept
	if accept == "" {
		accept = defaultAccept
	}
	contentType := r.contentType
	if contentType == "" && !r.omitContentType {
		contentType = "application/json"
	}
	var resp *http.Response
	var err error
	for attempt := 0; ; attempt++ {
		var body io.Reader
		if r.body != nil {
			body = bytes.NewReader(r.body)
		}
		req, rerr := http.NewRequestWithContext(ctx, r.method, u, body)
		if rerr != nil {
			return nil, rerr
		}
		req.SetBasicAuth(c.username, c.password)
		req.Header.Set("Accept", accept)
		if contentType != "" {
			req.Header.Set("Content-Type", contentType)
		}
		resp, err = c.httpClient.Do(req)
		if !retryable(resp, err) || attempt >= c.retries {
			break
		}
		if resp != nil {
			io.Copy(io.Discard, resp.Body)
			resp.Body.Close()
		}
		select {
		case <-time.After(500 * time.Millisecond << attempt):
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
	return resp, err
}

// parse reads and decodes a response: JSON into map[string]any /
// []any, text/plain into int (if numeric) or string, empty bodies
// (204 No Content) into nil. 4xx/5xx statuses become an *APIError
// (3xx such as 304 Not Modified are not errors, mirroring requests).
func parse(resp *http.Response) (any, error) {
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 400 {
		return nil, &APIError{
			StatusCode: resp.StatusCode,
			Status:     resp.Status,
			URL:        resp.Request.URL.String(),
			Body:       string(body),
		}
	}
	if len(body) == 0 {
		return nil, nil
	}
	ct := resp.Header.Get("Content-Type")
	if strings.Contains(ct, "application/json") {
		var v any
		if err := json.Unmarshal(body, &v); err != nil {
			return nil, fmt.Errorf("opennms: decoding response: %w", err)
		}
		return v, nil
	}
	if strings.Contains(ct, "text/plain") {
		text := strings.TrimSpace(string(body))
		if n, err := strconv.Atoi(text); err == nil {
			return n, nil
		}
		return text, nil
	}
	var v any
	if err := json.Unmarshal(body, &v); err == nil {
		return v, nil
	}
	return string(body), nil
}

// send performs the request and parses the response.
func (c *Client) send(ctx context.Context, r request) (any, error) {
	resp, err := c.do(ctx, r)
	if err != nil {
		return nil, err
	}
	return parse(resp)
}

// sendRaw performs the request and returns the raw body, raising an
// *APIError on non-2xx statuses.
func (c *Client) sendRaw(ctx context.Context, r request) ([]byte, error) {
	resp, err := c.do(ctx, r)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 400 {
		return nil, &APIError{
			StatusCode: resp.StatusCode,
			Status:     resp.Status,
			URL:        resp.Request.URL.String(),
			Body:       string(body),
		}
	}
	return body, nil
}

func marshalJSON(v any) ([]byte, error) {
	if v == nil {
		return nil, nil
	}
	return json.Marshal(v)
}

// get sends a GET request and returns the parsed response.
func (c *Client) get(ctx context.Context, path string, params url.Values, v2 bool) (any, error) {
	return c.send(ctx, request{method: "GET", path: path, params: params, v2: v2})
}

// post sends a JSON POST request and returns the parsed response.
func (c *Client) post(ctx context.Context, path string, jsonBody any, params url.Values, v2 bool) (any, error) {
	body, err := marshalJSON(jsonBody)
	if err != nil {
		return nil, err
	}
	return c.send(ctx, request{method: "POST", path: path, body: body, params: params, v2: v2})
}

// postForm sends a form-encoded POST request (the OpenNMS acks
// endpoint rejects JSON on every version).
func (c *Client) postForm(ctx context.Context, path string, form url.Values, params url.Values, v2 bool) (any, error) {
	return c.send(ctx, request{
		method: "POST", path: path, params: params, v2: v2,
		body:        []byte(form.Encode()),
		contentType: "application/x-www-form-urlencoded",
	})
}

// put sends a JSON PUT request and returns the parsed response.
func (c *Client) put(ctx context.Context, path string, jsonBody any, params url.Values, v2 bool) (any, error) {
	body, err := marshalJSON(jsonBody)
	if err != nil {
		return nil, err
	}
	return c.send(ctx, request{method: "PUT", path: path, body: body, params: params, v2: v2})
}

// putForm sends a form-encoded PUT request (alarm/notification acks,
// group and user updates are form-encoded on every OpenNMS version).
func (c *Client) putForm(ctx context.Context, path string, form url.Values, params url.Values, v2 bool) (any, error) {
	return c.send(ctx, request{
		method: "PUT", path: path, params: params, v2: v2,
		body:        []byte(form.Encode()),
		contentType: "application/x-www-form-urlencoded",
	})
}

// del sends a DELETE request. Pass accept for endpoints that cannot
// produce JSON — some (e.g. GraphML) return 500 against the default
// Accept header and the default Content-Type on bodyless requests,
// so both are overridden when accept is given.
func (c *Client) del(ctx context.Context, path string, params url.Values, jsonBody any, v2 bool, accept string) (any, error) {
	body, err := marshalJSON(jsonBody)
	if err != nil {
		return nil, err
	}
	return c.send(ctx, request{
		method: "DELETE", path: path, body: body, params: params, v2: v2,
		accept: accept, omitContentType: accept != "",
	})
}

// patch sends a JSON PATCH request and returns the parsed response.
func (c *Client) patch(ctx context.Context, path string, jsonBody any, params url.Values, v2 bool) (any, error) {
	body, err := marshalJSON(jsonBody)
	if err != nil {
		return nil, err
	}
	return c.send(ctx, request{method: "PATCH", path: path, body: body, params: params, v2: v2})
}

// getText sends a GET request and returns the raw response text.
// Pass accept for endpoints that only produce a non-JSON content
// type (e.g. application/xml).
func (c *Client) getText(ctx context.Context, path string, params url.Values, v2 bool, accept string) (string, error) {
	body, err := c.sendRaw(ctx, request{
		method: "GET", path: path, params: params, v2: v2,
		accept: accept, omitContentType: accept != "",
	})
	if err != nil {
		return "", err
	}
	return string(body), nil
}

// getBytes sends a GET request with Accept: */* and returns the raw
// body — binary endpoints (PNG images, rendered reports) return 406
// against the default JSON Accept header.
func (c *Client) getBytes(ctx context.Context, path string, params url.Values, v2 bool) ([]byte, error) {
	return c.sendRaw(ctx, request{
		method: "GET", path: path, params: params, v2: v2,
		accept: "*/*", omitContentType: true,
	})
}

// postBytes sends a JSON POST request with Accept: */* and returns
// the raw response body.
func (c *Client) postBytes(ctx context.Context, path string, jsonBody any, params url.Values, v2 bool) ([]byte, error) {
	body, err := marshalJSON(jsonBody)
	if err != nil {
		return nil, err
	}
	return c.sendRaw(ctx, request{
		method: "POST", path: path, body: body, params: params, v2: v2,
		accept: "*/*",
	})
}

// fileUpload is one part of a multipart POST.
type fileUpload struct {
	filename string
	content  []byte
}

// postFiles sends a multipart/form-data POST request. Keys are form
// field names.
func (c *Client) postFiles(ctx context.Context, path string, files map[string]fileUpload, params url.Values, v2 bool) (any, error) {
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	for field, f := range files {
		part, err := w.CreateFormFile(field, f.filename)
		if err != nil {
			return nil, err
		}
		if _, err := part.Write(f.content); err != nil {
			return nil, err
		}
	}
	if err := w.Close(); err != nil {
		return nil, err
	}
	return c.send(ctx, request{
		method: "POST", path: path, body: buf.Bytes(),
		params: params, v2: v2, contentType: w.FormDataContentType(),
	})
}

// putText sends a PUT request with a raw text body.
func (c *Client) putText(ctx context.Context, path, data, contentType string, v2 bool) (any, error) {
	return c.send(ctx, request{
		method: "PUT", path: path, body: []byte(data),
		contentType: contentType, v2: v2,
	})
}

// postText sends a POST request with a raw text body. Pass accept
// for endpoints that cannot produce JSON.
func (c *Client) postText(ctx context.Context, path, data, contentType string, v2 bool, accept string, params url.Values) (any, error) {
	return c.send(ctx, request{
		method: "POST", path: path, body: []byte(data),
		contentType: contentType, v2: v2, accept: accept, params: params,
	})
}
