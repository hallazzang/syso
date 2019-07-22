package rsrc

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/hallazzang/syso/pkg/ico"
)

func TestAddIconByName(t *testing.T) {
	f, err := os.Open(filepath.Join("..", "..", "testdata", "golang.ico"))
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	icons, err := ico.DecodeAll(f)
	if err != nil {
		t.Fatal(err)
	}

	r := New()
	if err := r.AddResourceByName(IconResource, "ICON", icons.Images[0]); err != nil {
		t.Fatal(err)
	}

	b := new(bytes.Buffer)
	n, err := r.WriteTo(b)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("wrote %d bytes", n)
}
