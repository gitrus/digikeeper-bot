package noteparser_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/gitrus/digikeeper-bot/internal/note"
	"github.com/gitrus/digikeeper-bot/pkg/noteparser"
)

func TestParseDateTimeAndTags(t *testing.T) {
	createdAt := time.Date(2025, time.January, 1, 10, 0, 0, 0, time.Local)
	input := "Meeting with team 2025-06-01 09:30 #work"

	n, err := noteparser.Parse(createdAt, input)
	assert.NoError(t, err)
	assert.Equal(t, createdAt, n.CreatedAt)
	assert.Equal(t, []string{"work"}, n.Tags)
	expected := time.Date(2025, 6, 1, 9, 30, 0, 0, time.Local)
	assert.Equal(t, expected, n.Payload.EventAt)
	assert.Equal(t, "Meeting with team", n.Payload.Text)
}

func TestParseTimeOnly(t *testing.T) {
	createdAt := time.Date(2025, time.May, 5, 8, 0, 0, 0, time.Local)
	input := "Call mom 15:04 #family"

	n, err := noteparser.Parse(createdAt, input)
	assert.NoError(t, err)
	assert.Equal(t, []string{"family"}, n.Tags)
	expected := time.Date(2025, 5, 5, 15, 4, 0, 0, time.Local)
	assert.Equal(t, expected, n.Payload.EventAt)
	assert.Equal(t, "Call mom", n.Payload.Text)
}

func TestParseCyrillicTag(t *testing.T) {
	createdAt := time.Date(2025, time.July, 7, 12, 0, 0, 0, time.Local)
	input := "Обед #обед"

	n, err := noteparser.Parse(createdAt, input)
	assert.NoError(t, err)
	assert.Equal(t, []string{"обед"}, n.Tags)
	assert.Equal(t, createdAt, n.Payload.EventAt)
	assert.Equal(t, "Обед", n.Payload.Text)
}
