// Copyright 2017 John Scherff
//
// Licensed under the Apache License, Version 2.0 (the `License`);
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an `AS IS` BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	`errors`
	`fmt`
	`log`
	`os`

	`github.com/google/gousb`
	`github.com/jscherff/gocmdb/usbci`
	`github.com/jscherff/goutils`
)

var (
	conf *Config
	slog, clog, elog *log.Logger
)

func init() {

	var err error

	// Build systemwide configuration from config file.

	if conf, err = NewConfig(`config.json`); err != nil {
		log.Fatalf(`%v`, goutils.ErrorDecorator(err))
	}

	// Process command-line actions and options.

	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, `You must specify an action.`)
		fsAction.Usage()
		os.Exit(1)
	}

	// Initialized loggers.

	slog, clog, elog = NewLoggers()

	// Parse action flag.

	fsAction.Parse(os.Args[1:2])

	// Parse option flags associated with selected action flag.

	switch {

	case *fActionReport:
		if fsReport.Parse(os.Args[2:]); fsReport.NFlag() == 0 {
			fmt.Fprintln(os.Stderr, `You must specify an option.`)
			fsReport.Usage()
			os.Exit(1)
		}

	case *fActionSerial:
		if fsSerial.Parse(os.Args[2:]); fsSerial.NFlag() == 0 {
			fmt.Fprintln(os.Stderr, `You must specify an option.`)
			fsSerial.Usage()
			os.Exit(1)
		}
	}
}

func main() {

	// Instantiate context to enumerate attached USB devices.

	context := gousb.NewContext()
	defer context.Close()

	// Open devices that match selection criteria in the Include.ProductID
	// and Include.VendorID maps from the configuration file.

	devices, _ := context.OpenDevices(func(desc *gousb.DeviceDesc) bool {

		vid, pid := desc.Vendor.String(), desc.Product.String()

		if val, ok := conf.Include.ProductID[vid][pid]; ok {return val}
		if val, ok := conf.Include.VendorID[vid]; ok {return val}

		return conf.Include.Default
	})

	// Log and exit if no relevant devices found.

	if len(devices) == 0 {
		elog.Fatalf(`%v`, goutils.ErrorDecorator(errors.New(`no devices found`)))
	}

	// Pass devices to relevant device handlers.

	for _, device := range devices {

		defer device.Close()

		slog.Printf(`found USB device, VID %s, VID %s`,
			device.Desc.Vendor.String(),
			device.Desc.Product.String(),
		)

		switch uint16(device.Desc.Vendor) {

		case usbci.MagtekVendorID:

			if d, err := usbci.NewMagtek(device); err != nil {
				elog.Printf("%v", goutils.ErrorDecorator(err))
			} else {
				slog.Printf(`identified USB device as %s`, d.Type())
				magtekRouter(d)
			}

		default:

			if d, err := usbci.NewGeneric(device); err != nil {
				elog.Printf("%v", goutils.ErrorDecorator(err))
			} else {
				slog.Printf(`identified USB device as %s`, d.Type())
				genericRouter(d)
			}
		}
	}
}
