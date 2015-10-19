package memorycache

import (
	. "gopkg.in/check.v1"
	"testing"
)

func TestCreateEntry(t *testing.T) {
	TestingT(t)
}

type CreateEntryTestsSuite struct{}

var _ = Suite(&CreateEntryTestsSuite{})

func (s *CreateEntryTestsSuite) Test_CreateEntry(c *C) {
	//c.Skip("Not now")
	k := NewKey("asdasd")
	b := CreateEntry(k, "", Nothing)
	c.Assert(b, NotNil)
}

func (s *CreateEntryTestsSuite) Test_CreateEntry_2(c *C) {
	//c.Skip("Not now")
	k := NewKey("asdasd")
	b := CreateEntry(k, "")
	c.Assert(b, NotNil)
}
