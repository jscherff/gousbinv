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

import `flag`

var (
	fsAction = flag.NewFlagSet("action", flag.ExitOnError)
	fActionAudit = fsAction.Bool("audit", false, "Audit devices")
	fActionCheckin = fsAction.Bool("checkin", false, "Check devices in")
	fActionReport = fsAction.Bool("report", false, "Report actions")
	fActionReset = fsAction.Bool("reset", false, "Reset device")
	fActionSerial = fsAction.Bool("serial", false, "Set serial number")
	fActionState = fsAction.Bool("state", false, "Show device state")
	fActionVersion = fsAction.Bool("version", false, "Display version")

	fsReport = flag.NewFlagSet("report", flag.ExitOnError)
	fReportFolder = fsReport.String("folder", "", "Write reports to `<path>`")
	fReportFormat = fsReport.String("format", "json", "Report `<format>` {csv|nvp|xml|json}")
	fReportConsole = fsReport.Bool("console", false, "Write reports to console")

	fsSerial = flag.NewFlagSet("serial", flag.ExitOnError)
	fSerialDefault = fsSerial.Bool("default", false, "Set serial number to default")
	fSerialErase = fsSerial.Bool("erase", false, "Erase current serial number")
	fSerialForce = fsSerial.Bool("force", false, "Force serial number change")
	fSerialFetch = fsSerial.Bool("fetch", false, "Fetch serial number from server")
	fSerialSet = fsSerial.String("set", "", "Set serial number to `<string>`")
)
