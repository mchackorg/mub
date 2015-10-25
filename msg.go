package main

import (
	"fmt"
	"time"
)

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

func members(time time.Time, members string) {
	timestr := time.Format("15:04:05")
	fmt.Printf("%v Members: %v\n", timestr, members)
}
