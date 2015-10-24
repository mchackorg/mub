// Elsa is a simple IRC bot who logs everything that's said on a
// channel and sends a backlog when asked.
package main

import "bufio"
import "fmt"
import "time"
import "os"
import "os/exec"

import irc "github.com/fluffle/goirc/client"

// Our log file.
var file *os.File

// Log text
func logmsg(time time.Time, nick string, text string) {
	line := time.UTC().Format("2006-01-02 15:04:05")
	line += " <" + nick + "> " + text + "\n"

	_, err := file.WriteString(line)
	if err != nil {
		fmt.Println(err)
	}
}

// Parse a numer from text. Returns number.
func parsenumber(text string) (lines int) {
	n, err := fmt.Sscanf(text, "%d", &lines)
	if err != nil {
		panic(err)
	}

	if n != 1 {
		lines = 0
	}

	return lines
}

// Wait for command cmd to exit so we can remove all its resources.
func wait(cmd *exec.Cmd) {
	cmd.Wait()
}

// Send the last lines to requester.
func printlast(conn *irc.Conn, lines int, nick string) {
	if lines <= 0 {
		return
	}

	linesarg := fmt.Sprintf("%v", lines)
	cmd := exec.Command("tail", "-n", linesarg, "./botlog.txt")

	output, err := cmd.StdoutPipe()
	if err != nil {
		panic(err)
	}

	err = cmd.Start()
	if err != nil {
		panic(err)
	}

	go wait(cmd)

	scanner := bufio.NewScanner(output)
	for scanner.Scan() {
		conn.Privmsg(nick, scanner.Text())
	}
}

// PRIVMSG handler.
func handlemsg(channel string, conn *irc.Conn, line *irc.Line) {
	switch line.Text() {
	case "!last":
		printlast(conn, 10, line.Nick)
		return
	case "!version":
		conn.Privmsg(line.Target(), "0.1")
		return
	case "!help":
		conn.Privmsg(line.Target(), "!version, !help, !last [N], where N is number of lines.")
		return
	}

	if line.Text()[:5] == "!last" {
		lines := parsenumber(line.Text()[6:])
		// Don't give more than 200 lines backlog.
		if lines > 200 {
			lines = 200
			conn.Privmsg(line.Nick, "I can't give you more than the last 200 lines:")
		}
		printlast(conn, lines, line.Nick)
		return
	}

	// Only log if this was said to the channel.
	if line.Target() == channel {
		logmsg(line.Time, line.Nick, line.Text())
	}
}

// Join a channel and do something.
func joinchannel(channel string, conn *irc.Conn, line *irc.Line) {
	conn.Join(channel)
	conn.Privmsg(channel, "Slå dig lös, slå dig fri!")
}

func main() {
	var err error

	channel := "#hacklunch"
	filename := "botlog.txt"

	file, err = os.OpenFile(filename, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
	if err != nil {
		panic(err)
	}

	cfg := irc.NewConfig("elsa")
	cfg.SSL = true
	cfg.Server = "chat.hack.org:4713"
	cfg.NewNick = func(n string) string { return n + "^" }
	cfg.Me.Ident = "elsabot"
	cfg.Me.Name = "Elsa"
	c := irc.Client(cfg)

	// Join channel on connect.
	c.HandleFunc("connected",
		func(conn *irc.Conn, line *irc.Line) { joinchannel(channel, conn, line) })
	// And a signal on disconnect
	quit := make(chan bool)
	c.HandleFunc("disconnected",
		func(conn *irc.Conn, line *irc.Line) { quit <- true })

	// Tell client to connect.
	if err := c.Connect(); err != nil {
		fmt.Printf("Connection error: %s\n", err)
	}

	// Handle messages.
	c.HandleFunc("PRIVMSG",
		func(conn *irc.Conn, line *irc.Line) { handlemsg(channel, conn, line) })

	// Wait for disconnect
	<-quit
}
