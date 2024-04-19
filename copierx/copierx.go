package copierx

import (
	"database/sql"
	"errors"
	"time"

	"github.com/jinzhu/copier"
	"github.com/vmee/go-tools/constant"
)

var (
	timeToStringConverter = copier.TypeConverter{
		SrcType: time.Time{},
		DstType: copier.String,
		Fn: func(src interface{}) (interface{}, error) {
			s, ok := src.(time.Time)
			if !ok {
				return nil, errors.New("src type not matching")
			}
			return s.Format(constant.DateTime), nil
		},
	}

	stringToSqlNullStringConverter = copier.TypeConverter{
		SrcType: copier.String,
		DstType: sql.NullString{},
		Fn: func(src interface{}) (interface{}, error) {
			s, ok := src.(string)
			if !ok {
				return nil, errors.New("src type not matching")
			}

			if s == "" {
				return sql.NullString{
					Valid: false,
				}, nil
			}

			return sql.NullString{
				String: s,
				Valid:  true,
			}, nil
		},
	}

	stringToSqlNullTimeConverter = copier.TypeConverter{
		SrcType: copier.String,
		DstType: sql.NullTime{},
		Fn: func(src interface{}) (interface{}, error) {
			s, ok := src.(string)
			if !ok {
				return nil, errors.New("src type not matching")
			}

			if s == "" {
				return sql.NullTime{
					Valid: false,
				}, nil
			}

			// 解析日期时间字符串为时间对象
			dateTime, err := time.Parse(constant.DateTime, s)
			if err == nil {
				return sql.NullTime{
					Time:  dateTime,
					Valid: true,
				}, nil
			}

			// 解析日期字符串为时间对象
			date, err := time.Parse(constant.DateOnly, s)
			if err == nil {
				return sql.NullTime{
					Time:  date,
					Valid: true,
				}, nil
			}

			return sql.NullTime{
				Valid: false,
			}, nil
		},
	}

	SqlNullTimeToStringConverter = copier.TypeConverter{
		SrcType: sql.NullTime{},
		DstType: copier.String,
		Fn: func(src interface{}) (interface{}, error) {
			s, ok := src.(sql.NullTime)
			if !ok {
				return nil, errors.New("src type not matching")
			}
			if !s.Valid {
				return "", nil
			}
			return s.Time.Format(constant.DateTime), nil
		},
	}
)

// Copy copy things
func Copy(toValue interface{}, fromValue interface{}) (err error) {
	return copier.CopyWithOption(toValue, fromValue, copier.Option{
		Converters: []copier.TypeConverter{
			timeToStringConverter,
			stringToSqlNullStringConverter,
			SqlNullTimeToStringConverter,
			stringToSqlNullTimeConverter,
		},
	})
}
