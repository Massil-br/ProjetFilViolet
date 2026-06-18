package main

import (
	"ProjetFilViolet/backend/api/config"
	"ProjetFilViolet/backend/game/server"
	"fmt"
	"log"
	"net/http"
)

func main() {
	// Initialiser la configuration et la connexion DB (partagée avec backend/api)
	config.Init()

	if config.DB == nil {
		log.Fatal("❌ Impossible de démarrer : la base de données n'est pas initialisée")
	}
	log.Println("✅ Connexion à la base de données initialisée pour le serveur de jeu")

	// Créer et démarrer le serveur de jeu
	gameServer := server.NewGameServer()
	gameServer.Run()

	http.Handle("/ws", gameServer)

	port := ":8082"
	fmt.Printf("♠️ ♥️ Server de jeu démarré sur le port %s ♦️ ♣️\n", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal("ListenAndServe error:", err)
	}
}
