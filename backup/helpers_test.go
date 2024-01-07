package backup

import "testing"

func TestLastPathElement(t *testing.T) {
	// Given:
	path := "/some/path/to/last"
	expected := "last"

	// When:
	last := lastPathElement(path)

	// Then
	if last != expected {
		t.Errorf("Expecting %s but got %s", expected, last)
	}
}

func TestLastPathElementTrailingSlash(t *testing.T) {
	// Given:
	path := "/some/path/to/last/"
	expected := "last"

	// When:
	last := lastPathElement(path)

	// Then
	if last != expected {
		t.Errorf("Expecting %s but got %s", expected, last)
	}
}

func TestLastPathElementWithEmptyPath(t *testing.T) {
	// Given:
	path := "/"
	expected := "/"

	// When:
	last := lastPathElement(path)

	// Then
	if last != expected {
		t.Errorf("Expecting %s but got %s", expected, last)
	}
}
