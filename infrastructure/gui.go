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
	"math"
	"math/rand"
	"sync"
	"time"
)

func IniciarGUI(servicio *application.ServicioEstacionamiento) {
	rand.Seed(time.Now().UnixNano()) // Inicializa la semilla de aleatoriedad

	a := app.New()
	w := a.NewWindow("Simulación de Estacionamiento")
	w.Resize(fyne.NewSize(600, 700))

	// Crear una etiqueta de estado
	statusLabel := widget.NewLabel("Simulación de Estacionamiento en curso...")

	// Crear el canvas para el estacionamiento y el contenedor de cajones
	cajones := make([]*canvas.Rectangle, servicio.Estacionamiento.Capacidad)
	cajonMutexes := make([]sync.Mutex, servicio.Estacionamiento.Capacidad) // Mutex por cada cajón
	parkingGrid := container.NewGridWithColumns(5)

	// Canal para manejar cajones disponibles
	cajonesDisponibles := make(chan int, servicio.Estacionamiento.Capacidad)

	for i := range cajones {
		cajones[i] = canvas.NewRectangle(color.RGBA{200, 200, 200, 255}) // Color gris para espacio libre
		cajones[i].SetMinSize(fyne.NewSize(100, 100))                    // Tamaño de cada cajón (100x100)
		parkingGrid.Add(cajones[i])
		cajonesDisponibles <- i // Agregar todos los cajones al canal como disponibles
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

		// Obtener un cajón del canal (bloquea hasta que haya uno disponible)
		destCajon := <-cajonesDisponibles

		// Bloquear el cajón específico mientras el vehículo lo ocupa
		go func(cajon int) {
			cajonMutexes[cajon].Lock()
			defer cajonMutexes[cajon].Unlock()

			// Centrar el vehículo en el cajón seleccionado
			destX := float32((cajon % 5) * 100 + 25)
			destY := float32((cajon / 5) * 100 + 125)

			// Mover el vehículo hacia el cajón destino con precisión
			for {
				dx := destX - initialPosX
				dy := destY - initialPosY
				distance := float32(math.Sqrt(float64(dx*dx + dy*dy)))

				// Si el vehículo está cerca del destino, detener el movimiento
				if distance < 2 {
					break
				}

				// Mover el vehículo paso a paso
				initialPosX += 2 * (dx / distance)
				initialPosY += 2 * (dy / distance)

				canvasCar.Move(fyne.NewPos(initialPosX, initialPosY))
				canvasCar.Refresh()
				time.Sleep(10 * time.Millisecond)
			}

			// Marcar el cajón como ocupado y refrescarlo
			cajones[cajon].FillColor = color.RGBA{255, 0, 0, 255} // Cambiar a rojo
			cajones[cajon].Refresh()                              // Refrescar para aplicar el cambio de color inmediatamente

			// Espera mientras el vehículo permanece estacionado
			time.Sleep(duracion)

			// Liberar el cajón inmediatamente después de que el vehículo comienza a salir
			cajonesDisponibles <- cajon

			// Cambiar el color del cajón a gris para indicar que está libre
			cajones[cajon].FillColor = color.RGBA{200, 200, 200, 255} // Cambiar a gris
			cajones[cajon].Refresh()                                  // Refrescar para mostrar el cambio de color

			// Mover el vehículo de regreso
			go func() {
				for {
					dx := 50 - initialPosX
					dy := 550 - initialPosY
					distance := float32(math.Sqrt(float64(dx*dx + dy*dy)))

					// Si el vehículo está cerca de la posición inicial, detener el movimiento
					if distance < 2 {
						break
					}

					initialPosX += 2 * (dx / distance)
					initialPosY += 2 * (dy / distance)

					canvasCar.Move(fyne.NewPos(initialPosX, initialPosY))
					canvasCar.Refresh()
					time.Sleep(10 * time.Millisecond)
				}

				// Eliminar el vehículo de la interfaz cuando se va
				vehicleContainer.Remove(canvasCar)
			}()
		}(destCajon)
	}

	// Máximo de vehículos a crear
	maxVehicles := 30 // Ajusta a 30 para facilitar la visualización

	// Crear todos los vehículos de manera concurrente sin pausa entre ellos
	for i := 0; i < maxVehicles; i++ {
		go createVehicle(vehicleID)
		vehicleID++
	}

	// Asignar el contenido y ejecutar la aplicación
	w.SetContent(content)
	w.ShowAndRun()
}
