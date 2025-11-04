package svg

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSvgSnapshot(t *testing.T) {
	result, err := RenderAbcToSvg("anything", "X:1\nT:Title\nM:4/4\nL:1/4\nK:C cleff=perc stafflines=0\nP:A\nA4 | A4")
	assert.NoError(t, err)

	bytes, err := os.ReadFile("testdata/expected.svg")
	assert.NoError(t, err)
	assert.Equal(t, result, strings.TrimSpace(string(bytes)))
}
