package main

import (
	"main/application"
	"main/infrastructure"
)

func main() {
	// Inicializar el servicio de estacionamiento
	servicio := application.NewServicioEstacionamiento(5)

	// Inicia la lógica de simulación en paralelo (100 vehículos)
	go servicio.NuevaSimulacion(0)

	// Iniciar la interfaz gráfica
	infrastructure.IniciarGUI(servicio)
}
