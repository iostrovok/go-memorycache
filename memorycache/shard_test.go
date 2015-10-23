package memorycache

import (
	"fmt"
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

func (s *CreateShardTestsSuite) Test_Put(c *C) {
	//c.Skip("Not now")
	sh := NewShard(10)
	c.Assert(sh, NotNil)
	mes := NewRequest(TypePut, NewKey("key"))
	sh.put(mes)

	c.Check(len(sh.entries), Equals, 1)
}

func (s *CreateShardTestsSuite) Test_Put_10(c *C) {
	//c.Skip("Not now")
	sh := NewShard(10)
	c.Assert(sh, NotNil)

	for i := 0; i < 10; i++ {
		mes := NewRequest(TypePut, NewKey(fmt.Sprintf("key:%d", i)))
		sh.put(mes)
	}

	c.Check(len(sh.entries), Equals, 10)
}

func (s *CreateShardTestsSuite) Test_setTTL(c *C) {
	//c.Skip("Not now")
	sh := NewShard(10)
	c.Assert(sh, NotNil)

	for i := 0; i < 10; i++ {
		mes := NewRequest(TypePut, NewKey(fmt.Sprintf("key:%d", i)))
		sh.put(mes)
	}

	c.Check(len(sh.entries), Equals, 10)
}

func (s *CreateShardTestsSuite) Test_trim(c *C) {
	//c.Skip("Not now")
	sh := NewShard(10)
	c.Assert(sh, NotNil)

	for i := 0; i < 10; i++ {
		mes := NewRequest(TypePut, NewKey(fmt.Sprintf("key:%d", i)))
		sh.put(mes)
	}

	c.Check(len(sh.entries), Equals, 10)
}
