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
func NewMysqlRepository(conn *sql.DB) domain.RadacctRepository {
	return &mysqlRepository{Conn: conn}
}

func (m *mysqlRepository) fetch(ctx context.Context, query string, args ...interface{}) (result []domain.Radacct, err error) {
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

	result = make([]domain.Radacct, 0)
	for rows.Next() {
		t := domain.Radacct{}
		err = rows.Scan(
			&t.Radacctid,
			&t.Acctsessionid,
			&t.Acctuniqueid,
			&t.Username,
			&t.Realm,
			&t.Nasipaddress,
			&t.Nasportid,
			&t.Nasporttype,
			&t.Acctstarttime,
			&t.Acctupdatetime,
			&t.Acctstoptime,
			&t.Acctinterval,
			&t.Acctsessiontime,
			&t.Acctauthentic,
			&t.ConnectinfoStart,
			&t.ConnectinfoStop,
			&t.Acctinputoctets,
			&t.Acctoutputoctets,
			&t.Calledstationid,
			&t.Callingstationid,
			&t.Acctterminatecause,
			&t.Servicetype,
			&t.Framedprotocol,
			&t.Framedipaddress,
			&t.Secret,
		)

		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		result = append(result, t)
	}

	return result, nil
}

func (m *mysqlRepository) Get(ctx context.Context, radacct domain.Radacct) (res []domain.Radacct, err error) {
	query := "SELECT radacct.*, nas.secret FROM radacct INNER JOIN nas ON nas.nasname = radacct.nasipaddress "
	addWhere := false

	args := make([]interface{}, 0)

	if radacct.Radacctid != nil {
		if addWhere == false {
			addWhere = true
			query += "WHERE radacct.radacctid = ? "
		} else {
			query += "AND radacct.radacctid = ? "
		}
		args = append(args, *radacct.Radacctid)
	}

	if radacct.Username != nil {
		if addWhere == false {
			addWhere = true
			query += "WHERE radacct.username = ? "
		} else {
			query += "AND radacct.username = ? "
		}
		args = append(args, *radacct.Username)
	}

	query += " ORDER by radacct.radacctid DESC "
	res, err = m.fetch(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	return
}

func (m *mysqlRepository) FetchWithUsernameBatch(ctx context.Context, usernameList string) (res []domain.Radacct, err error) {
	query := "SELECT radacct.*, nas.secret FROM radacct INNER JOIN nas ON nas.nasname = radacct.nasipaddress WHERE acctstoptime is NULL AND username IN(" + usernameList + ")"

	res, err = m.fetch(ctx, query)
	if err != nil {
		return nil, err
	}

	return
}
