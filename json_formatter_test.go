package debug

import (
	"encoding/json"
	"errors"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrorNotLost(t *testing.T) {
	SetFormatter(&JSONFormatter{})

	s := formatter.Format(Debug("error_not_lost").WithField("error", errors.New("wild walrus")), "hi")

	entry := make(map[string]interface{})
	err := json.Unmarshal([]byte(s), &entry)

	assert.Nil(t, err, "Unable to unmarshal formatted entry")

	assert.Equal(t, entry["error"], "wild walrus")
	assert.Equal(t, entry["msg"], "hi")
	assert.NotEmpty(t, entry["delta"])
	assert.NotEmpty(t, entry["time"])
	assert.NotEmpty(t, entry["namespace"])
}

func TestErrorNotLostOnFieldNotNamedError(t *testing.T) {
	SetFormatter(&JSONFormatter{})

	s := formatter.Format(Debug("mapped_field_error").WithField("omg", errors.New("wild walrus")), "hi")

	entry := make(map[string]interface{})
	err := json.Unmarshal([]byte(s), &entry)

	assert.Nil(t, err, "Unable to unmarshal formatted entry")

	assert.Equal(t, entry["omg"], "wild walrus")
	assert.Equal(t, entry["msg"], "hi")
	assert.NotEmpty(t, entry["delta"])
	assert.NotEmpty(t, entry["time"])
	assert.NotEmpty(t, entry["namespace"])
}

func TestFieldClashWithTime(t *testing.T) {
	SetFormatter(&JSONFormatter{})

	s := formatter.Format(Debug("clash_time").WithField("time", "right now!"), "hi")

	entry := make(map[string]interface{})
	err := json.Unmarshal([]byte(s), &entry)

	assert.Nil(t, err, "Unable to unmarshal formatted entry")
	assert.NotEqual(t, entry["fields.time"], "right now!", "fields.time not set to original time field")

	r := regexp.MustCompile(`\d\d:\d\d:\d\d\.(\d*)`)
	ss := (entry["time"]).(string)
	assert.Equal(t, len(r.FindStringSubmatch(ss)), 2, "time check")
}

func TestFieldClashWithMsg(t *testing.T) {
	SetFormatter(&JSONFormatter{})

	s := formatter.Format(Debug("clash_msg").WithField("msg", errors.New("wild walrus")), "hi")

	entry := make(map[string]interface{})
	err := json.Unmarshal([]byte(s), &entry)

	assert.Nil(t, err, "Unable to unmarshal formatted entry")
	assert.Equal(t, entry["msg"], "hi")
}

func TestFieldClashWithNamespace(t *testing.T) {
	SetFormatter(&JSONFormatter{})

	s := formatter.Format(Debug("clash_namespace").WithField("namespace", errors.New("wild walrus")), "hi")

	entry := make(map[string]interface{})
	err := json.Unmarshal([]byte(s), &entry)

	assert.Nil(t, err, "Unable to unmarshal formatted entry")
	assert.Equal(t, entry["namespace"], "clash_namespace", "namespace is correct")
}

func TestJSONEntryEndsWithNewline(t *testing.T) {
	SetFormatter(&JSONFormatter{})

	s := formatter.Format(Debug("newline").WithField("dog", errors.New("wild walrus")), "hi")

	entry := make(map[string]interface{})
	err := json.Unmarshal([]byte(s), &entry)

	assert.Nil(t, err, "Unable to unmarshal formatted entry")
	assert.Equal(t, "\n", string(s[len(s)-1]), "Expected JSON log entry to end with a newline")
}

func TestJSONPretty(t *testing.T) {
	HAS_TIME = false
	SetFormatter(&JSONFormatter{PrettyPrint: true})

	s := formatter.Format(Debug("pretty").WithField("dog", "wild walrus"), "hi")
	expected := "{\n  \"dog\": \"wild walrus\",\n  \"msg\": \"hi\",\n  \"namespace\": \"pretty\"\n}\n"

	assert.Equal(t, expected, s, "is pretty")
}
