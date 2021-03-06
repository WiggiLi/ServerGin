package app

import (
	"log"
)

// IncomeRegistration is an interface for accepting income requesrs for neccassery operations  from Web Server
type IncomeRegistration interface {
	GiveEvents(*Event) *AllEvents
	GiveCount(*Event) int
	GiveCsv()
}

// DataAccessLayer is an interface for DAL usage from Application
type DataAccessLayer interface {
	ReadEvent(*Event) (*AllEvents, error)
	ReadCountRows(*Event) (int, error)
	GetCsv() error
}

// Application is responsible for all logics and communicates with other layers
type Application struct {
	DB   DataAccessLayer
	errc chan<- error
}

// RegisterEvent sends Event to DAL for saving/registration
func (app *Application) GiveEvents(currentData *Event) *AllEvents {
	allEv := GetEvents()
	allEv, err := app.DB.ReadEvent(currentData)

	if err != nil {
		app.errc <- err
		return nil
	}
	log.Print("Events readed from DB...")
	return allEv
}

func (app *Application) GiveCount(currentData *Event) int {
	count, err := app.DB.ReadCountRows(currentData)

	if err != nil {
		app.errc <- err
		return 0
	}

	return count
}
func (app *Application) GiveCsv() {
	err := app.DB.GetCsv()

	if err != nil {
		app.errc <- err
		return
	}
	return
}

// NewApplication constructs Application
func NewApplication(db DataAccessLayer, errchannel chan<- error) *Application {
	res := &Application{}

	res.DB = db
	res.errc = errchannel

	return res
}
