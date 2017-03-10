// Small cross-platform IRC client.

package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	irc "github.com/fluffle/goirc/client"
	yaml "gopkg.in/yaml.v2"
)

// Config is the IRC client configuration.
type Config struct {
	Server          string
	Nick            string
	LogFile         string
	Ident           string
	RealName        string
	TLS             bool
	Pass            string
	BlockedCommands map[string]string
	AllowServers    map[string]string
}

var (
	file       *os.File
	conf       *Config
	conn       *irc.Conn
	currtarget string
	quitclient bool
)

// Parse the configuration file. Returns the configuration.
func parseconfig(filename string) (conf *Config, err error) {
	conf = new(Config)

	conf.Ident = "mub"
	conf.RealName = "unknown"
	contents, err := ioutil.ReadFile(filename)
	if err != nil {
		return
	}

	if err = yaml.Unmarshal(contents, &conf); err != nil {
		return
	}

	return
}

// Log text
func logmsg(time time.Time, nick string, target string, text string, action bool) {
	line := time.UTC().Format("2006-01-02 15:04:05")
	if target != "" {
		line += " " + target + " "
	}
	text = strings.TrimRight(text, "\r\n")
	if action {
		line += fmt.Sprintf("[%v %v]\n", nick, text)
	} else {
		line += fmt.Sprintf("<%v> %v\n", nick, text)
	}

	if file != nil {
		_, err := file.WriteString(line)
		if err != nil {
			msg := fmt.Sprintf("Couldn't write to log file: %s", err)
			errormsg(msg)
		}
	}
}

func connect(server string, nickname string, pass string, usetls bool) bool {
	var tlsconfig tls.Config

	// Check if we're allowed to connect to this server
	if conf.AllowServers != nil {
		if _, val := conf.AllowServers[server]; !val {
			info("This server is blocked by configuration.")
			return false
		}
	}

	cfg := irc.NewConfig(nickname)
	// Don't recover any crashes.
	cfg.Recover = func(c *irc.Conn, l *irc.Line) {}
	cfg.SSL = usetls
	tlsconfig.InsecureSkipVerify = true
	cfg.SSLConfig = &tlsconfig
	cfg.Server = server
	cfg.NewNick = func(n string) string { return n + "^" }
	cfg.Me.Ident = conf.Ident
	cfg.Me.Name = conf.RealName
	cfg.Pass = pass

	conn = irc.Client(cfg)

	conn.HandleFunc("connected",
		func(conn *irc.Conn, line *irc.Line) {
			connected(conn.Me().Nick)
		})

	conn.HandleFunc("disconnected",
		func(conn *irc.Conn, line *irc.Line) { quitclient = true })

	conn.HandleFunc("privmsg",
		func(conn *irc.Conn, line *irc.Line) {

			msg(line.Time, line.Nick, line.Args[0], line.Text(), false)
			logmsg(line.Time, line.Nick, line.Args[0], line.Text(), false)
		})

	conn.HandleFunc("action",
		func(conn *irc.Conn, line *irc.Line) {
			msg(line.Time, line.Nick, line.Args[0], line.Text(), true)
			logmsg(line.Time, line.Nick, line.Args[0], line.Text(), true)
		})

	conn.HandleFunc("join",
		func(conn *irc.Conn, line *irc.Line) {
			joined(line.Nick, line.Target())
		})

	conn.HandleFunc("part",
		func(conn *irc.Conn, line *irc.Line) {
			parted(line.Nick, line.Target())
		})

	conn.HandleFunc("quit",
		func(conn *irc.Conn, line *irc.Line) {
			quit(line.Nick)
		})

	conn.HandleFunc("nick",
		func(conn *irc.Conn, line *irc.Line) {
			nick(line.Nick, line.Args[0])
		})

	conn.HandleFunc("notice",
		func(conn *irc.Conn, line *irc.Line) {
			notice(line.Nick, line.Args[1])
		})

	conn.HandleFunc("353",
		func(conn *irc.Conn, line *irc.Line) {
			members(line.Args[2], line.Args[3])
		})

	conn.HandleFunc("311",
		func(conn *irc.Conn, line *irc.Line) {
			whois(line.Args[1], line.Args[5], line.Args[2], line.Args[3])
		})

	connecting(server, usetls)

	if err := conn.Connect(); err != nil {
		connectionerror(err)
		return false
	}

	return true
}

func greeting() {
	fmt.Printf("Welcome to IRC!\n")
	fmt.Printf("Type /help to list available commands.\n")
	fmt.Printf("Use /tlsconnect server nick [password] to connect to an IRC server.\n")
	fmt.Printf("Then use /join #channelname to join the chat rooms.\n")
}

func main() {
	var err error
	var configfile = flag.String("config", "mub.yaml", "Path to configuration file")
	var subprocess = flag.Bool("sub", false, "Run as subprocess without prompt and readline.")

	flag.Parse()

	greeting()

	conf, err = parseconfig(*configfile)
	if err != nil {
		fmt.Printf("Couldn't parse configuration file %v\n", *configfile)
	}

	if conf != nil && conf.LogFile != "" {
		file, err = os.OpenFile(conf.LogFile, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
		if err != nil {
			cantopenfile(conf.LogFile, err)
		}
	}

	rl, bio := initUI(*subprocess)

	if conf != nil && conf.Server != "" && conf.Nick != "" {
		connect(conf.Server, conf.Nick, conf.Pass, conf.TLS)
	}

	ui(*subprocess, rl, bio)

	disconnected()
}
