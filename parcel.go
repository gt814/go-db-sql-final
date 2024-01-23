package main

import (
	"database/sql"
)

type ParcelStore struct {
	db *sql.DB
}

func NewParcelStore(db *sql.DB) ParcelStore {
	return ParcelStore{db: db}
}

func (s ParcelStore) Add(p Parcel) (int, error) {
	res, err := s.db.Exec("INSERT INTO parcel (client, status, address, created_at) VALUES (:client, :status, :address, :created_at)",
		sql.Named("client", p.Client),
		sql.Named("status", p.Status),
		sql.Named("address", p.Address),
		sql.Named("created_at", p.CreatedAt))
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()

	return int(id), err
}

func (s ParcelStore) Get(number int) (Parcel, error) {
	row := s.db.QueryRow("SELECT client, status, address, created_at FROM parcel WHERE number = :id", sql.Named("id", number))

	p := Parcel{}
	err := row.Scan(&p.Client, &p.Status, &p.Address, &p.CreatedAt)

	return p, err
}

func (s ParcelStore) GetByClient(client int) ([]Parcel, error) {
	var res []Parcel

	rows, err := s.db.Query("SELECT number, status, address, created_at FROM parcel WHERE client = :id", sql.Named("id", client))
	if err != nil {
		return res, err
	}
	defer rows.Close()

	var p Parcel
	for rows.Next() {
		p.Client = client
		err := rows.Scan(&p.Number, &p.Status, &p.Address, &p.CreatedAt)
		if err != nil {
			return nil, err
		}
		res = append(res, p)
	}

	return res, nil
}

func (s ParcelStore) SetStatus(number int, status string) error {
	_, err := s.db.Exec("UPDATE parcel SET status = :new_status WHERE number = :id",
		sql.Named("new_status", status),
		sql.Named("id", number))

	return err
}

func (s ParcelStore) SetAddress(number int, address string) error {
	p, err := s.Get(number)
	if err != nil {
		return err
	}

	if p.Status == ParcelStatusRegistered {
		_, err = s.db.Exec("UPDATE parcel SET address = :new_address WHERE number = :id",
			sql.Named("new_address", address),
			sql.Named("id", number))
	}

	return err
}

func (s ParcelStore) Delete(number int) error {
	p, err := s.Get(number)
	if err != nil {
		return err
	}

	if p.Status == ParcelStatusRegistered {
		_, err = s.db.Exec("DELETE FROM parcel WHERE number = :id", sql.Named("id", number))
	}

	return err
}
