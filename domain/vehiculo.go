package domain

import (
	"math/rand"
	"time"
)

type Vehiculo struct {
	ID        int
	Duracion  time.Duration
	PosicionX float32
	PosicionY float32
}

func NewVehiculo(id int) *Vehiculo {
	return &Vehiculo{
		ID:       id,
		Duracion: time.Duration(rand.Intn(3)+3) * time.Second,
		PosicionX: 50,
		PosicionY: 550,
	}
}
