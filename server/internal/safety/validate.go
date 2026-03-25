package safety

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/url"
	"strconv"
	"strings"
)

type ErrorCode string

const (
	ErrorInvalidURL        ErrorCode = "invalid_url"
	ErrorUnsupportedScheme ErrorCode = "unsupported_scheme"
	ErrorMissingHost       ErrorCode = "missing_host"
	ErrorBlockedHost       ErrorCode = "blocked_host"
	ErrorBlockedIP         ErrorCode = "blocked_ip"
	ErrorResolutionFailed  ErrorCode = "resolution_failed"
	ErrorPortNotAllowed    ErrorCode = "port_not_allowed"
	ErrorTooManyRedirects  ErrorCode = "too_many_redirects"
)

type ValidationError struct {
	Code    ErrorCode
	Message string
	URL     string
	Host    string
	Port    string
	IP      net.IP
	Err     error
}

func AsValidationError(err error, target **ValidationError) bool {
	return errors.As(err, target)
}

func (e *ValidationError) Error() string {
	if e == nil {
		return "<nil>"
	}

	if e.URL != "" {
		return fmt.Sprintf("%s: %s (%s)", e.Code, e.Message, e.URL)
	}

	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func (e *ValidationError) Unwrap() error {
	if e == nil {
		return nil
	}

	return e.Err
}

type Resolver interface {
	LookupIPAddr(ctx context.Context, host string) ([]net.IPAddr, error)
}

type Options struct {
	AllowedSchemes []string
	AllowedPorts   []int
	MaxRedirects   int
	Resolver       Resolver
}

type ValidationResult struct {
	URL          *url.URL
	ResolvedIPs  []net.IP
	ResolvedHost string
	Port         int
}

type ChainResult struct {
	Hops []ValidationResult
}

type defaultResolver struct{}

func (defaultResolver) LookupIPAddr(ctx context.Context, host string) ([]net.IPAddr, error) {
	return net.DefaultResolver.LookupIPAddr(ctx, host)
}

func DefaultOptions() Options {
	return Options{
		AllowedSchemes: []string{"http", "https"},
		AllowedPorts:   []int{80, 443},
		MaxRedirects:   3,
		Resolver:       defaultResolver{},
	}
}

func ValidateURL(ctx context.Context, rawURL string, opts Options) (ValidationResult, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	normalized := normalizeOptions(opts)
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return ValidationResult{}, &ValidationError{
			Code:    ErrorInvalidURL,
			Message: "invalid URL",
			URL:     rawURL,
			Err:     err,
		}
	}

	if parsed.Scheme == "" || !containsString(normalized.AllowedSchemes, parsed.Scheme) {
		return ValidationResult{}, &ValidationError{
			Code:    ErrorUnsupportedScheme,
			Message: "unsupported URL scheme",
			URL:     rawURL,
		}
	}

	host := parsed.Hostname()
	if host == "" {
		return ValidationResult{}, &ValidationError{
			Code:    ErrorMissingHost,
			Message: "URL host is required",
			URL:     rawURL,
		}
	}

	if blockedHost(host) {
		return ValidationResult{}, &ValidationError{
			Code:    ErrorBlockedHost,
			Message: "destination host is blocked",
			URL:     rawURL,
			Host:    host,
		}
	}

	resolvedIPs, err := resolveIPs(ctx, normalized.Resolver, host)
	if err != nil {
		return ValidationResult{}, err
	}

	for _, ip := range resolvedIPs {
		if blockedIP(ip) {
			return ValidationResult{}, &ValidationError{
				Code:    ErrorBlockedIP,
				Message: "destination IP is blocked",
				URL:     rawURL,
				Host:    host,
				Port:    parsed.Port(),
				IP:      cloneIP(ip),
			}
		}
	}

	portString := parsed.Port()
	port, hasExplicitPort, err := resolvePort(parsed)
	if err != nil {
		return ValidationResult{}, err
	}
	if hasExplicitPort && !containsInt(normalized.AllowedPorts, port) {
		return ValidationResult{}, &ValidationError{
			Code:    ErrorPortNotAllowed,
			Message: "destination port is not allowed",
			URL:     rawURL,
			Host:    host,
			Port:    portString,
		}
	}

	return ValidationResult{
		URL:          parsed,
		ResolvedIPs:  cloneIPs(resolvedIPs),
		ResolvedHost: host,
		Port:         port,
	}, nil
}

