package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	neturl "net/url"
	"os"
	"path/filepath"
	"strings"
)

type CustomJarsList struct {
	Results []CustomJar `json:"results"`
}

type CustomJar struct {
	Name string `json:"name,omitempty"`
	// FilePath is the local path to the .jar file to upload. When set,
	// UploadCustomJar performs a multipart/form-data upload of the file
	// (as required by the real API). When empty, it falls back to a JSON
	// metadata upload (used by the mock API in acceptance tests).
	FilePath string `json:"-"`
	// SourceURL is an http(s) URL to download the .jar file from. When set
	// (and FilePath is empty), UploadCustomJar downloads the file from the URL
	// and uploads it via multipart/form-data, just like a local file.
	SourceURL string `json:"-"`
}

// UnmarshalJSON supports both response shapes for the custom-jars endpoint:
//   - The real API returns an array of per-server objects, each with a "jars"
//     list of installed jar names, e.g. [{"jars":["a.jar","b.jar"]}, ...].
//   - The mock API returns an object with a "results" list of jar objects,
//     e.g. {"results":[{"name":"a.jar"}]}.
func (l *CustomJarsList) UnmarshalJSON(data []byte) error {
	trimmed := bytes.TrimSpace(data)
	if len(trimmed) > 0 && trimmed[0] == '[' {
		var servers []struct {
			Jars []string `json:"jars"`
		}
		if err := json.Unmarshal(trimmed, &servers); err != nil {
			return err
		}
		seen := make(map[string]bool)
		l.Results = nil
		for _, s := range servers {
			for _, name := range s.Jars {
				if !seen[name] {
					seen[name] = true
					l.Results = append(l.Results, CustomJar{Name: name})
				}
			}
		}
		return nil
	}

	type alias CustomJarsList
	var a alias
	if err := json.Unmarshal(trimmed, &a); err != nil {
		return err
	}
	*l = CustomJarsList(a)
	return nil
}

func (c *Client) GetCustomJars(accountName, deploymentID string) (*CustomJarsList, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/account/%s/deployment/%s/solr/custom-jars/", c.HostURL, accountName, deploymentID), nil)
	if err != nil {
		return nil, err
	}
	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}
	out := CustomJarsList{}
	if err := json.Unmarshal(body, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// UploadCustomJar uploads a custom jar to a deployment.
// The source of the .jar file is chosen in this order:
//   - jar.FilePath set: upload the local file via multipart/form-data.
//   - jar.SourceURL set: download the file from the http(s) URL, then upload
//     it via multipart/form-data.
//   - neither set: send a JSON metadata payload (used by the mock API in
//     acceptance tests).
func (c *Client) UploadCustomJar(accountName, deploymentID string, jar CustomJar) error {
	url := fmt.Sprintf("%s/account/%s/deployment/%s/solr/custom-jars/", c.HostURL, accountName, deploymentID)

	switch {
	case jar.FilePath != "":
		return c.uploadCustomJarFile(url, jar.FilePath)
	case jar.SourceURL != "":
		return c.uploadCustomJarURL(url, jar.SourceURL, jar.Name)
	default:
		return c.uploadCustomJarJSON(url, jar)
	}
}

// uploadCustomJarURL downloads the .jar file from an http(s) URL and uploads
// it to the deployment via multipart/form-data.
func (c *Client) uploadCustomJarURL(url, sourceURL, name string) error {
	u, err := neturl.Parse(sourceURL)
	if err != nil {
		return fmt.Errorf("parsing custom jar source_url %q: %w", sourceURL, err)
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return fmt.Errorf("custom jar source_url must be an http or https URL, got %q", sourceURL)
	}

	getReq, err := http.NewRequest("GET", sourceURL, nil)
	if err != nil {
		return err
	}
	getResp, err := c.HTTPClient.Do(getReq)
	if err != nil {
		return fmt.Errorf("downloading custom jar from %q: %w", sourceURL, err)
	}
	defer getResp.Body.Close()
	if getResp.StatusCode < 200 || getResp.StatusCode >= 300 {
		return fmt.Errorf("downloading custom jar from %q: unexpected status %d", sourceURL, getResp.StatusCode)
	}

	filename := name
	if filename == "" {
		filename = filepath.Base(u.Path)
	}

	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)
	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return err
	}
	if _, err := io.Copy(part, getResp.Body); err != nil {
		return err
	}
	if err := writer.Close(); err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, &buf)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	if _, err := c.doRequest(req); err != nil {
		return err
	}
	return nil
}

// uploadCustomJarFile performs the multipart/form-data file upload.
func (c *Client) uploadCustomJarFile(url, filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("opening custom jar %q: %w", filePath, err)
	}
	defer file.Close()

	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)
	part, err := writer.CreateFormFile("file", filepath.Base(filePath))
	if err != nil {
		return err
	}
	if _, err := io.Copy(part, file); err != nil {
		return err
	}
	if err := writer.Close(); err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, &buf)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	if _, err := c.doRequest(req); err != nil {
		return err
	}
	return nil
}

// uploadCustomJarJSON sends a JSON metadata payload (mock API path).
func (c *Client) uploadCustomJarJSON(url string, jar CustomJar) error {
	rb, err := json.Marshal(jar)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", url, strings.NewReader(string(rb)))
	if err != nil {
		return err
	}
	// The real API returns the updated jar list; a non-2xx status is already an
	// error, so reaching here means the jar was uploaded.
	if _, err := c.doRequest(req); err != nil {
		return err
	}
	return nil
}

func (c *Client) DeleteCustomJar(accountName, deploymentID, jarName string) error {
	// The real API returns a 500 (not a 404) when asked to delete a jar that
	// is no longer installed, so first check whether the jar is still present.
	// If it is already gone, deletion is a no-op (matches the Python module,
	// which GETs the jar list and returns success when the jar is not found).
	if jars, err := c.GetCustomJars(accountName, deploymentID); err == nil {
		installed := false
		for _, j := range jars.Results {
			if j.Name == jarName {
				installed = true
				break
			}
		}
		if !installed {
			return nil
		}
	}

	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/account/%s/deployment/%s/solr/custom-jars/%s/", c.HostURL, accountName, deploymentID, jarName), nil)
	if err != nil {
		return err
	}
	if _, err := c.doRequest(req); err != nil {
		// A 404 means the jar is already gone; deletion is idempotent.
		if isNotFound(err) {
			return nil
		}
		return err
	}
	// The real API returns a 2xx response whose body does not include a
	// "deleted" flag, so any successful response means the jar was removed.
	return nil
}
