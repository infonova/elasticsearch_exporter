package collector

import (
	"fmt"
	"github.com/go-kit/kit/log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestIlm(t *testing.T) {
	ti := map[string]string{
		"7.3.2": `{"indices":{"foo_11":{"index":"foo_1","managed":true,"policy":"foo_policy","lifecycle_date_millis":1575630854324,"phase":"hot","phase_time_millis":1575605054674,"action":"complete","action_time_millis":1575630855862,"step":"complete","step_time_millis":1575630855862,"phase_execution":{"policy":"foo_policy","phase_definition":{"min_age":"0ms","actions":{"rollover":{"max_size":"15gb","max_age":"1d"},"set_priority":{"priority":100}}},"version":7,"modified_date_in_millis":1573070716617}}}}`,
	}
	for ver, out := range ti {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, out)
		}))
		defer ts.Close()

		u, err := url.Parse(ts.URL)
		if err != nil {
			t.Fatalf("Failed to parse URL: %s", err)
		}
		i := NewIlm(log.NewNopLogger(), http.DefaultClient, u)
		ilm, err := i.fetchAndDecodeIlm()
		if err != nil {
			t.Fatalf("Failed to fetch or decode ILM stats: %s", err)
		}
		t.Logf("[%s] ILM Response: %+v", ver, ilm)
		for ilmIndex, stats := range ilm.Indices {
			t.Logf(
				"Index: %s - Managed: %t - Action: %s - Phase: %s - Step: %s",
				ilmIndex,
				stats.Managed,
				stats.Action,
				stats.Phase,
				stats.Step,
			)
			if stats.Managed != true {
				t.Errorf("Wrong parsed value of managed flag")
			}
		}

	}
}
