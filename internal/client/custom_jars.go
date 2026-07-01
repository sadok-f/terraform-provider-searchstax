package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
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
// When jar.FilePath is set, the actual .jar file is uploaded using a
// multipart/form-data request with a "file" form field, as required by the
// SearchStax API. When FilePath is empty, a JSON metadata payload is sent
// instead (used by the mock API in acceptance tests).
func (c *Client) UploadCustomJar(accountName, deploymentID string, jar CustomJar) error {
	url := fmt.Sprintf("%s/account/%s/deployment/%s/solr/custom-jars/", c.HostURL, accountName, deploymentID)

	if jar.FilePath == "" {
		return c.uploadCustomJarJSON(url, jar)
	}
	return c.uploadCustomJarFile(url, jar.FilePath)
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
	body, err := c.doRequest(req)
	if err != nil {
		return err
	}
	var resp struct {
		Uploaded bool `json:"uploaded"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		return err
	}
	if !resp.Uploaded {
		return fmt.Errorf("custom jar not uploaded")
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
