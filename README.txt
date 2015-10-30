# mub - a Small IRC client

This is a small, bare-bones IRC client.

It's marginally usable in itself with a readline-like user interface.

If it is started with the -sub flag it can be used to create the core
of an IRC client with a richer user interfaces in front of it,
connected to stdout and stdin.

To connect to a server, use

  /tlsconnect server:port nick

or /connect for unencrypted connections.

The server clause in the configuration file limits what servers are
allowed to connect to.
