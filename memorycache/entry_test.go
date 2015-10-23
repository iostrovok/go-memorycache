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
	b := CreateEntry(k, "", Nothing, []string{}, 10*time.Millisecond, map[string]Press{})
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
	b := CreateEntry(k, "", Nothing, []string{"seller", "my_tag"}, time.Duration(0), map[string]Press{})
	c.Assert(b, NotNil)

	c.Check(len(b.Tags), Equals, 2)
	c.Check(b.CheckTime, Equals, false)
	c.Check(b.Valid(), Equals, true)
}

func myPress(interface{}) (interface{}, error) {
	return `qazwsx`, nil
}
func myPress2(interface{}) (interface{}, error) {
	return `qazwsx2`, nil
}

func (s *CreateEntryTestsSuite) Test_CreateEntry_Press(c *C) {
	//c.Skip("Not now")
	k := NewKey("asdasd")

	PF := map[string]Press{
		"": myPress,
	}

	b := CreateEntry(k, "", Nothing, []string{}, 10*time.Millisecond, PF)
	c.Check(len(b.Tags), Equals, 0)
	c.Check(b.CheckTime, Equals, true)
	c.Assert(b, NotNil)

	c.Check(b.Valid(), Equals, true)

	time.Sleep(11 * time.Millisecond)

	c.Check(b.Valid(), Equals, false)

	c.Check(b.Data, Equals, `qazwsx`)
}

func (s *CreateEntryTestsSuite) Test_CreateEntry_2_Press(c *C) {
	//c.Skip("Not now")
	k := NewKey("asdasd")
	PF := map[string]Press{
		"": myPress,
	}
	b := CreateEntry(k, "", Nothing, []string{"seller", "my_tag"}, time.Duration(0), PF)
	c.Assert(b, NotNil)

	c.Check(len(b.Tags), Equals, 2)
	c.Check(b.CheckTime, Equals, false)
	c.Check(b.Valid(), Equals, true)

	c.Check(b.Data, Equals, `qazwsx`)
}

func (s *CreateEntryTestsSuite) Test_CreateEntry_3_Press(c *C) {
	//c.Skip("Not now")
	k := NewKey("asdasd")
	PF := map[string]Press{
		"seller": myPress,
		"my_tag": myPress2,
	}
	b := CreateEntry(k, "", Nothing, []string{"seller", "my_tag"}, time.Duration(0), PF)
	c.Assert(b, NotNil)

	c.Check(len(b.Tags), Equals, 2)
	c.Check(b.CheckTime, Equals, false)
	c.Check(b.Valid(), Equals, true)

	c.Check(b.Data, Equals, `qazwsx`)
}
