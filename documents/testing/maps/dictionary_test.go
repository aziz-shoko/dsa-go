package maps

import (
	"testing"
)

func TestSearch(t *testing.T) {
	dictionary := Dictionary{
		"test": "this is just a test",
	}
	t.Run("known word", func(t *testing.T) {
		got, _ := dictionary.Search("test")
		want := "this is just a test"
		assertStrings(t, got, want)
	})

	t.Run("unknown word", func(t *testing.T) {
		// "unknown" is a nonexistent word
		_, err := dictionary.Search("unknown")
		if err == nil {
			t.Fatalf("expected an error but got no errors")
		}
		assertErrors(t, err, ErrNotFound)
	})
}

func TestAdd(t *testing.T) {
	t.Run("add word", func(t *testing.T) {
		dictionary := Dictionary{}
		word := "example"
		definition := "definition of example"

		err := dictionary.Add(word, definition)

		assertErrors(t, err, nil)
		assertDefinition(t, dictionary, word, definition)
	})

	t.Run("existing word", func(t *testing.T) {
		word := "test"
		definition := "this is just a test"
		dictionary := Dictionary{word: definition}
		err := dictionary.Add(word, "new test")

		assertErrors(t, err, ErrWordExists)
		assertDefinition(t, dictionary, word, definition)
	})
}

func TestUpdate(t *testing.T) {
	word := "example"
	definition := "test definition for test"
	dictionary := Dictionary{word: definition}
	t.Run("test updating existing word", func(t *testing.T) {
		updatedDefinition := "definition of example"
		err := dictionary.Update(word, updatedDefinition)
		assertErrors(t, err, nil)
		assertDefinition(t, dictionary, word, updatedDefinition)
	})

	t.Run("test updating none existing word", func(t *testing.T) {
		word := "test"	
		updatedDefinition := "definition of test"

		err := dictionary.Update(word, updatedDefinition)
		assertErrors(t, err, ErrWordDoesNotExist)
	})
}

func TestDelete(t *testing.T) {
	word := "example"
	definition := "test definition for test"
	dictionary := Dictionary{word: definition}

	dictionary.Delete(word)
	_, err := dictionary.Search(word)
	assertErrors(t, err, ErrNotFound)
}

func assertStrings(t testing.TB, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}

func assertErrors(t testing.TB, got, want error) {
	t.Helper()

	if got != want {
		t.Errorf("got error %q want %q", got, want)
	}
}

func assertDefinition(t testing.TB, dictionary Dictionary, word, definition string) {
	t.Helper()

	got, err := dictionary.Search(word)
	if err != nil {
		t.Fatal("Should find added word:", err)
	}
	assertStrings(t, got, definition)
}
