# waffle
https://waffle-chat-demo.herokuapp.com/

<img src="https://user-images.githubusercontent.com/2475286/51789238-bc33ed80-2154-11e9-9a37-9874ab7adf77.png" width="400">

## Overview
A simple chat application that enables a user to chat with other anonymous
users.

Users are identified by a randomly-generated client ID (similar to a session ID)
that is saved in Local Storage

Only the most recent 4096 messages are saved. Once this capacity is reached,
earlier message will be truncated

#### Project layout
The project is roughly structured according to the standard Go project layout as
defined in this project: https://github.com/golang-standards/project-layout

- `cmd/` - Main entry point/runner
- `pkg/` - Application packages

#### Name
It is thought that the origin of the word "gopher" may be rooted in the French
word "gaufre", meaning "waffle". The backend for the project is written in Go,
and as such I felt it was appropriate to reference the Go mascot in the name.


#### Stack
I viewed this project as a good opportunity to learn about a couple topics and
technologies that I did not have much prior experience with:

- Go
- Server Sent Events (SSE)
- Docker

## Server
The chat server is written in [Go](https://golang.org/) and utilizes
[Server Sent Events (SSE)](https://developer.mozilla.org/en-US/docs/Web/API/Server-sent_events) to
broadcast new chat messages to connected clients. While WebSockets are typically
used for browser chat applications as they support two-way messaging, I planned
to incorporate some sort of persistence for broadcasted messages, and figured it would be a
fine application to create the messages via a POST endpoint, then broadcast
using SSEs.

I considered using GraphQL for the API layer, however ultimately decided this
was overkill as the server just needed to serve 2 endpoints that both had a very
well defined, simple schema that did not require a lot of flexibility. Also, due
to time constraints, it made sense to keep the client code as simple and
lightweight as possible.

#### Notes
- This app in its current state is not able to be distributed. Messages are
  stored in memory for the duration of the currently running process. To scale
  the application, a shared and durable data store would need to be introduced, as
  well as a message queue.
- SSE supports multiple named message types and message IDs. TODO: classify
  events (ie. 'message.created') to allow for handling them cleanly on the
  client

## User Interface
The UI is a simple, vanilla JS class written using ES6 language features.
There is no transpiling, polyfilling, build process, or other trickery in the front
end, so the browser must compatible with ES6 and SSE. Twitter Bootstrap is used
for some basic styling, and the UI is meant to resemble messaging apps like
iMessage or Facebook Messenger.

### Notes
- Browser must support
  - Server Sent Events (https://caniuse.com/#feat=eventsource)
  - Web Cryptography (SubtleCrypto) (https://caniuse.com/#feat=cryptography)
  - ES6 classes (https://caniuse.com/#feat=es6-class)

## Development

### Prerequisites
- A Linux based system
- Docker and/or Go with dep installed

### Setup
If you have Docker installed and just want to run the app locally:
1. `git clone https://github.com/thejchap/waffle`
2. `cd waffle`
3. `make`

For local development
1. Make sure your `$GOPATH` is set and you have [dep](https://github.com/golang/dep) installed
2. `go get github.com/thejchap/waffle/cmd/waffle`
3. `cd $GOPATH/src/github.com/thejchap/waffle`
4. `dep ensure`
5. `go run cmd/waffle/main.go`

### Documentation
1. The doc server can be run using `godoc -http=:6060` (assuming go is
   installed) in the application directory
2. Visit http://localhost:6060/pkg/github.com/thejchap/waffle/

## License

Copyright 2019 Justin Chapman

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
