package repository

import "gorm.io/gorm"

// tx 在事务中执行 fn，供仓储内部需要多表原子写入时复用。
func tx(db *gorm.DB, fn func(tx *gorm.DB) error) error {
	return db.Transaction(fn)
}
