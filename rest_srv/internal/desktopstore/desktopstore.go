package desktopstore

import (
	"fmt"
	"sync"
	"time"
)

type Desktop struct {
	Id    int       `json:"id"`
	CName string    `json:"cname"`
	CIp   string    `json:"cip"`
	User  string    `json:"user"`
	At    time.Time `json:"at"`
}

type DesktopStore struct {
	sync.Mutex

	desktops map[int]Desktop
	nextId   int
}

func New() *DesktopStore {
	ts := &DesktopStore{}
	ts.desktops = make(map[int]Desktop)
	ts.nextId = 0
	return ts
}

func (ds *DesktopStore) CreateDesktop(cname string, cip string, user string, at time.Time) int {
	ds.Lock()
	defer ds.Unlock()

	desktop := Desktop{
		Id:    ds.nextId,
		CName: cname,
		CIp:   cip,
		User:  user,
		At:    at}

	ds.desktops[ds.nextId] = desktop
	ds.nextId++
	return desktop.Id
}

func (ds *DesktopStore) DeleteDesktop(id int) error {
	ds.Lock()
	defer ds.Unlock()

	if _, ok := ds.desktops[id]; !ok {
		return fmt.Errorf("desktop with id=%d not found", id)
	}

	delete(ds.desktops, id)
	return nil
}

func (ds *DesktopStore) UpdateDesktop(id int, cname string, cip string, user string, at time.Time) error {
	ds.Lock()
	defer ds.Unlock()

	if _, ok := ds.desktops[id]; !ok {
		return fmt.Errorf("desktop with id=%d not found", id)
	}

	ds.desktops[id] = Desktop{
		Id:    id,
		CName: cname,
		CIp:   cip,
		User:  user,
		At:    at}
	return nil
}
