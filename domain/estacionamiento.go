package domain

import (
	"sync"
	"time"
)

type Estacionamiento struct {
	Capacidad        int
	Ocupados         int
	Cajones          []bool // True si el cajón está ocupado
	Bloqueados       []bool // True si el cajón está bloqueado temporalmente
	Mu               sync.Mutex
	ultimoCajonUsado int
}

// NuevoEstacionamiento crea una nueva instancia de Estacionamiento
func NuevoEstacionamiento(capacidad int) *Estacionamiento {
	return &Estacionamiento{
		Capacidad:        capacidad,
		Cajones:          make([]bool, capacidad),
		Bloqueados:       make([]bool, capacidad),
		ultimoCajonUsado: -1,
	}
}

// Entrar permite a un vehículo ingresar si hay espacio
func (e *Estacionamiento) Entrar() int {
	e.Mu.Lock()
	defer e.Mu.Unlock()
	if e.Ocupados >= e.Capacidad {
		return -1 // No hay espacio
	}

	// Busca el siguiente cajón disponible que no esté ocupado ni bloqueado
	for i := 1; i <= e.Capacidad; i++ {
		cajon := (e.ultimoCajonUsado + i) % e.Capacidad
		if !e.Cajones[cajon] && !e.Bloqueados[cajon] { // Cajón disponible y no bloqueado
			e.Cajones[cajon] = true
			e.Ocupados++
			e.ultimoCajonUsado = cajon // Actualiza el último cajón usado
			return cajon               // Retorna el número del cajón asignado
		}
	}
	return -1 // No hay espacio disponible
}

// Salir permite que un vehículo deje el estacionamiento y bloquea el cajón temporalmente
func (e *Estacionamiento) Salir(cajon int) {
	e.Mu.Lock()
	defer e.Mu.Unlock()

	// Asegura que el cajón especificado esté realmente ocupado
	if cajon >= 0 && cajon < e.Capacidad && e.Cajones[cajon] {
		e.Cajones[cajon] = false  // Libera el cajón
		e.Ocupados--

		// Bloquea el cajón temporalmente por 5 segundos para simular tiempo de salida
		e.Bloqueados[cajon] = true
		go func() {
			time.Sleep(1* time.Microsecond)
			e.Mu.Lock()
			defer e.Mu.Unlock()
			e.Bloqueados[cajon] = false // Desbloquea el cajón después de 5 segundos
		}()
	}
}
