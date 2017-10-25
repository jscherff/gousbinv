// Copyright 2017 John Scherff
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	 http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	`fmt`
	`io/ioutil`
	`os`
	`path/filepath`
	`github.com/jscherff/cmdb/ci/peripheral/usb`
)

// audit performs a change audit against properties from the last checkin.
func audit(dev usb.Auditer) (err error) {

	var ch [][]string

	if dev.SN() == `` {
		sl.Printf(`device %s-%s skipping audit, no SN`, dev.VID(), dev.PID())
		return err
	}

	sl.Printf(`device %s-%s-%s fetching previous state from server`,
		dev.VID(), dev.PID(), dev.SN(),
	)

	var j []byte

	if j, err = usbCiCheckoutV1(dev); err != nil {
		return err
	}

	if ch, err = dev.CompareJSON(j); err != nil {
		return err
	}

	sl.Printf(`device %s-%s-%s saving current state to server`,
		dev.VID(), dev.PID(), dev.SN(),
	)

	if err := usbCiCheckinV1(dev); err != nil {
		el.Print(err)
	}

	if len(ch) == 0 {
		sl.Printf(`device %s-%s-%s detected no changes`,
			dev.VID(), dev.PID(), dev.SN(),
		)
		return nil
	}

	sl.Printf(`device %s-%s-%s recording changes in change log`,
		dev.VID(), dev.PID(), dev.SN(),
	)

	for _, c := range ch {
		cl.Printf(`device %s-%s-%s modified: %q was %q, now %q`,
			dev.VID(), dev.PID(), dev.SN(), c[0], c[1], c[2],
		)
	}

	sl.Printf(`device %s-%s-%s reporting changes to server`,
		dev.VID(), dev.PID(), dev.SN(),
	)

	dev.SetChanges(ch)

	return usbCiAuditV1(dev)
}

// report processes options and writes report to the selected destination.
func report(dev usb.Reporter) (err error) {

	var b []byte

	switch *fReportFormat {

	case `csv`:
		b, err = dev.CSV()

	case `nvp`:
		b, err = dev.NVP()

	case `xml`:
		b, err = dev.PrettyXML()

	case `json`:
		b, err = dev.PrettyJSON()

	default:
		err = fmt.Errorf(`invalid format %q`, *fReportFormat)
	}

	if err != nil {
		return err
	}

	if *fReportConsole {
		fmt.Fprintln(os.Stdout, string(b))
		return nil
	}

	f := fmt.Sprintf(`%s-%s.%s`, dev.SN(), dev.Conn(), *fReportFormat)

	if *fReportFolder != `` {
		f = filepath.Join(*fReportFolder, f)
	} else {
		f = filepath.Join(conf.Paths.ReportDir, f)
	}

	return ioutil.WriteFile(f, b, FileMode)
}

// serial processes options and configures the the serial number.
func serial(dev usb.Serializer) (err error) {

	var s string

	if *fSerialErase {

		sl.Printf(`device %s-%s erasing serial %q`,
			dev.VID(), dev.PID(), dev.SN(),
		)

		if err = dev.EraseDeviceSN(); err != nil {
			return err
		}
	}

	if !*fSerialForce && dev.SN() != `` {

		return fmt.Errorf(`device %s-%s serial already set to %q`,
			dev.VID(), dev.PID(), dev.SN(),
		)

	}

	switch {

	case *fSerialSet != ``:

		sl.Printf(`device %s-%s setting serial to %q`,
			dev.VID(), dev.PID(), *fSerialSet,
		)

		err = dev.SetDeviceSN(*fSerialSet)

	case *fSerialDefault:

		sl.Printf(`device %s-%s setting serial to default`,
			dev.VID(), dev.PID(),
		)

		err = dev.SetDefaultSN()

	case *fSerialFetch:

		if s, err = usbCiNewSnV1(dev); err != nil {
			break
		}

		sl.Printf(`device %s-%s setting serial to %q`,
			dev.VID(), dev.PID(), s,
		)

		if err = dev.SetDeviceSN(s); err != nil {
			break
		}

		sl.Printf(`device %s-%s-%s checking in with server`,
			dev.VID(), dev.PID(), dev.SN(),
		)

		err = usbCiCheckinV1(dev)
	}

	return err
}