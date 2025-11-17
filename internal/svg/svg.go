package svg

import (
	"errors"
	"lesheets/internal/logger"
)

func AbcToHtml(sourceFile string, defaultLength string, abcInput string) (string, error) {
	abc := `
%%topspace 0
%%musicfont
%%pagewidth 700px
%%scale 1.1
%%topmargin      0px
%%botmargin      0px
%%leftmargin     0px
%%rightmargin    0px
%%titlespace     0px
` + abcInput
	svg, err := AbcToSvg(sourceFile, abc)
	if err != nil {
		return "", err
	}

	return svg, nil
}

func InlineAbcToHtml(sourceFile string, defaultLength string, abcInput string) (string, error) {
	abc := `
%%topspace 0
%%musicfont
%%pagewidth 300px
%%scale 1.1
%%topmargin      0px
%%botmargin      0px
%%leftmargin     20px
%%rightmargin    0px
%%titlespace     0px
%%map all2A * print=G
X:1
M:none
L:` + defaultLength + `
K:none clef=none stafflines=0 stem=up
%%voicemap all2A
` + abcInput
	svg, err := AbcToSvg(sourceFile, abc)
	if err != nil {
		return "", err
	}

	return svg, nil
}

func AbcToSvg(sourceFile string, abcInput string) (string, error) {
	defer logger.LogElapsedTime("RenderSvg")()
	res, err := RenderAbcToSvg(sourceFile, abcInput)
	if err != nil {
		return "", errors.New("error rendering abc to svg: " + err.Error())
	}
	return res, nil
}
