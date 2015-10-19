package memorycache

import (
	. "gopkg.in/check.v1"
	"testing"
)

func TestCreateKey(t *testing.T) {
	TestingT(t)
}

type CreateKeyTestsSuite struct{}

var _ = Suite(&CreateKeyTestsSuite{})

func (s *CreateKeyTestsSuite) Test_New(c *C) {
	//c.Skip("Not now")
	k := NewKey("asdasdas", "tag1")
	c.Assert(k, NotNil)
}

func (s *CreateKeyTestsSuite) Test_New_2(c *C) {
	//c.Skip("Not now")
	k := NewKey("asdasdas")
	c.Assert(k, NotNil)
}

func (s *CreateKeyTestsSuite) Test_New_3(c *C) {
	//c.Skip("Not now")
	k := NewKey("asdasdas", "tag1", "tag1", "tag1", "tag1")
	c.Assert(k, NotNil)
}
