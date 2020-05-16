package gorm

var templateDB = `// Code generated by gen_batch.go; DO NOT EDIT.
// GENERATED FILE DO NOT EDIT

package {{ .Package }}

import (
	"errors"
	"github.com/jinzhu/gorm"
	"time"
)

var db *gorm.DB

// InitDB - init db handler
func InitDB(handler *gorm.DB) error {
	if handler == nil {
		return errors.New("invalid db handler")
	}

	db = handler
	db.DB().SetMaxIdleConns({{ .MaxIdleConns }})
	db.DB().SetMaxOpenConns({{ .MaxOpenConns }})

	return nil
}

// Model - mode base struct define
type Model struct {
	ID          uint 
	CreatedTime time.Time
	UpdatedTime time.Time
}
`

var templateHeader = `// Code generated by gen_batch.go; DO NOT EDIT.
// GENERATED FILE DO NOT EDIT

package {{ .Package }} 

import (
	"github.com/jinzhu/gorm"
)

`

var templateBody = `
// {{ .Name }} - {{ .Name | ToLower }} detail info
{{- .Detail }}

// TableName Set {{ .Name }}'s table name to be '{{ .Name | ToSnake }}'
func (m *{{ .Name }}) TableName() string {
	return "{{ .Name | ToSnake }}"
}

// Add{{ .Name }} - add an new {{ .Name | ToSnake }}
func Add{{ .Name }}(val *{{ .Name }}) (int64, error) {
	txn := db.Begin()
	if err := txn.Create(val).Error; err != nil {
		txn.Rollback()
		return 0, err
	}

	var id []int64
	if err := txn.Raw("select LAST_INSERT_ID() as id").Pluck("id", &id).Error; err != nil {
		txn.Rollback()
		return 0, err
	}

	if err := txn.Commit().Error; err != nil {
		txn.Rollback()
		return 0, err
	}

	return id[0], nil
}

// Add{{ .Name }}Batch - batch add {{ .Name | ToSnake }}
func Add{{ .Name }}Batch(items []*{{ .Name }}) (int64, error) {
	txn := db.Begin()
	for i := 0; i < len(items); i++ {
		val := items[i]
		if err := txn.Create(val).Error; err != nil {
			txn.Rollback()
			return 0, err
		}
	}

	var id []int64
	if err := txn.Raw("select LAST_INSERT_ID() as id").Pluck("id", &id).Error; err != nil {
		txn.Rollback()
		return 0, err
	}

	if err := txn.Commit().Error; err != nil {
		txn.Rollback()
		return 0, err
	}

	return id[0], nil
}

// Get{{ .Name }} - query {{ .Name | ToSnake }} info by id
func Get{{ .Name }}(id int64) (*{{ .Name }}, error) {
	val := &{{ .Name }}{}
	v := db.Where("id = ?", id).Find(val)
	if v.Error != nil {
		return nil, v.Error
	}
	return val, nil
}

// Edit{{ .Name }} - modify {{ .Name | ToSnake }} info by id
func Edit{{ .Name }}(val *{{ .Name }}, params map[string]interface{}) (int64, error) {
	v := db.Model(val).Update(params)
	if v.Error != nil {
		return 0, v.Error
	}
	return v.RowsAffected, nil
}

// Delete{{ .Name }} - remove {{ .Name | ToSnake }} info by id
func Delete{{ .Name }} (id int64) (int64, error) {
	v := db.Where("id = ?", id).Delete(&{{ .Name }}{})
	if v.Error != nil {
		return 0, v.Error
	}
	return v.RowsAffected, nil
}

// Count{{ .Name }} - query the count of {{ .Name | ToSnake }} info
func Count{{ .Name }}(params map[string]interface{}) (int64, error) {
	var count int64
	v := db.Model(&{{ .Name }}{}).Where(params).Count(&count)
	if v.Error != nil {
		return 0, v.Error
	}
	return count, nil
}

// List{{ .Name }} - query list of {{ .Name | ToSnake }} info
func List{{ .Name }}(offset, limit int64, order string, params map[string]interface{}) ([]{{ .Name }}, error) {
	var list []{{ .Name }}
	var err error

	if offset >= 0 && limit > 0 {
		if order != "" {
			err = db.Where(params).Order(order).Offset(offset).Limit(limit).Find(&list).Error
		} else {
			err = db.Where(params).Offset(offset).Limit(limit).Find(&list).Error
		}
	} else {
		err = db.Where(params).Order(order).Find(&list).Error
	}

	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return list, nil
}`
