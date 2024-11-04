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

// NuevaSimulacion crea una nueva simulación de estacionamiento concurrente
func (s *ServicioEstacionamiento) NuevaSimulacion(vehiculos int) {
	for i := 0; i < vehiculos; i++ {
		go func(id int) {
			// Genera una duración aleatoria entre 3 y 5 segundos
			duracion := time.Duration(rand.Intn(3)+3) * time.Second
			vehiculo := domain.Vehiculo{ID: id, Duracion: duracion}

			for {
				// Intenta entrar al estacionamiento
				if cajon := s.Estacionamiento.Entrar(); cajon != -1 {
					fmt.Printf("Vehículo %d ocupa el cajón %d\n", id, cajon) // Imprime el cajón ocupado
					time.Sleep(vehiculo.Duracion)                           // Simula el tiempo estacionado
					s.Estacionamiento.Salir(cajon)
					fmt.Printf("Vehículo %d deja libre el cajón %d\n", id, cajon) // Imprime el cajón liberado
					return
				}
				// Si no se puede entrar, espera 1 segundo antes de intentar nuevamente
				time.Sleep(1 * time.Second)
			}
		}(i)
	}
}
