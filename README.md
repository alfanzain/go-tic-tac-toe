# go-tic-tac-toe

A simple real-time multiplayer Tic Tac Toe game built with Go using WebSockets.
Designed for clarity, simplicity, and to help understand the fundamentals of real-time communication in a live web application.

---

## 🚀 Getting Started

```bash
git clone https://github.com/alfanzain/go-tic-tac-toe.git

cd go-tic-tac-toe
```

Run the Go files like this:

```bash
go run *.go
```

or with Makefile and Docker:

```bash
make run
```

Then, open `localhost:3000` in two separate browser windows or tabs.

---

## 🧠 Why I Built This

The main goal of this project is to **demonstrate real-time communication using WebSocket in Go**.
Rather than building something complex, this game serves as a minimal but complete example of:

* WebSocket connection handling
* Real-time data synchronization
* Managing game state across multiple clients

It's a teaching tool just as much as a game.

---

## 🛠️ How I Built It

* **Language**: [Go](https://go.dev/)
* **Communication**: WebSocket (via [Gorilla WebSocket](https://github.com/gorilla/websocket))
* **Game State**: Managed in-memory
* **Architecture**: Plain and messy
* **Frontend**: Plain HTML + CSS (Tailwind) + JS

You can see the version of the libraries on [go.mod](./go.mod)

---

## ⚙️ How It Works

The application uses **Go WebSocket** to manage real-time communication between players. Here's a breakdown of the flow:

### 🧩 Core Components

* **`main.go`**

  * Starts the HTTP server
  * Exposes REST API (if any) and WebSocket endpoints

* **Room List**

  * The server keeps a registry of active game rooms
  * Each room handles its own state and player communication

* **Room Lifecycle**

  * When the frontend creates a room, the server instantiates a new `Room` object
  * The room maintains:

    * A list of connected players
    * Game state (board, turn, etc.)
    * A goroutine to **broadcast messages** (game info, updates)

* **Player Connection**

  * When a player joins a room, a WebSocket connection is established
  * The server creates a `Player` object for them
  * Each player runs goroutines to:

    * **Read messages** from the WebSocket
    * **Write messages** to the WebSocket

### 📡 Channels

Two main Go channels handle server-client communication:

1. **Game Info Channel** – broadcasts general server messages (e.g., player joined, game start)
2. **Game Action Channel** – relays player actions like moves, quits

### 🧠 Game Logic

* All logic (turns, win/draw detection, move validation) is handled **on the server**
* The frontend is kept **dumb**: it sends player actions and renders what the server responds with
* After each player action, the server checks:

  * Has a player won? Checking the winning state every movement
  * Is the game a draw? After checking the winning state on last turn (9th) and still doesn't return a win state, automatically change the state into draw

### 🔚 Connection Handling

* When a player **closes or refreshes** the browser tab, the WebSocket disconnects and the server cleans up
* If **inactive for 15 minutes**, the server automatically closes the connection and clears the room


---

## 📌 Future Ideas

Although this project is for learning Websocket purpose only, I want to list the ideas that possibly I don't have time to make it. But who knows? 

* Better architecture

    Hexagonal architecture perhaps, learning purpose too

* Add player names

* Game lobby support

    The idea is to make the server can create the player for the first time than listed them on lobby. They can see the available room too. I believe it supports persistent player whenever the player changes the page from main menu to the game room and back

* Add chat functionality via WebSocket
* Add game history (stored temporarily in memory)
* Improve UI with animations or visual effects

* Player rematch options

    The player is not persistent. The frontend doesn't record the player info. When the player refresh the game room page, the Websocket creates a new one, makes previous session destroyed. It means the game can't be continued. Currently, the game stucks because the player will not listed as room players and the game doesn't continue automatically. I think it should not because can be exploited

    Before make this feature, I need to make sure the player is persistent

* Score tracking

    Win lose tracking so you can screenshot this and upload on your Linkedin as achievement

* AI opponent mode

    LOL this never gonna happen, I suppose.

--- 

## 🤝 Want to Contribute?

I don't expect you guys to contribute because this is not that worth to contribute

Unless, you insist want to learn with me. Just open issue, open PR, or contact me by email or Linkedin (the links on my profile)

---


## 📄 License

IDK if this necessary, I think I should put this on my future ideas list too

---
