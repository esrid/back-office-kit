package ui

import "strings"

// Tone classes are written out in full, never concatenated, so Tailwind's
// content scanner finds them in this file. `"badge-" + tone` would compile
// but the CSS would be missing at runtime.
// tone is the single place deciding what a semantic state looks like, in
// every shape it is drawn: a pill, a banner, a trend figure, a timeline dot.
type tone struct {
	Badge string // StatusBadge
	Alert string // Flashes
	Text  string // StatRow change figures
	Dot   string // AuditTimeline markers
}

var tones = map[string]tone{
	"success": {Badge: "badge-success", Alert: "alert-success", Text: "text-success", Dot: "status-success"},
	"warning": {Badge: "badge-warning", Alert: "alert-warning", Text: "text-warning", Dot: "status-warning"},
	"error":   {Badge: "badge-error", Alert: "alert-error", Text: "text-error", Dot: "status-error"},
	"info":    {Badge: "badge-info", Alert: "alert-info", Text: "text-info", Dot: "status-info"},
	// Neutral carries no trend colour on purpose: an unremarkable number
	// should not be tinted.
	"neutral": {Badge: "badge-neutral", Alert: "alert-info", Text: "", Dot: "status-neutral"},
}

// statusTones maps a lowercased status to a tone key.
var statusTones = map[string]string{
	"active": "success", "actif": "success", "enabled": "success",
	"paid": "success", "payé": "success", "done": "success",
	"completed": "success", "terminé": "success", "online": "success",

	"pending": "warning", "en attente": "warning", "processing": "warning",
	"queued": "warning", "review": "warning", "trial": "warning",

	"failed": "error", "error": "error", "échec": "error",
	"rejected": "error", "canceled": "error", "cancelled": "error",
	"annulé": "error", "overdue": "error", "banned": "error",

	"draft": "neutral", "brouillon": "neutral", "inactive": "neutral",
	"disabled": "neutral", "archived": "neutral", "offline": "neutral",
}

// ToneFor returns the tone key for a status string.
// Pure; case- and whitespace-insensitive. Unknown statuses are neutral.
func ToneFor(status string) string {
	if t, ok := statusTones[strings.ToLower(strings.TrimSpace(status))]; ok {
		return t
	}
	return "neutral"
}

func toneOf(key string) tone {
	if t, ok := tones[key]; ok {
		return t
	}
	return tones["neutral"]
}
