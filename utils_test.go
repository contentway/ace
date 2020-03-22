package ace

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConcat(t *testing.T) {
	a := assert.New(t)

	a.Equal("a", concat("a"))
	a.Equal("ab", concat("a", "b"))
	a.Equal("abc", concat("a", "b", "c"))
}

func TestFileExists(t *testing.T) {
	assert.Equal(t, FileExists("/etc/hosts"), true, "/etc/hosts must exist")
	assert.Equal(t, FileExists("/etc"), false, "/etc is a folder and not a file")
	assert.Equal(t, FileExists("/paththatdoesntexist"), false, "This path shouldn't exist")
}
