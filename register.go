package simple_orm

import (
	"errors"
	"reflect"
	"sync"
	"unicode"
)

type field struct {
	columnName string
}

type tableModel struct {
	tableName string            // 表名
	fieldMap  map[string]*field // 表中各个字段
}

// Registry 注册中心，存储表信息
type Registry struct {
	lock   sync.RWMutex // 防止读写冲突
	models map[reflect.Type]*tableModel
}

func (r *Registry) get(val any) (*tableModel, error) {
	r.lock.RLock()
	typ := reflect.TypeOf(val)
	model, ok := r.models[typ]
	r.lock.RUnlock()
	if ok {
		return model, nil
	}
	r.lock.Lock()
	defer r.lock.Unlock()
	model, err := r.parseModel(typ)
	if err != nil {
		return nil, err
	}
	r.models[typ] = model
	return model, nil
}

func (r *Registry) parseModel(typ reflect.Type) (*tableModel, error) {
	for typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	if typ.Kind() != reflect.Struct {
		return nil, errors.New("type is wrong")
	}
	fieldMap := map[string]*field{}
	for i := 0; i < typ.NumField(); i++ {
		fd := typ.Field(i)
		fdName := fd.Name
		// 支持标签，标签的优先级高于字段名
		/*
			type student struct {
				name string `orm:"title"`
				age  int
			}
			fieldMap中的key&value是title，若没配置orm则key&value是name
		*/
		tag := fd.Tag.Get("orm")
		if tag != "" {
			fdName = tag
		}
		fieldMap[fdName] = &field{
			columnName: underscoreName(fdName),
		}
	}
	return &tableModel{
		tableName: underscoreName(typ.Name()),
		fieldMap:  fieldMap,
	}, nil
}

// underscoreName 驼峰转字符串命名
func underscoreName(tableName string) string {
	var buf []byte
	for i, v := range tableName {
		if unicode.IsUpper(v) {
			if i != 0 {
				buf = append(buf, '_')
			}
			buf = append(buf, byte(unicode.ToLower(v)))
		} else {
			buf = append(buf, byte(v))
		}
	}
	return string(buf)
}
