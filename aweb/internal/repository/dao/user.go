package dao

import (
	"context"
	"errors"

	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

var (
	ErrDuplicateEmail = errors.New("邮箱冲突")
	ErrRecordNotFound = gorm.ErrRecordNotFound
)

type UserDAO struct {
	db *gorm.DB
}

func NewUserDAO(db *gorm.DB) *UserDAO {
	return &UserDAO{
		db: db,
	}
}

func (dao *UserDAO) Insert(ctx context.Context, u User) error {
	// now := time.Now().UnixMilli()
	// u.Ctime = now
	// u.Utime = now
	err := dao.db.WithContext(ctx).Create(&u).Error
	if me, ok := err.(*mysql.MySQLError); ok {
		const duplicateErr uint16 = 1062
		if me.Number == duplicateErr {
			// 用户冲突，邮箱冲突
			return ErrDuplicateEmail
		}
	}
	return err
}

func (dao *UserDAO) FindByEmail(ctx context.Context, email string) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Where("email=?", email).First(&u).Error
	return u, err
}

func (dao *UserDAO) UpdateById(ctx context.Context, userId int64, updateFields map[string]interface{}) error {
	return dao.db.WithContext(ctx).Model(&User{}).Where("id = ?", userId).Updates(updateFields).Error
}

func (dao *UserDAO) FindById(ctx context.Context, userId int64) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Where("id=?", userId).First(&u).Error
	return u, err
}

type User struct {
	Id       int64  `gorm:"primaryKey,autoIncrement"`
	Email    string `gorm:"unique"`
	Password string

	Age int `gorm:"default:0;index:idx_age"`
	// 时区，UTC 0 的毫秒数
	// 创建时间
	Ctime int64 `gorm:"autoCreateTime"`
	// 更新时间
	Utime int64 `gorm:"autoUpdateTime"`

	Nickname string `gorm:"default:'';type=varchar(128)"`

	Birthday int64 `gorm:"default:0"`

	Description string `gorm:"default:'';type=varchar(4096)"`

	// json 存储
	//Addr string
}

//type Address struct {
//	Uid
//}
