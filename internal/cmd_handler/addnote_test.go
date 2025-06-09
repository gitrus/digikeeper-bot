package cmdhandler_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/gitrus/digikeeper-bot/internal/cmd_handler"
	"github.com/gitrus/digikeeper-bot/internal/note"
)

func TestHandleAddNote_ReturnsHandler(t *testing.T) {
	svc := note.NewInMemoryService()
	h := cmdhandler.HandleAddNote(svc)
	assert.NotNil(t, h)
}
