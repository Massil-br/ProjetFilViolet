import React, { useEffect, useState } from 'react'

const API_URL = `http://${window.location.hostname}:8081`;

export default function Lobby({ token, user, onJoinTable }) {
  const [saloons, setSaloons] = useState([])
  const [tables, setTables] = useState([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState(null)

  // Creation State
  const [showCreateForm, setShowCreateForm] = useState(false)
  const [createData, setCreateData] = useState({
    name: '',
    minBet: 10,
    maxBet: 200,
    playerCount: 7
  })
  const [createError, setCreateError] = useState(null)

  const fetchData = async () => {
    setLoading(true)
    setError(null)
    try {
      // Fetch Saloons
      const saloonsRes = await fetch(`${API_URL}/api/saloons`, {
        headers: { 'Authorization': `Bearer ${token}` }
      })
      if (!saloonsRes.ok) {
        const errorText = await saloonsRes.text()
        throw new Error(`Erreur lors de la récupération des salons (${saloonsRes.status} : ${errorText})`)
      }
      const saloonsData = await saloonsRes.json()
      setSaloons(saloonsData || [])

      // Fetch Tables
      const tablesRes = await fetch(`${API_URL}/api/tables`, {
        headers: { 'Authorization': `Bearer ${token}` }
      })
      if (!tablesRes.ok) {
        const errorText = await tablesRes.text()
        throw new Error(`Erreur lors de la récupération des tables (${tablesRes.status} : ${errorText})`)
      }
      const tablesData = await tablesRes.json()
      setTables(tablesData || [])

    } catch (err) {
      setError(err.message)
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    fetchData()
  }, [token])

  const handleCreateSaloon = async (e) => {
    e.preventDefault()
    setCreateError(null)

    try {
      const response = await fetch(`${API_URL}/api/saloons`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${token}`
        },
        body: JSON.stringify({
          name: createData.name,
          player_count: parseInt(createData.playerCount),
          min_bet: parseInt(createData.minBet),
          max_bet: parseInt(createData.maxBet)
        })
      })

      const data = await response.json()
      if (!response.ok) {
        throw new Error(data.error || 'Seuls les modérateurs et administrateurs peuvent créer des salons.')
      }

      setShowCreateForm(false)
      setCreateData({ name: '', minBet: 10, maxBet: 200, playerCount: 7 })
      fetchData() // Refresh list
    } catch (err) {
      setCreateError(err.message)
    }
  }

  // Helper to map tables to saloons
  const getTablesForSaloon = (saloonId) => {
    return tables.filter(t => t.saloon_id === saloonId)
  }

  return (
    <div className="lobby-container">
      <div style={{ flex: 1 }}>
        <div className="lobby-header">
          <h2>♦️ Salons de Poker Disponibles</h2>
          <button className="btn btn-secondary" onClick={fetchData}>
            Actualiser
          </button>
        </div>

        {error && (
          <div style={{ background: 'rgba(224, 79, 95, 0.1)', color: 'var(--accent-red)', padding: '1rem', borderRadius: '8px', marginBottom: '1.5rem', border: '1px solid var(--accent-red)' }}>
            {error}
          </div>
        )}

        {loading ? (
          <p style={{ textAlign: 'center', color: 'var(--text-muted)', padding: '2rem' }}>
            Chargement des tables de jeu...
          </p>
        ) : saloons.length === 0 ? (
          <div className="glass-panel" style={{ textAlign: 'center', padding: '3rem' }}>
            <p style={{ color: 'var(--text-muted)', marginBottom: '1rem' }}>
              Aucun salon actif n'a été trouvé.
            </p>
            {user.role >= 1 && (
              <button className="btn btn-primary" onClick={() => setShowCreateForm(true)} style={{ margin: '0 auto' }}>
                Créer un Salon
              </button>
            )}
          </div>
        ) : (
          <div className="saloon-list">
            {saloons.map((saloon) => {
              const saloonTables = getTablesForSaloon(saloon.id)
              return (
                <div key={saloon.id} className="glass-panel saloon-card">
                  <div className="saloon-info">
                    <h3 className="text-gold">{saloon.name}</h3>
                    <div className="saloon-details">
                      <span className="saloon-detail">
                        💰 Blinds: {saloon.min_bet} / {saloon.max_bet} Jetons
                      </span>
                      <span className="saloon-detail">
                        👥 Max Joueurs: {saloon.player_count}
                      </span>
                    </div>
                    
                    {/* Render matching tables */}
                    <div style={{ marginTop: '1rem', display: 'flex', flexDirection: 'column', gap: '0.5rem' }}>
                      {saloonTables.length === 0 ? (
                        <span style={{ fontSize: '0.85rem', color: 'var(--text-muted)' }}>
                          Aucune table associée.
                        </span>
                      ) : (
                        saloonTables.map((t) => (
                          <div 
                            key={t.id} 
                            style={{ 
                              display: 'flex', 
                              justifyContent: 'space-between', 
                              alignItems: 'center',
                              background: 'rgba(0,0,0,0.15)',
                              padding: '0.5rem 0.8rem',
                              borderRadius: '6px',
                              border: '1px solid rgba(255,255,255,0.03)'
                            }}
                          >
                            <span style={{ fontSize: '0.9rem', fontWeight: '500' }}>
                              Table #{t.id} ({t.slots_available} places max)
                            </span>
                            <button 
                              className="btn btn-primary" 
                              style={{ padding: '0.4rem 0.8rem', fontSize: '0.85rem' }}
                              onClick={() => onJoinTable(t.id)}
                            >
                              Rejoindre
                            </button>
                          </div>
                        ))
                      )}
                    </div>
                  </div>
                </div>
              )
            })}
          </div>
        )}
      </div>

      {/* Side Creation Panel */}
      <div className="glass-panel" style={{ height: 'fit-content' }}>
        <h3 style={{ marginBottom: '1rem', borderBottom: '1px solid var(--border-color)', paddingBottom: '0.5rem' }}>
          Mon Compte
        </h3>
        <p style={{ fontWeight: '700', fontSize: '1.1rem', marginBottom: '0.25rem' }}>
          {user.nick_name}
        </p>
        <p style={{ color: 'var(--text-muted)', fontSize: '0.85rem', marginBottom: '1rem' }}>
          {user.email}
        </p>
        
        <div className="chip-badge" style={{ justifyContent: 'center', marginBottom: '1.5rem', width: '100%' }}>
          <span className="chip-icon"></span>
          <span>{user.money} Jetons</span>
        </div>

        {user.role >= 1 ? (
          <div>
            {!showCreateForm ? (
              <button className="btn btn-secondary" style={{ width: '100%' }} onClick={() => setShowCreateForm(true)}>
                + Créer un Salon
              </button>
            ) : (
              <form onSubmit={handleCreateSaloon} style={{ display: 'flex', flexDirection: 'column', gap: '1rem' }}>
                <h4 style={{ color: 'var(--accent-gold)' }}>Nouveau Salon</h4>
                {createError && (
                  <div style={{ background: 'rgba(224, 79, 95, 0.1)', color: 'var(--accent-red)', padding: '0.5rem', borderRadius: '4px', fontSize: '0.8rem', border: '1px solid var(--accent-red)' }}>
                    {createError}
                  </div>
                )}
                <div className="form-group" style={{ marginBottom: 0 }}>
                  <label htmlFor="create-name" style={{ fontSize: '0.75rem' }}>Nom du salon</label>
                  <input 
                    type="text" 
                    id="create-name" 
                    className="form-input" 
                    required 
                    value={createData.name}
                    onChange={(e) => setCreateData({ ...createData, name: e.target.value })}
                    placeholder="Ex: Table des High Rollers"
                  />
                </div>
                <div style={{ display: 'flex', gap: '0.5rem' }}>
                  <div className="form-group" style={{ marginBottom: 0, flex: 1 }}>
                    <label htmlFor="create-min" style={{ fontSize: '0.75rem' }}>Min Blind</label>
                    <input 
                      type="number" 
                      id="create-min" 
                      className="form-input" 
                      required 
                      value={createData.minBet}
                      onChange={(e) => setCreateData({ ...createData, minBet: e.target.value })}
                    />
                  </div>
                  <div className="form-group" style={{ marginBottom: 0, flex: 1 }}>
                    <label htmlFor="create-max" style={{ fontSize: '0.75rem' }}>Max Blind</label>
                    <input 
                      type="number" 
                      id="create-max" 
                      className="form-input" 
                      required 
                      value={createData.maxBet}
                      onChange={(e) => setCreateData({ ...createData, maxBet: e.target.value })}
                    />
                  </div>
                </div>
                <div className="form-group" style={{ marginBottom: 0 }}>
                  <label htmlFor="create-slots" style={{ fontSize: '0.75rem' }}>Places Joueurs</label>
                  <input 
                    type="number" 
                    id="create-slots" 
                    className="form-input" 
                    required 
                    value={createData.playerCount}
                    onChange={(e) => setCreateData({ ...createData, playerCount: e.target.value })}
                  />
                </div>
                <div style={{ display: 'flex', gap: '0.5rem', marginTop: '0.5rem' }}>
                  <button type="submit" className="btn btn-primary" style={{ flex: 1, padding: '0.5rem' }}>
                    Valider
                  </button>
                  <button type="button" className="btn btn-secondary" style={{ flex: 1, padding: '0.5rem' }} onClick={() => setShowCreateForm(false)}>
                    Annuler
                  </button>
                </div>
              </form>
            )}
          </div>
        ) : (
          <p style={{ fontSize: '0.8rem', color: 'var(--text-muted)', textAlign: 'center', background: 'rgba(255,255,255,0.02)', padding: '0.75rem', borderRadius: '8px' }}>
            Compte Joueur standard. Pour créer des salons, vous devez disposer du rôle Modérateur.
          </p>
        )}
      </div>
    </div>
  )
}
