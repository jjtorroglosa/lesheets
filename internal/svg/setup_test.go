package svg

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	LoadJsRuntime(os.DirFS("./../../"))
	os.Exit(m.Run())
}
