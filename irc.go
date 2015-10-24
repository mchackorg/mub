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

import irc "github.com/fluffle/goirc/client"

type Config struct {
	Nick     string
	RealName string
	Server   string
	Port     int
	LogFile  string
}

var file *os.File

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
func logmsg(time time.Time, nick string, text string) {
	line := time.UTC().Format("2006-01-02 15:04:05")
	line += " <" + nick + "> " + text + "\n"

	_, err := file.WriteString(line)
	if err != nil {
		fmt.Println(err)
	}
}

// PRIVMSG handler.
func handlemsg(channel string, conn *irc.Conn, line *irc.Line) {
	time := line.Time.Format("15:04:05")
	fmt.Printf("%v <%v> %v\n", time, line.Nick, line.Text())

	// Only log if this was said to the channel.
	if line.Target() == channel {
		logmsg(line.Time, line.Nick, line.Text())
	}
}

func connected(conn *irc.Conn, line *irc.Line) {
	fmt.Printf("Connected.\n")
	fmt.Printf("Joining #larshack\n")
	conn.Join("#larshack")
}

func main() {
	var err error
	var configfile = flag.String("config", "mub.yaml", "Path to configuration file")

	flag.Parse()

	conf, err := parseconfig(*configfile)
	if err != nil {
		log.Fatal("Couldn't parse configuration file")
	}

	channel := "#larshack"

	file, err = os.OpenFile(conf.LogFile, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
	if err != nil {
		panic(err)
	}

	cfg := irc.NewConfig(conf.Nick)
	cfg.SSL = false
	cfg.Server = conf.Server
	cfg.NewNick = func(n string) string { return n + "^" }
	cfg.Me.Ident = "elsabot"
	cfg.Me.Name = conf.RealName
	c := irc.Client(cfg)

	// Join channel on connect.
	c.HandleFunc("connected",
		func(conn *irc.Conn, line *irc.Line) {
			connected(conn, line)
		})
	// And a signal on disconnect
	quit := make(chan bool)
	c.HandleFunc("disconnected",
		func(conn *irc.Conn, line *irc.Line) { quit <- true })

	// Tell client to connect.
	fmt.Printf("Connecting to %v...\n", conf.Server)
	if err := c.Connect(); err != nil {
		fmt.Printf("Connection error: %s\n", err)
	}

	// Handle messages.
	c.HandleFunc("PRIVMSG",
		func(conn *irc.Conn, line *irc.Line) { handlemsg(channel, conn, line) })

	for {
		fmt.Printf("[%v] ", channel)
		bio := bufio.NewReader(os.Stdin)
		line, _, _ := bio.ReadLine()
		if string(line) != "" {
			c.Privmsg(channel, string(line))
		}
	}

	// Wait for disconnect
	<-quit
	fmt.Printf("Disconnected from server.\n")
}
