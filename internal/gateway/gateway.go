package gateway

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"github.com/uforg/ufogateway/internal/util/randutil"
)

// Route represents a routing rule that maps an endpoint prefix to a destination URL.
type Route struct {
	ID                string // is the unique identifier for the route
	Endpoint          string // is the prefix to match incoming requests
	OriginURL         string // is the destination URL to proxy requests to
	TLSClientCert     string // is the content of the PEM file (optional)
	TLSClientKey      string // is the content of the key file (optional)
	TLSCaCert         string // is the content of the CA certificate (optional)
	TLSSkipCertVerify bool   // is a flag to skip TLS verification
}

// RouteProvider defines an interface to obtain the current list of routes.
type RouteProvider interface {
	// Routes returns the list of current routing rules.
	Routes() ([]Route, error)
}

// LogStorer defines an interface for storing request and response logs.
type LogStorer interface {
	// StoreRequestLog stores the log entry for a request.
	StoreRequestLog(reqLog RequestLog)
	// StoreResponseLog stores the log entry for a response.
	StoreResponseLog(respLog ResponseLog)
}

// RequestLog represents the data to be logged for an incoming request.
type RequestLog struct {
	RouteID           string              // Identifier of the route handling the request
	Timestamp         time.Time           // Timestamp when the request was received
	RequestID         string              // Unique identifier for the request
	RequestIP         string              // IP address of the client making the request
	RequestMethod     string              // HTTP method of the request
	RequestGatewayURL string              // URL of the gateway receiving the request
	RequestOriginURL  string              // URL of the origin server handling the request
	RequestHeaders    map[string][]string // Headers of the request
	RequestBody       io.Reader           // Body of the request
}

// ResponseLog represents the data to be logged for an outgoing response.
type ResponseLog struct {
	RouteID         string              // Identifier of the route handling the request
	Timestamp       time.Time           // Timestamp when the response was sent
	Duration        time.Duration       // Time taken to process the request
	RequestID       string              // Unique identifier matching the request
	ResponseHeaders map[string][]string // Headers of the response
	ResponseBody    io.Reader           // Body of the response
}

// Gateway represents the main gateway instance that routes and proxies requests.
type Gateway struct {
	routeProvider RouteProvider // Provider for obtaining the current routes
	logStorer     LogStorer     // Storer for logging requests and responses
}

// NewGateway creates a new gateway instance with the given route provider and log storer.
func NewGateway(routeProvider RouteProvider, logStorer LogStorer) *Gateway {
	return &Gateway{
		routeProvider: routeProvider,
		logStorer:     logStorer,
	}
}

// ServeHTTP handles incoming HTTP requests and proxies them to the appropriate backend.
// It implements the http.Handler interface.
func (g *Gateway) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	requestIP, err := getRequestIP(r)
	if err != nil {
		http.Error(w, "Gateway Error: failed to get request IP", http.StatusInternalServerError)
		return
	}

	if g.routeProvider == nil {
		http.Error(w, "Gateway Error: routeProvider is nil", http.StatusInternalServerError)
		return
	}
	if g.logStorer == nil {
		http.Error(w, "Gateway Error: logStorer is nil", http.StatusInternalServerError)
		return
	}

	routes, err := g.routeProvider.Routes()
	if err != nil {
		http.Error(w, "Gateway Error: failed to get routes", http.StatusInternalServerError)
		return
	}

	route, found := findRoute(routes, r.URL.Path)
	if !found {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	destURL, err := url.Parse(route.OriginURL)
	if err != nil {
		http.Error(w, "Gateway Error: failed to parse destination URL", http.StatusInternalServerError)
		return
	}

	requestID := randutil.GenerateIDForPocketBase()
	startTime := time.Now()

	var reqBody bytes.Buffer
	r.Body, err = readAndRestoreBody(r.Body, &reqBody)
	if err != nil {
		http.Error(w, "Gateway Error: failed to read request body", http.StatusInternalServerError)
		return
	}

	requestGatewayURL, requestOriginURL := getRequestURL(r, route)
	g.logStorer.StoreRequestLog(RequestLog{
		RouteID:           route.ID,
		Timestamp:         startTime,
		RequestID:         requestID,
		RequestIP:         requestIP,
		RequestMethod:     r.Method,
		RequestGatewayURL: requestGatewayURL,
		RequestOriginURL:  requestOriginURL,
		RequestHeaders:    cloneHeaderMap(r.Header),
		RequestBody:       bytes.NewReader(reqBody.Bytes()),
	})

	proxy := httputil.NewSingleHostReverseProxy(destURL)

	// Configure TLS if CertPEM is available
	if route.TLSClientCert != "" {
		tlsConfig, err := configureTLS(route.TLSClientCert, route.TLSClientKey, route.TLSCaCert, route.TLSSkipCertVerify)
		if err != nil {
			http.Error(w, "Gateway Error: failed to configure TLS", http.StatusInternalServerError)
			return
		}
		proxy.Transport = &http.Transport{TLSClientConfig: tlsConfig}
	}

	customWriter := newResponseWriter(w)
	r.URL.Path = gatewayToOriginPath(r.URL.Path, route.Endpoint)
	r.Host = destURL.Host
	proxy.ServeHTTP(customWriter, r)

	g.logStorer.StoreResponseLog(ResponseLog{
		RouteID:         route.ID,
		Timestamp:       time.Now(),
		Duration:        time.Since(startTime),
		RequestID:       requestID,
		ResponseHeaders: cloneHeaderMap(w.Header()),
		ResponseBody:    bytes.NewReader(customWriter.getBody()),
	})
}
