// Code created by gotmpl. DO NOT MODIFY.
// source: internal/shared/otlp/otlptrace/otlpconfig/options.go.tmpl

// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package otlpconfig // import "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp/internal/otlpconfig"

import (
	"crypto/tls"
	"fmt"
	"path"
	"strings"
	"time"

	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp/internal/retry"
)

const (
	// DefaultTracesPath is a default URL path for endpoint that
	// receives spans.
	DefaultTracesPath string = "/v1/traces"
	// DefaultTimeout is a default max waiting time for the backend to process
	// each span batch.
	DefaultTimeout time.Duration = 10 * time.Second
)

type (
	SignalConfig struct {
		Endpoint    string
		Insecure    bool
		TLSCfg      *tls.Config
		Headers     map[string]string
		Compression Compression
		Timeout     time.Duration
		URLPath     string
	}

	Config struct {
		// Signal specific configurations
		Traces SignalConfig

		RetryConfig retry.Config
	}
)

// NewConfig returns a new Config with all settings applied from opts and
// any unset setting using the default HTTP config values.
func NewConfig(opts ...Option) Config {
	cfg := Config{
		Traces: SignalConfig{
			Endpoint:    fmt.Sprintf("%s:%d", DefaultCollectorHost, DefaultCollectorHTTPPort),
			URLPath:     DefaultTracesPath,
			Compression: NoCompression,
			Timeout:     DefaultTimeout,
		},
		RetryConfig: retry.DefaultConfig,
	}
	cfg = ApplyHTTPEnvConfigs(cfg)
	for _, opt := range opts {
		cfg = opt.Apply(cfg)
	}
	cfg.Traces.URLPath = cleanPath(cfg.Traces.URLPath, DefaultTracesPath)
	return cfg
}

// cleanPath returns a path with all spaces trimmed and all redundancies
// removed. If urlPath is empty or cleaning it results in an empty string,
// defaultPath is returned instead.
func cleanPath(urlPath string, defaultPath string) string {
	tmp := path.Clean(strings.TrimSpace(urlPath))
	if tmp == "." {
		return defaultPath
	}
	if !path.IsAbs(tmp) {
		tmp = fmt.Sprintf("/%s", tmp)
	}
	return tmp
}

// Option applies an option to the driver.
type Option interface {
	Apply(Config) Config

	// A private method to prevent users implementing the
	// interface and so future additions to it will not
	// violate compatibility.
	private()
}

type option struct {
	fn func(Config) Config
}

func (g option) Apply(cfg Config) Config {
	return g.fn(cfg)
}

func (option) private() {}

func NewOption(fn func(cfg Config) Config) Option {
	return &option{fn: fn}
}

func WithEndpoint(endpoint string) Option {
	return NewOption(func(cfg Config) Config {
		cfg.Traces.Endpoint = endpoint
		return cfg
	})
}

func WithCompression(compression Compression) Option {
	return NewOption(func(cfg Config) Config {
		cfg.Traces.Compression = compression
		return cfg
	})
}

func WithURLPath(urlPath string) Option {
	return NewOption(func(cfg Config) Config {
		cfg.Traces.URLPath = urlPath
		return cfg
	})
}

func WithRetry(rc retry.Config) Option {
	return NewOption(func(cfg Config) Config {
		cfg.RetryConfig = rc
		return cfg
	})
}

func WithTLSClientConfig(tlsCfg *tls.Config) Option {
	return NewOption(func(cfg Config) Config {
		cfg.Traces.TLSCfg = tlsCfg.Clone()
		return cfg
	})
}

func WithInsecure() Option {
	return NewOption(func(cfg Config) Config {
		cfg.Traces.Insecure = true
		return cfg
	})
}

func WithSecure() Option {
	return NewOption(func(cfg Config) Config {
		cfg.Traces.Insecure = false
		return cfg
	})
}

func WithHeaders(headers map[string]string) Option {
	return NewOption(func(cfg Config) Config {
		cfg.Traces.Headers = headers
		return cfg
	})
}

func WithTimeout(duration time.Duration) Option {
	return NewOption(func(cfg Config) Config {
		cfg.Traces.Timeout = duration
		return cfg
	})
}
