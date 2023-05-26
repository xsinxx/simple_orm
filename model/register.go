package model

import (
	"errors"
	"reflect"
	"unicode"
)

func NewRegistry() *Registry {
	return &Registry{
		TableModels: map[reflect.Type]*TableModel{},
	}
}

func (r *Registry) Get(val any) (*TableModel, error) {
	r.lock.RLock()
	typ := reflect.TypeOf(val)
	tableModel, ok := r.TableModels[typ]
	r.lock.RUnlock()
	if ok {
		return tableModel, nil
	}
	r.lock.Lock()
	defer r.lock.Unlock()
	tableModel, err := r.parseModel(typ)
	if err != nil {
		return nil, err
	}
	r.TableModels[typ] = tableModel
	return tableModel, nil
}

func (r *Registry) parseModel(typ reflect.Type) (*TableModel, error) {
	for typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	if typ.Kind() != reflect.Struct {
		return nil, errors.New("type is wrong")
	}
	tag2Field := map[string]*Field{}
	col2Field := map[string]*Field{}
	for i := 0; i < typ.NumField(); i++ {
		fd := typ.Field(i)
		fdName := fd.Name
		/*
			type student struct {
				name string `orm:"title"`
				age  int
			}
			fieldMap中的key是title，若没配置orm则key是name
		*/
		tag := fd.Tag.Get("orm")
		// 若不配置标签默认取typeName
		if tag == "" {
			tag = fdName
		}
		field := &Field{
			ColumnName: underscoreName(fdName),
			Typ:        fd.Type,
			TypName:    fd.Name,
			Offset:     fd.Offset,
		}
		tag2Field[tag] = field
		col2Field[fdName] = field
	}
	return &TableModel{
		TableName: underscoreName(typ.Name()),
		Tag2Field: tag2Field,
		Col2Field: col2Field,
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
