package application

import (
	"fmt"
	"main/domain"
	"math/rand"
	"time"
)

type ServicioEstacionamiento struct {
	Estacionamiento *domain.Estacionamiento
}

// NewServicioEstacionamiento crea un nuevo servicio para gestionar el estacionamiento
func NewServicioEstacionamiento(capacidad int) *ServicioEstacionamiento {
	return &ServicioEstacionamiento{
		Estacionamiento: domain.NuevoEstacionamiento(capacidad),
	}
}

// NuevaSimulacion inicia la simulación del estacionamiento con el número de vehículos especificado
func (s *ServicioEstacionamiento) NuevaSimulacion(vehiculos int) {
	for i := 0; i < vehiculos; i++ {
		go func(id int) {
			// Genera una duración aleatoria entre 3 y 5 segundos
			duracion := time.Duration(rand.Intn(3)+3) * time.Second

			for {
				// Intenta ocupar un cajón en el estacionamiento
				cajon := s.Estacionamiento.Entrar()
				if cajon != -1 {
					fmt.Printf("Vehículo %d ocupa el cajón %d por %v\n", id, cajon, duracion)
					
					// Mantiene el vehículo en el cajón durante la duración especificada
					time.Sleep(duracion)
					
					// Libera el cajón y muestra un mensaje
					s.Estacionamiento.Salir(cajon)
					fmt.Printf("Vehículo %d deja libre el cajón %d\n", id, cajon)
					break
				} else {
					// Si no hay espacio, espera 1 segundo antes de reintentar
					fmt.Printf("Vehículo %d está en espera...\n", id)
					time.Sleep(1 * time.Second)
				}
			}
		}(i)
	}
}
