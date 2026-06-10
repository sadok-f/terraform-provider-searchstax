package client

import (
	"encoding/json"
	"fmt"
	"strings"
)

// FlexStringList unmarshals JSON string arrays or space-delimited strings.
type FlexStringList []string

func (f *FlexStringList) UnmarshalJSON(data []byte) error {
	if string(data) == "null" || len(data) == 0 {
		*f = nil
		return nil
	}
	var arr []string
	if err := json.Unmarshal(data, &arr); err == nil {
		*f = arr
		return nil
	}
	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		s = strings.TrimSpace(s)
		if s == "" {
			*f = nil
			return nil
		}
		*f = strings.Fields(s)
		return nil
	}
	return fmt.Errorf("FlexStringList: unsupported JSON value")
}

// DeploymentSpecifications holds deployment sizing metadata from the API.
type DeploymentSpecifications struct {
	JVMHeapMemory  string `json:"jvm_heap_memory"`
	DiskSpace      string `json:"disk_space"`
	PhysicalMemory string `json:"physical_memory"`
}

// FlexSpecifications unmarshals a specifications object and ignores non-object values.
type FlexSpecifications DeploymentSpecifications

func (f *FlexSpecifications) UnmarshalJSON(data []byte) error {
	if string(data) == "null" || len(data) == 0 {
		return nil
	}
	if data[0] != '{' {
		return nil
	}
	var specs DeploymentSpecifications
	if err := json.Unmarshal(data, &specs); err != nil {
		return err
	}
	*f = FlexSpecifications(specs)
	return nil
}
