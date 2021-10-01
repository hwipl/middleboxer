# middleboxer

middleboxer is a command line tool for testing middleboxes like firewalls. It
uses cooperating client instances that send packets through a middlebox and
checks which packets pass through the middlebox and which do not.

## Overview

A test setup using middleboxer is depicted in the figure below:

```
+----------+     +-----------+     +----------+
| Client 1 |---->| Middlebox |---->| Client 2 |
+----------+     +-----------+     +----------+
     |_______________    _______________|
                     |  |
                  +--------+
                  | Server |
                  +--------+
```

The `middlebox` is the device under test, e.g., a firewall. Two middleboxer
clients, `client 1` and `client 2`, are connected to a middleboxer `server`.
Both clients are placed in such a way that they communicate with each other
through the middlebox.

The server manages the properties of packets, e.g., IP addresses, layer 4
protocol and ports, that should be sent through the middlebox by the clients.
It instructs one client to send those packets and the other client to receive
them.

The clients report back to the server for each packet if the packet passed
through the middlebox or if they received error messages like ICMP errors or
TCP resets.

The server collects all results and prints them to the console or writes them
to a file.

## Usage

```
Usage of middleboxer:
  -address string
        set address to connect to (client mode) or listen on (server mode)
  -diffs
        show packet diffs in results
  -id uint
        set id of the client (default 1)
  -out string
        set output file
  -ports string
        set port range to be tested (default "1:65535")
  -prot uint
        set layer 4 protocol to 6 (tcp) or 17 (udp) (default 6)
  -rdev string
        set device of the receiving client
  -rdip string
        set destination IP of the receiving client
  -rdmac string
        set destination MAC of the receiving client
  -rid uint
        set id of the receiving client (default 2)
  -rsip string
        set source IP of the receiving client
  -rsmac string
        set source MAC of the receiving client
  -sdev string
        set device of the sending client
  -sdip string
        set destination IP of the sending client
  -sdmac string
        set destination MAC of the sending client
  -server
        run as server (default: run as client)
  -sid uint
        set id of the sending client (default 1)
  -ssip string
        set source IP of the sending client
  -ssmac string
        set source MAC of the sending client
  -ssport uint
        set source port of the sending client
```
