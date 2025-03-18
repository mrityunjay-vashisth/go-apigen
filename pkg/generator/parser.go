package generator

import (
	"context"
	"fmt"
	"os"

	"github.com/getkin/kin-openapi/openapi3"
)

// ParseOpenAPIFile reads an OpenAPI YAML or JSON file from disk and returns the parsed doc.
func ParseOpenAPIFile(filename string) (*openapi3.T, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("cannot read file %s: %w", filename, err)
	}

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true

	doc, err := loader.LoadFromData(data)
	if err != nil {
		return nil, fmt.Errorf("failed to parse OpenAPI data: %w", err)
	}

	if err := doc.Validate(context.Background()); err != nil {
		return nil, fmt.Errorf("OpenAPI validation error: %w", err)
	}

	return doc, nil
}
