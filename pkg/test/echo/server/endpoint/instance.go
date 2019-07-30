// Copyright 2019 Istio Authors
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

package endpoint

import (
	"fmt"
	"io"

	"istio.io/istio/pilot/pkg/model"
	"istio.io/istio/pkg/config"
	"istio.io/istio/pkg/test/echo/common"
)

// IsServerReadyFunc is a function that indicates whether the server is currently ready to handle traffic.
type IsServerReadyFunc func() bool

// OnReadyFunc is a callback function that informs the server that the endpoint is ready.
type OnReadyFunc func()

// Config for a single endpoint Instance.
type Config struct {
	IsServerReady IsServerReadyFunc
	Version       string
	TLSCert       string
	TLSKey        string
	UDSServer     string
	Dialer        common.Dialer
	Port          *model.Port
}

// Instance of an endpoint that serves the Echo application on a single port/protocol.
type Instance interface {
	io.Closer
	Start(onReady OnReadyFunc) error
}

// New creates a new endpoint Instance.
func New(cfg Config) (Instance, error) {
	if cfg.Port != nil {
		switch cfg.Port.Protocol {
		case config.ProtocolTCP, config.ProtocolHTTP, config.ProtocolHTTPS:
			return newHTTP(cfg), nil
		case config.ProtocolHTTP2, config.ProtocolGRPC:
			return newGRPC(cfg), nil
		default:
			return nil, fmt.Errorf("unsupported protocol: %s", cfg.Port.Protocol)
		}
	}

	if len(cfg.UDSServer) > 0 {
		return newHTTP(cfg), nil
	}

	return nil, fmt.Errorf("either port or UDS must be specified")
}