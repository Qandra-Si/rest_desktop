package desktopstore

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"os"
	"sync"
	"time"
)

type DesktopStore struct {
	sync.Mutex
}

func New() *DesktopStore {
	ts := &DesktopStore{}
	return ts
}

func (ds *DesktopStore) RefreshDesktop(cname string, user string, cip string, at time.Time) (int, error) {
	ds.Lock()
	defer ds.Unlock()

	// DBDSN=postgres://testuser:testpassword@localhost/testdb?sslmode=disable
	db, err := sql.Open("postgres", os.Getenv("DBURL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		return 0, err
	}
	defer db.Close()

	var id int
	err = db.QueryRow(
		`insert into public.rest_srv_table(cname,cip,"user","at")
values($1,$2,$3,$4) 
on conflict on constraint unq_rest_srv_cname do update set
cip=$2,"user"=$3,"at"=$4 returning id;`, cname, cip, user, at).Scan(&id)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Refresh failed: %v\n", err)
		return 0, err
	}
	return id, nil
}

func (ds *DesktopStore) DeleteDesktop(cname string, user string, cip string, at time.Time) error {
	ds.Lock()
	defer ds.Unlock()

	// DBDSN=postgres://testuser:testpassword@localhost/testdb?sslmode=disable
	db, err := sql.Open("postgres", os.Getenv("DBURL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		return err
	}
	defer db.Close()

	var id int
	err = db.QueryRow("delete from public.rest_srv_table where cname=$1 returning id;", cname).Scan(&id)
	if err != nil {
		return fmt.Errorf("desktop with cname=%s not found", cname)
	}
	return nil
}
