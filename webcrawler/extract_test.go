package webcrawler

import "testing"

func TestNewAbsoluteURLNoTrailing(t *testing.T) {
	base := "http://test.com"
	link := "http://otherdomain.com"
	parsed := parseReturnedLink(base, link)
	if parsed != link {
		t.Errorf("expected %s, got %s", link, parsed)
	}
}

func TestNewAbsoluteURLTrailingSlash(t *testing.T) {
	base := "http://test.com"
	link := "http://otherdomain.com/"
	expected := "http://otherdomain.com"
	parsed := parseReturnedLink(base, link)
	if parsed != expected {
		t.Errorf("expected %s, got %s", expected, parsed)
	}
}

func TestNewRelativeURLNoTrailingSlash(t *testing.T) {
	base := "http://test.com"
	link := "/otherpage"
	expected := "http://test.com/otherpage"
	parsed := parseReturnedLink(base, link)
	if parsed != expected {
		t.Errorf("expected %s, got %s", expected, parsed)
	}
}

func TestNewRelativeURLTrailingSlash(t *testing.T) {
	base := "http://test.com"
	link := "/otherpage/"
	expected := "http://test.com/otherpage"
	parsed := parseReturnedLink(base, link)
	if parsed != expected {
		t.Errorf("expected %s, got %s", expected, parsed)
	}
}
