// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package source // import "github.com/open-telemetry/opentelemetry-collector-contrib/extension/headerssetterextension/internal/source"

import (
	"context"
	"fmt"
	"strings"

	"go.opentelemetry.io/collector/client"
)

var _ Source = (*ContextSource)(nil)

type ContextSource struct {
	Key          string
	DefaultValue string
}

func (ts *ContextSource) Get(ctx context.Context) (string, error) {
	const authPrefix = "auth."

	cl := client.FromContext(ctx)
	ss := cl.Metadata.Get(ts.Key)

	if strings.HasPrefix(ts.Key, authPrefix) {
		attrName := strings.TrimPrefix(ts.Key, authPrefix)
		attr := cl.Auth.GetAttribute(attrName)

		switch a := attr.(type) {
		case string:
			ss = []string{a}
		case []string:
			ss = a
		}
	}

	if len(ss) == 0 {
		return ts.DefaultValue, nil
	}

	if len(ss) > 1 {
		return "", fmt.Errorf("%d source keys found in the context, can't determine which one to use", len(ss))
	}

	return ss[0], nil
}
