package mysql

import (
	"context"
	"database/sql"
	"strconv"
	"time"

	"github.com/novriyantoAli/go-insinyur-radius-v1/domain"
	"github.com/sirupsen/logrus"
)

type mysqlRepository struct {
	Conn *sql.DB
}

func NewMysqlRepository(conn *sql.DB) domain.ResellerRepository {
	return &mysqlRepository{Conn: conn}
}

// Transaction(ctx context.Context, idUsers int64) (res Transaction, err error)
func (m *mysqlRepository) Transaction(ctx context.Context, idUsers int64, username string, password string, packages domain.Package) (idTransaction int64, err error) {
	tx, err := m.Conn.BeginTx(ctx, nil)

	// pertama buatlah username dan password
	// kedua buatlah radpackage atas username tersebut
	// hasil dari id radpackage letakkan di transaksi

	// radcheck
	query := "INSERT INTO radcheck(username, attribute, op, value) VALUES(?,?,?,?),(?,?,?,?)"
	_, err = tx.ExecContext(
		ctx,
		query,
		username, "Cleartext-Password", ":=", password,
		username, "User-Profile", ":=", *packages.Profile,
	)
	if err != nil {
		logrus.Error(err)

		tx.Rollback()
		return -1, err
	}

	// radpackage
	query = "INSERT INTO radpackage(id_package, username) VALUES(?,?)"
	rows, err := tx.ExecContext(ctx, query, packages.ID, username)
	if err != nil {
		logrus.Error(err)

		tx.Rollback()
		return -1, err
	}

	id, err := rows.LastInsertId()
	if err != nil {
		logrus.Error(err)

		tx.Rollback()
		return -1, err
	}

	t := time.Now()
	uniqTime := strconv.FormatInt(int64(time.Nanosecond)*t.UnixNano()/int64(time.Microsecond), 10)
	// transaction
	query = "INSERT INTO transaction(id_reseller, id_radpackage, transaction_code, status, value) VALUES(?,?,?,?,?)"
	rowsTransaction, err := tx.ExecContext(
		ctx,
		query,
		idUsers, id, uniqTime, "out", *packages.Price,
	)
	if err != nil {
		logrus.Error(err)

		tx.Rollback()
		return -1, err
	}

	idTransaction, err = rowsTransaction.LastInsertId()
	if err != nil {
		logrus.Error(err)

		tx.Rollback()
		return -1, err
	}

	err = tx.Commit()
	if err != nil {
		logrus.Error(err)
		return -1, err
	}

	return
}
