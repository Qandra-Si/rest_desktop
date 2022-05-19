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
	desktops []Desktop
	nextId   int
}

func New() *DesktopStore {
	ts := &DesktopStore{}
	ts.desktops = []Desktop{}
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
		At:    at,
	}
	ds.desktops = append(ds.desktops, desktop)
	ds.nextId++
	return desktop.Id
}

func (ds *DesktopStore) DeleteDesktop(cname string, cip string, user string, at time.Time) error {
	ds.Lock()
	defer ds.Unlock()

	for i, desktop := range ds.desktops {
		if desktop.CName == cname {
			ds.desktops = append(ds.desktops[:i], ds.desktops[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("desktop with cname=%s not found", cname)
}

func (ds *DesktopStore) UpdateDesktop(cname string, cip string, user string, at time.Time) error {
	ds.Lock()
	defer ds.Unlock()

	for _, desktop := range ds.desktops {
		if desktop.CName == cname {
			desktop.CIp = cip
			desktop.User = user
			desktop.At = at
			return nil
		}
	}
	return fmt.Errorf("desktop with cname=%s not found", cname)
}
