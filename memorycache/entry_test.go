package memorycache

import (
	. "gopkg.in/check.v1"
	"testing"
	"time"
)

func TestCreateEntry(t *testing.T) {
	TestingT(t)
}

type CreateEntryTestsSuite struct{}

var _ = Suite(&CreateEntryTestsSuite{})

func (s *CreateEntryTestsSuite) Test_CreateEntry(c *C) {
	//c.Skip("Not now")
	k := NewKey("asdasd")
	b := CreateEntry(k, "", Nothing, []string{}, 10*time.Millisecond)
	c.Check(len(b.Tags), Equals, 0)
	c.Check(b.CheckTime, Equals, true)
	c.Assert(b, NotNil)

	c.Check(b.Valid(), Equals, true)

	time.Sleep(11 * time.Millisecond)

	c.Check(b.Valid(), Equals, false)

}

func (s *CreateEntryTestsSuite) Test_CreateEntry_2(c *C) {
	//c.Skip("Not now")
	k := NewKey("asdasd")
	b := CreateEntry(k, "", Nothing, []string{"seller", "my_tag"}, time.Duration(0))
	c.Assert(b, NotNil)

	c.Check(len(b.Tags), Equals, 2)
	c.Check(b.CheckTime, Equals, false)
	c.Check(b.Valid(), Equals, true)
}
