package admin_api_test

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/APTrust/registry/app"
	"github.com/APTrust/registry/constants"
	"github.com/APTrust/registry/middleware"
	"github.com/APTrust/registry/pgmodels"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

// These tests keep admin_api_v3.yml honest. The Admin API has no generated
// documentation, so the spec is hand-written, and hand-written specs rot as
// soon as someone adds a route or a filter without updating them. Rather than
// check the prose, these tests check the facts that can be checked: that every
// route is documented, that nothing is documented that doesn't exist, and that
// the query params we tell callers about match the ones the Registry actually
// accepts.
//
// If one of these fails, fix admin_api_v3.yml. Don't delete the assertion.

const specPath = "../../../admin_api_v3.yml"

// specParam is one entry in an operation's parameter list. Entries are either
// inline (Name is set) or a reference to components/parameters (Ref is set).
type specParam struct {
	Name string `yaml:"name"`
	In   string `yaml:"in"`
	Ref  string `yaml:"$ref"`
}

type specOperation struct {
	OperationID string      `yaml:"operationId"`
	Parameters  []specParam `yaml:"parameters"`
}

type apiSpec struct {
	Components struct {
		Parameters map[string]specParam `yaml:"parameters"`
	} `yaml:"components"`
	Paths map[string]map[string]specOperation `yaml:"paths"`
}

// paramName returns the name of a parameter, resolving it through
// components/parameters if the entry is a $ref.
func (spec *apiSpec) paramName(p specParam) string {
	if p.Name != "" {
		return p.Name
	}
	key := strings.TrimPrefix(p.Ref, "#/components/parameters/")
	return spec.Components.Parameters[key].Name
}

func loadSpec(t *testing.T) *apiSpec {
	path, err := filepath.Abs(specPath)
	require.Nil(t, err, "Can't resolve path to admin_api_v3.yml")
	data, err := os.ReadFile(path)
	require.Nil(t, err, "Can't read %s", path)
	spec := &apiSpec{}
	require.Nil(t, yaml.Unmarshal(data, spec), "Can't parse %s", path)
	require.NotEmpty(t, spec.Paths)
	return spec
}

// adminRoutes returns the Admin API routes the Gin router actually serves,
// keyed by "method path", with the name of the handler function as the value.
// Path params are converted from Gin's :id and *id forms to OpenAPI's {id}.
func adminRoutes(t *testing.T) map[string]string {
	engine := app.InitAppEngine(true)
	routes := make(map[string]string)
	for _, route := range engine.Routes() {
		if !strings.HasPrefix(route.Path, constants.APIPrefixAdmin) {
			continue
		}
		path := strings.NewReplacer(":", "{", "*", "{").Replace(route.Path)
		parts := strings.Split(path, "/")
		for i, part := range parts {
			if strings.HasPrefix(part, "{") {
				parts[i] = part + "}"
			}
		}
		handler := route.Handler
		if i := strings.LastIndex(handler, "."); i > -1 {
			handler = handler[i+1:]
		}
		key := strings.ToLower(route.Method) + " " + strings.Join(parts, "/")
		routes[key] = handler
	}
	require.NotEmpty(t, routes)
	return routes
}

func specOperations(spec *apiSpec) map[string]specOperation {
	ops := make(map[string]specOperation)
	for path, methods := range spec.Paths {
		for method, op := range methods {
			ops[method+" "+path] = op
		}
	}
	return ops
}

// TestAdminAPISpecCoversAllRoutes fails when a route is added to or removed
// from the Admin API without a matching change to admin_api_v3.yml.
func TestAdminAPISpecCoversAllRoutes(t *testing.T) {
	spec := loadSpec(t)
	routes := adminRoutes(t)
	ops := specOperations(spec)

	undocumented := make([]string, 0)
	for key := range routes {
		if _, ok := ops[key]; !ok {
			undocumented = append(undocumented, key)
		}
	}
	sort.Strings(undocumented)
	assert.Empty(t, undocumented, "These Admin API routes are not documented in admin_api_v3.yml")

	nonexistent := make([]string, 0)
	for key := range ops {
		if _, ok := routes[key]; !ok {
			nonexistent = append(nonexistent, key)
		}
	}
	sort.Strings(nonexistent)
	assert.Empty(t, nonexistent, "admin_api_v3.yml documents these routes, but the Registry does not serve them")
}

// TestAdminAPISpecOperationIDs checks that each operationId is the name of the
// Go handler function that serves the route. Keeping these in sync means a
// reader can jump from the docs straight to the code that produced them.
func TestAdminAPISpecOperationIDs(t *testing.T) {
	spec := loadSpec(t)
	routes := adminRoutes(t)

	for key, op := range specOperations(spec) {
		handler, ok := routes[key]
		if !ok {
			continue // reported by TestAdminAPISpecCoversAllRoutes
		}
		assert.Equal(t, handler, op.OperationID, "operationId for %s should be the handler function name", key)
	}
}

// TestAdminAPISpecDocumentsAllFilters checks that the query params documented
// for each index route match the filters the Registry actually accepts.
//
// This matters more than it looks. Request.ValidateFilters rejects any query
// param it doesn't recognize with a 400, so an undocumented filter is a feature
// callers can't discover, and a documented one that doesn't exist is a call
// that fails.
//
// This keys off operationId rather than the router, so it needs no database.
// TestAdminAPISpecOperationIDs is what guarantees operationId is the handler
// name, which is what makes the AuthMap lookup below correct.
func TestAdminAPISpecDocumentsAllFilters(t *testing.T) {
	spec := loadSpec(t)

	// Every index route accepts these in addition to its own filters.
	commonParams := []string{"page", "per_page", "sort"}

	tested := 0
	for key, op := range specOperations(spec) {
		handler := op.OperationID
		if !strings.HasSuffix(handler, "Index") {
			continue
		}
		authMeta, ok := middleware.AuthMap[handler]
		require.True(t, ok, "No AuthMap entry for handler %s", handler)

		expected := append(commonParams, pgmodels.FiltersFor(authMeta.ResourceType)...)
		sort.Strings(expected)

		actual := make([]string, 0)
		for _, param := range op.Parameters {
			if name := spec.paramName(param); name != "" {
				actual = append(actual, name)
			}
		}
		sort.Strings(actual)

		assert.Equal(t, expected, actual, "Documented query params for %s don't match the filters the Registry accepts", key)
		tested++
	}
	assert.NotZero(t, tested, "Found no index routes to check. Has the route table or handler naming changed?")
}
