// Program to program HPSDR boards from the command line
//
// by David R. Larsen KV0S, Copyright 2014-11-24
//
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/kv0s/openhpsdr"
)

const version string = "0.1.4"
const update string = "2020-10-24"

// function to point users to the command list
func usage() {
	fmt.Printf("    For a list of commands use -help \n\n")
}

// Function to print the program name info
func program() {
	fmt.Printf("CmdHPSDRProgrammer  version:(%s)\n", version)
	fmt.Printf("    By Dave KV0S, 2014-11-24, GPL2 \n")
	fmt.Printf("    Last Updated: %s \n\n", update)
}

// Listboard is a convenience function to print board data
func Listboard(str openhpsdr.Hpsdrboard) {
	if str.Macaddress != "0:0:0:0:0:0" {
		fmt.Printf("       HPSDR Board: (%s)\n", str.Macaddress)
		fmt.Printf("              IPV4: %s\n", str.Baddress)
		fmt.Printf("              Port: %s\n", str.Bport)
		fmt.Printf("              Type: %s\n", str.Board)
		fmt.Printf("          Firmware: %s\n", str.Firmware)
		fmt.Printf("            Status: %s\n\n", str.Status)
		fmt.Printf("            PC    : %s\n\n", str.Pcaddress)
	}
}

//Listinterface is a convenience function to print interface data
func Listinterface(itr openhpsdr.Intface) {
	fmt.Printf("          Computer: (%v)\n", itr.MAC)
	fmt.Printf("                OS: %s (%s) %d CPU(s)\n", runtime.GOOS, runtime.GOARCH, runtime.NumCPU())
	fmt.Printf("              IPV4: %v\n", itr.Ipv4)
	fmt.Printf("              Mask: %d\n", itr.Mask)
	fmt.Printf("           Network: %v\n", itr.Network)
	fmt.Printf("              IPV6: %v\n\n", itr.Ipv6)
}

//Listflags is a convienience function to print flag data
func Listflags(fg flagsettings) {
	fmt.Printf("    Saved Settings: \n")
	fmt.Printf("         Interface: %v\n", fg.Intface)
	fmt.Printf("          Filename: %v\n", fg.Filename)
	fmt.Printf("      Selected MAC: (%v)\n", fg.SelectMAC)
	fmt.Printf("            SetRBF: %v\n", fg.SetRBF)
	fmt.Printf("             Debug: %v\n", fg.Debug)
	fmt.Printf("            Ddelay: %d\n\n", fg.Ddelay)
	fmt.Printf("            Edelay: %d\n\n", fg.Edelay)
}

//Listflagstemp is a convienience function to print temporary flag data
func Listflagstemp(fgt flagtemp) {
	fmt.Printf("     Temp settings: \n")
	fmt.Printf("          Settings: %v\n", fgt.Settings)
	fmt.Printf("             SetIP: %v\n", fgt.SetIP)
	fmt.Printf("              Save: %v\n", fgt.Save)
	fmt.Printf("              Load: %v\n\n", fgt.Load)
}

// Initflags is a convienience function to initialize command line flags
func Initflags(fg *flagsettings) {
	fg.Intface = "none"
	fg.Filename = "none"
	fg.SelectMAC = "none"
	fg.SetRBF = "none"
	fg.Debug = "none"
	fg.Ddelay = 2
	fg.Edelay = 60
}

//flagsetting is a local structure to work with command line flags
type flagsettings struct {
	Filename  string
	Intface   string
	SelectMAC string
	SetRBF    string
	Debug     string
	Ddelay    int
	Edelay    int
}

//flagtemp is a local structure to work with command line temp flags
type flagtemp struct {
	SetIP    string
	Settings string
	Save     string
	Load     string
}

//Initflagstemp is a function to initialize the temp flags
func Initflagstemp(fgt *flagtemp) {
	fgt.SetIP = "none"
	fgt.Settings = "none"
	fgt.Save = "none"
	fgt.Load = "none"
}

