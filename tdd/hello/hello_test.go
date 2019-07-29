package main

import "testing"

func TestHello(t *testing.T) {
	assertCorrectMessage := func(t *testing.T, got, want string) {
		t.Helper()
		if got != want {
			t.Errorf("Error testing Hello function, got: %q, want: %q", got, want)
		}
	}
	t.Run("Classic hello world", func(t *testing.T) {
		got := Hello("", "")
		want := "Hello, World"
		assertCorrectMessage(t, got, want)
	})

	t.Run("Hello with a name", func(t *testing.T) {
		got := Hello("John", "")
		want := "Hello, John"
		assertCorrectMessage(t, got, want)
	})

	t.Run("Hello in spanish with no name", func(t *testing.T) {
		got := Hello("", "in Spanish")
		want := "Hola, World"
		assertCorrectMessage(t, got, want)
	})

	t.Run("Hello in spanish with name", func(t *testing.T) {
		got := Hello("Rosita", "in Spanish")
		want := "Hola, Rosita"
		assertCorrectMessage(t, got, want)
	})

	t.Run("Hello in french with name", func(t *testing.T) {
		got := Hello("Alexia", "in French")
		want := "Bonjour, Alexia"
		assertCorrectMessage(t, got, want)
	})

}
