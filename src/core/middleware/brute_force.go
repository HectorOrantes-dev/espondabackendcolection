package middleware

import (
	"fmt"
	"sync"
	"time"
)

// loginState rastrea el estado de intentos fallidos por email.
type loginState struct {
	failedAttempts int
	blockedUntil   time.Time
	// afterFirstBlock indica que ya se cumplió el bloqueo inicial de 5 min
	// y ahora solo se permite 1 intento antes del bloqueo largo de 20 min.
	afterFirstBlock bool
}

// BruteForceGuard protege el login contra ataques de fuerza bruta.
// Lógica: 3 fallos → bloqueo 5 min → 1 fallo → bloqueo 20 min → reinicio.
type BruteForceGuard struct {
	mu     sync.Mutex
	states map[string]*loginState
}

func NewBruteForceGuard() *BruteForceGuard {
	return &BruteForceGuard{
		states: make(map[string]*loginState),
	}
}

// CheckAndRecord verifica si el email puede intentar login.
// Debe llamarse ANTES del intento. Retorna error si está bloqueado.
func (g *BruteForceGuard) CheckBlocked(email string) error {
	g.mu.Lock()
	defer g.mu.Unlock()

	s := g.getOrCreate(email)

	if time.Now().Before(s.blockedUntil) {
		remaining := time.Until(s.blockedUntil).Round(time.Second)
		return fmt.Errorf("cuenta bloqueada temporalmente, intenta en %s", remaining)
	}

	// Si el bloqueo ya expiró, limpiar el estado de bloqueo activo
	// pero preservar afterFirstBlock para saber en qué fase estamos.
	if !s.blockedUntil.IsZero() && time.Now().After(s.blockedUntil) {
		s.blockedUntil = time.Time{}
	}

	return nil
}

// RecordFailure registra un intento fallido y aplica el bloqueo si corresponde.
func (g *BruteForceGuard) RecordFailure(email string) (blocked bool, blockDuration time.Duration) {
	g.mu.Lock()
	defer g.mu.Unlock()

	s := g.getOrCreate(email)
	s.failedAttempts++

	if s.afterFirstBlock {
		// Fase 2: solo 1 intento permitido → bloqueo de 20 minutos
		s.blockedUntil = time.Now().Add(20 * time.Minute)
		s.failedAttempts = 0
		s.afterFirstBlock = false // reinicia la fase para cuando expire el bloqueo de 20 min
		return true, 20 * time.Minute
	}

	if s.failedAttempts >= 3 {
		// Fase 1: 3 fallos → bloqueo de 5 minutos
		s.blockedUntil = time.Now().Add(5 * time.Minute)
		s.failedAttempts = 0
		s.afterFirstBlock = true // el siguiente ciclo es de 1 intento
		return true, 5 * time.Minute
	}

	return false, 0
}

// RecordSuccess reinicia el estado al lograr un login exitoso.
func (g *BruteForceGuard) RecordSuccess(email string) {
	g.mu.Lock()
	defer g.mu.Unlock()
	delete(g.states, email)
}

// AttemptsLeft retorna cuántos intentos quedan antes del bloqueo.
func (g *BruteForceGuard) AttemptsLeft(email string) int {
	g.mu.Lock()
	defer g.mu.Unlock()

	s, ok := g.states[email]
	if !ok {
		return 3
	}
	if s.afterFirstBlock {
		return 1 - s.failedAttempts
	}
	return 3 - s.failedAttempts
}

func (g *BruteForceGuard) getOrCreate(email string) *loginState {
	if _, ok := g.states[email]; !ok {
		g.states[email] = &loginState{}
	}
	return g.states[email]
}
