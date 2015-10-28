package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

func errormsg(msg string) {
	message(msg)
}

func info(msg string) {
	message(msg)
}

func warn(msg string) {
	message(msg)
}

func showmsg(nick string, target string, text string) {
	timestr := time.Now().Format("15:04:05")
	fmt.Printf("%v <%v→%v> %v\n", timestr, nick, target, text)
}

func message(msg string) {
	timestr := time.Now().Format("15:04:05")
	fmt.Printf("%v %s\n", timestr, msg)
}

func parsecommand(line string) {
	fields := strings.Fields(line)

	switch fields[0] {
	case "/tlsconnect":
		if len(fields) != 3 {
			warn("Use /connect server:port nick")
			return
		}
		connect(fields[1], fields[2], true)

	case "/connect":
		if len(fields) != 3 {
			warn("Use /connect server:port nick")
			return
		}
		connect(fields[1], fields[2], false)

	case "/nick":
		if conn == nil {
			noconnection()
			break
		}

		conn.Nick(fields[1])

	case "/join":
		if conn == nil {
			noconnection()
			break
		}

		if len(fields) != 2 {
			warn("Use /join #channel")
			return
		}

		currtarget = fields[1]
		conn.Join(currtarget)

	case "/part":
		if conn == nil {
			noconnection()
			break
		}

		if len(fields) != 2 {
			warn("Use /part #channel")
			return
		}

		conn.Part(fields[1])
		currtarget = ""

	case "/names":
		if conn == nil {
			noconnection()
			break
		}

		namescmd := fmt.Sprintf("NAMES %v", currtarget)
		conn.Raw(namescmd)

	case "/whois":
		if conn == nil {
			noconnection()
			break
		}

		if len(fields) != 2 {
			warn("Use /whois <nick>")
			return
		}

		conn.Whois(fields[1])

	case "/query":
		if conn == nil {
			noconnection()
			break
		}

		if len(fields) != 2 {
			warn("Use /query <nick/channel>")
			return
		}

		currtarget = fields[1]

	case "/quit":
		iquit()
		if conn != nil {
			if len(fields) == 2 {
				conn.Quit(fields[1])
			} else {
				conn.Quit()
			}
		}

		quitclient = true

	default:
		warn("Unknown command: " + fields[0])
	}
}

func ui() {
	quitclient = false
	for !quitclient {
		bio := bufio.NewReader(os.Stdin)
		line, err := bio.ReadString('\n')

		if err != nil {
			log.Fatal("Couldn't get input.\n")
		}

		if line != "\n" && line != "\r\n" {
			if line[0] == '/' {
				// A command
				parsecommand(line)
			} else {
				// Send line to target.
				if currtarget == "" {
					notarget()
				} else {
					conn.Privmsg(currtarget, line)
					logmsg(time.Now(), conn.Me().Nick, currtarget, line)
				}
			}
		}
	}
}