//Parseflagstruct is a function to parse input flags
func Parseflagstruct(fg *flagsettings, fgt *flagtemp, ifn string, stmac string, stip string, strbf string, db string, ss string, sv string, ld string, dd int, ed int) {

	Initflags(fg)
	Initflagstemp(fgt)

	if (ld == "default") || (ld == "Default") {
		fg.Filename = "CmdHPSDRProgrammer.json"
	} else if ld != "none" {
		fg.Filename = ld
	}

	if ld != "none" {

		dta, _ := ioutil.ReadFile(fg.Filename)
		err := json.Unmarshal(dta, &fg)
		if err != nil {
			fmt.Println("error:", err)
		}
	}

	if ifn != "none" {
		fg.Intface = ifn
	}
	if stmac != "none" {
		fg.SelectMAC = stmac
	}
	if strbf != "none" {
		fg.SetRBF = strbf
	}
	if db != "none" {
		fg.Debug = db
	}
	if ed != 20 {
		fg.Edelay = ed
	}
	if dd != 2 {
		fg.Ddelay = dd
	}
	if ed != 2 {
		fg.Edelay = ed
	}
	if stip != "none" {
		fgt.SetIP = stip
	}
	if ss != "none" {
		fgt.Settings = ss
	}
	if sv == "default" {
		fgt.Save = sv
		fg.Filename = "CmdHPSDRProgrammer.json"
	} else if sv != "none" {
		fg.Filename = sv
		fgt.Save = sv
	} else {
		fgt.Save = sv
	}
	if ld == "default" {
		fgt.Load = ld
		fg.Filename = "CmdHPSDRProgrammer.json"
	} else if ld != "none" {
		fg.Filename = ld
		fgt.Load = ld
	} else {
		fgt.Load = ld
	}

	if fgt.Save != "none" {

		f, err := os.Create(fg.Filename)
		if err != nil {
			panic(err)
		}

		b, err := json.MarshalIndent(fg, "", "\t")
		if err != nil {
			fmt.Println("error:", err)
		}

		fmt.Fprintf(f, "%s\n", b)
	}

	if ss != "none" {
		Listflags(*fg)
		Listflagstemp(*fgt)
	}

}

