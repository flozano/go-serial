//
// Copyright 2014-2026 Cristian Maglie. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

// portlist is a tool to list all the available serial ports.
// It will print the port name and, when available, the USB VID/PID and other details.

package main

import (
	"flag"
	"fmt"
	"log"

	"go.bug.st/serial/enumerator"
)

func main() {
	probe := flag.Bool("probe", false, "actively probe USB devices to retrieve manufacturer, product and configuration strings")
	flag.Parse()

	var filters []func(vid, pid string) bool
	if *probe {
		filters = append(filters, enumerator.All)
	}

	ports, err := enumerator.GetDetailedPortsList(filters...)
	if err != nil {
		log.Fatal(err)
	}
	if len(ports) == 0 {
		return
	}
	for _, port := range ports {
		fmt.Printf("Port: %s\n", port.Name)
		if port.IsUSB {
			fmt.Printf("   USB VID/PID      : %s:%s\n", port.VID, port.PID)
			fmt.Printf("   USB serial no.   : %s\n", port.SerialNumber)
			fmt.Printf("   USB manufacturer : %s\n", port.Manufacturer)
			fmt.Printf("   USB product      : %s\n", port.Product)
			fmt.Printf("   USB config       : %s\n", port.Configuration)
		}
	}
}
