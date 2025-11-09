//go:build js && wasm
// +build js,wasm

package svg

func RenderAbcToSvg(file string, abcInput string) (string, error) {
	return "svg from wasm", nil
}
