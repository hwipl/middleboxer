# middleboxer

## Overview

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
