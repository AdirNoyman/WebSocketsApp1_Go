package handlers

import (
	"fmt"
	"github.com/CloudyKit/jet/v6"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

// webSocketChannel - The channel will except only data type WsPayload
var webSocketChannel = make(chan WsPayload)

// clients - a map of websocket connected clients (every client has his own websocket connection)
var clients = make(map[WebSocketConnection]string)

var views = jet.NewSet(

	// Tell jet html rendering engine where are my html templates files restored
	jet.NewOSFileSystemLoader("./html"),
	// Util that replicates nodemon
	jet.InDevelopmentMode(),
)

// upgradeConnection function - Upgrades regular web connection to a websocket connection
var upgradeConnection = websocket.Upgrader{

	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// for security, set checking request's origin to true
	CheckOrigin: func(r *http.Request) bool { return true },
}

// Home renders and displays the home page
func Home(w http.ResponseWriter, r *http.Request) {

	err := renderPage(w, "home.jet", nil)

	if err != nil {

		log.Println(err)
	}

}

// WebSocketConnection data structure that holds the data about the websocket connection
type WebSocketConnection struct {
	*websocket.Conn
}

// WsJsonResponse defines the response data structure sent back from websocket
type WsJsonResponse struct {
	Action      string `json:"action"`
	Message     string `json:"message"`
	MessageType string `json:"message_type"`
}

// WsPayload is the data structure of the data we send to the websocket
type WsPayload struct {
	// Action like 'Joined' (the chat) for example
	Action   string `json:"action"`
	Username string `json:"username"`
	Message  string `json:"message"`
	// When you put "-" it means don't show it in the json file
	Conn WebSocketConnection `json:"-"`
}

// WebSocketEndPoint handler that upgrades connection to web socket and sends a response back to the client
func WebSocketEndPoint(w http.ResponseWriter, r *http.Request) {

	ws, err := upgradeConnection.Upgrade(w, r, nil)

	if err != nil {

		log.Println("Error upgrading to connection to websocket connection", err)

	}

	log.Println("Client connected to endpoint successfully 😎🤟")

	var response WsJsonResponse
	response.Message = `<em><small>Connected to server successfully</small></em>`

	// Add the user who connects to the websockets connections map
	conn := WebSocketConnection{Conn: ws}
	clients[conn] = ""
	// Marshall the response to JSON and send it back to the client
	err = ws.WriteJSON(response)
	if err != nil {

		log.Println("Error connecting to websocket", err)

	}

	// Listen for WsPayload, and if WsPayload is sent to the channel we are listening for, we send it to the webSocketChannel
	go ListenForWS(&conn)

}

// ListenForWS is a Go routine. This function will run on a concurrently on separate thread
func ListenForWS(conn *WebSocketConnection) {
	// defer function is recovery function, that is called if the parent function('ListenForWS') dies for whatever reason
	defer func() {

		if err := recover(); err != nil {

			log.Println("Error", fmt.Sprintf("%v", err))
		}

	}()

	var payload WsPayload

	// Run forever and keep reading the payload that is sent by the user requests to the websocket connection
	for {

		err := conn.ReadJSON(payload)
		if err != nil {
			// do nothing
		} else {

			payload.Conn = *conn
			// Send the payload to my webSocketChannel
			webSocketChannel <- payload
		}
	}
}

func ListenToWebSocketChannel() {

	var response WsJsonResponse

	for {

		event := <-webSocketChannel
		response.Action = "Got here 😎🤟"
		response.Message = fmt.Sprintf("Some message 🙄... and action was %s", event.Action)
		// Broadcast a response to all the connected users. simultaneously
		broadcastMsgToAllConnectedClients(response)
	}
}

func broadcastMsgToAllConnectedClients(response WsJsonResponse) {

	for client := range clients {

		err := client.WriteJSON(response)
		// Error example -> user already left the connection
		if err != nil {

			log.Println("Websocket error")
			// Close there connection
			_ = client.Close()
			// Delete that user's connection from the map
			delete(clients, client)

		}
	}

}

// Render HTML pages function controller/handler
// tmpl = html template we want to render
// data = the model data we want mount on our html template. This is optional (you don't have to pass data to the template)
func renderPage(w http.ResponseWriter, tmpl string, data jet.VarMap) error {

	view, err := views.GetTemplate(tmpl)

	if err != nil {
		log.Println("Error #1 rendering the html template:", err)
		return err
	}

	err = view.Execute(w, data, nil)
	if err != nil {
		log.Println("Error #2 rendering the html template:", err)
		return err
	}
	// There was no error
	return nil
}
