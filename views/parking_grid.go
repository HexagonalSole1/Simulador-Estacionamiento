package views

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"image/color"
	"sync"
)

// CreateParkingGrid crea la interfaz de una cuadrícula de estacionamiento con cajones
func CreateParkingGrid(capacidad int) (*fyne.Container, []*canvas.Rectangle, []sync.Mutex) {
	cajones := make([]*canvas.Rectangle, capacidad)
	cajonMutexes := make([]sync.Mutex, capacidad)
	parkingGrid := container.NewGridWithColumns(5)

	for i := range cajones {
		cajones[i] = canvas.NewRectangle(color.RGBA{200, 200, 200, 255}) // Color gris para espacio libre
		cajones[i].SetMinSize(fyne.NewSize(100, 100))                    // Tamaño de cada cajón
		parkingGrid.Add(cajones[i])
	}

	return parkingGrid, cajones, cajonMutexes
}
