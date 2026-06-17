package server

import (
	"ProjetFilViolet/backend/api/config"
	"ProjetFilViolet/backend/api/models"
	"ProjetFilViolet/backend/game/deck"
	"ProjetFilViolet/backend/game/engine"
	"ProjetFilViolet/backend/game/money"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Autoriser toutes les origines pour le dev
	},
}

type Client struct {
	UserID   uint
	NickName string
	Conn     *websocket.Conn
	TableID  uint
}

type GameServer struct {
	clients    map[uint]*Client // userID -> Client
	tables     map[uint]*engine.Table
	tablesMu   sync.RWMutex
	register   chan *Client
	unregister chan *Client
	broadcast  chan uint // tableID
}

func NewGameServer() *GameServer {
	return &GameServer{
		clients:    make(map[uint]*Client),
		tables:     make(map[uint]*engine.Table),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan uint),
	}
}

func (gs *GameServer) Run() {
	// Goroutine pour gérer l'enregistrement/désenregistrement des clients
	go func() {
		for {
			select {
			case client := <-gs.register:
				gs.tablesMu.Lock()
				gs.clients[client.UserID] = client
				gs.tablesMu.Unlock()
				log.Printf("Client %s (User %d) connecté à la table %d", client.NickName, client.UserID, client.TableID)
				gs.broadcastTableState(client.TableID)

			case client := <-gs.unregister:
				gs.tablesMu.Lock()
				if _, ok := gs.clients[client.UserID]; ok {
					delete(gs.clients, client.UserID)
					client.Conn.Close()
					log.Printf("Client %s déconnecté", client.NickName)

					// Si le joueur est assis à une table, on le retire et on le recrédite
					if table, exists := gs.tables[client.TableID]; exists {
						chipsRefunded, err := table.RemovePlayer(client.UserID)
						if err == nil && chipsRefunded > 0 {
							// Effectuer le cash-out en DB
							_, dbErr := money.AccountMoneyManager(client.UserID, int64(chipsRefunded))
							if dbErr != nil {
								log.Printf("Erreur lors du cash-out de %s : %v", client.NickName, dbErr)
							} else {
								log.Printf("Cash-out réussi pour %s : +%d jetons", client.NickName, chipsRefunded)
							}
						}
					}
				}
				gs.tablesMu.Unlock()
				gs.broadcastTableState(client.TableID)
			}
		}
	}()

	// Goroutine pour vérifier les timeouts toutes les secondes
	go func() {
		for {
			time.Sleep(1 * time.Second)
			gs.tablesMu.Lock()
			for tableID, table := range gs.tables {
				if table.CheckTimeouts() {
					log.Printf("Joueur timed out à la table %d", tableID)
					// Sauvegarder les jetons en DB si la partie s'est arrêtée (Showdown)
					if table.Stage == engine.Showdown {
						gs.persistShowdownResults(table)
					}
					go gs.triggerBroadcast(tableID)
				}
			}
			gs.tablesMu.Unlock()
		}
	}()
}

func (gs *GameServer) triggerBroadcast(tableID uint) {
	gs.broadcastTableState(tableID)
}

func (gs *GameServer) persistShowdownResults(table *engine.Table) {
	// A la fin de la main (Showdown), on met à jour les montants des joueurs en base de données.
	// Dans notre modèle, Player.Money en DB représente la somme de jetons en cours sur la table.
	for _, p := range table.Players {
		var dbPlayer models.Player
		// Trouver l'entrée de joueur en DB
		err := config.DB.Where("user_id = ? AND table_id = ?", p.ID, table.ID).First(&dbPlayer).Error
		if err == nil {
			dbPlayer.Money = p.Chips
			config.DB.Save(&dbPlayer)
		}
	}
}

