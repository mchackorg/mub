package main

import (
	"bufio"
	"fmt"
	"github.com/chzyer/readline"
	"io"
	"log"
	"os"
	"reflect"
	"strings"
	"time"
)

type command struct {
	Name      string
	Prototype interface{}
	Desc      string
}

type nickname string
type channel string
type nickorchan string

type nocommand struct{}

type helpcommand struct{}

type tlsconnectcommand struct {
	Server string "IRC server"
	Nick   string "Your nickname"
}

type connectcommand struct {
	Server string "IRC server"
	Nick   string "Your nickname"
}

type quitcommand struct{}

type querycommand struct {
	Target nickorchan "channel or nick"
}

type joincommand struct {
	Channel channel "channel"
}

type whoiscommand struct {
	Nick nickname "nick"
}

type nickcommand struct {
	Nick nickname "nick"
}

type partcommand struct {
	Channel channel "channel"
}

type mecommand struct {
	Action string
}

type namescommand struct{}

// Internal state of completer.
type CommandState struct {
	FoundCmd int
	Channels map[string]string
	NickMap  map[string]string
}

type Commands struct {
	Commands []command
	State    *CommandState
}

var (
	output io.Writer // All output should go here, not to stdout.

	commands = Commands{
		Commands: []command{
			{"", nocommand{}, "No command given."},
			{"/help", helpcommand{}, "Give this help"},
			{"/tlsconnect", tlsconnectcommand{}, "Connect to IRC server using TLS."},
			{"/connect", connectcommand{}, "Connect to IRC server."},
			{"/quit", quitcommand{}, "Quit the IRC client."},
			{"/query", querycommand{}, "Start talking to a nick or channel."},
			{"/x", querycommand{}, "Shorthand for /query."},
			{"/join", joincommand{}, "Join a channel."},
			{"/part", partcommand{}, "Leave a channel."},
			{"/whois", whoiscommand{}, "Show information about someone."},
			{"/me", mecommand{}, "Show a string describing you doing something."},
			{"/nick", nickcommand{}, "Change your nickname."},
			{"/names", namescommand{}, "List members on current channel."}},
	}
)

// Completer for readline
func (c Commands) Do(line []rune, pos int) (newLine [][]rune, length int) {
	var linestr string = string(line)
	var matches int

	// Find out what command this is:
	space := strings.IndexRune(linestr, ' ')
	if space == -1 {
		if len(line) != 0 && linestr[0] == '/' {
			// This is a command completion.
			for i, cmd := range c.Commands {
				if strings.HasPrefix(cmd.Name, strings.ToLower(linestr)) {
					newLine = append(newLine, []rune(cmd.Name[pos:]))
					matches++
					c.State.FoundCmd = i
				}
			}
			if matches != 1 {
				c.State.FoundCmd = 0
			}
		} else {
			// Nick completion.
			newLine = findmap(linestr[space+1:], c.State.NickMap, pos)
		}
	} else {
		// Argument completion.

		// The line so far is...
		head := linestr[:space] + " "
		// ...and our position in the this word is:
		wordpos := pos - len(head)

		//msg := fmt.Sprintf("foundcmd: %v", c.State.FoundCmd)
		//info(msg)

		switch c.Commands[c.State.FoundCmd].Prototype.(type) {
		case querycommand:
			// if len(linestr) < space {
			// 	return
			// }

			if strings.HasPrefix(linestr[space+1:], "#") {
				// Complete a channel.
				newLine = findmap(linestr[space+1:], c.State.Channels, wordpos)
			} else {
				// Complete a nickname.
				newLine = findmap(linestr[space+1:], c.State.NickMap, wordpos)
			}

		case whoiscommand:
			newLine = findmap(linestr[space+1:], c.State.NickMap, wordpos)
		case joincommand:
			newLine = findmap(linestr[space+1:], c.State.Channels, wordpos)
		case partcommand:
			newLine = findmap(linestr[space+1:], c.State.Channels, wordpos)
		}
	}

	length = len(linestr)

	return
}

