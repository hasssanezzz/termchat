# Termchat

A simple terminal-based chat application with rooms based on a TCP server-client architecture.

## Features
* Multiple rooms for chatting
* Commands for managing rooms and interactions
* Help message for command reference

## Usage
1. __Starting the Server__ Run the server using the following command:
```bash
go run server.go
```
2. __Connecting to the Server__ Clients can connect to the server using TCP connection. For example, using telnet:

```bash
telnet localhost 8080
```

3. __Available Commands__ The following commands are available for clients:

```
	:help                   display this help message
	:info                   display client's name and joined rooms
	:name <name>            set client's name
	:join <room>            join a room
	:leave <room>           leave a room
	:exit                   disconnect
	to:<room> <message>     send a message to a room
```
Example usage:


```
:name Alice
:join room1
you have joined room "room1"
to:room1 Hello, everyone!
[Alice]:[room1] >> Hello, everyone!
:info
============ Client info ==============
Name: "Alice"
Rooms:
room1
============ Client info ==============
:leave room1
you have left room "room1"
```

### Notes
* The server uses a simple TCP connection for communication.
* Clients can interact with the server using commands prefixed with a colon (:) or the to: syntax for sending messages to specific rooms.
* Use the :help command for detailed command usage and information.
* The app maybe full of bugs :7

## Contributing
I welcome contributions to improve this chat application. If you're interested in contributing, here are some areas where help is needed:
* Code Clean-Up: There are several TODO comments in the code that need attention. Feel free to address them by implementing the required functionality or optimizing existing code.
* Enhancements: If you have ideas for new features or improvements, please share them via GitHub issues or contribute directly to the codebase.
* Validation: My string manipulation might not be my best thing :)