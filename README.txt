# mub - a Small IRC client

This is a small, bare-bones IRC client. It can be used to create the
core of an IRC client with more rich user interfaces in front of it,
connected to stdout and stdin.

To connect to a server, use

  /connect server:port nick

TLS will be used if turned on in the configuration file.

The server clause in the configuration file limits what servers are
allowed to connect to.
