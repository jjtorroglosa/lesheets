//go:build !js || !wasm
// +build !js !wasm

package svg

import (
	"crypto/sha256"
	"embed"
	"encoding/hex"
	"lesheets/internal/timer"
	"log"
	"sync"

	"github.com/fastschema/qjs"
	lru "github.com/hashicorp/golang-lru/v2"
)

//go:embed abc2svg/user.js abc2svg/tosvg.js abc2svg/abc2svg-1.js
var abc2svg embed.FS

func loadFile(ctx *qjs.Context, filename string) (*qjs.Value, func()) {
	filepath := "abc2svg/" + filename

	code, err := abc2svg.ReadFile(filepath)
	if err != nil {
		panic(err)
	}
	result, err := ctx.Eval(filename, qjs.Code(string(code)))
	if err != nil {
		log.Fatal(err)
	}
	return result, func() {
		result.Free()
	}
}

var renderAbcToSvg func(string, string) (string, error)

// cache entry
type renderResult struct {
	svg string
	err error
}

var (
	renderCache, _ = lru.New[string, renderResult](128)
	mu             sync.Mutex
)

func makeKey(file, data string) string {
	sum := sha256.Sum256([]byte(file + data))
	return hex.EncodeToString(sum[:])
}

func RenderAbcToSvg(file, data string) (string, error) {
	key := makeKey(file, data)

	if res, ok := renderCache.Get(key); ok {
		return res.svg, res.err
	}

	svg, err := renderAbcToSvg(file, data)

	mu.Lock()
	renderCache.Add(key, renderResult{svg, err})
	mu.Unlock()

	return svg, err
}

func init() {
	loadJsRuntime()
}

func loadJsRuntime() func() {
	defer timer.LogElapsedTime("LoadQjs")()
	rt, err := qjs.New()
	if err != nil {
		log.Fatal(err)
	}

	ctx := rt.Context()

	for _, f := range []string{"user.js", "abc2svg-1.js"} {
		_, cleanup := loadFile(ctx, f)
		defer cleanup()
	}
	result, _ := loadFile(ctx, "tosvg.js")

	if err != nil {
		log.Fatal("Eval error:", err)
	}
	jsRenderFunction := result.GetPropertyStr("tosvg")
	goRenderFunc, err := qjs.JsFuncToGo[func(string, string) (string, error)](jsRenderFunction)
	if err != nil {
		log.Fatal("Func conversion error:", err)
	}
	renderAbcToSvg = goRenderFunc
	return func() {
		jsRenderFunction.Free()
		result.Free()
		rt.Close()
	}
}
