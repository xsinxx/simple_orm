package accessLog

import (
	"context"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/simple_orm"
	"github.com/simple_orm/model"
	"testing"
)

func TestMiddlewareBuilder_Build(t *testing.T) {
	logMiddleWare := NewLogMiddleWare()
	options := simple_orm.DBWithMiddleWare(logMiddleWare.Build())
	mockDB, _, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = mockDB.Close() }()
	db, err := simple_orm.OpenDB(mockDB, options)
	if err != nil {
		t.Fatal(err)
	}
	selectQuery := simple_orm.NewSelector[model.TestModel](db)
	selectQuery.Get(context.Background())
}
