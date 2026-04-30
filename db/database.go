package db

import (
    "database/sql"
		"context"
    _ "modernc.org/sqlite"
		_ "embed"
)

//go:embed schema.sql
var schema string

type Store struct {
    DB *sql.DB
}

func NewStore(dbPath string) (*Store, error) {
    db, err := sql.Open("sqlite", dbPath)
		if err != nil {
return nil, err
		}
		_, err = db.Exec(schema)
    return &Store{DB: db}, err
}

func (s *Store) AddSignature(ctx context.Context, name string, email *string) error {
    _, err := s.DB.ExecContext(ctx, "INSERT INTO signatures (name, email) VALUES (?, ?)", name, email)
    if err != nil {
      return err
    }
    return nil
}


func (s *Store) GetSignatureCount(ctx context.Context) (int, error) {
    var count int
    err := s.DB.QueryRowContext(ctx, "SELECT COUNT(*) FROM signatures").Scan(&count)
    return count, err
}
