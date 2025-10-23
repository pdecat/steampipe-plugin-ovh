package ovh

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httputil"
	"strings"

	"github.com/hashicorp/go-hclog"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin"
)

type transport struct {
	name      string
	ctx       context.Context
	transport http.RoundTripper
	debug     bool
}

func IsDebugOrHigher(ctx context.Context) bool {
	level := plugin.Logger(ctx).GetLevel()
	return level == hclog.Debug || level == hclog.Trace
}

func (t *transport) RoundTrip(req *http.Request) (*http.Response, error) {
	if t.debug {
		reqData, err := httputil.DumpRequestOut(req, true)
		if err == nil {
			plugin.Logger(t.ctx).Debug(fmt.Sprintf(logReqMsg, t.name, prettyPrintJsonLines(reqData)))
		} else {
			plugin.Logger(t.ctx).Error(fmt.Sprintf("%s API Request error: %#v", t.name, err))
		}
	}

	resp, err := t.transport.RoundTrip(req)
	if err != nil {
		return resp, err
	}

	if t.debug {
		respData, err := httputil.DumpResponse(resp, true)
		if err == nil {
			plugin.Logger(t.ctx).Debug(fmt.Sprintf(logRespMsg, t.name, prettyPrintJsonLines(respData)))
		} else {
			plugin.Logger(t.ctx).Error(fmt.Sprintf("%s API Response error: %#v", t.name, err))
		}
	}

	return resp, nil
}

func NewTransport(name string, ctx context.Context, t http.RoundTripper, debug bool) *transport {
	if debug {
		plugin.Logger(ctx).Debug("[DEBUG] Enabling HTTP requests/responses tracing")
	}
	return &transport{name, ctx, t, debug}
}

// prettyPrintJsonLines iterates through a []byte line-by-line,
// transforming any lines that are complete json into pretty-printed json.
func prettyPrintJsonLines(b []byte) string {
	parts := strings.Split(string(b), "\n")
	for i, p := range parts {
		if b := []byte(p); json.Valid(b) {
			var out bytes.Buffer
			_ = json.Indent(&out, b, "", " ")
			parts[i] = out.String()
		}
	}
	return strings.Join(parts, "\n")
}

const logReqMsg = `%s API Request Details:
---[ REQUEST ]---------------------------------------
%s
-----------------------------------------------------`

const logRespMsg = `%s API Response Details:
---[ RESPONSE ]--------------------------------------
%s
-----------------------------------------------------`
