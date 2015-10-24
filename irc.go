// Small cross-platform IRC client.

package main

import "bufio"
import "fmt"
import "time"
import "os"
import "gopkg.in/yaml.v2"
import "io/ioutil"
import "flag"
import "log"
import "strings"

import irc "github.com/fluffle/goirc/client"

type Config struct {
	Nick     string
	RealName string
	Server   string
	LogFile  string
	TLS      bool
}

var (
	file       *os.File
	conn       *irc.Conn
	target     string
	quitclient bool
)

const (
	IRC_privmsg = iota
	IRC_join
	IRC_part
	IRC_quit
	IRC_names
)

// Parse the configuration file. Returns the configuration.
func parseconfig(filename string) (conf *Config, err error) {
	contents, err := ioutil.ReadFile(filename)
	if err != nil {
		return
	}

	conf = new(Config)

	if err = yaml.Unmarshal(contents, &conf); err != nil {
		return
	}

	return
}

// Log text
func logmsg(time time.Time, nick string, target string, text string) {
	line := time.UTC().Format("2006-01-02 15:04:05")
	if target != "" {
		line += " " + target
	}
	text = strings.TrimRight(text, "\r\n")
	line += " <" + nick + "> " + text + "\n"

	if file != nil {
		_, err := file.WriteString(line)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func handler(msgtype int, line *irc.Line) {
	time := line.Time.Format("15:04:05")

	switch msgtype {
	case IRC_privmsg:
		if line.Target() != target {
			fmt.Printf("%v %v <%v> %v\n", time, line.Target(), line.Nick, line.Text())
		} else {
			fmt.Printf("%v <%v> %v\n", time, line.Nick, line.Text())
		}

		logmsg(line.Time, line.Nick, line.Target(), line.Text())

	case IRC_join:
		fmt.Printf("%v %v joined %v\n", time, line.Nick, line.Target())
	case IRC_part:
		fmt.Printf("%v %v left %v\n", time, line.Nick, line.Target())
	case IRC_names:
		fmt.Printf("Members: %v\n", line.Args[3])
	case IRC_quit:
		fmt.Printf("%v %v quit IRC.\n", time, line.Nick)
	default:
		fmt.Printf("Unknown message to handler.\n")
	}

}

func connected(conn *irc.Conn, line *irc.Line) {
	fmt.Printf("Connected.\n")
}

func parsecommand(line string) {
	fields := strings.Fields(line)

	switch fields[0] {
	case "/nick":
		conn.Nick(fields[1])
	case "/join":
		if len(fields) != 2 {
			fmt.Printf("Use /join #channel\n")
			return
		}

		target = fields[1]
		fmt.Printf("Now talking on %v\n", target)
		conn.Join(target)
	case "/part":
		if len(fields) != 2 {
			fmt.Printf("Use /part #channel\n")
			return
		}

		fmt.Printf("Leaving channel %v\n", fields[1])
		conn.Part(fields[1])
		target = ""
	case "/names":
		conn.Raw("NAMES")
	case "/quit":
		fmt.Printf("Quitting.\n")
		if len(fields) == 2 {
			conn.Quit(fields[1])
		} else {
			conn.Quit()
		}
		quitclient = true
	}
}

func ui() {
	quitclient = false
	for !quitclient {
		fmt.Printf("[%v] ", target)
		bio := bufio.NewReader(os.Stdin)
		line, err := bio.ReadString('\n')
		logmsg(time.Now(), conn.Me().Nick, target, line)

		if err != nil {
			log.Fatal("Couldn't get input.\n")
		}

		if line != "\n" {
			if line[0] == '/' {
				// A command
				parsecommand(line)
			} else {
				// Send line to target.
				conn.Privmsg(target, line)
			}
		}
	}
}

func main() {
	var err error
	var configfile = flag.String("config", "mub.yaml", "Path to configuration file")

	flag.Parse()

	conf, err := parseconfig(*configfile)
	if err != nil {
		log.Fatal("Couldn't parse configuration file")
	}

	file, err = os.OpenFile(conf.LogFile, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
	if err != nil {
		fmt.Printf("Could not open file '%s', will not write to log\n", conf.LogFile)
		fmt.Printf("Error Message: %s\n", err)
	}

	cfg := irc.NewConfig(conf.Nick)
	cfg.SSL = conf.TLS
	cfg.Server = conf.Server
	cfg.NewNick = func(n string) string { return n + "^" }
	cfg.Me.Ident = "elsabot"
	cfg.Me.Name = conf.RealName
	conn = irc.Client(cfg)

	// Join channel on connect.
	conn.HandleFunc("connected",
		func(conn *irc.Conn, line *irc.Line) {
			connected(conn, line)
		})

	conn.HandleFunc("disconnected",
		func(conn *irc.Conn, line *irc.Line) { quitclient = true })

	// Tell client to connect.
	fmt.Printf("Connecting to %v...\n", conf.Server)
	if err := conn.Connect(); err != nil {
		fmt.Printf("Connection error: %s\n", err)
	}

	conn.HandleFunc("privmsg",
		func(conn *irc.Conn, line *irc.Line) {
			handler(IRC_privmsg, line)
		})

	conn.HandleFunc("join",
		func(conn *irc.Conn, line *irc.Line) {
			handler(IRC_join, line)
		})

	conn.HandleFunc("part",
		func(conn *irc.Conn, line *irc.Line) {
			handler(IRC_part, line)
		})

	conn.HandleFunc("quit",
		func(conn *irc.Conn, line *irc.Line) {
			handler(IRC_quit, line)
		})

	conn.HandleFunc("353",
		func(conn *irc.Conn, line *irc.Line) {
			handler(IRC_names, line)
		})

	ui()

	fmt.Printf("Disconnected from server.\n")
}
