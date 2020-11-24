package databases

import (
	"github.com/jinzhu/gorm"
)

const DefaultPageSize = 10

/**
传入的参数值
pageData {
	page int  // 页码
	pageSize int // 条数
}
*/
func FindPage(db *gorm.DB, pageData map[string]uint64, out interface{}, count *uint64) error {

	if err := db.Count(count).Error; err != nil {
		return err
	}

	pageSize, has := pageData["pageSize"]

	if has {
		db = db.Limit(pageSize)
	} else {
		db = db.Limit(DefaultPageSize)
	}

	// 不传页码为通过最后一条 ID 查询
	if page, has := pageData["page"]; has {
		if page < 0 {
			return nil
		}
		db = db.Offset((page - 1) * pageSize)
	}

	return db.Find(out).Error
}
