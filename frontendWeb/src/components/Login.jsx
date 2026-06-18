import React, { useState } from 'react'

const API_URL = `http://${window.location.hostname}:8081`;

export default function Login({ onAuthSuccess }) {
  const [isLogin, setIsLogin] = useState(true)
  const [formData, setFormData] = useState({
    firstName: '',
    lastName: '',
    nickName: '',
    email: '',
    confirmEmail: '',
    password: '',
    confirmPassword: ''
  })
  const [error, setError] = useState(null)
  const [loading, setLoading] = useState(false)
  const [successMessage, setSuccessMessage] = useState(null)

  const handleChange = (e) => {
    setFormData({ ...formData, [e.target.name]: e.target.value })
  }

  const handleSubmit = async (e) => {
    e.preventDefault()
    setError(null)
    setSuccessMessage(null)
    setLoading(true)

    try {
      if (isLogin) {
        // Handle Login
        const response = await fetch(`${API_URL}/api/login`, {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({
            email: formData.email,
            password: formData.password
          })
        })
        const data = await response.json()
        if (!response.ok) {
          throw new Error(data.error || 'Identifiants invalides')
        }
        // Save token and user details, call parent handler
        localStorage.setItem('auth_token', data.token)
        localStorage.setItem('auth_user', JSON.stringify(data.user))
        onAuthSuccess(data.token, data.user)
      } else {
        // Handle Register
        const response = await fetch(`${API_URL}/api/register`, {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({
            first_name: formData.firstName,
            last_name: formData.lastName,
            nick_name: formData.nickName,
            email: formData.email,
            confirmEmail: formData.confirmEmail,
            password: formData.password,
            confirmPassword: formData.confirmPassword
          })
        })
        const data = await response.json()
        if (!response.ok) {
          throw new Error(data.error || 'Erreur lors de la création du compte')
        }
        setSuccessMessage('Compte créé avec succès ! Connectez-vous maintenant.')
        setIsLogin(true)
        setFormData({
          ...formData,
          password: '',
          confirmPassword: ''
        })
      }
    } catch (err) {
      setError(err.message)
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="auth-page">
      <div className="glass-panel auth-card">
        <h2 style={{ textAlign: 'center', marginBottom: '1.5rem', fontSize: '1.75rem' }}>
          {isLogin ? '♠️ Connexion' : '♥️ Créer un compte'}
        </h2>

        <div className="auth-tabs">
          <button 
            className={`auth-tab ${isLogin ? 'active' : ''}`}
            onClick={() => { setIsLogin(true); setError(null); }}
          >
            Se Connecter
          </button>
          <button 
            className={`auth-tab ${!isLogin ? 'active' : ''}`}
            onClick={() => { setIsLogin(false); setError(null); }}
          >
            S'inscrire
          </button>
        </div>

        {error && (
          <div style={{
            background: 'rgba(224, 79, 95, 0.1)',
            border: '1px solid var(--accent-red)',
            color: 'var(--accent-red)',
            padding: '0.75rem',
            borderRadius: '8px',
            marginBottom: '1.25rem',
            fontSize: '0.9rem',
            fontWeight: '600'
          }}>
            {error}
          </div>
        )}

        {successMessage && (
          <div style={{
            background: 'rgba(24, 210, 110, 0.1)',
            border: '1px solid var(--accent-green-glow)',
            color: 'var(--accent-green-glow)',
            padding: '0.75rem',
            borderRadius: '8px',
            marginBottom: '1.25rem',
            fontSize: '0.9rem',
            fontWeight: '600'
          }}>
            {successMessage}
          </div>
        )}

        <form onSubmit={handleSubmit}>
          {!isLogin && (
            <>
              <div className="form-group">
                <label htmlFor="firstName">Prénom</label>
                <input 
                  type="text" 
                  id="firstName" 
                  name="firstName" 
                  className="form-input" 
                  required 
                  value={formData.firstName}
                  onChange={handleChange}
                  placeholder="Ex: Jean"
                />
              </div>

              <div className="form-group">
                <label htmlFor="lastName">Nom</label>
                <input 
                  type="text" 
                  id="lastName" 
                  name="lastName" 
                  className="form-input" 
                  required 
                  value={formData.lastName}
                  onChange={handleChange}
                  placeholder="Ex: Dupont"
                />
              </div>

              <div className="form-group">
                <label htmlFor="nickName">Pseudo</label>
                <input 
                  type="text" 
                  id="nickName" 
                  name="nickName" 
                  className="form-input" 
                  required 
                  value={formData.nickName}
                  onChange={handleChange}
                  placeholder="Ex: PokerKing99"
                />
              </div>
            </>
          )}

          <div className="form-group">
            <label htmlFor="email">Adresse E-mail</label>
            <input 
              type="email" 
              id="email" 
              name="email" 
              className="form-input" 
              required 
              value={formData.email}
              onChange={handleChange}
              placeholder="Ex: contact@poker.fr"
            />
          </div>

          {!isLogin && (
            <div className="form-group">
              <label htmlFor="confirmEmail">Confirmer l'e-mail</label>
              <input 
                type="email" 
                id="confirmEmail" 
                name="confirmEmail" 
                className="form-input" 
                required 
                value={formData.confirmEmail}
                onChange={handleChange}
                placeholder="Ex: contact@poker.fr"
              />
            </div>
          )}

          <div className="form-group">
            <label htmlFor="password">Mot de passe</label>
            <input 
              type="password" 
              id="password" 
              name="password" 
              className="form-input" 
              required 
              value={formData.password}
              onChange={handleChange}
              placeholder="••••••••"
            />
          </div>

          {!isLogin && (
            <div className="form-group">
              <label htmlFor="confirmPassword">Confirmer le mot de passe</label>
              <input 
                type="password" 
                id="confirmPassword" 
                name="confirmPassword" 
                className="form-input" 
                required 
                value={formData.confirmPassword}
                onChange={handleChange}
                placeholder="••••••••"
              />
            </div>
          )}

          <button 
            type="submit" 
            className="btn btn-primary" 
            style={{ width: '100%', marginTop: '1.5rem', height: '45px' }}
            disabled={loading}
          >
            {loading ? 'Traitement en cours...' : (isLogin ? 'Se Connecter' : "S'inscrire")}
          </button>
        </form>
      </div>
    </div>
  )
}
