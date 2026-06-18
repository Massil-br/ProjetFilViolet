import React, { useEffect, useState, useRef } from 'react'
import { ArrowLeft, User, DollarSign, Play, LogOut, Check, ChevronUp } from 'lucide-react'

// Map suit integer to visual symbol and class
const SUITS = [
  { char: '♠', name: 'spades', colorClass: '' },
  { char: '♥', name: 'hearts', colorClass: 'red' },
  { char: '♦', name: 'diamonds', colorClass: 'red' },
  { char: '♣', name: 'clubs', colorClass: '' }
]

// Map rank integer to name
const RANKS = {
  11: 'J',
  12: 'Q',
  13: 'K',
  14: 'A'
}

function getCardRankStr(rank) {
  return RANKS[rank] || String(rank)
}

function Card({ card, hidden }) {
  if (hidden || !card) {
    return <div className="poker-card card-back" />
  }

  const suitInfo = SUITS[card.suit] || { char: '?', colorClass: '' }
  const rankStr = getCardRankStr(card.rank)

  return (
    <div className={`poker-card ${suitInfo.colorClass}`}>
      <div style={{ textAlign: 'left', lineHeight: 1 }}>{rankStr}</div>
      <div className="suit-icon">{suitInfo.char}</div>
    </div>
  )
}

