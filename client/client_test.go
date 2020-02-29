package client

import (
	"reflect"
	"testing"

	"github.com/BurntSushi/xgb/xproto"
)

func TestNew(t *testing.T) {
	window := xproto.Window(50)
	t.Run("TypeNormal", func(t *testing.T) {
		x11 := &mockX11{t: t}
		cfg := &Config{}
		_, err := New(x11, cfg, window, TypeNormal)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if len(x11.reparentedWins) != 1 {
			t.Fatalf("expected the window to be reparented")
		}

		got := x11.reparentedWins[0]
		want := mockReparented{w: window, p: 1, x: 0, y: 0}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got = %v, want = %v", got, want)
		}
	})
}
