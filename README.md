# chatwithme
A verrrry basic go command line chat app. Inspired by @corya14!

# Features
  - Clients and Host can chat with each other
  - Host supports multiple client connections
  - Host repeats messages sent from one client to all others
  - ...that's about it!

# Usage
`chatwithme -s <port>` to start a server

`chatwithme -c <host:port>` to start a client

# Known Bugs
This program was thrown together in a little over an hour. There are minimal comments and documentation. In addition to that...
  - Clients' name are not properly propogated to each when forwarded through the server (shows up as if the server sent the msg)
  - There is no concept of timeouts in any part of the program
  - ...many more!
