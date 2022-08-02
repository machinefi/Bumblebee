package roundtrippers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/iotexproject/Bumblebee/conf/log"
	"github.com/iotexproject/Bumblebee/x/misc/timer"
	"github.com/pkg/errors"
)

type LogRoundTripper struct {
	next http.RoundTripper
}

func NewLogRoundTripper() func(rt http.RoundTripper) http.RoundTripper {
	return func(rt http.RoundTripper) http.RoundTripper {
		return &LogRoundTripper{
			next: rt,
		}
	}
}

func (rt *LogRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	cost := timer.Start()
	ctx, logger := log.Start(req.Context(), "Request")
	defer logger.End()

	resp, err := rt.next.RoundTrip(req.WithContext(ctx))

	defer func() {
		l := logger.WithValues(
			"cost", fmt.Sprintf("%0.3fms", float64(cost()/time.Millisecond)),
			"method", req.Method,
			"url", req.URL.String(),
			"metadata", req.Header,
		)

		if err == nil {
			l.Info("success")
		} else {
			l.Warn(errors.Wrap(err, "http request failed"))
		}
	}()

	return resp, err
}
