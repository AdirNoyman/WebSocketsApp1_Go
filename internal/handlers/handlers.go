package handlers

import (
	"github.com/CloudyKit/jet/v6"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var views = jet.NewSet(

	// Tell jet html rendering engine where are my html templates files restored
	jet.NewOSFileSystemLoader("./html"),
	// Util that replicates nodemon
	jet.InDevelopmentMode(),
)

// upgradeConnestion function - Upgrades regular web connection to a websocket connection
var upgradeConnestion = websocket.Upgrader{

	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// for security, set checking request's origin to true
	CheckOrigin: func(r *http.Request) bool { return true },
}

// Home renders the home page
func Home(w http.ResponseWriter, r *http.Request) {

	err := renderPage(w, "home.jet", nil)

	if err != nil {

		log.Println(err)
	}

}

// WsJsonResponse defines the response sent back from websocket
type WsJsonResponse struct {
	Action      string `json:"action"`
	Message     string `json:"message"`
	MessageType string `json:"message_type"`
}

// WebSocketEndPoint handler that upgrades connection to web socket
func WebSocketEndPoint(w http.ResponseWriter, r *http.Request) {

	ws, err := upgradeConnestion.Upgrade(w, r, nil)

	if err != nil {

		log.Println("Error upgrading to connection to websocket connection", err)

	}

	log.Println("Client connected to endpoint successfully 😎🤟")

	var response WsJsonResponse
	response.Message = `<em><small>Connected to server successfully</small></em>`

	// Marshall the response to JSON and send it back to the client
	err = ws.WriteJSON(response)
	if err != nil {

		log.Println("Error connecting to websocket", err)

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
