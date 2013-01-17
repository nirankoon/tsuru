// Copyright 2013 tsuru authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package juju

import (
	"github.com/globocom/commandmocker"
	"github.com/globocom/tsuru/heal"
	. "launchpad.net/gocheck"
)

func (s *S) TestBootstrapShouldBeRegistered(c *C) {
	h, err := heal.Get("bootstrap")
	c.Assert(err, IsNil)
	c.Assert(h, FitsTypeOf, &BootstrapHealer{})
}

func (s *S) TestBootstrapNeedsHeal(c *C) {
	tmpdir, err := commandmocker.Add("juju", collectOutputBootstrapDown)
	c.Assert(err, IsNil)
	defer commandmocker.Remove(tmpdir)
	h := BootstrapHealer{}
	c.Assert(h.NeedsHeal(), Equals, true)
}

func (s *S) TestBootstrapDontNeedsHeal(c *C) {
	tmpdir, err := commandmocker.Add("juju", collectOutput)
	c.Assert(err, IsNil)
	defer commandmocker.Remove(tmpdir)
	h := BootstrapHealer{}
	c.Assert(h.NeedsHeal(), Equals, false)
}

func (s *S) TestBootstrapHeal(c *C) {
	jujuTmpdir, err := commandmocker.Add("juju", collectOutputBootstrapDown)
	c.Assert(err, IsNil)
	defer commandmocker.Remove(jujuTmpdir)
	sshTmpdir, err := commandmocker.Add("ssh", "$*")
	c.Assert(err, IsNil)
	defer commandmocker.Remove(sshTmpdir)
	jujuOutput := []string{
		"status", // for verify if heal is need
		"status", // for juju status that gets the output
	}
	sshOutput := []string{
		"-o",
		"StrictHostKeyChecking no",
		"-q",
		"-l",
		"ubuntu",
		"10.10.10.96",
		"sudo",
		"restart",
		"juju-machine-agent",
	}
	h := BootstrapHealer{}
	err = h.Heal()
	c.Assert(err, IsNil)
	c.Assert(commandmocker.Ran(jujuTmpdir), Equals, true)
	c.Assert(commandmocker.Parameters(jujuTmpdir), DeepEquals, jujuOutput)
	c.Assert(commandmocker.Ran(sshTmpdir), Equals, true)
	c.Assert(commandmocker.Parameters(sshTmpdir), DeepEquals, sshOutput)
}

func (s *S) TestBootstrapOnlyHealsWhenItIsNeeded(c *C) {
	tmpdir, err := commandmocker.Add("juju", collectOutput)
	c.Assert(err, IsNil)
	defer commandmocker.Remove(tmpdir)
	cmdOutput := []string{
		"status", // for verify if heal is need
	}
	h := BootstrapHealer{}
	err = h.Heal()
	c.Assert(err, IsNil)
	c.Assert(commandmocker.Ran(tmpdir), Equals, true)
	c.Assert(commandmocker.Parameters(tmpdir), DeepEquals, cmdOutput)
}
