package sntracker

import (
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/nikstoyanov/image-recognition/probability"
	sp "github.com/snowplow/snowplow-golang-tracker/v2/tracker"
)

// SetupTracker for a POST request on http.
func SetupTracker() *sp.Tracker {
	tracker := initTracker("0.0.0.0:9090", "image-app-id", "POST", "http")
	return tracker
}

// initTracker creates the emitter and tracker objects.
func initTracker(collector string, appid string, method string, protocol string) *sp.Tracker {
	// Create emitter and log success/failures.
	emitter := sp.InitEmitter(
		sp.RequireCollectorUri(collector),
		sp.OptionRequestType(method),
		sp.OptionProtocol(protocol),
		sp.OptionStorage(sp.InitStorageMemory()),
		sp.OptionCallback(func(s []sp.CallbackResult, f []sp.CallbackResult) {
			log.Println("Successes: " + strconv.Itoa(len(s)))
			log.Println("Failure: " + strconv.Itoa(len(f)))
		}),
	)

	subject := sp.InitSubject()

	tracker := sp.InitTracker(
		sp.RequireEmitter(emitter),
		sp.OptionSubject(subject),
		sp.OptionAppId(appid),
	)

	return tracker
}

// TrackNewImage creates a timing event when a new image classification is requested.
func TrackNewImage(tracker *sp.Tracker, n string) {
	tracker.TrackStructEvent(sp.StructuredEvent{
		Category:      sp.NewString("image"),
		Action:        sp.NewString("new-image"),
		Label:         sp.NewString(n),
		TrueTimestamp: sp.NewInt64(time.Now().Unix()),
	})
}

// TrackWrongExt stores the wrong extensions users input.
func TrackWrongExt(tracker *sp.Tracker, ext string) {
	tracker.TrackStructEvent(sp.StructuredEvent{
		Category:      sp.NewString("image"),
		Action:        sp.NewString("wrong-extension"),
		Label:         sp.NewString(ext),
		TrueTimestamp: sp.NewInt64(time.Now().Unix()),
	})
}

// TrackLabels collects the top 5 probability labels from an image classification.
func TrackLabels(tracker *sp.Tracker, prob []probability.LabelResult) {
	var labels strings.Builder

	for _, elem := range prob {
		labels.WriteString(elem.Label)
		labels.WriteString(",")
	}

	tracker.TrackStructEvent(sp.StructuredEvent{
		Category:      sp.NewString("image"),
		Action:        sp.NewString("label"),
		Label:         sp.NewString(labels.String()),
		TrueTimestamp: sp.NewInt64(time.Now().Unix()),
	})
}
