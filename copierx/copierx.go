package copierx

import (
	"database/sql"
	"errors"
	"time"

	"github.com/jinzhu/copier"
	"github.com/vmee/go-tools/constant"
)

const (
	Uint64 uint64 = 0
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

	stringToTimeConveter = copier.TypeConverter{
		SrcType: copier.String,
		DstType: time.Time{},
		Fn: func(src interface{}) (interface{}, error) {
			s, ok := src.(string)
			if !ok {
				return nil, errors.New("src type not matching")
			}
			if s == "" {
				return time.Time{}, nil
			}
			t, err := time.Parse(constant.DateTime, s)
			if err != nil {
				return nil, err
			}
			return t, nil
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

	uint64ToSqlNullInt64Converter = copier.TypeConverter{
		SrcType: Uint64,
		DstType: sql.NullInt64{},
		Fn: func(src interface{}) (interface{}, error) {
			s, ok := src.(uint64)
			if !ok {
				return nil, errors.New("src type not matching")
			}

			if s == 0 {
				return sql.NullInt64{
					Valid: false,
				}, nil
			}

			return sql.NullInt64{
				Int64: int64(s),
				Valid: true,
			}, nil
		},
	}

	sqlNullInt64ToUint64Converter = copier.TypeConverter{
		SrcType: sql.NullInt64{},
		DstType: Uint64,
		Fn: func(src interface{}) (interface{}, error) {
			s, ok := src.(sql.NullInt64)
			if !ok {
				return nil, errors.New("src type not matching")
			}
			if !s.Valid {
				return 0, nil
			}
			return uint64(s.Int64), nil
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
			stringToTimeConveter,
			stringToSqlNullStringConverter,
			SqlNullTimeToStringConverter,
			stringToSqlNullTimeConverter,
			uint64ToSqlNullInt64Converter,
			sqlNullInt64ToUint64Converter,
		},
	})
}
