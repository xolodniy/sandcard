package application

import "github.com/gorilla/websocket"

type Application struct {
	tables map[int]*Table
}

func New() *Application {
	return &Application{
		tables: make(map[int]*Table),
	}
}

func (a *Application) CreateTable() (int, error) {
	t := NewTable()
	go t.Start()
	a.tables[t.ID] = t
	return t.ID, nil
}

func (a *Application) JoinTable(c *websocket.Conn, tableID int) error {
	if err := a.tables[tableID].Join(c); err != nil {
		return err
	}
	return nil
}
