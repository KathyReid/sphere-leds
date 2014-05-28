package agent

import (
	"strings"

	. "launchpad.net/gocheck"
)

type LoadBusSuite struct {
	bus        *Bus
	sampleJson string
}

var _ = Suite(&LoadBusSuite{})

func (s *LoadBusSuite) SetUpTest(c *C) {

	s.bus = &Bus{}

}

func (s *LoadBusSuite) TestEncode(c *C) {

}

func (s *LoadBusSuite) TestDecode(c *C) {

}

func trim(str string) string {
	return strings.Trim(str, "\n\r")
}
