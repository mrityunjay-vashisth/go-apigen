package generator

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gorilla/mux"
)

// RouteDefinition holds both a Handler and optional Middlewares for each operationId.
type RouteDefinition struct {
	Handler     http.HandlerFunc
	Middlewares []mux.MiddlewareFunc
}

// OperationMap is a map of operationId → RouteDefinition.
// The user populates it with the desired handler and optional middlewares.
type OperationMap map[string]RouteDefinition

// GenerateMuxRouter takes a parsed OpenAPI doc and an OperationMap (mapping operationId->handler),
// then returns a fully-wired Gorilla Mux router. If an operationId has no matching handler,
// it registers a default 501 (Not Implemented) handler.
func GenerateMuxRouter(doc *openapi3.T, ops OperationMap, globalMiddlewares ...mux.MiddlewareFunc) (*mux.Router, error) {
	r := mux.NewRouter()

	// Optionally attach global middlewares at the root router level
	if len(globalMiddlewares) > 0 {
		r.Use(globalMiddlewares...)
	}

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
			routeDef, found := ops[opID]
			if !found {
				// No route def for this opID => use a default "not implemented" handler
				routeDef = RouteDefinition{Handler: notImplementedHandler(opID)}
			}

			// 1) Create a route that matches this path & method
			route := r.NewRoute().Path(path).Methods(strings.ToUpper(method))

			// 2) Convert it to a subrouter so we can attach middlewares
			subrouter := route.Subrouter()

			// 3) Attach any middlewares for this operation to the subrouter
			if len(routeDef.Middlewares) > 0 {
				subrouter.Use(routeDef.Middlewares...)
			}

			// 4) Register the actual handler function on the subrouter
			//    Notice we use HandleFunc("") because the subrouter already has the path set.
			subrouter.HandleFunc("", routeDef.Handler)
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
