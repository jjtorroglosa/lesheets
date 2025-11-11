package svg

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSvg(t *testing.T) {
	abcInput := `
X:1
T:Example Tune
M:4/4
L:1/4
K:C
C D E F | G A B c |
`
	svgOutput, err := AbcToSvg("testSvg", abcInput)
	assert.NoError(t, err)

	assert.True(t, strings.HasPrefix(`<svg xmlns="http://www.w3.org/2000/svg" version="1.1"`, svgOutput[0:53]))
}
