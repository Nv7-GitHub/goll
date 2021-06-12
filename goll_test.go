package goll

import (
	_ "embed"
	"os"
	"testing"
)

//go:embed test.goll
var src string

func TestCompileSrc(t *testing.T) {
	out, err := CompileSrc(src)
	if err != nil {
		t.Fatal(err)
	}

	err = os.WriteFile("test.ll", []byte(out.String()), os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}
}