func main() {
	var fg flagsettings
	var fgt flagtemp
	var erstat openhpsdr.Erasestatus

	// Create the command line flags
	ifn := flag.String("interface", "none", "Select one interface")
	stmac := flag.String("selectMAC", "none", "Select Board by MAC address")
	stip := flag.String("setIP", "none", "Set IP address, unused number from your subnet or 0.0.0.0 for DHCP")
	strbf := flag.String("setRBF", "none", "Select the RBF file to write to the board")
	dd := flag.Int("ddelay", 2, "Discovery delay before a rediscovery")
	ed := flag.Int("edelay", 60, "Discovery delay before a rediscovery")
	db := flag.String("debug", "none", "Turn debugging and output type, (none, dec, hex)")
	ss := flag.String("settings", "none", "Show the settings values (show)")
	sv := flag.String("save", "none", "Save these current flags for future use in default or a named file")
	ld := flag.String("load", "none", "Load a saved command file from default or a named file")
	cadr := flag.Bool("checkaddress", true, "check if new address is in subdomain and not restricted space")
	cbad := flag.Bool("checkboard", true, "check if new RBF file name has the same name as the board type")

	flag.Parse()

	if flag.NFlag() < 1 {
		program()
		usage()
	}

	Parseflagstruct(&fg, &fgt, *ifn, *stmac, *stip, *strbf, *db, *ss, *sv, *ld, *dd, *ed)

	intf, err := openhpsdr.Interfaces()
	if err != nil {
		fmt.Println(err)
	}

	if flag.NFlag() < 1 {
		fmt.Printf("Interfaces on this Computer: \n")
	}
	for i := range intf {
		if flag.NFlag() < 1 {
			// if no flags list the interfaces in short form
			fmt.Printf("    %s (%s)\n", intf[i].Intname, intf[i].MAC)
		} else if (flag.NFlag() == 1) && (fg.Intface == "none") {
			if fg.Debug == "none" {
				// if one flag and it is debug = none, list the interface in short form
				fmt.Printf("    %s (%s)\n", intf[i].Intname, intf[i].MAC)
			} else {
				// if one flag and it is debug = dec or hex, list the interface in long form
				fmt.Printf("    %s (%s) %s  %s\n", intf[i].Intname, intf[i].MAC, intf[i].Ipv4, intf[i].Ipv6)
			}
		}

		// if ifn flag matches the current interface
		if fg.Intface == intf[i].Intname {
			if len(intf[i].Ipv4) != 0 {
				if fg.Debug == "none" {
					//list the sending computer information
					Listinterface(intf[i])
				}

				var adr string
				var bcadr string
				adr = intf[i].Ipv4 + ":1024"
				bcadr = intf[i].Ipv4Bcast + ":1024"

				// perform a discovery
				str, err := openhpsdr.Discover(adr, bcadr, fg.Ddelay, fg.Debug)
				if err != nil {
					fmt.Println("Error ", err)
				}

				var bdid int
				bdid = 0
				//loop throught the list of discovered HPSDR boards
				for i := 0; i < len(str); i++ {
					if fg.SelectMAC == str[i].Macaddress {
						// if a MAC is selected
						fmt.Printf("      Selected MAC: (%s) %s\n", fg.SelectMAC, str[i].Board)
						bdid = i
					}
				}

				if (len(str) > 0 && fgt.SetIP != str[bdid].Baddress) && (fgt.SetIP != "none") {
					//If the IPV4 changes
					if strings.Contains(*stip, "255.255.255.255") {
						fmt.Printf("     Changing IP address from %s to DHCP address\n\n", str[bdid].Baddress)
					} else {
						fmt.Printf("     Changing IP address from %s to %s\n\n", str[bdid].Baddress, *stip)
					}

					str2, err := openhpsdr.Setip(fgt.SetIP, str[bdid], fg.Debug, *cadr)
					if err != nil {
						fmt.Printf("Error %v", str2)
						panic(err)
					}

					// perform a rediscovery
					time.Sleep(time.Duration(fg.Ddelay) * time.Second)
					str, err = openhpsdr.Discover(adr, bcadr, fg.Ddelay, fg.Debug)
					if err != nil {
						fmt.Println("Error ", err)
					}

					//loop throught the list of discovered HPSDR boards
					for i := 0; i < len(str); i++ {
						if fg.SelectMAC == str[i].Macaddress {
							// if a MAC is selected
							//fmt.Printf("      Selected MAC: %s\n", fg.SelectMAC)
							bdid = i
						}
					}
				}
				if *strbf != "none" {
					if *cbad && (fg.SelectMAC != "none") && (fg.SelectMAC == str[bdid].Macaddress) {
						if strings.Contains(strings.ToLower(*strbf), strings.ToLower(str[bdid].Board)) {
							// erasy the board flash memory
							erstat, err = openhpsdr.Eraseboard(str[bdid], fg.SetRBF, fg.Edelay, fg.Debug)
							if err != nil {
								panic(err)
							} else {
								fmt.Printf(" %v %v\n", erstat.Seconds, erstat.State)
								// send the RBF to the flash memory
								openhpsdr.Programboard(str[bdid], fg.SetRBF, fg.Debug)
							}
						} else {
							fmt.Printf("\n      Input Check: RBF name \"%s\" and selectedMAC board name \"%s\" (%s) do not match!\n", *strbf, str[bdid].Board, str[bdid].Macaddress)
							fmt.Println("       Please correct to program the board.\n")
						}
					} else {
						// easy the board flash memory
						erstat, err = openhpsdr.Eraseboard(str[bdid], fg.SetRBF, fg.Edelay, fg.Debug)
						if err != nil {
							panic(err)
						} else {
							fmt.Printf(" %v %v\n", erstat.Seconds, erstat.State)
							// send the RBF to the flash memory
							openhpsdr.Programboard(str[bdid], fg.SetRBF, fg.Debug)
						}

					}
				}

				if fg.Debug == "none" {
					if (fg.SelectMAC != "none") && (fg.SelectMAC == str[bdid].Macaddress) {
						// list all the HPSDR Board information or the select HPSDR Board information

						Listboard(str[bdid])
					} else if fg.SelectMAC == "none" && len(str) > 0 {
						//loop throught the list of discovered HPSDR boards
						for i := 0; i < len(str); i++ {
							Listboard(str[i])
						}
					} else if len(str) == 0 {
						fmt.Printf("      No HPSDR Boards found on interface \"%s\"! \n", *ifn)
					}
				}
			} else {
				fmt.Printf("      Interface not active! \n")
			}
		}
	}
}
