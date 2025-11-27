package blogposts_test

import (
	"errors"
	blogposts "github.com/aziz-shoko/dsa-go/blogposts"
	"io/fs"
	"reflect"
	"sort"
	"testing"
	"testing/fstest"
)

func TestNewBlogPosts(t *testing.T) {
	const (
		firstBody = `Title: Post 1
Description: Description 1
Tags: tdd, go
---
Hello
World`
		secondBody = `Title: Post 2
Description: Description 2
Tags: rust, borrow-checker
---
B
L
M`
	)

	fs := fstest.MapFS{
		"hello_word.md":   {Data: []byte(firstBody)},
		"hello-world2.md": {Data: []byte(secondBody)},
	}

	posts, err := blogposts.NewPostFromFS(fs)
	if err != nil {
		t.Fatal(err)
	}

	if len(posts) != 2 {
		t.Errorf("got %d posts, wanted %d posts", len(posts), 2)
	}

	sort.Slice(posts, func(i, j int) bool {
		return posts[i].Title < posts[j].Title
	})

	assertPost(t, posts[0], blogposts.Post{
		Title:       "Post 1",
		Description: "Description 1",
		Tags:        []string{"tdd", "go"},
		Body: `Hello
World`,
	})
}

func assertPost(t *testing.T, got blogposts.Post, want blogposts.Post) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got post %+v, wanted %+v", got, want)
	}
}

func TestNewBlogPosts_ErrorHandling(t *testing.T) {
	_, err := blogposts.NewPostFromFS(StubFailingFS{})

	if err == nil {
		t.Error("expected an error but didn't get one")
	}
}

type StubFailingFS struct{}

func (s StubFailingFS) Open(name string) (fs.File, error) {
	return nil, errors.New("oh no, i always fail")
}
