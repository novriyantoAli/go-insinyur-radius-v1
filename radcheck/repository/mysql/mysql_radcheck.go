package mysql

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/novriyantoAli/go-insinyur-radius-v1/domain"
	"github.com/sirupsen/logrus"
)

type mysqlRepository struct {
	Conn *sql.DB
}

func NewMysqlRepository(conn *sql.DB) domain.RadcheckRepository {
	return &mysqlRepository{Conn: conn}
}

func (m *mysqlRepository) fetch(ctx context.Context, query string, args ...interface{}) (result []domain.Radcheck, err error) {
	rows, err := m.Conn.QueryContext(ctx, query, args...)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	defer func() {
		errRow := rows.Close()
		if errRow != nil {
			logrus.Error(errRow)
		}
	}()

	result = make([]domain.Radcheck, 0)
	for rows.Next() {
		t := domain.Radcheck{}
		err = rows.Scan(
			&t.ID,
			&t.Username,
			&t.Attribute,
			&t.OP,
			&t.Value,
		)

		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		result = append(result, t)
	}

	return result, nil
}

func (m *mysqlRepository) Get(ctx context.Context, radcheck domain.Radcheck) (res []domain.Radcheck, err error) {
	query := "SELECT * FROM radcheck "
	args := make([]interface{}, 0)

	addWhere := false

	if radcheck.ID != nil {
		if addWhere {
			query += " AND id = ? "
		} else {
			addWhere = true
			query += " WHERE id = ? "
		}
		args = append(args, *radcheck.ID)
	}

	if radcheck.Username != nil {
		if addWhere {
			query += " AND username = ? "
		} else {
			addWhere = true
			query += " WHERE username = ? "
		}
		args = append(args, *radcheck.Username)
	}

	if radcheck.Attribute != nil {
		if addWhere {
			query += " AND attribute = ? "
		} else {
			addWhere = true
			query += " WHERE attribute = ? "
		}
		args = append(args, *radcheck.Attribute)
	}

	if radcheck.OP != nil {
		if addWhere {
			query += " AND op = ? "
		} else {
			addWhere = true
			query += " WHERE op = ? "
		}
		args = append(args, *radcheck.OP)
	}

	if radcheck.Value != nil {
		if addWhere {
			query += " AND value = ? "
		} else {
			addWhere = true
			query += " WHERE value = ? "
		}
		args = append(args, *radcheck.Value)
	}

	res, err = m.fetch(ctx, query, args...)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	return
}

func (m *mysqlRepository) FetchWithValueExpiration(ctx context.Context, delete bool) (res []domain.Radcheck, err error) {
	query := `SELECT id, username, attribute, op, value FROM radcheck WHERE attribute='Expiration' AND STR_TO_DATE(value, "%d %b %Y") <= CURDATE()`

	if delete {
		query += " - INTERVAL 7 DAY"
	}

	res, err = m.fetch(ctx, query)

	return
}

// Update(ctx context.Context, radcheck Radcheck) (res Users, err error)
func (m *mysqlRepository) Update(ctx context.Context, radcheck domain.Radcheck) (err error) {
	tx, err := m.Conn.BeginTx(ctx, nil)
	if err != nil {
		logrus.Error(err)
		return err
	}

	query := "UPDATE radcheck SET username = ?, attribute = ?, op = ?, value = ? WHERE id = ?"
	_, err = tx.ExecContext(
		ctx, query,
		*radcheck.Username, *radcheck.Attribute, *radcheck.OP, *radcheck.Value, *radcheck.ID,
	)

	if err != nil {
		logrus.Error(err)

		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		logrus.Error(err)
		return err
	}

	return
}

func (m *mysqlRepository) DeleteWithUsername(ctx context.Context, username string) (err error) {
	query := "DELETE FROM radcheck WHERE username = ?"

	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return
	}

	res, err := stmt.ExecContext(ctx, username)
	if err != nil {
		logrus.Error(err)
		return
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return
	}

	if rowsAffected < 1 {
		err = fmt.Errorf("Weird  Behavior. Total Affected: %d", rowsAffected)
		return
	}

	return
}
