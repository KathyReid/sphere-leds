package agent

import . "launchpad.net/gocheck"

type LoadAgentSuite struct {
	agent *Agent
}

var _ = Suite(&LoadAgentSuite{})

func (s *LoadAgentSuite) SetUpTest(c *C) {

	conf := &Config{}

	s.agent = createAgent(conf)
}

func (s *LoadAgentSuite) TestConfig(c *C) {

}
