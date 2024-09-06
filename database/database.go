package database

import (
	"time"

	"github.com/jmoiron/sqlx"

	_ "github.com/lib/pq"
)

type Card struct {
	Id            int       `db:"id"`
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

func UpdateCard(db *sqlx.DB, cardId int, front, back string) error {
	_, err := db.Exec("update cards set front = $1, back = $2 where id = $3", front, back, cardId)
	return err
}

func UpdateDate(db *sqlx.DB, cardId int, dateExpired time.Time, currentStrike int) error {
	_, err := db.Exec("update cards set date_expired = $1, current_strike = $2 where id = $3", dateExpired, currentStrike, cardId)
	return err
}

func GetList(db *sqlx.DB, user_id int64, to_train bool) ([]Card, error) {
	cards := []Card{}
	if to_train {
		err := db.Select(&cards, "select * from cards where user_id = $1 and date_expired <= $2", user_id, time.Now())
		return cards, err
	}
	err := db.Select(&cards, "select * from cards where user_id = $1", user_id)
	return cards, err
}

func Get(db *sqlx.DB, card_id int) (Card, error) {
	card := Card{}
	err := db.Get(&card, "select * from cards where id = $1 limit 1", card_id)
	return card, err
}

func Delete(db *sqlx.DB, card_id int) error {
	_, err := db.Exec("delete from cards * where id = $1", card_id)
	return err
}

func Init() (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", "postgres://card_admin:password@localhost/cards?sslmode=disable")
	return db, err
}
