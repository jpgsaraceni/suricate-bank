package cpfvalidator

import (
	"testing"
)

func TestIsValid(t *testing.T) {
	t.Run("Global invalid cpf test, will be detailed later", func(t *testing.T) {
		t.Parallel()
		const invalidCpf Cpf = "12345678901"

		got, _ := invalidCpf.IsValid()
		want := false

		if got != want {
			t.Errorf("got %v expected %v", got, want)
		}
	})

	t.Run("Global valid cpf test, will be detailed later", func(t *testing.T) {
		t.Parallel()
		const validCpf Cpf = "045.591.180-00"

		got, _ := validCpf.IsValid()
		want := true

		if got != want {
			t.Errorf("got %v expected %v", got, want)
		}
	})

}
