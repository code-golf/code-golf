package pager

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPage(t *testing.T) {
	assert.Exactly(t, Page(1), 1)
	assert.Exactly(t, Page(13), 1)
	assert.Exactly(t, Page(PerPage), 1)
	assert.Exactly(t, Page(PerPage+1), 2)
}
