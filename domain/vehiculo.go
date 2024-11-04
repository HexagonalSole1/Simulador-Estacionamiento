package domain

import (
	"time"
)

type Vehiculo struct {
	ID        int
	Duracion  time.Duration // Tiempo que estará estacionado
	PosicionX float32       // Posición X del vehículo en la interfaz
	PosicionY float32       // Posición Y del vehículo en la interfaz
}

func (v *Vehiculo) Mover(posX, posY float32) {
	v.PosicionX = posX
	v.PosicionY = posY
}
