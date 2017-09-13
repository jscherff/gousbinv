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
	"encoding/json"
	"path/filepath"
	"os"
)

// The filename of the JSON configuration file.
const configFile string = "config.json"

// Config holds the application configuration settings. The struct tags
// must match the field names in the JSON configuration file.
type Config struct {

	AppPath		string

	LogDir		string
	AuditDir	string
	StateDir	string
	ReportDir	string

	AppLog		string
	ChangeLog	string
	LegacyReport	string

	ServerURL	string
	AuditPath	string
	ChangesPath	string
	CheckinPath	string
	FetchSnPath	string

	IncludeVID	map[string]bool
	IncludePID	map[string]map[string]bool

	DefaultInclude	bool
	DefaultFormat	string
}

// GetConfig retrieves the settings in the JSON configuration file and
// populates the fields in the runtime configuration. It also creates
// subdirectories in the application path if they do not exist.
func getConfig() (c *Config, err error) {

	c = new(Config)

	ep := filepath.Dir(os.Args[0])
	fp := filepath.Join(ep, configFile)

	fh, err := os.Open(fp)
	defer fh.Close()

	// Decode JSON from configuration file.

	if err == nil {
		jd := json.NewDecoder(fh)
		err = jd.Decode(&c)
	}

	// If app path is empty, set it to executable path.

	if err == nil {
		if len(c.AppPath) == 0 {
			c.AppPath = ep
		}
	}

	// Configure and create log directory.

	if err == nil {

		d, sd := filepath.Split(c.LogDir)

		if len(d) == 0 {
			c.LogDir = filepath.Join(c.AppPath, sd)
		}

		err = os.MkdirAll(c.LogDir, 0755)
	}

	// Configure and create audit directory.

	if err == nil {

		d, sd := filepath.Split(c.AuditDir)

		if len(d) == 0 {
			c.LogDir = filepath.Join(c.AppPath, sd)
		}

		err = os.MkdirAll(c.AuditDir, 0755)
	}

	// Configure and create report directory.

	if err == nil {

		d, sd := filepath.Split(c.ReportDir)

		if len(d) == 0 {
			c.LogDir = filepath.Join(c.AppPath, sd)
		}

		err = os.MkdirAll(c.ReportDir, 0755)
	}

	return c, err
}
