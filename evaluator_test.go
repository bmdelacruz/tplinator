package tplinator_test

import (
	"testing"

	"github.com/bmdelacruz/tplinator"
	"golang.org/x/net/html"
)

func TestTryEvaluateOnContext(t *testing.T) {
	type song struct {
		artistName string
		length     string
	}

	divNode := tplinator.CreateNode(html.ElementNode, "div", nil, false)
	divNode.SetContextParams(tplinator.EvaluatorParams{
		`isSong`:   true,
		`songName`: `Someday`,
		`songDetails`: song{
			artistName: `HalfNoise`,
			length:     `3:24`,
		},
	})

	deps := tplinator.NewDefaultExtensionDependencies()
	evaluator := deps.Get(tplinator.EvaluatorExtDepKey).(tplinator.Evaluator)

	t.Run(`bool`, func(t *testing.T) {
		hasIsSong, isSong, err := tplinator.TryEvaluateBoolOnContext(divNode, evaluator, `isSong`)
		if err != nil {
			t.Error("unexpected error:", err)
		} else if !hasIsSong || !isSong {
			t.Error("expecting `isSong` param with value of `true`")
		}
	})
	t.Run(`not bool`, func(t *testing.T) {
		_, _, err := tplinator.TryEvaluateBoolOnContext(divNode, evaluator, `songName`)
		if err == nil {
			t.Error("expecting an error because `songName` is not a boolean type")
		}
	})
	t.Run(`bool (invalid input string)`, func(t *testing.T) {
		_, _, err := tplinator.TryEvaluateBoolOnContext(divNode, evaluator, `isS--ong`)
		if err == nil {
			t.Error("expecting an error because input string is not a valid expression")
		}
	})
	t.Run(`bool (missing param)`, func(t *testing.T) {
		_, _, err := tplinator.TryEvaluateBoolOnContext(divNode, evaluator, `isSpecialSong`)
		if err == nil {
			t.Error("expecting an error because the param is not present on the evaluator params")
		}
	})
	t.Run(`string`, func(t *testing.T) {
		hasSongName, songName, err := tplinator.TryEvaluateStringOnContext(divNode, evaluator, `songName`)
		if err != nil {
			t.Error("unexpected error:", err)
		} else if !hasSongName || songName != `Someday` {
			t.Error("expecting `songName` param with value of `Someday`")
		}
	})
	t.Run(`not string`, func(t *testing.T) {
		_, _, err := tplinator.TryEvaluateStringOnContext(divNode, evaluator, `isSong`)
		if err == nil {
			t.Error("expecting an error because `isSong` is not a string type")
		}
	})
	t.Run(`string (invalid input string)`, func(t *testing.T) {
		_, _, err := tplinator.TryEvaluateStringOnContext(divNode, evaluator, `songNa-- me`)
		if err == nil {
			t.Error("expecting an error because input string is not a valid expression")
		}
	})
	t.Run(`string (missing param)`, func(t *testing.T) {
		_, _, err := tplinator.TryEvaluateStringOnContext(divNode, evaluator, `songLongName`)
		if err == nil {
			t.Error("expecting an error because the param is not present on the evaluator params")
		}
	})
	t.Run(`custom struct`, func(t *testing.T) {
		hasSongDetails, songDetails, err := tplinator.TryEvaluateOnContext(divNode, evaluator, `songDetails`)
		if err != nil {
			t.Error("unexpected error:", err)
			return
		} else if !hasSongDetails {
			t.Error("expecting `songDetails` param of type `song`")
			return
		}
		songDeets, isSongDeets := songDetails.(song)
		if !isSongDeets {
			t.Error("expecting `songDetails` param of type `song`")
		} else if songDeets.artistName != `HalfNoise` || songDeets.length != `3:24` {
			t.Error("songDeets has incorrect member values")
		}
	})
	t.Run(`custom struct (invalid input string)`, func(t *testing.T) {
		_, _, err := tplinator.TryEvaluateOnContext(divNode, evaluator, `songDe-t ails`)
		if err == nil {
			t.Error("expecting an error because input string is not a valid expression")
		}
	})
	t.Run(`custom struct (missing param)`, func(t *testing.T) {
		_, _, err := tplinator.TryEvaluateOnContext(divNode, evaluator, `songDeets`)
		if err == nil {
			t.Error("expecting an error because the param is not present on the evaluator params")
		}
	})
}
