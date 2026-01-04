package app

import (
	"github.com/ArthurGuatsaev/smarthome/internal/storage"
)

type App struct {
	Devices  *storage.DeviceRepo
	States   *storage.StateRepo
	Commands *storage.CommandRepo
}

func New(db *storage.DB) *App {
	return &App{
		Devices:  storage.NewDeviceRepo(db.DB),
		States:   storage.NewStateRepo(db.DB),
		Commands: storage.NewCommandRepo(db.DB),
	}
}
