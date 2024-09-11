package mysql

import (
	"database/sql"

	"github.com/novriyantoAli/go-insinyur-radius-v1/domain"
	"gorm.io/gorm"
)

// mysqlRepository ...
type mysqlRepository struct {
	Conn *sql.DB
	Gorm *gorm.DB
}

// NewMysqlRepository ...
func NewMysqlRepository(conn *sql.DB, gorm *gorm.DB) domain.UsersRepository {
	return &mysqlRepository{Conn: conn, Gorm: gorm}
}

func (m *mysqlRepository) Find(user *domain.Usr) (res []domain.Usr, err error) {
	// err = m.Gorm.Find(&res, user).Error
	err = m.Gorm.Model(&domain.Usr{}).Preload("UserLevel").Find(&res, user).Error
	return
}

func (m *mysqlRepository) LevelFind(level *domain.Level) (res []domain.Level, err error) {
	err = m.Gorm.Find(&res, level).Error
	return
}

func (m *mysqlRepository) First(user *domain.Usr) (res domain.Usr, err error) {
	err = m.Gorm.Where(&user).First(&res).Error
	return
}

func (m *mysqlRepository) LevelFirst(level *domain.Level) (res domain.Level, err error) {
	err = m.Gorm.Where(&level).First(&res).Error
	return
}

func (m *mysqlRepository) Save(user *domain.Usr) (err error) {
	err = m.Gorm.Save(&user).Error
	return
}

func (m *mysqlRepository) LevelSave(level *domain.Level) (err error) {
	err = m.Gorm.Save(&level).Error
	return
}

func (m *mysqlRepository) Delete(user *domain.Usr) (err error) {
	err = m.Gorm.Delete(&user).Error
	return
}

func (m *mysqlRepository) LevelDelete(level *domain.Level) (err error) {
	err = m.Gorm.Delete(&level).Error
	return
}

// func (m *mysqlRepository) fetch(c context.Context, query string, args ...interface{}) (res []domain.Users, er error) {
// 	rows, err := m.Conn.QueryContext(c, query, args...)
// 	if err != nil {
// 		logrus.Error(err)
// 		return nil, err
// 	}

// 	defer func() {
// 		errRow := rows.Close()
// 		if errRow != nil {
// 			logrus.Error(errRow)
// 		}
// 	}()

// 	res = make([]domain.Users, 0)
// 	for rows.Next() {
// 		t := domain.Users{}
// 		err = rows.Scan(
// 			&t.ID,
// 			&t.Username,
// 			&t.Password,
// 			&t.Level,
// 			&t.CreatedAt,
// 		)

// 		if err != nil {
// 			logrus.Error(err)
// 			return nil, err
// 		}
// 		res = append(res, t)
// 	}

// 	return res, nil
// }

// func (m *mysqlRepository) Search(ctx context.Context, user domain.Users) (res []domain.Users, err error) {
// 	query := `SELECT * FROM users `
// 	args := make([]interface{}, 0)

// 	addWhere := false

// 	if user.ID != nil {
// 		if addWhere == false {
// 			query += " WHERE id LIKE '%?%' "
// 			addWhere = true
// 		} else {
// 			query += " OR id LIKE '%?%' "
// 		}
// 		args = append(args, *user.ID)
// 	}

// 	if user.Username != nil {
// 		if addWhere == false {
// 			query += " WHERE username LIKE '%?%'"
// 			addWhere = true
// 		} else {
// 			query += " OR username LIKE '%?%' "
// 		}
// 		args = append(args, *user.Username)
// 	}

// 	if user.Level != nil {
// 		if addWhere == false {
// 			query += " WHERE level LIKE '%?%'"
// 			addWhere = true
// 		} else {
// 			query += " OR level LIKE '%?%' "
// 		}
// 		args = append(args, *user.Level)
// 	}

// 	res, err = m.fetch(ctx, query, args)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return
// }

// Update(ctx context.Context, user Users) (res Users, err error)
// func (m *mysqlRepository) Update(ctx context.Context, user domain.Users) (err error) {
// 	tx, err := m.Conn.BeginTx(ctx, nil)

// 	query := "UPDATE users SET  username = ?, password = ?, level = ? WHERE id = ?"
// 	_, err = tx.ExecContext(
// 		ctx, query, *user.Username, *user.Password, *user.Level, *user.ID,
// 	)

// 	if err != nil {
// 		logrus.Error(err)

// 		tx.Rollback()
// 		return err
// 	}

// 	err = tx.Commit()
// 	if err != nil {
// 		logrus.Error(err)
// 		return err
// 	}

// 	return
// }

// func (m *mysqlRepository) Insert(ctx context.Context, user *domain.Users) (err error) {
// 	tx, err := m.Conn.BeginTx(ctx, nil)

// 	query := "INSERT INTO users( nik, password, level) VALUES(?,?,?,?)"
// 	res, err := tx.ExecContext(
// 		ctx, query, *user.Username, *user.Password, *user.Level,
// 	)

// 	if err != nil {
// 		logrus.Error(err)

// 		tx.Rollback()
// 		return err
// 	}

// 	lastID, err := res.LastInsertId()
// 	if err != nil {
// 		logrus.Error(err)

// 		tx.Rollback()
// 		return
// 	}

// 	user.ID = &lastID

// 	err = tx.Commit()
// 	if err != nil {
// 		logrus.Error(err)

// 		return err
// 	}

// 	return nil
// }

// func (m *mysqlRepository) Delete(ctx context.Context, id int64) (err error) {
// 	query := "DELETE FROM users WHERE id = ?"

// 	stmt, err := m.Conn.PrepareContext(ctx, query)
// 	if err != nil {
// 		logrus.Error(err)
// 		return
// 	}

// 	res, err := stmt.ExecContext(ctx, id)
// 	if err != nil {
// 		logrus.Error(err)
// 		return
// 	}

// 	rowsAffected, err := res.RowsAffected()
// 	if err != nil {
// 		logrus.Error(err)
// 		return
// 	}

// 	if rowsAffected != 1 {
// 		err = fmt.Errorf("Weird  Behavior. Total Affected: %d", rowsAffected)
// 		logrus.Error(err)
// 		return
// 	}

// 	return
// }
