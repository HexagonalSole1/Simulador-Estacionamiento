package infrastructure

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"main/application"
	"main/views"
)

func IniciarGUI(servicio *application.ServicioEstacionamiento) {
	a := app.New()
	w := a.NewWindow("Simulación de Estacionamiento")
	w.Resize(fyne.NewSize(600, 700))

	// Crear una etiqueta de estado
	statusLabel := widget.NewLabel("Simulación de Estacionamiento en curso...")

	// Crear el estacionamiento y el contenedor de vehículos
	parkingGrid, cajones, cajonMutexes := views.CreateParkingGrid(servicio.Estacionamiento.Capacidad)
	vehicleContainer := container.NewWithoutLayout()

	// Contenedor principal
	content := container.NewVBox(
		statusLabel,
		parkingGrid,
		vehicleContainer,
	)

	// Crear vehículos concurrentemente
	views.LaunchVehicles(servicio, vehicleContainer, cajones, cajonMutexes)

	// Asignar el contenido y ejecutar la aplicación
	w.SetContent(content)
	w.ShowAndRun()
}
