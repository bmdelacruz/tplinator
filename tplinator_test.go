package tplinator_test

import (
	"strings"
	"testing"

	"github.com/bmdelacruz/tplinator"
)

func TestTplinate(t *testing.T) {
	t.Run(`ok`, func(t *testing.T) {
		_, err := tplinator.Tplinate(
			strings.NewReader(`<div></div>`),
		)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
	t.Run(`improperly closed tag`, func(t *testing.T) {
		_, err := tplinator.Tplinate(
			strings.NewReader(`<div></div`),
		)
		if err == nil {
			t.Error("expecting an error because" +
				" the tag was improperly closed")
		}
	})
}
