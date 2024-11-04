package infrastructure

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"image/color"
	"main/application"
	"main/domain"
	"math"
	"math/rand"
	"sync"
	"time"
)

func IniciarGUI(servicio *application.ServicioEstacionamiento) {
	a := app.New()
	w := a.NewWindow("Simulación de Estacionamiento")
	w.Resize(fyne.NewSize(600, 700))

	// Crear una etiqueta de estado
	statusLabel := widget.NewLabel("Simulación de Estacionamiento en curso...")

	// Crear el canvas para el estacionamiento y el contenedor de cajones
	cajones := make([]*canvas.Rectangle, servicio.Estacionamiento.Capacidad)
	cajonMutexes := make([]sync.Mutex, servicio.Estacionamiento.Capacidad) // Mutex por cada cajón
	parkingGrid := container.NewGridWithColumns(5)

	for i := range cajones {
		cajones[i] = canvas.NewRectangle(color.RGBA{200, 200, 200, 255}) // Color gris para espacio libre
		cajones[i].SetMinSize(fyne.NewSize(100, 100))                    // Tamaño de cada cajón (100x100)
		parkingGrid.Add(cajones[i])
	}

	// Contenedor para todos los vehículos
	vehicleContainer := container.NewWithoutLayout()

	// Contenedor principal
	content := container.NewVBox(
		statusLabel,
		parkingGrid,      // Agregar el estacionamiento directamente al VBox
		vehicleContainer, // Contenedor donde se añadirán los vehículos
	)

	// Inicializa el ID para los vehículos
	vehicleID := 0

	// Función para manejar la lógica de cada vehículo
	createVehicle := func(id int) {
		fmt.Printf("Creando vehículo con ID: %d\n", id)
		canvasCar := canvas.NewRectangle(color.RGBA{0, 128, 255, 255}) // Carrito azul
		canvasCar.SetMinSize(fyne.NewSize(50, 50))
		canvasCar.Resize(fyne.NewSize(50, 50))

		initialPosX, initialPosY := float32(50), float32(550) // Posición inicial fuera del estacionamiento
		canvasCar.Move(fyne.NewPos(initialPosX, initialPosY))

		// Añadir el vehículo al contenedor
		vehicleContainer.Add(canvasCar)

		// Generar una duración aleatoria de entre 3 y 5 segundos para el tiempo de estacionamiento
		duracion := time.Duration(rand.Intn(3)+3) * time.Second
		vehiculo := &domain.Vehiculo{
			ID:        id,
			Duracion:  duracion,
			PosicionX: initialPosX,
			PosicionY: initialPosY,
		}

		for {
			destCajon := servicio.Estacionamiento.Entrar()
			if destCajon != -1 {
				// Cajón encontrado, continuar con el estacionamiento
				fmt.Printf("Vehículo %d ocupa el cajón %d por %v\n", id, destCajon, vehiculo.Duracion)
				destX := float32((destCajon % 5) * 100 + 25)  // Centra el vehículo dentro del cajón
				destY := float32((destCajon / 5) * 100 + 125) // Centra el vehículo dentro del cajón

				// Asegurarse de que el cajón específico esté bloqueado mientras el vehículo lo ocupa
				go func(cajon int) {
					cajonMutexes[cajon].Lock()
					defer cajonMutexes[cajon].Unlock()

					// Mover el vehículo hacia el cajón destino con precisión
					for {
						dx := destX - vehiculo.PosicionX
						dy := destY - vehiculo.PosicionY
						distance := float32(math.Sqrt(float64(dx*dx + dy*dy)))

						// Si el vehículo está cerca del destino, detenemos el movimiento
						if distance < 2 {
							vehiculo.PosicionX = destX
							vehiculo.PosicionY = destY
							break
						}

						// Mover el vehículo paso a paso
						if dx != 0 {
							vehiculo.PosicionX += 2 * (dx / distance)
						}
						if dy != 0 {
							vehiculo.PosicionY += 2 * (dy / distance)
						}

						canvasCar.Move(fyne.NewPos(vehiculo.PosicionX, vehiculo.PosicionY))
						canvasCar.Refresh()
						time.Sleep(10 * time.Millisecond)
					}

					// Imprimir mensaje cuando el vehículo está en el cajón
					fmt.Printf("Vehículo %d ha llegado al cajón %d\n", id, cajon)

					// Marcar el cajón como ocupado y refrescarlo
					cajones[cajon].FillColor = color.RGBA{255, 0, 0, 255} // Cambiar a rojo
					cajones[cajon].Refresh()                              // Refrescar para aplicar el cambio de color inmediatamente

					// Espera mientras el vehículo permanece estacionado
					time.Sleep(vehiculo.Duracion)

					// Salir del cajón y regresar a la posición inicial
					servicio.Estacionamiento.Salir(cajon)
					fmt.Printf("Vehículo %d deja libre el cajón %d\n", id, cajon)
					cajones[cajon].FillColor = color.RGBA{200, 200, 200, 255} // Cambiar a gris
					cajones[cajon].Refresh()                                  // Refrescar para mostrar el cambio de color

					// Mover el vehículo de regreso
					go func() {
						for {
							dx := initialPosX - vehiculo.PosicionX
							dy := initialPosY - vehiculo.PosicionY
							distance := float32(math.Sqrt(float64(dx*dx + dy*dy)))

							// Si el vehículo está cerca de la posición inicial, detenemos el movimiento
							if distance < 2 {
								vehiculo.PosicionX = initialPosX
								vehiculo.PosicionY = initialPosY
								break
							}

							if dx != 0 {
								vehiculo.PosicionX += 2 * (dx / distance)
							}
							if dy != 0 {
								vehiculo.PosicionY += 2 * (dy / distance)
							}

							canvasCar.Move(fyne.NewPos(vehiculo.PosicionX, vehiculo.PosicionY))
							canvasCar.Refresh()
							time.Sleep(10 * time.Millisecond)
						}

						// Eliminar el vehículo de la interfaz cuando se va
						vehicleContainer.Remove(canvasCar)
					}()
				}(destCajon)
				break
			} else {
				// Si no hay cajón disponible, el vehículo espera y vuelve a intentar
				fmt.Printf("Vehículo %d está en espera...\n", id)
				time.Sleep(1 * time.Second) // Esperar antes de reintentar
			}
		}
	}

	// Máximo de vehículos a crear
	maxVehicles := 100

	// Crear todos los vehículos de manera concurrente sin pausa entre ellos
	for i := 0; i < maxVehicles; i++ {
		go createVehicle(vehicleID)
		vehicleID++
	}

	// Asignar el contenido y ejecutar la aplicación
	w.SetContent(content)
	w.ShowAndRun()
}