export default function PokerTable({ token, user, tableId, onLeaveTable }) {
  const [tableState, setTableState] = useState(null)
  const [logs, setLogs] = useState([])
  const [wsConnected, setWsConnected] = useState(false)
  const [showBuyInModal, setShowBuyInModal] = useState(false)
  const [showWinnerModal, setShowWinnerModal] = useState(false)
  const [buyInAmount, setBuyInAmount] = useState(500)
  const [raiseAmount, setRaiseAmount] = useState(10)
  
  const wsRef = useRef(null)
  const logsEndRef = useRef(null)

  // Trigger winner modal on showdown
  useEffect(() => {
    if (tableState?.stage === 'Showdown') {
      setShowWinnerModal(true)
    } else {
      setShowWinnerModal(false)
    }
  }, [tableState?.stage])

  // Auto scroll logs
  useEffect(() => {
    logsEndRef.current?.scrollIntoView({ behavior: 'smooth' })
  }, [logs])

  // Setup WebSocket connection
  useEffect(() => {
    const wsUrl = `ws://${window.location.hostname}:8082/ws?token=${token}&table_id=${tableId}`
    const ws = new WebSocket(wsUrl)
    wsRef.current = ws

    ws.onopen = () => {
      setWsConnected(true)
      addLog("Connecté au serveur de jeu. Rejoignez un siège pour jouer !")
    }

    ws.onmessage = (event) => {
      try {
        const data = JSON.parse(event.data)
        if (data.type === 'state') {
          setTableState(data.state)
          
          // Add table logs from last action message
          if (data.state.last_action_message) {
            addLog(data.state.last_action_message)
          }
        } else if (data.type === 'error') {
          addLog(`Erreur: ${data.error}`, true)
        }
      } catch (err) {
        console.error("Error parsing WS message:", err)
      }
    }

    ws.onclose = () => {
      setWsConnected(false)
      addLog("Connexion au serveur de jeu fermée.")
    }

    ws.onerror = (err) => {
      console.error("WS error:", err)
      addLog("Erreur de connexion avec le serveur de jeu.", true)
    }

    return () => {
      ws.close()
    }
  }, [token, tableId])

  const addLog = (message, isError = false) => {
    if (!message) return
    setLogs(prev => {
      // Avoid duplicates of the last message
      if (prev.length > 0 && prev[prev.length - 1].text === message) {
        return prev
      }
      const newLogs = [...prev, { text: message, isError, time: new Date().toLocaleTimeString() }]
      // Garder uniquement les 30 derniers logs pour éviter d'agrandir la table
      if (newLogs.length > 30) {
        return newLogs.slice(newLogs.length - 30)
      }
      return newLogs
    })
  }

  const sendAction = (action, amount = 0) => {
    if (!wsRef.current || wsRef.current.readyState !== WebSocket.OPEN) {
      addLog("Impossible d'envoyer l'action : déconnecté du serveur", true)
      return
    }
    wsRef.current.send(JSON.stringify({
      action,
      amount: parseInt(amount, 10)
    }))
  }

  const handleBuyInSubmit = (e) => {
    e.preventDefault()
    sendAction('buy_in', buyInAmount)
    setShowBuyInModal(false)
  }

  const handleAction = (actionName) => {
    if (actionName === 'raise') {
      sendAction('raise', raiseAmount)
    } else {
      sendAction(actionName)
    }
  }

  // Calculate seat placements so that current user is ALWAYS at seat-0 (bottom)
  const getArrangedPlayers = () => {
    if (!tableState || !tableState.players) return []
    
    const playersList = tableState.players
    const myPlayerIndex = playersList.findIndex(p => p.id === user.id)

    // Arrange the list of players to start from current user
    const arranged = Array(7).fill(null)
    
    playersList.forEach((player, originalIdx) => {
      let seatIdx = originalIdx
      if (myPlayerIndex >= 0) {
        // Shift seats so my player is at seat-0
        seatIdx = (originalIdx - myPlayerIndex + 7) % 7
      }
      arranged[seatIdx] = player
    })

    return arranged
  }

  const arrangedPlayers = getArrangedPlayers()
  const myPlayerState = tableState?.players?.find(p => p.id === user.id)
  const isMyTurn = tableState?.game_in_progress && tableState?.current_turn_idx !== -1 && 
    tableState?.players[tableState.current_turn_idx]?.id === user.id

  // Calculate bet/blind markers
  const getBlindBadge = (playerIndex, table) => {
    if (playerIndex === table.dealer_idx) return <span className="role-badge d">D</span>
    if (playerIndex === table.small_blind_idx) return <span className="role-badge sb">SB</span>
    if (playerIndex === table.big_blind_idx) return <span className="role-badge bb">BB</span>
    return null
  }

  // Pots calculation
  const totalPotAmount = tableState?.pots?.reduce((sum, pot) => sum + pot.amount, 0) || 0

  // Raise controls constraints
  const minRaise = tableState?.min_raise || 10
  const maxRaise = myPlayerState?.chips || 1000
  
  useEffect(() => {
    if (tableState) {
      // Align default raise amount to min required
      setRaiseAmount(Math.min(minRaise, maxRaise))
    }
  }, [tableState, minRaise, maxRaise])

  return (
    <div className="game-container">
      {/* Board Felts */}
      <div className="board-area">
        {/* Header / Table status bar */}
        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', background: 'rgba(255,255,255,0.02)', padding: '0.75rem 1.25rem', borderRadius: '12px', border: '1px solid var(--border-color)' }}>
          <div style={{ display: 'flex', alignItems: 'center', gap: '0.75rem' }}>
            <button className="btn btn-secondary btn-icon-only" onClick={onLeaveTable} title="Quitter la table">
              <ArrowLeft size={18} />
            </button>
            <div>
              <h3 style={{ fontSize: '1.05rem', fontWeight: '700' }}>Table #{tableId}</h3>
              <span style={{ fontSize: '0.8rem', color: 'var(--text-muted)' }}>
                Blinds: {tableState?.small_blind_amount || 10} / {tableState?.big_blind_amount || 20}
              </span>
            </div>
          </div>

          <div style={{ display: 'flex', alignItems: 'center', gap: '1rem' }}>
            <span style={{ 
              display: 'inline-flex', 
              alignItems: 'center', 
              gap: '0.4rem', 
              fontSize: '0.85rem', 
              fontWeight: '600',
              color: wsConnected ? 'var(--accent-green-glow)' : 'var(--accent-red)' 
            }}>
              <span style={{ 
                width: 8, 
                height: 8, 
                borderRadius: '50%', 
                background: wsConnected ? 'var(--accent-green-glow)' : 'var(--accent-red)',
                boxShadow: wsConnected ? '0 0 8px var(--accent-green-glow)' : 'none'
              }}></span>
              {wsConnected ? 'En Ligne' : 'Hors Ligne'}
            </span>

            {myPlayerState ? (
              <button className="btn btn-danger" style={{ padding: '0.4rem 0.8rem', fontSize: '0.85rem' }} onClick={() => handleAction('leave')}>
                Quitter le Siège
              </button>
            ) : (
              <button className="btn btn-primary" style={{ padding: '0.4rem 0.8rem', fontSize: '0.85rem' }} onClick={() => setShowBuyInModal(true)}>
                Prendre une place
              </button>
            )}
          </div>
        </div>

        {/* Felt Green Canvas */}
        <div className="poker-felt">
          
          {/* Render 7 Seats */}
          {arrangedPlayers.map((player, seatIdx) => {
            const seatClass = `seat-${seatIdx}`
            
            // Map arranged seat index back to actual table player index
            const actualPlayerIndex = tableState && player
              ? tableState.players.findIndex(p => p.id === player.id)
              : -1

            const isCurrentTurn = tableState?.game_in_progress && tableState?.current_turn_idx === actualPlayerIndex

            if (!player) {
              return (
                <div key={seatIdx} className={`player-seat ${seatClass}`}>
                  <button className="empty-seat-btn" onClick={() => setShowBuyInModal(true)}>
                    <User size={18} style={{ marginBottom: '0.2rem' }} />
                    <span>Siège Vide</span>
                  </button>
                </div>
              )
            }

            return (
              <div 
                key={player.id} 
                className={`player-seat ${seatClass} ${player.is_active ? '' : 'folded'} ${isCurrentTurn ? 'active-turn' : ''}`}
              >
                {/* Visual Avatar Bubble */}
                <div className="player-avatar-card">
                  {isCurrentTurn && <div className="turn-timer-ring" />}
                  <span className="nick">{player.nick_name}</span>
                  <span className="chips">{player.chips}</span>
                  {actualPlayerIndex !== -1 && getBlindBadge(actualPlayerIndex, tableState)}
                </div>

                {/* Hand Cards */}
                {player.hole_cards && player.hole_cards.length > 0 && (
                  <div className="player-hand">
                    <Card card={player.hole_cards[0]} hidden={player.id !== user.id && tableState?.stage !== 'Showdown'} />
                    <Card card={player.hole_cards[1]} hidden={player.id !== user.id && tableState?.stage !== 'Showdown'} />
                  </div>
                )}
              </div>
            )
          })}

          {/* Render Bets (CurrentBets) placed on felt */}
          {tableState?.players?.map((player, idx) => {
            if (player.current_bet === 0) return null
            
            // Find which arranged seat this player sits in
            const arrangedIdx = arrangedPlayers.findIndex(p => p && p.id === player.id)
            if (arrangedIdx === -1) return null

            const betClass = `bet-${arrangedIdx}`
            return (
              <div key={`bet-${player.id}`} className={`seat-bet ${betClass}`}>
                <DollarSign size={10} />
                <span>{player.current_bet}</span>
              </div>
            )
          })}

          {/* Table Center (Board cards, pots, timer) */}
          <div className="table-center">
            
            {/* Main Pot */}
            <div className="pot-display">
              💰 Total Pot: {totalPotAmount}
            </div>

            {/* Community Cards */}
            <div className="community-cards">
              <Card card={tableState?.community_cards[0]} />
              <Card card={tableState?.community_cards[1]} />
              <Card card={tableState?.community_cards[2]} />
              <Card card={tableState?.community_cards[3]} />
              <Card card={tableState?.community_cards[4]} />
            </div>

            {/* Game Stage Label and Status */}
            {tableState?.game_in_progress && (
              <div style={{ background: 'rgba(0,0,0,0.6)', padding: '0.3rem 0.8rem', borderRadius: '4px', fontSize: '0.8rem', fontWeight: '600' }}>
                Étape: {tableState.stage} 
                {tableState.seconds_remaining > 0 && ` | Temps restant: ${tableState.seconds_remaining}s`}
              </div>
            )}

            {!tableState?.game_in_progress && (
              <div style={{ background: 'rgba(0,0,0,0.6)', padding: '0.4rem 1rem', borderRadius: '4px', fontSize: '0.85rem', fontWeight: '700', color: 'var(--accent-gold)' }}>
                En attente de joueurs...
              </div>
            )}
          </div>
        </div>

        {/* User control Action bar */}
        <div className="action-panel-bar">
          <div>
            {isMyTurn && (
              <span className="text-gold" style={{ fontSize: '0.95rem', fontWeight: '700', display: 'flex', alignItems: 'center', gap: '0.25rem' }}>
                ⚡ C'est votre tour de jouer !
              </span>
            )}
            {!isMyTurn && myPlayerState && (
              <span className="text-muted" style={{ fontSize: '0.85rem' }}>
                En attente des autres joueurs...
              </span>
            )}
            {!myPlayerState && (
              <span className="text-muted" style={{ fontSize: '0.85rem' }}>
                Vous regardez la partie en tant que spectateur.
              </span>
            )}
          </div>

          <div className="action-buttons-group">
            {/* Start game button */}
            {myPlayerState && !tableState?.game_in_progress && tableState?.players?.length >= 2 && (
              <button className="btn btn-green" onClick={() => handleAction('start_game')}>
                <Play size={16} /> Démarrer la partie
              </button>
            )}

            {/* Standard play action buttons */}
            {isMyTurn && (
              <>
                <button className="btn btn-secondary" onClick={() => handleAction('fold')}>
                  Se Coucher (Fold)
                </button>

                {/* Check: only allowed if current player's bet matches the table's current bet */}
                {myPlayerState.current_bet === tableState.current_bet ? (
                  <button className="btn btn-secondary" onClick={() => handleAction('check')}>
                    Parole (Check)
                  </button>
                ) : (
                  <button className="btn btn-primary" onClick={() => handleAction('call')}>
                    Suivre (Call) ({tableState.current_bet - myPlayerState.current_bet})
                  </button>
                )}

                {/* Raise inputs */}
                {myPlayerState.chips > 0 && (
                  <div style={{ display: 'flex', alignItems: 'center', gap: '0.5rem', background: 'rgba(0,0,0,0.2)', padding: '0.2rem 0.5rem', borderRadius: '8px', border: '1px solid var(--border-color)' }}>
                    <div className="bet-slider-container">
                      <input 
                        type="range" 
                        className="custom-slider" 
                        min={minRaise} 
                        max={maxRaise}
                        value={raiseAmount}
                        onChange={(e) => setRaiseAmount(parseInt(e.target.value, 10))}
                      />
                      <input 
                        type="number"
                        className="form-input bet-amount-input" 
                        min={minRaise} 
                        max={maxRaise}
                        value={raiseAmount}
                        onChange={(e) => setRaiseAmount(Math.min(maxRaise, Math.max(minRaise, parseInt(e.target.value, 10) || 0)))}
                      />
                    </div>
                    <button className="btn btn-primary" style={{ padding: '0.5rem 1rem' }} onClick={() => handleAction('raise')}>
                      Relancer (Raise) ({raiseAmount})
                    </button>
                  </div>
                )}
              </>
            )}
          </div>
        </div>
      </div>

      {/* Logs / Chat Feed sidebar */}
      <div className="side-column">
        <div className="glass-panel chat-panel">
          <div className="chat-header">
            <h3>📝 Log des Actions</h3>
          </div>
          
          <div className="chat-feed">
            {logs.map((log, index) => (
              <div key={index} className={`log-entry ${log.isError ? 'error-log' : ''}`}>
                <span className="text-muted" style={{ fontSize: '0.75rem', marginRight: '0.4rem' }}>
                  [{log.time}]
                </span>
                <span>{log.text}</span>
              </div>
            ))}
            <div ref={logsEndRef} />
          </div>

          <div style={{ background: 'rgba(0,0,0,0.1)', padding: '0.75rem', borderRadius: '8px', border: '1px solid var(--border-color)' }}>
            <h4 style={{ fontSize: '0.85rem', color: 'var(--accent-gold)', marginBottom: '0.25rem' }}>Détails de la Main</h4>
            <div style={{ fontSize: '0.8rem', display: 'flex', flexDirection: 'column', gap: '0.25rem' }}>
              <div style={{ display: 'flex', justifyContent: 'space-between' }}>
                <span className="text-muted">Étape active:</span>
                <span>{tableState?.stage || 'Attente'}</span>
              </div>
              <div style={{ display: 'flex', justifyContent: 'space-between' }}>
                <span className="text-muted">Mise en cours:</span>
                <span>{tableState?.current_bet || 0} jetons</span>
              </div>
              <div style={{ display: 'flex', justifyContent: 'space-between' }}>
                <span className="text-muted">Relance Min:</span>
                <span>{tableState?.min_raise || 0} jetons</span>
              </div>
            </div>
          </div>
        </div>
      </div>

      {/* Buy In modal */}
      {showBuyInModal && (
        <div className="modal-overlay">
          <div className="glass-panel modal-content">
            <div className="modal-header">
              <h3>💵 Choisir votre Cavée (Buy-In)</h3>
              <button 
                style={{ background: 'transparent', border: 'none', color: '#fff', fontSize: '1.2rem', cursor: 'pointer' }}
                onClick={() => setShowBuyInModal(false)}
              >
                &times;
              </button>
            </div>
            <form onSubmit={handleBuyInSubmit}>
              <div className="form-group">
                <label>Montant en jetons (Max: {user.money})</label>
                <input 
                  type="number" 
                  className="form-input" 
                  min="20" 
                  max={user.money}
                  value={buyInAmount}
                  onChange={(e) => setBuyInAmount(Math.min(user.money, Math.max(20, parseInt(e.target.value, 10) || 0)))}
                />
              </div>
              <div className="modal-actions">
                <button type="button" className="btn btn-secondary" onClick={() => setShowBuyInModal(false)}>
                  Annuler
                </button>
                <button type="submit" className="btn btn-primary">
                  S'asseoir à la Table
                </button>
              </div>
            </form>
          </div>
        </div>
      )}

      {/* Winner Announcement Modal */}
      {showWinnerModal && (
        <div className="modal-overlay">
          <div className="glass-panel modal-content" style={{ textAlign: 'center', maxWidth: '450px', border: '2px solid var(--accent-gold)' }}>
            <div style={{ fontSize: '3rem', marginBottom: '0.5rem' }}>🏆</div>
            <h2 className="text-gold" style={{ marginBottom: '1rem', fontSize: '1.8rem' }}>Fin de la Main !</h2>
            
            <div style={{ 
              background: 'rgba(0, 0, 0, 0.3)', 
              padding: '1.25rem', 
              borderRadius: '12px', 
              marginBottom: '1.5rem',
              lineHeight: '1.5',
              fontSize: '1.05rem',
              fontWeight: '600'
            }}>
              {tableState?.last_action_message}
            </div>

            <div className="modal-actions" style={{ justifyContent: 'center', gap: '1rem', marginTop: '1rem' }}>
              <button className="btn btn-secondary" onClick={() => setShowWinnerModal(false)}>
                Voir la table
              </button>
              {myPlayerState && !tableState?.game_in_progress && tableState?.players?.length >= 2 && (
                <button className="btn btn-green" onClick={() => { sendAction('start_game'); setShowWinnerModal(false); }}>
                  Main suivante
                </button>
              )}
            </div>
          </div>
        </div>
      )}
    </div>
  )
}
