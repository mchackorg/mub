// All messages to the user should be defined here.
package main

import (
	"fmt"
	"time"
)

func notarget() {
	warn("Join a channel first.")
}

func noconnection() {
	warn("Not connected to a server.")
}

func connecting(server string) {
	info("Connecting to " + server + "...")
}

func nick(oldnick string, newnick string) {
	info(oldnick + " is now known as " + newnick)
}

func connectionerror(err error) {
	line := fmt.Sprintf("Connection error: %v", err)
	errormsg(line)
}

func disconnected() {
	warn("Disconnected from server.")
}

func connected(nick string) {
	info("Connected as " + nick)
}

func msg(time time.Time, nick string, target string, text string, action bool) {
	showmsg(nick, target, text, action)
}

func iquit() {
	info("Quitting.")
}

func quit(nick string) {
	line := fmt.Sprintf("%v quit IRC.", nick)
	info(line)
}

func notice(nick string, msg string) {
	line := fmt.Sprintf("NOTICE: %v %v", nick, msg)
	info(line)
}

func whois(nick string, realname string, ident string, host string) {
	line := fmt.Sprintf("%s is %s <%s@%s>", nick, realname, ident, host)
	info(line)
}

func joined(nick string, channel string) {
	line := fmt.Sprintf("%s joined %s", nick, channel)
	info(line)
}

func parted(nick string, channel string) {
	line := fmt.Sprintf("%v parted %v", nick, channel)
	info(line)
}

func cantopenfile(filename string, err error) {
	fmt.Printf("Could not open file '%s', will not write to log\n", filename)
	fmt.Printf("Error Message: %s\n", err)
}

func members(channel string, members string) {
	line := fmt.Sprintf("Members on %v: %v", channel, members)
	info(line)
}
