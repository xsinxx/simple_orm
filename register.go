package simple_orm

import (
	"errors"
	"reflect"
	"sync"
	"unicode"
)

type field struct {
	ColumnName string // 对应的数据库中表的列
	Typ        reflect.Type
	Offset     uintptr
}

type TableModel struct {
	TableName string            // 表名
	Tag2Field map[string]*field // 标签名到字段的映射
	Col2Field map[string]*field // 列名到字段的映射
}

// Registry 注册中心，存储表信息
type Registry struct {
	lock   sync.RWMutex // 防止读写冲突
	models map[reflect.Type]*TableModel
}

func (r *Registry) Get(val any) (*TableModel, error) {
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

func (r *Registry) parseModel(typ reflect.Type) (*TableModel, error) {
	for typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	if typ.Kind() != reflect.Struct {
		return nil, errors.New("type is wrong")
	}
	tag2Field := map[string]*field{}
	col2Field := map[string]*field{}
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
		field := &field{
			ColumnName: underscoreName(fdName),
			Typ:        fd.Type,
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