// ServeHTTP gère les connexions WebSockets et authentifie avec JWT
func (gs *GameServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	tokenStr := r.URL.Query().Get("token")
	tableIDStr := r.URL.Query().Get("table_id")

	if tokenStr == "" || tableIDStr == "" {
		http.Error(w, "Missing token or table_id", http.StatusBadRequest)
		return
	}

	tableIDVal, err := strconv.ParseUint(tableIDStr, 10, 32)
	if err != nil {
		http.Error(w, "Invalid table_id", http.StatusBadRequest)
		return
	}
	tableID := uint(tableIDVal)

	// Valider le JWT
	secret := os.Getenv("JWT_SECRET")
	token, err := jwt.ParseWithClaims(tokenStr, jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte(secret), nil
	})
	if err != nil || !token.Valid {
		http.Error(w, "Unauthorized: invalid token", http.StatusUnauthorized)
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		http.Error(w, "Unauthorized: invalid claims", http.StatusUnauthorized)
		return
	}

	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		http.Error(w, "Unauthorized: invalid user_id", http.StatusUnauthorized)
		return
	}
	userID := uint(userIDFloat)

	// Récupérer le pseudo depuis la DB
	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Upgrade websocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Upgrade error: %v", err)
		return
	}

	client := &Client{
		UserID:   userID,
		NickName: user.NickName,
		Conn:     conn,
		TableID:  tableID,
	}

	gs.tablesMu.Lock()
	// Récupérer ou créer la table de jeu
	if _, exists := gs.tables[tableID]; !exists {
		// Récupérer les blinds du salon depuis la DB si possible
		var dbTable models.Table
		var minBet uint64 = 10
		var maxBet uint64 = 20
		if err := config.DB.Preload("Saloon").First(&dbTable, tableID).Error; err == nil {
			var saloon models.Saloon
			if err := config.DB.First(&saloon, dbTable.SaloonID).Error; err == nil {
				minBet = saloon.MinBet
				maxBet = saloon.MaxBet
				if minBet == 0 {
					minBet = 10
				}
				if maxBet == 0 {
					maxBet = minBet * 2
				}
			}
		}
		gs.tables[tableID] = engine.NewTable(tableID, minBet, maxBet)
	}
	gs.tablesMu.Unlock()

	gs.register <- client

	// Boucle de lecture des messages
	go gs.readPump(client)
}

type ActionMessage struct {
	Action string `json:"action"`
	Amount uint64 `json:"amount"`
}

func (gs *GameServer) readPump(c *Client) {
	defer func() {
		gs.unregister <- c
	}()

	for {
		var msg ActionMessage
		err := c.Conn.ReadJSON(&msg)
		if err != nil {
			log.Printf("Erreur de lecture de %s: %v", c.NickName, err)
			break
		}

		gs.tablesMu.Lock()
		table, exists := gs.tables[c.TableID]
		if !exists {
			gs.tablesMu.Unlock()
			continue
		}

		var actionErr error
		switch msg.Action {
		case "buy_in":
			// Déduire l'argent en base de données de l'utilisateur
			// On appelle AccountMoneyManager avec un montant négatif
			if msg.Amount == 0 {
				msg.Amount = 500 // Montant par défaut
			}
			dbBalance, dbErr := money.AccountMoneyManager(c.UserID, -int64(msg.Amount))
			if dbErr != nil {
				actionErr = fmt.Errorf("solde insuffisant en DB : %v", dbErr)
			} else {
				log.Printf("%s a débité %d de sa DB (nouveau solde DB: %d)", c.NickName, msg.Amount, dbBalance)
				actionErr = table.AddPlayer(c.UserID, c.NickName, msg.Amount)
				if actionErr != nil {
					// Rembourser en DB si l'ajout à la table a échoué
					money.AccountMoneyManager(c.UserID, int64(msg.Amount))
				} else {
					// Créer le Player en DB pour l'historique/le suivi de la table
					var dbPlayer models.Player
					dbPlayer.UserID = c.UserID
					dbPlayer.TableID = c.TableID
					dbPlayer.InitialBet = msg.Amount
					dbPlayer.Money = msg.Amount
					config.DB.Save(&dbPlayer)
				}
			}

		case "start_game":
			actionErr = table.StartHand()

		case "leave":
			chips, err := table.RemovePlayer(c.UserID)
			if err == nil {
				if chips > 0 {
					// Rembourser l'argent en DB
					money.AccountMoneyManager(c.UserID, int64(chips))
				}
				// Retirer le Player de la table en DB
				config.DB.Where("user_id = ? AND table_id = ?", c.UserID, c.TableID).Delete(&models.Player{})
			}
			actionErr = err

		default:
			// Actions de poker classiques (fold, check, call, raise)
			actionErr = table.PlayAction(c.UserID, msg.Action, msg.Amount)
			// Si Showdown atteint après cette action, enregistrer les résultats
			if actionErr == nil && table.Stage == engine.Showdown {
				gs.persistShowdownResults(table)
			}
		}

		if actionErr != nil {
			log.Printf("Erreur d'action pour %s (%s) : %v", c.NickName, msg.Action, actionErr)
			c.Conn.WriteJSON(map[string]interface{}{
				"type":  "error",
				"error": actionErr.Error(),
			})
		} else {
			log.Printf("Action réussie pour %s : %s", c.NickName, msg.Action)
		}

		gs.tablesMu.Unlock()
		gs.broadcastTableState(c.TableID)
	}
}

