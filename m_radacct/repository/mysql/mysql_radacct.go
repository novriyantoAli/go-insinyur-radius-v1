package mysql

import (
	"context"
	"database/sql"

	"github.com/novriyantoAli/go-insinyur-radius-v1/domain"
	"gorm.io/gorm"

	"github.com/sirupsen/logrus"
)

type mysqlRepository struct {
	Conn  *sql.DB
	Conng *gorm.DB
}

// NewMysqlRepository ...
func NewMysqlRepository(conn *sql.DB, conng *gorm.DB) domain.RadacctRepository {
	return &mysqlRepository{Conn: conn, Conng: conng}
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

func (m *mysqlRepository) FindUsernameIn(ctx context.Context, usernamein []string) (res []domain.Rdacct, err error) {
	err = m.Conng.Model(domain.Rdacct{}).Preload("Nas").Where("acctstoptime IS NULL AND username IN ? ", usernamein).Find(&res).Error
	return
}

func (m *mysqlRepository) FindToday(ctx context.Context, todaystart string, todayend string) (res []domain.Rdacct, err error) {
	err = m.Conng.Where("acctstarttime BETWEEN ? AND ?", todaystart, todayend).Group("username").Find(&res).Error
	return
}

func (m *mysqlRepository) FindWeek(ctx context.Context, weekbreakdown int) (res []domain.Rdacct, err error) {
	err = m.Conng.Model(domain.Rdacct{}).Where("WEEK(acctstarttime) = WEEK(acctstarttime) - ?", weekbreakdown).Group("username").Find(&res).Error
	return
}

func (m *mysqlRepository) FindMonth(ctx context.Context, targetmonth int, targetyear int) (res []domain.Rdacct, err error) {
	err = m.Conng.Model(domain.Rdacct{}).Where("MONTH(acctstarttime) = ? AND YEAR(acctstarttime) = ? ", targetmonth, targetyear).Group("username").Find(&res).Error
	return
}

func (m *mysqlRepository) Find(ctx context.Context, param *domain.Rdacct) (res []domain.Rdacct, err error) {
	err = m.Conng.Find(&res, param).Error
	return
}

func (m *mysqlRepository) Save(ctx context.Context, param *domain.Rdacct) (err error) {
	err = m.Conng.Save(param).Error
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
