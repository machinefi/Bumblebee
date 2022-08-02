package roundtrippers_test

import (
	"net/http"
	"testing"

	"github.com/iotexproject/Bumblebee/kit/httptransport"
	. "github.com/iotexproject/Bumblebee/kit/httptransport/client/roundtrippers"
)

func TestLogRoundTripper(t *testing.T) {
	mgr := httptransport.NewRequestTsfmFactory(nil, nil)
	mgr.SetDefault()

	req, _ := mgr.NewRequest(http.MethodGet, "https://github.com", nil)
	_, _ = NewLogRoundTripper()(http.DefaultTransport).RoundTrip(req)
}