func findmatch(arg string, args []string, wordpos int) (newLine [][]rune) {
	for _, n := range args {
		//msg := fmt.Sprintf("comparing %v to %v", arg, n)
		//info(msg)
		if strings.HasPrefix(n, strings.ToLower(arg)) {
			newLine = append(newLine, []rune(n[wordpos:]))
		}
	}

	return
}

func findmap(arg string, args map[string]string, wordpos int) (newLine [][]rune) {
	for _, n := range args {
		//msg := fmt.Sprintf("comparing %v to %v", arg, n)
		//info(msg)
		if strings.HasPrefix(n, strings.ToLower(arg)) {
			newLine = append(newLine, []rune(n[wordpos:]))
		}
	}

	return
}

func errormsg(msg string) {
	message(msg)
}

func info(msg string) {
	message(msg)
}

func warn(msg string) {
	message(msg)
}

func showmsg(nick string, target string, text string, action bool) {
	var str string

	if action {
		str = fmt.Sprintf("* On %v: %v %v", target, nick, text)
	} else {
		str = fmt.Sprintf("<%v→%v> %v", nick, target, text)
	}

	message(str)
}

// Sanitize string msg from ESC and control characters.
func sanitizestring(msg string) (out string) {
	for _, c := range msg {
		if c == 127 || (c < 32 && c != '\t') {
			out = out + "?"
		} else {
			out = out + string(c)
		}
	}

	return out
}

func message(msg string) {
	timestr := time.Now().Format("15:04:05")
	msg = sanitizestring(msg)
	fmt.Fprintf(output, "%v %s\n", timestr, msg)
}

func printhelp() {
	for _, cmd := range commands.Commands {
		msg := cmd.Name
		prototype := reflect.TypeOf(cmd.Prototype)
		for i := 0; i < prototype.NumField(); i++ {
			msg += " <" + strings.ToLower(prototype.Field(i).Name) + ">"
		}

		message(msg + " - " + cmd.Desc)
	}
}

func parsecommand(line string) {
	fields := strings.Fields(line)

	switch fields[0] {
	case "/help":
		printhelp()

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
		commands.State.Channels[currtarget] = currtarget
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
		// Forget about this channel
		delete(commands.State.Channels, currtarget)
	case "/me":
		if conn == nil {
			noconnection()
			break
		}

		if len(fields) < 2 {
			warn("Use /me action text")
			return
		}

		actiontext := line[strings.Index(line, " ")+1:]
		conn.Action(currtarget, actiontext)
		logmsg(time.Now(), conn.Me().Nick, currtarget, actiontext, true)

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

	case "/x":
		fallthrough
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

func ui(subprocess bool) {
	var state CommandState
	var rl *readline.Instance
	var line string
	var err error
	var bio *bufio.Reader

	// Internal state for command completer.
	commands.State = &state
	commands.State.NickMap = make(map[string]string)
	commands.State.Channels = make(map[string]string)

	if subprocess {
		// We're running as a subprocess. Just read from stdin.
		bio = bufio.NewReader(os.Stdin)
		output = os.Stdout
	} else {
		// Slightly smarter UI is used.
		rl, err = readline.NewEx(&readline.Config{
			HistoryFile:  "/tmp/mub.tmp",
			AutoComplete: commands,
		})
		if err != nil {
			panic(err)
		}
		defer rl.Close()

		// Send output to readline's handler so prompt can
		// refresh.
		output = rl.Stdout()
	}

	quitclient = false
	for !quitclient {
		if subprocess {
			line, err = bio.ReadString('\n')
			if err != nil {
				log.Fatal("Couldn't get input.\n")
			}
		} else {
			rl.SetPrompt("\033[33m[" + currtarget + "] \033[0m")
			line, err = rl.Readline()
			if err != nil {
				break
			}
		}

		line = strings.TrimSpace(line)
		if line != "" && line != "\n" && line != "\r\n" {
			if line[0] == '/' {
				// A command
				parsecommand(line)
			} else {
				// Send line to target.
				if currtarget == "" {
					notarget()
				} else {
					conn.Privmsg(currtarget, line)
					logmsg(time.Now(), conn.Me().Nick, currtarget, line, false)
				}
			}
		}
	}
}
