/*
 Licensed to the Apache Software Foundation (ASF) under one
 or more contributor license agreements.  See the NOTICE file
 distributed with this work for additional information
 regarding copyright ownership.  The ASF licenses this file
 to you under the Apache License, Version 2.0 (the
 "License"); you may not use this file except in compliance
 with the License.  You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing,
 software distributed under the License is distributed on an
 "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 KIND, either express or implied.  See the License for the
 specific language governing permissions and limitations
 under the License.
*/

package main

import (
	"fmt"
	"os"

	"github.com/joanbono/color"
	"github.com/integrii/flaggy"

	"github.com/joanbono/Gurp/modules/commander"
	"github.com/joanbono/Gurp/modules/configure"
	"github.com/joanbono/Gurp/modules/nmap"
)

// Defining colors
var yellow = color.New(color.Bold, color.FgYellow).SprintfFunc()
var red = color.New(color.Bold, color.FgRed).SprintfFunc()
var cyan = color.New(color.Bold, color.FgCyan).SprintfFunc()
var green = color.New(color.Bold, color.FgGreen).SprintfFunc()

var VERSION = `1.1.1`

//var BurpAPI, username, password, ApiToken string
var target, port string = "127.0.0.1", "1337"
var export string
var username, password string = "", ""
var key string = ""
var description string = ""
var metrics, issues bool = false, false
var scan, scan_id, scanList, nmapScan string
var version, description_names bool = false, false

func init() {
	flaggy.SetName("Gurp")
	flaggy.SetDescription("Interact with Burp API")
	flaggy.DefaultParser.ShowVersionWithVersionFlag = false

	flaggy.String(&target, "t", "target", "Burp Address. Default 127.0.0.1")
	flaggy.String(&port, "p", "port", "Burp API Port. Default 1337")

	flaggy.String(&username, "U", "username", "Username for an authenticated scan")
	flaggy.String(&password, "P", "password", "Password for an authenticated scan")

	flaggy.String(&scan, "s", "scan", "URLs to scan")
	flaggy.String(&scan_id, "S", "scan-id", "Scanned URL identifier")

	flaggy.String(&nmapScan, "sn", "scan-nmap", "Nmap xml file to scan")
	flaggy.String(&scanList, "sl", "scan-list", "File with hosts/Ip's to scan")

	flaggy.Bool(&metrics, "M", "metrics", "Provides metrics for a given task")
	flaggy.String(&description, "D", "description", "Provides description for a given issue")
	flaggy.Bool(&description_names, "d", "description-names", "Returns vulnerability names from PortSwigger")
	flaggy.Bool(&issues, "I", "issues", "Provides issues for a given task")
	flaggy.String(&export, "e", "export", "Export issues' json.")

	flaggy.String(&key, "k", "key", "Api Key")
	flaggy.Bool(&version, "v", "version", "Gurp version")
}

func main() {
	flaggy.Parse()

	// Check how many args are provided
	if len(os.Args) < 2 {
		fmt.Fprintf(color.Output, "\n %v No argument provided. Try with %v.\n\n", cyan("[i] INFO:"), green("gurp -h"))
		os.Exit(0)
	}

	if version == true {
		fmt.Fprintf(color.Output, "%v Gurp %v.\n", cyan(" [i] INFO:"), VERSION)
		os.Exit(0)
	}
	if configure.CheckBurp(target, port, key) == true {
		fmt.Fprintf(color.Output, "%v Found Burp API endpoint on %v.\n", green(" [+] SUCCESS:"), target+":"+port)
	} else {
		fmt.Fprintf(color.Output, "%v No Burp API endpoint found on %v.\n", red(" [-] ERROR:"), target+":"+port)
		os.Exit(0)
	}

	if nmapScan != "" {

		scanList, err := nmap.ParseNmap(nmapScan)
		if err != nil {
			fmt.Fprintf(color.Output, "%v  %v.\n", red(" [-] ERROR:"), err)
			os.Exit(0)
		}
		for _, scan := range scanList {
			Location := configure.ScanConfig(target, port, scan, username, password, key)
			if Location != "" {
				fmt.Fprintf(color.Output, "%v Scanning %v over %v.\n", green(" [+] SUCCESS:"), scan, Location)
			} else {
				fmt.Fprintf(color.Output, "%v Can't start scan .\n", red(" [-] ERROR:"))
				os.Exit(0)
			}
		}

	}

	if scanList != "" {

		targets := nmap.ParseFile(scanList)
		for _, scan := range targets {
			Location := configure.ScanConfig(target, port, scan, username, password, key)
			if Location != "" {
				fmt.Fprintf(color.Output, "%v Scanning %v over %v.\n", green(" [+] SUCCESS:"), scan, Location)
			} else {
				fmt.Fprintf(color.Output, "%v Can't start scan over %s .\n", red(" [-] ERROR:"), scan)
			}
		}
	}

	if scan != "" {
		//fmt.Println(configure.ScanConfig(target, port, scan))
		Location := configure.ScanConfig(target, port, scan, username, password, key)
		if Location != "" {
			fmt.Fprintf(color.Output, "%v Scanning %v over %v.\n", green(" [+] SUCCESS:"), scan, Location)
		} else {
			fmt.Fprintf(color.Output, "%v Can't start scan .\n", red(" [-] ERROR:"))
			os.Exit(0)
		}

	}

	if scan == "" && scan_id != "" && metrics == true && issues == false {
		//commander.GetScan(target, port, scan_id)
		commander.GetMetrics(target, port, scan_id, key)
	} else if scan == "" && scan_id != "" && metrics == true && issues == true {
		commander.GetScan(target, port, scan_id, export, key)
		commander.GetMetrics(target, port, scan_id, key)
	} else if scan == "" && scan_id != "" && metrics == false {
		commander.GetScan(target, port, scan_id, export, key)
	}

	if description != "" {
		configure.GetDescription(target, port, description, key)
	}
	if description_names == true {
		configure.GetNames(target, port, key)
	}
}
