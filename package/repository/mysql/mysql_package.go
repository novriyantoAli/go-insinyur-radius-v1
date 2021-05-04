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

// NewMysqlRepository ...
func NewMysqlRepository(conn *sql.DB) domain.PackageRepository {
	return &mysqlRepository{Conn: conn}
}

func (m *mysqlRepository) fetch(ctx context.Context, query string, args ...interface{}) (result []domain.Package, err error) {
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

	result = make([]domain.Package, 0)
	for rows.Next() {
		t := domain.Package{}
		err = rows.Scan(
			&t.ID,
			&t.Name,
			&t.ValidityValue,
			&t.ValidityUnit,
			&t.Price,
			&t.Margin,
			&t.Profile,
			&t.CreatedAt,
		)

		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		result = append(result, t)
	}

	return result, nil
}

func (m *mysqlRepository) Fetch(ctx context.Context, id int64, limit int64) (res []domain.Package, err error) {
	args := make([]interface{}, 0)

	query := `SELECT * FROM package `

	if id != 0 {
		query += "WHERE id <= ? "
		args = append(args, id)
	}

	query += "ORDER BY id DESC LIMIT ?"
	args = append(args, limit)

	res, err = m.fetch(ctx, query, args...)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	return
}

func (m *mysqlRepository) CountPage(ctx context.Context, packages domain.Package) (res int64, err error) {
	query := "SELECT COUNT(id) as total_page FROM package "
	args := make([]interface{}, 0)

	addWhere := false

	if packages.ID != nil {
		if addWhere {
			query += " AND id = ? "
		} else {
			addWhere = true
			query += " WHERE id = ? "
		}
		args = append(args, *packages.ID)
	}

	if packages.Name != nil {
		if addWhere {
			query += " AND name = ? "
		} else {
			addWhere = true
			query += " WHERE name = ? "
		}
		args = append(args, *packages.Name)
	}

	if packages.ValidityValue != nil {
		if addWhere {
			query += " AND validity_value = ? "
		} else {
			addWhere = true
			query += " WHERE validity_value = ? "
		}
		args = append(args, *packages.ValidityValue)
	}

	if packages.ValidityUnit != nil {
		if addWhere {
			query += " AND validity_unit = ? "
		} else {
			addWhere = true
			query += " WHERE validity_unit = ? "
		}
		args = append(args, *packages.ValidityUnit)
	}

	if packages.Price != nil {
		if addWhere {
			query += " AND price = ? "
		} else {
			addWhere = true
			query += " WHERE price = ? "
		}
		args = append(args, *packages.Price)
	}

	if packages.Margin != nil {
		if addWhere {
			query += " AND margin = ? "
		} else {
			addWhere = true
			query += " WHERE margin = ? "
		}
		args = append(args, *packages.Margin)
	}

	if packages.Profile != nil {
		if addWhere {
			query += " AND profile = ? "
		} else {
			addWhere = true
			query += " WHERE profile = ? "
		}
		args = append(args, *packages.Profile)
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

func (m *mysqlRepository) Get(ctx context.Context, packages domain.Package) (res []domain.Package, err error) {
	query := "SELECT * FROM package "
	args := make([]interface{}, 0)

	addWhere := false

	if packages.ID != nil {
		if addWhere {
			query += " AND id = ? "
		} else {
			addWhere = true
			query += " WHERE id = ? "
		}
		args = append(args, *packages.ID)
	}

	if packages.Name != nil {
		if addWhere {
			query += " AND name = ? "
		} else {
			addWhere = true
			query += " WHERE name = ? "
		}
		args = append(args, *packages.Name)
	}

	if packages.ValidityValue != nil {
		if addWhere {
			query += " AND validity_value = ? "
		} else {
			addWhere = true
			query += " WHERE validity_value = ? "
		}
		args = append(args, *packages.ValidityValue)
	}

	if packages.ValidityUnit != nil {
		if addWhere {
			query += " AND validity_unit = ? "
		} else {
			addWhere = true
			query += " WHERE validity_unit = ? "
		}
		args = append(args, *packages.ValidityUnit)
	}

	if packages.Price != nil {
		if addWhere {
			query += " AND price = ? "
		} else {
			addWhere = true
			query += " WHERE price = ? "
		}
		args = append(args, *packages.Price)
	}

	if packages.Margin != nil {
		if addWhere {
			query += " AND margin = ? "
		} else {
			addWhere = true
			query += " WHERE margin = ? "
		}
		args = append(args, *packages.Margin)
	}

	if packages.Profile != nil {
		if addWhere {
			query += " AND profile = ? "
		} else {
			addWhere = true
			query += " WHERE profile = ? "
		}
		args = append(args, *packages.Profile)
	}

	res, err = m.fetch(ctx, query, args...)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	return
}
