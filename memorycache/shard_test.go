package memorycache

import (
	. "gopkg.in/check.v1"
	"testing"
)

func TestCreateShard(t *testing.T) {
	TestingT(t)
}

type CreateShardTestsSuite struct{}

var _ = Suite(&CreateShardTestsSuite{})

func (s *CreateShardTestsSuite) Test_New(c *C) {
	//c.Skip("Not now")
	sh := NewShard(10)
	c.Assert(sh, NotNil)

	c.Check(len(sh.entries), Equals, 0)
	c.Check(sh.Len, Equals, 10)
}
