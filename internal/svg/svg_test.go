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

	svgOutput := AbcToSvg(abcInput)

	assert.True(t, strings.HasPrefix(`<svg xmlns="http://www.w3.org/2000/svg" version="1.1"
 xmlns:xlink="http://www.w3.org/1999/xlink"
 viewBox="0 0 794 146">
<svg xmlns="http://www.w3.org/2000/svg" version="1.1"
 xmlns:xlink="http://www.w3.org/1999/xlink"
 fill="currentColor" stroke-width=".7" class="f1 tune0"
 width="794px" height="76.00px" viewBox="0 0 794 76.00"
 y="0">
<style>
.f0{font:20.0px text,serif}
.f1{font:24.0px music}
@font-face{
 font-family:music;
 src:url("data:application/octet-stream;base64`, svgOutput))
}
