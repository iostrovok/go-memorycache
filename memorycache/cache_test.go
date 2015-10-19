package memorycache

import (
	"fmt"
	. "gopkg.in/check.v1"
	"testing"
	"time"
)

func TestCreateCache(t *testing.T) {
	TestingT(t)
}

type CreateCacheTestsSuite struct{}

var _ = Suite(&CreateCacheTestsSuite{})

func (s *CreateCacheTestsSuite) Test_New(c *C) {
	//c.Skip("Not now")
	st := New(5, 100)
	c.Assert(st, NotNil)

	c.Check(len(st.Shards), Equals, 5)
}

func (s *CreateCacheTestsSuite) Test_Close(c *C) {
	//c.Skip("Not now")
	st := New(5, 100)
	c.Assert(st, NotNil)
	st.Close()
}

func (s *CreateCacheTestsSuite) Test_Put(c *C) {
	//c.Skip("Not now")
	st := New(5, 100)
	st.Put(100, "==> key")
	st.Close()
}

func (s *CreateCacheTestsSuite) Test_Put_10000(c *C) {
	//c.Skip("Not now")
	st := New(5, 100)
	for i := 0; i < 10000; i++ {
		st.Put(i, fmt.Sprintf("==> %d", i))
	}

	st.Close()
}

func (s *CreateCacheTestsSuite) Test_Put_Over(c *C) {
	//c.Skip("Not now")

	total := 100
	shards := 1
	perShard := 1 + total/shards

	st := New(1, 100)

	for i := 0; i < 10000; i++ {
		st.Put(i, fmt.Sprintf("==> %d", i))
	}

	time.Sleep(time.Millisecond * 1000)

	c.Check(len(st.Shards[0].entries), Equals, perShard)
	st.Close()
}

func (s *CreateCacheTestsSuite) Test_Get_1(c *C) {
	//c.Skip("Not now")
	st := New(1, 100)
	c.Assert(st, NotNil)

	c.Check(len(st.Shards), Equals, 1)

	st.Put("some data", "key")
	res, ok := st.Get("key")

	c.Check(ok, Equals, true)
	c.Check(res, Equals, "some data")

	st.Close()
}
