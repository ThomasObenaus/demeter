package alertmanager

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	saa "github.com/thomasobenaus/sokar/scaleAlertAggregator"
)

func TestNewConnector(t *testing.T) {

	cfg := Config{}
	connector := cfg.New()

	assert.NotNil(t, connector)
}

func Test_GenReceiver(t *testing.T) {
	emitter := genEmitterName("hello")
	assert.Equal(t, "AM.hello", emitter)

	emitter = genEmitterName("")
	assert.Equal(t, "AM", emitter)
}

func Test_FireScaleAlert(t *testing.T) {

	cfg := Config{}
	connector := cfg.New()
	require.NotNil(t, connector)

	subscriber := make(chan saa.ScaleAlertPacket)

	connector.Subscribe(subscriber)

	var alertsAll []saa.ScaleAlert

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		for pkg := range subscriber {
			alertsAll = append(alertsAll, pkg.ScaleAlerts...)
		}
		defer wg.Done()
	}()

	sentAlerts := make([]saa.ScaleAlert, 0)
	sentAlerts = append(sentAlerts, saa.ScaleAlert{Firing: true, Name: "A"})
	sentAlerts = append(sentAlerts, saa.ScaleAlert{Firing: true, Name: "B"})
	sentAlerts = append(sentAlerts, saa.ScaleAlert{Firing: true, Name: "C"})
	pkg := saa.ScaleAlertPacket{ScaleAlerts: sentAlerts}
	connector.fireScaleAlertPacket(pkg)

	sentAlerts = make([]saa.ScaleAlert, 0)
	sentAlerts = append(sentAlerts, saa.ScaleAlert{Firing: false, Name: "A"})
	sentAlerts = append(sentAlerts, saa.ScaleAlert{Firing: false, Name: "B"})
	sentAlerts = append(sentAlerts, saa.ScaleAlert{Firing: false, Name: "C"})
	pkg = saa.ScaleAlertPacket{ScaleAlerts: sentAlerts}
	connector.fireScaleAlertPacket(pkg)

	close(subscriber)

	wg.Wait()
	assert.Equal(t, 6, len(alertsAll))
}
func Test_HandleScaleAlert_Invalid(t *testing.T) {

	cfg := Config{}
	connector := cfg.New()
	require.NotNil(t, connector)

	req := httptest.NewRequest("GET", "http://example.com/foo", nil)
	w := httptest.NewRecorder()
	connector.HandleScaleAlerts(w, req, httprouter.Params{})

	resp := w.Result()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	buf := bytes.NewBufferString("INVALID RESPONSE")
	req = httptest.NewRequest("POST", "http://example.com/foo", buf)
	w = httptest.NewRecorder()
	connector.HandleScaleAlerts(w, req, httprouter.Params{})

	resp = w.Result()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func Test_HandleScaleAlert_Success(t *testing.T) {

	cfg := Config{}
	connector := cfg.New()
	require.NotNil(t, connector)
	alertName := "ABC"
	startsAt := time.Now()

	labels := map[string]string{"alertname": alertName}
	alerts := make([]alert, 0)
	alerts = append(alerts, alert{
		Status:   "Firing",
		Labels:   labels,
		StartsAt: startsAt,
	})

	data, err := json.Marshal(response{Alerts: alerts})

	require.NoError(t, err)
	require.NotEmpty(t, data)
	buf := bytes.NewReader(data)
	req := httptest.NewRequest("POST", "http://example.com/foo", buf)
	w := httptest.NewRecorder()

	subscriber := make(chan saa.ScaleAlertPacket)
	connector.Subscribe(subscriber)
	var alertsAll []saa.ScaleAlert
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		for pkg := range subscriber {
			alertsAll = append(alertsAll, pkg.ScaleAlerts...)
		}
		defer wg.Done()
	}()

	connector.HandleScaleAlerts(w, req, httprouter.Params{})
	resp := w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	close(subscriber)
	wg.Wait()
	assert.Equal(t, 1, len(alertsAll))
	assert.Equal(t, alertName, alertsAll[0].Name)
	assert.True(t, alertsAll[0].Firing)
	assert.True(t, alertsAll[0].StartedAt.Equal(startsAt))
}
