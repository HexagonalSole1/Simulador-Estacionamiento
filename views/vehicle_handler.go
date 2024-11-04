package views

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"image/color"
	"main/application"
	"main/domain"
	"math"
	"sync"
	"time"
)

// LaunchVehicles lanza los vehículos en la simulación
func LaunchVehicles(servicio *application.ServicioEstacionamiento, vehicleContainer *fyne.Container, cajones []*canvas.Rectangle, cajonMutexes []sync.Mutex) {
	for i := 0; i < 100; i++ {
		go func(id int) {
			// Crear y configurar el vehículo
			canvasCar := canvas.NewRectangle(color.RGBA{0, 128, 255, 255}) // Carrito azul
			canvasCar.SetMinSize(fyne.NewSize(50, 50))
			vehicleContainer.Add(canvasCar)

			vehiculo := domain.NewVehiculo(id)
			moveVehicleToSlot(servicio, vehiculo, canvasCar, cajones, cajonMutexes, vehicleContainer)
		}(i)
	}
}

func moveVehicleToSlot(servicio *application.ServicioEstacionamiento, vehiculo *domain.Vehiculo, canvasCar *canvas.Rectangle, cajones []*canvas.Rectangle, cajonMutexes []sync.Mutex, vehicleContainer *fyne.Container) {
	initialPosX, initialPosY := vehiculo.PosicionX, vehiculo.PosicionY

	for {
		destCajon := servicio.Estacionamiento.Entrar()
		if destCajon != -1 {
			go func(cajon int) {
				cajonMutexes[cajon].Lock()
				defer cajonMutexes[cajon].Unlock()
				animateMovementToSlot(vehiculo, canvasCar, cajones[cajon], cajon)
				time.Sleep(vehiculo.Duracion)

				// Libera el cajón después del tiempo de permanencia
				servicio.Estacionamiento.Salir(cajon)
				cajones[cajon].FillColor = color.RGBA{200, 200, 200, 255}
				cajones[cajon].Refresh()

				// Mover el vehículo de regreso a su posición inicial
				animateMovementToExit(vehiculo, canvasCar, vehicleContainer, initialPosX, initialPosY)
			}(destCajon)
			break
		} else {
			time.Sleep(1 * time.Second) // Espera antes de reintentar
		}
	}
}

func animateMovementToSlot(vehiculo *domain.Vehiculo, canvasCar *canvas.Rectangle, cajon *canvas.Rectangle, cajonID int) {
	destX := float32((cajonID % 5) * 100 + 25)
	destY := float32((cajonID / 5) * 100 + 125)

	for {
		dx := destX - vehiculo.PosicionX
		dy := destY - vehiculo.PosicionY
		distance := float32(math.Sqrt(float64(dx*dx + dy*dy)))
		if distance < 2 {
			vehiculo.PosicionX = destX
			vehiculo.PosicionY = destY
			break
		}

		// Movimiento gradual
		vehiculo.PosicionX += 2 * (dx / distance)
		vehiculo.PosicionY += 2 * (dy / distance)
		canvasCar.Move(fyne.NewPos(vehiculo.PosicionX, vehiculo.PosicionY))
		canvasCar.Refresh()
		time.Sleep(10 * time.Millisecond)
	}
}

func animateMovementToExit(vehiculo *domain.Vehiculo, canvasCar *canvas.Rectangle, vehicleContainer *fyne.Container, initialPosX, initialPosY float32) {
	for {
		dx := initialPosX - vehiculo.PosicionX
		dy := initialPosY - vehiculo.PosicionY
		distance := float32(math.Sqrt(float64(dx*dx + dy*dy)))

		if distance < 2 {
			vehiculo.PosicionX = initialPosX
			vehiculo.PosicionY = initialPosY
			break
		}

		vehiculo.PosicionX += 2 * (dx / distance)
		vehiculo.PosicionY += 2 * (dy / distance)
		canvasCar.Move(fyne.NewPos(vehiculo.PosicionX, vehiculo.PosicionY))
		canvasCar.Refresh()
		time.Sleep(10 * time.Millisecond)
	}

	// Remueve el vehículo una vez que llega a su posición inicial
	vehicleContainer.Remove(canvasCar)
}
