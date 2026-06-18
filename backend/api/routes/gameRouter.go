package routes

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

// Configuration du WebSocket
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Attention: en production, il faudra restreindre ça
	},
}

// Fonction pour initialiser la route WebSocket
func InitGameRoutes(e *echo.Echo) {
	// On ne met pas le AuthMiddleware classique ici, on passera le token dans la requête WS
	e.GET("/ws/game", HandleWebSocket)
}

func HandleWebSocket(c echo.Context) error {
	// Upgrade la connexion HTTP en WebSocket
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		log.Println("Erreur lors de l'upgrade WebSocket:", err)
		return err
	}
	defer ws.Close()

	// Boucle d'écoute pour ce joueur
	for {
		// Lire le message envoyé par le mobile
		var msg map[string]interface{}
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Println("Joueur déconnecté ou erreur de lecture:", err)
			break
		}

		log.Printf("Message reçu du mobile : %v\n", msg)

		// Exemple de réponse du serveur : "J'ai bien reçu ton action"
		response := map[string]interface{}{
			"type":    "SERVER_ACK",
			"message": "Action reçue, en attente des autres joueurs...",
		}

		err = ws.WriteJSON(response)
		if err != nil {
			log.Println("Erreur d'écriture:", err)
			break
		}
	}
	return nil
}