// Structures sécurisées pour masquer les cartes des autres joueurs
type SecurePlayerState struct {
	ID                uint        `json:"id"`
	NickName          string      `json:"nick_name"`
	Chips             uint64      `json:"chips"`
	HoleCards         []deck.Card `json:"hole_cards,omitempty"` // Rempli uniquement pour le joueur concerné
	CurrentBet        uint64      `json:"current_bet"`
	TotalContribution uint64      `json:"total_contribution"`
	IsActive          bool        `json:"is_active"`
	IsAllIn           bool        `json:"is_all_in"`
	HasActed          bool        `json:"has_acted"`
}

type SecureTableState struct {
	ID                uint                `json:"id"`
	Players           []SecurePlayerState `json:"players"`
	CommunityCards    []deck.Card         `json:"community_cards"`
	Stage             string              `json:"stage"`
	DealerIdx         int                 `json:"dealer_idx"`
	SmallBlindIdx     int                 `json:"small_blind_idx"`
	BigBlindIdx       int                 `json:"big_blind_idx"`
	CurrentTurnIdx    int                 `json:"current_turn_idx"`
	CurrentBet        uint64              `json:"current_bet"`
	MinRaise          uint64              `json:"min_raise"`
	SmallBlindAmount  uint64              `json:"small_blind_amount"`
	BigBlindAmount    uint64              `json:"big_blind_amount"`
	GameInProgress    bool                `json:"game_in_progress"`
	Pots              []engine.Pot        `json:"pots"`
	SecondsRemaining  int                 `json:"seconds_remaining"`
	LastActionMessage string              `json:"last_action_message"`
}

func (gs *GameServer) broadcastTableState(tableID uint) {
	gs.tablesMu.Lock()
	table, exists := gs.tables[tableID]
	if !exists {
		gs.tablesMu.Unlock()
		return
	}

	// Préparer la liste des clients connectés à cette table
	var tableClients []*Client
	for _, client := range gs.clients {
		if client.TableID == tableID {
			tableClients = append(tableClients, client)
		}
	}
	gs.tablesMu.Unlock()

	// Envoyer l'état personnalisé à chaque client
	for _, client := range tableClients {
		secureState := gs.buildSecureState(table, client.UserID)
		err := client.Conn.WriteJSON(map[string]interface{}{
			"type":  "state",
			"state": secureState,
		})
		if err != nil {
			log.Printf("Erreur lors de l'envoi de l'état à %s: %v", client.NickName, err)
		}
	}
}

func (gs *GameServer) buildSecureState(table *engine.Table, forUserID uint) SecureTableState {
	securePlayers := make([]SecurePlayerState, 0, len(table.Players))
	for _, p := range table.Players {
		var cards []deck.Card
		// Révéler les cartes privées uniquement au propriétaire ou à la fin lors du Showdown
		if p.ID == forUserID || table.Stage == engine.Showdown {
			cards = p.HoleCards
		}

		securePlayers = append(securePlayers, SecurePlayerState{
			ID:                p.ID,
			NickName:          p.NickName,
			Chips:             p.Chips,
			HoleCards:         cards,
			CurrentBet:        p.CurrentBet,
			TotalContribution: p.TotalContribution,
			IsActive:          p.IsActive,
			IsAllIn:           p.IsAllIn,
			HasActed:          p.HasActed,
		})
	}

	secsRemaining := 0
	if table.GameInProgress && table.CurrentTurnIdx != -1 {
		elapsed := time.Since(table.TurnStartTime).Seconds()
		secsRemaining = table.ActionTimeoutSecs - int(elapsed)
		if secsRemaining < 0 {
			secsRemaining = 0
		}
	}

	return SecureTableState{
		ID:                table.ID,
		Players:           securePlayers,
		CommunityCards:    table.CommunityCards,
		Stage:             table.Stage.String(),
		DealerIdx:         table.DealerIdx,
		SmallBlindIdx:     table.SmallBlindIdx,
		BigBlindIdx:       table.BigBlindIdx,
		CurrentTurnIdx:    table.CurrentTurnIdx,
		CurrentBet:        table.CurrentBet,
		MinRaise:          table.MinRaise,
		SmallBlindAmount:  table.SmallBlindAmount,
		BigBlindAmount:    table.BigBlindAmount,
		GameInProgress:    table.GameInProgress,
		Pots:              table.Pots,
		SecondsRemaining:  secsRemaining,
		LastActionMessage: table.LastActionMessage,
	}
}
