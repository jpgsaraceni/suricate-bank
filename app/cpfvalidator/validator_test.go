package cpfvalidator

import (
	"testing"
)

func TestIsValid(t *testing.T) {
	t.Run("Global test, will be detailed later", func(t *testing.T) {
		t.Parallel()

		got, _ := IsValid("12345678901")
		want := false

		if got != want {
			t.Errorf("got %v expected %v", got, want)
		}
	})

	t.Run("Global test, will be detailed later", func(t *testing.T) {
		t.Parallel()

		got, _ := IsValid("045.591.180-00")
		want := true

		if got != want {
			t.Errorf("got %v expected %v", got, want)
		}
	})

}
