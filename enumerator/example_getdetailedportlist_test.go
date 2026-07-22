//
// Copyright 2014-2026 Cristian Maglie. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package enumerator_test

import (
	"fmt"
	"log"

	"github.com/flozano/go-serial/enumerator"
)

func ExampleGetDetailedPortsList() {
	// Passing enumerator.All actively probes every USB device to retrieve the
	// Manufacturer, Product and Configuration fields. This may interfere with
	// the normal operation of some devices, so in production code you should
	// prefer a filter that only allows probing specific VID/PID pairs, e.g.:
	//
	//   enumerator.GetDetailedPortsList(func(vid, pid string) bool {
	//       return vid == "2341" // only probe Arduino devices
	//   })
	ports, err := enumerator.GetDetailedPortsList(enumerator.All)
	if err != nil {
		log.Fatal(err)
	}
	if len(ports) == 0 {
		fmt.Println("No serial ports found!")
		return
	}
	for _, port := range ports {
		fmt.Printf("Found port: %s\n", port.Name)
		if port.IsUSB {
			fmt.Printf("   USB ID     %s:%s\n", port.VID, port.PID)
			fmt.Printf("   USB vendor %s\n", port.Manufacturer)
			fmt.Printf("   USB prod.  %s\n", port.Product)
			fmt.Printf("   USB serial %s\n", port.SerialNumber)
			fmt.Printf("   USB config %s\n", port.Configuration)
		}
	}
}