func ValidateRedirectChain(ctx context.Context, urls []string, opts Options) (ChainResult, error) {
	normalized := normalizeOptions(opts)
	if len(urls) == 0 {
		return ChainResult{}, &ValidationError{
			Code:    ErrorInvalidURL,
			Message: "redirect chain is empty",
		}
	}

	if normalized.MaxRedirects >= 0 && len(urls)-1 > normalized.MaxRedirects {
		return ChainResult{}, &ValidationError{
			Code:    ErrorTooManyRedirects,
			Message: "redirect limit exceeded",
		}
	}

	result := ChainResult{Hops: make([]ValidationResult, 0, len(urls))}
	for _, rawURL := range urls {
		hop, err := ValidateURL(ctx, rawURL, normalized)
		if err != nil {
			return ChainResult{}, err
		}

		result.Hops = append(result.Hops, hop)
	}

	return result, nil
}

func normalizeOptions(opts Options) Options {
	normalized := opts
	if len(normalized.AllowedSchemes) == 0 {
		normalized.AllowedSchemes = DefaultOptions().AllowedSchemes
	}
	if len(normalized.AllowedPorts) == 0 {
		normalized.AllowedPorts = DefaultOptions().AllowedPorts
	}
	if normalized.MaxRedirects == 0 {
		normalized.MaxRedirects = DefaultOptions().MaxRedirects
	}
	if normalized.Resolver == nil {
		normalized.Resolver = defaultResolver{}
	}

	return normalized
}

func resolvePort(parsed *url.URL) (int, bool, error) {
	if portString := parsed.Port(); portString != "" {
		port, err := strconv.Atoi(portString)
		if err != nil || port <= 0 || port > 65535 {
			return 0, true, &ValidationError{
				Code:    ErrorPortNotAllowed,
				Message: "destination port is invalid",
				URL:     parsed.String(),
				Host:    parsed.Hostname(),
				Port:    portString,
			}
		}

		return port, true, nil
	}

	switch strings.ToLower(parsed.Scheme) {
	case "http":
		return 80, false, nil
	case "https":
		return 443, false, nil
	default:
		return 0, false, &ValidationError{
			Code:    ErrorUnsupportedScheme,
			Message: "unsupported URL scheme",
			URL:     parsed.String(),
		}
	}
}

func resolveIPs(ctx context.Context, resolver Resolver, host string) ([]net.IP, error) {
	if ip := net.ParseIP(host); ip != nil {
		return []net.IP{cloneIP(ip)}, nil
	}

	addrs, err := resolver.LookupIPAddr(ctx, host)
	if err != nil {
		return nil, &ValidationError{
			Code:    ErrorResolutionFailed,
			Message: "failed to resolve destination host",
			Host:    host,
			Err:     err,
		}
	}

	ips := make([]net.IP, 0, len(addrs))
	for _, addr := range addrs {
		if addr.IP == nil {
			continue
		}
		ips = append(ips, cloneIP(addr.IP))
	}

	if len(ips) == 0 {
		return nil, &ValidationError{
			Code:    ErrorResolutionFailed,
			Message: "host resolved to no usable IPs",
			Host:    host,
		}
	}

	return ips, nil
}

func blockedHost(host string) bool {
	lower := strings.TrimSuffix(strings.ToLower(host), ".")
	return lower == "localhost" || strings.HasSuffix(lower, ".localhost")
}

func blockedIP(ip net.IP) bool {
	if ip == nil {
		return true
	}

	if ip.IsLoopback() || ip.IsPrivate() || ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() || ip.IsUnspecified() || ip.IsMulticast() {
		return true
	}

	metadataIP := net.ParseIP("169.254.169.254")
	return metadataIP != nil && metadataIP.Equal(ip)
}

func cloneIP(ip net.IP) net.IP {
	if ip == nil {
		return nil
	}

	return append(net.IP(nil), ip...)
}

func cloneIPs(ips []net.IP) []net.IP {
	if len(ips) == 0 {
		return nil
	}

	cloned := make([]net.IP, 0, len(ips))
	for _, ip := range ips {
		cloned = append(cloned, cloneIP(ip))
	}

	return cloned
}

func containsString(values []string, target string) bool {
	for _, value := range values {
		if strings.EqualFold(value, target) {
			return true
		}
	}

	return false
}

func containsInt(values []int, target int) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}

	return false
}
