// cmd/main.go
package main

import (
	"main/application"
	"main/domain"
	"main/infrastructure"
)

func main() {
	estacionamiento := domain.NuevoEstacionamiento(20)
	servicio := application.ServicioEstacionamiento{Estacionamiento: estacionamiento}
	go servicio.NuevaSimulacion(0) // Inicia la simulación con 100 vehículos

	infrastructure.IniciarGUI(&servicio) // Inicia la interfaz gráfica
}
