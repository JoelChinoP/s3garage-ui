package router

import (
	"errors"
	"fmt"
	"khairul169/garage-webui/utils"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

var allowedProxyRoutes = map[string]map[string]struct{}{
	http.MethodGet: {
		"/v2/GetBucketInfo":    {},
		"/v2/GetClusterHealth": {},
		"/v2/GetClusterLayout": {},
		"/v2/GetClusterStatus": {},
		"/v2/GetKeyInfo":       {},
		"/v2/GetNodeInfo":      {},
		"/v2/ListBuckets":      {},
		"/v2/ListKeys":         {},
	},
	http.MethodPost: {
		"/v2/AddBucketAlias":      {},
		"/v2/AllowBucketKey":      {},
		"/v2/ApplyClusterLayout":  {},
		"/v2/ConnectClusterNodes": {},
		"/v2/CreateBucket":        {},
		"/v2/CreateKey":           {},
		"/v2/DeleteBucket":        {},
		"/v2/DeleteKey":           {},
		"/v2/DenyBucketKey":       {},
		"/v2/ImportKey":           {},
		"/v2/RemoveBucketAlias":   {},
		"/v2/RevertClusterLayout": {},
		"/v2/UpdateBucket":        {},
		"/v2/UpdateClusterLayout": {},
	},
}

func ProxyHandler(w http.ResponseWriter, r *http.Request) {
	if !isAllowedProxyRoute(r.Method, r.URL.Path) {
		utils.ResponseErrorStatus(w, errors.New("admin api endpoint is not allowed"), http.StatusForbidden)
		return
	}

	target, err := url.Parse(utils.Garage.GetAdminEndpoint())
	if err != nil {
		utils.ResponseError(w, err)
		return
	}

	proxy := &httputil.ReverseProxy{
		Rewrite: func(r *httputil.ProxyRequest) {
			r.SetURL(target)
			r.Out.URL.Path = strings.TrimPrefix(r.In.URL.Path, "/api")
			r.Out.Header.Set("Authorization", fmt.Sprintf("Bearer %s", utils.Garage.GetAdminKey()))
		},
	}

	proxy.ServeHTTP(w, r)
}

func isAllowedProxyRoute(method string, path string) bool {
	routes, ok := allowedProxyRoutes[method]
	if !ok {
		return false
	}

	_, ok = routes[path]
	return ok
}
