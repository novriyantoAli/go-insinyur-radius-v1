package mysql

import (
	"context"
	"database/sql"

	"github.com/novriyantoAli/go-insinyur-radius-v1/domain"
	"github.com/sirupsen/logrus"
)

type mysqlRepository struct {
	Conn *sql.DB
}

func NewMysqlRepository(conn *sql.DB) domain.TransactionRepository {
	return &mysqlRepository{Conn: conn}
}

func (m *mysqlRepository) fetch(ctx context.Context, query string, args ...interface{}) (result []domain.Transaction, err error) {
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

	result = make([]domain.Transaction, 0)
	for rows.Next() {
		t := domain.Transaction{}
		err = rows.Scan(
			&t.ID,
			&t.IDReseller,
			&t.IDRadpackage,
			&t.TransactionCode,
			&t.Status,
			&t.Value,
			&t.Information,
			&t.CreatedAt,
		)

		if err != nil {
			logrus.Error(err)
			return nil, err
		}

		// row reseller
		user := domain.Users{}
		err = m.Conn.QueryRowContext(ctx, "SELECT id, username, level, created_at FROM users WHERE id = ?", *t.IDReseller).Scan(
			&user.ID,
			&user.Username,
			&user.Level,
			&user.CreatedAt,
		)
		if err != nil {
			logrus.Error(err)
			return nil, err
		}

		t.Reseller = user

		// row radpackage
		if t.IDRadpackage != nil {
			pac := domain.Package{}
			rp := domain.Radpackage{}
			err = m.Conn.QueryRowContext(ctx, "SELECT radpackage.*, package.* FROM radpackage INNER JOIN package ON radpackage.id_package = package.id WHERE radpackage.id = ?", *t.IDRadpackage).Scan(
				&rp.ID,
				&rp.IDPackage,
				&rp.Username,
				&pac.ID,
				&pac.Name,
				&pac.ValidityValue,
				&pac.ValidityUnit,
				&pac.Price,
				&pac.Margin,
				&pac.Profile,
				&pac.CreatedAt,
			)
			if err != nil {
				logrus.Error(err)
				return nil, err
			}

			// row keluarga
			rowsRc, err := m.Conn.QueryContext(ctx, "SELECT * FROM radcheck WHERE username = ?", *rp.Username)
			if err != nil {
				logrus.Error(err)
				return nil, err
			}
			rc := make([]domain.Radcheck, 0)
			for rowsRc.Next() {
				kg := domain.Radcheck{}
				err = rowsRc.Scan(
					&kg.ID,
					&kg.Username,
					&kg.Attribute,
					&kg.OP,
					&kg.Value,
				)
				if err != nil {
					logrus.Error(err)
					return nil, err
				}
				rc = append(rc, kg)
			}
			rowsRc.Close()
			rp.Package = pac
			rp.Radcheck = rc

			t.Radpackage = rp
		}

		result = append(result, t)
	}

	return result, nil
}

func (m *mysqlRepository) Fetch(ctx context.Context, idUsers int64, id int64, limit int64) (res []domain.Transaction, err error) {
	args := make([]interface{}, 0)

	query := `SELECT * FROM transaction WHERE id_reseller = ? `
	args = append(args, idUsers)

	if id != 0 {
		query += " AND id <= ? "
		args = append(args, id)
	}

	query += " ORDER BY id DESC LIMIT ? "
	args = append(args, limit)

	res, err = m.fetch(ctx, query, args...)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	return
}

func (m *mysqlRepository) Get(ctx context.Context, transaction domain.Transaction) (res []domain.Transaction, err error) {
	query := "SELECT * FROM transaction "
	args := make([]interface{}, 0)

	addWhere := false

	if transaction.ID != nil {
		if addWhere {
			query += " AND id = ? "
		} else {
			addWhere = true
			query += " WHERE id = ? "
		}
		args = append(args, *transaction.ID)
	}

	if transaction.IDReseller != nil {
		if addWhere {
			query += " AND id_reseller = ? "
		} else {
			addWhere = true
			query += " WHERE id_reseller = ? "
		}
		args = append(args, *transaction.IDReseller)
	}

	if transaction.IDRadpackage != nil {
		if addWhere {
			query += " AND id_radpackage = ? "
		} else {
			addWhere = true
			query += " WHERE id_radpackage = ? "
		}
		args = append(args, *transaction.IDRadpackage)
	}

	if transaction.TransactionCode != nil {
		if addWhere {
			query += " AND transaction_code = ? "
		} else {
			addWhere = true
			query += " WHERE transaction_code = ? "
		}
		args = append(args, *transaction.TransactionCode)
	}

	if transaction.Status != nil {
		if addWhere {
			query += " AND status = ? "
		} else {
			addWhere = true
			query += " WHERE status = ? "
		}
		args = append(args, *transaction.Status)
	}

	if transaction.Value != nil {
		if addWhere {
			query += " AND value = ? "
		} else {
			addWhere = true
			query += " WHERE value = ? "
		}
		args = append(args, *transaction.Value)
	}

	if transaction.Information != nil {
		if addWhere {
			query += " AND information = ? "
		} else {
			addWhere = true
			query += " WHERE information = ? "
		}
		args = append(args, *transaction.Information)
	}

	res, err = m.fetch(ctx, query, args...)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	return
}

func (m *mysqlRepository) CountPage(ctx context.Context, transaction domain.Transaction) (res int64, err error) {
	query := "SELECT COUNT(id) as total_page FROM transaction "
	args := make([]interface{}, 0)

	addWhere := false

	if transaction.ID != nil {
		if addWhere {
			query += " AND id = ? "
		} else {
			addWhere = true
			query += " WHERE id = ? "
		}
		args = append(args, *transaction.ID)
	}

	if transaction.IDReseller != nil {
		if addWhere {
			query += " AND id_reseller = ? "
		} else {
			addWhere = true
			query += " WHERE id_reseller = ? "
		}
		args = append(args, *transaction.IDReseller)
	}

	if transaction.IDRadpackage != nil {
		if addWhere {
			query += " AND id_radpackage = ? "
		} else {
			addWhere = true
			query += " WHERE id_radpackage = ? "
		}
		args = append(args, *transaction.IDRadpackage)
	}

	if transaction.TransactionCode != nil {
		if addWhere {
			query += " AND transaction_code = ? "
		} else {
			addWhere = true
			query += " WHERE transaction_code = ? "
		}
		args = append(args, *transaction.TransactionCode)
	}

	if transaction.Status != nil {
		if addWhere {
			query += " AND status = ? "
		} else {
			addWhere = true
			query += " WHERE status = ? "
		}
		args = append(args, *transaction.Status)
	}

	if transaction.Value != nil {
		if addWhere {
			query += " AND value = ? "
		} else {
			addWhere = true
			query += " WHERE value = ? "
		}
		args = append(args, *transaction.Value)
	}

	if transaction.Information != nil {
		if addWhere {
			query += " AND information = ? "
		} else {
			addWhere = true
			query += " WHERE information = ? "
		}
		args = append(args, *transaction.Information)
	}

	rows, err := m.Conn.QueryContext(ctx, query, args...)
	if err != nil {
		logrus.Error(err)
		return 0, err
	}

	defer func() {
		errRow := rows.Close()
		if errRow != nil {
			logrus.Error(errRow)
		}
	}()

	for rows.Next() {
		err = rows.Scan(&res)
		if err != nil {
			logrus.Error(err)
			return 0, err
		}
	}

	return
}
