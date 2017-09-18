// Copyright 2017 John Scherff
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	`log`
	`io`
	`io/ioutil`
	`os`
	`strings`
	`github.com/RackSec/srslog`
	`github.com/jscherff/goutil`
)

func NewLoggers() (sl, cl, el *log.Logger) {

	var sw, cw, ew []io.Writer

	var newfl = func(f string) (h *os.File, err error) {

		if h, err = os.OpenFile(f, FileFlags, FileMode); err != nil {
			log.Printf(`%v`, goutil.ErrorDecorator(err))
		}

		return h, err
	}

	var newsl = func(prot, raddr, tag string, pri srslog.Priority) (s *srslog.Writer, err error) {

		if s, err = srslog.Dial(prot, raddr, pri, tag); err != nil {
			log.Printf(`%v`, goutil.ErrorDecorator(err))
		}

		return s, err
	}

	if conf.Logging.LogFiles {

		if f, err := newfl(conf.Files.SystemLog); err == nil {
			sw = append(sw, f)
		}
		if f, err := newfl(conf.Files.ChangeLog); err == nil {
			cw = append(cw, f)
		}
		if f, err := newfl(conf.Files.ErrorLog); err == nil {
			ew = append(ew, f)
		}
	}

	if conf.Logging.Console {
		sw = append(sw, os.Stdout)
		cw = append(cw, os.Stdout)
		ew = append(ew, os.Stderr)
	}

	if conf.Logging.Syslog {

		port, prot, addr := conf.Syslog.Port, conf.Syslog.Protocol, conf.Syslog.Address
		raddr := strings.Join([]string{addr, port}, `:`)

		if s, err := newsl(prot, raddr, `gocmdbcli`, PriInfo); err == nil {
			sw = append(sw, s)
		}
		if s, err := newsl(prot, raddr, `gocmdbcli`, PriInfo); err == nil {
			cw = append(cw, s)
		}
		if s, err := newsl(prot, raddr, `gocmdbcli`, PriErr); err == nil {
			ew = append(ew, s)
		}
	}

	if len(sw) == 0 {
		sw = append(sw, ioutil.Discard)
	}
	if len(cw) == 0 {
		cw = append(cw, ioutil.Discard)
	}
	if len(ew) == 0 {
		ew = append(ew, ioutil.Discard)
	}

	//sl = log.New(io.MultiWriter(sw...), `system: `, log.LstdFlags)
	//cl = log.New(io.MultiWriter(cw...), `change: `, log.LstdFlags)
	//el = log.New(io.MultiWriter(ew...), `error:  `, log.LstdFlags)
	sl = log.New(io.MultiWriter(sw...), ``, 0)
	cl = log.New(io.MultiWriter(cw...), ``, 0)
	el = log.New(io.MultiWriter(ew...), ``, 0)

	return sl, cl, el
}
