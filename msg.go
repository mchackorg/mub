// All messages to the user should be defined here.
package main

import (
	"fmt"
	"time"
)

func notarget() {
	fmt.Printf("Join a channel first.\n")
}

func connecting(server string) {
	fmt.Printf("Connecting to %v...\n", server)
}

func connectionerror(err error) {
	fmt.Printf("Connection error: %s\n", err)
}

func disconnected() {
	fmt.Printf("Disconnected from server.\n")
}

func connected(time time.Time) {
	timestr := time.Format("15:04:05")
	fmt.Printf("%v Connected.\n", timestr)
}

func msg(time time.Time, nick string, target string, text string) {
	timestr := time.Format("15:04:05")

	if target != currtarget {
		fmt.Printf("%v %v <%v> %v\n", timestr, target, nick, text)
	} else {
		fmt.Printf("%v <%v> %v\n", timestr, nick, text)
	}
}

func commanderror(help string) {
	fmt.Printf(help)
}

func iquit() {
	fmt.Printf("Quitting.\n")
}

func quit(time time.Time, nick string) {
	timestr := time.Format("15:04:05")
	fmt.Printf("%v %v quit IRC.\n", timestr, nick)
}

func whois(time time.Time, realname string, ident string, host string) {
	timestr := time.Format("15:04:05")
	fmt.Printf("%v %v <%v@%v>\n", timestr, realname, ident, host)
}

func joined(time time.Time, nick string, channel string) {
	timestr := time.Format("15:04:05")
	fmt.Printf("%v %v joined %v\n", timestr, nick, channel)
}

func parted(time time.Time, nick string, channel string) {
	timestr := time.Format("15:04:05")
	fmt.Printf("%v %v parted %v\n", timestr, nick, channel)
}

func cantopenfile(filename string, err error) {
	fmt.Printf("Could not open file '%s', will not write to log\n", filename)
	fmt.Printf("Error Message: %s\n", err)
}

func members(time time.Time, members string) {
	timestr := time.Format("15:04:05")
	fmt.Printf("%v Members: %v\n", timestr, members)
}
