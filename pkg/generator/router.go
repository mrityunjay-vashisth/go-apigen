package generator

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gorilla/mux"
)

// OperationMap is a map from operationId -> http.HandlerFunc
// The user of this package provides a function for each operationId found in the spec.
type OperationMap map[string]http.HandlerFunc

// GenerateMuxRouter takes a parsed OpenAPI doc and an OperationMap (mapping operationId->handler),
// then returns a fully-wired Gorilla Mux router. If an operationId has no matching handler,
// it registers a default 501 (Not Implemented) handler.
func GenerateMuxRouter(doc *openapi3.T, ops OperationMap) (*mux.Router, error) {
	r := mux.NewRouter()

	// For each path in the OpenAPI doc
	for path, pathItem := range doc.Paths.Map() {
		// For each method (GET, POST, etc.) in the path
		for method, operation := range pathItem.Operations() {
			opID := operation.OperationID
			if opID == "" {
				// If no operationId, skip or handle it differently as needed
				continue
			}

			// Decide which handler to attach
			handler, found := ops[opID]
			if !found {
				handler = notImplementedHandler(opID)
			}

			// Register the route
			r.HandleFunc(path, handler).Methods(strings.ToUpper(method))
		}
	}

	return r, nil
}

// notImplementedHandler is used if the user didn't provide a handler for an operationId.
func notImplementedHandler(opID string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotImplemented)
		fmt.Fprintf(w, "Operation %q is not implemented\n", opID)
	}
}

// Example helper for path params, if the user wants it:
func GetPathParam(r *http.Request, key string) string {
	return mux.Vars(r)[key]
}
