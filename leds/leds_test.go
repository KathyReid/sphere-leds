package leds

import (
	"testing"

	. "launchpad.net/gocheck"
)

type LoadLedSuite struct {
	ledArray *LedArray
}

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

var _ = Suite(&LoadLedSuite{})

func (s *LoadLedSuite) SetUpTest(c *C) {
	s.ledArray = CreateLedArray()
}

func (s *LoadLedSuite) TestConfig(c *C) {

	c.Assert(ValidLedName("power"), Equals, true)
	c.Assert(LedNameIndex("power"), Equals, 0)
}
