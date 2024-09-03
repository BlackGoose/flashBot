package database

import (
	"time"

	"github.com/jmoiron/sqlx"

	_ "github.com/lib/pq"
)

type Card struct {
	Id            int64     `db:"id"`
	Front         string    `db:"front"`
	Back          string    `db:"back"`
	UserId        int64     `db:"user_id"`
	DateExpired   time.Time `db:"date_expired"`
	CurrentStrike int       `db:"current_strike"`
}

func Create(db *sqlx.DB, front, back string, user_id int64) error {
	card := &Card{Front: front, Back: back, UserId: user_id, DateExpired: time.Now().AddDate(0, 0, 1), CurrentStrike: 1}
	_, err := db.NamedExec("insert into cards (front, back, user_id, date_expired, current_strike) values(:front, :back, :user_id, :date_expired, :current_strike)", &card)
	return err
}

func Update(db *sqlx.DB, front, back string, user_id int) error {
	_, err := db.Exec("update cards set back = $2 where front = $1 and user_id = $3", front, back, user_id)
	return err
}

func Get(db *sqlx.DB, user_id int64) ([]Card, error) {
	cards := []Card{}
	err := db.Select(&cards, "select * from cards where user_id = $1 and date_expired <= $2", user_id, time.Now())
	return cards, err
}

func Delete(db *sqlx.DB, front string, user_id int64) error {
	_, err := db.Exec("delete from cards * where front = $1 and user_id = $2", front, user_id)
	return err
}

func Init() (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", "postgres://card_admin:password@localhost/cards?sslmode=disable")
	return db, err
}
