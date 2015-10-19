package memorycache

import (
	. "gopkg.in/check.v1"
	"testing"
)

func TestCreateMessage(t *testing.T) {
	TestingT(t)
}

type CreateMessageTestsSuite struct{}

var _ = Suite(&CreateMessageTestsSuite{})

func (s *CreateMessageTestsSuite) Test_Res(c *C) {
	//c.Skip("Not now")
	b := Res{}
	c.Assert(b, NotNil)
}

func (s *CreateMessageTestsSuite) Test_Request(c *C) {
	//c.Skip("Not now")
	b := Request{}
	c.Assert(b, NotNil)
}
