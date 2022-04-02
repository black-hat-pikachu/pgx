package zeronull_test

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5/pgtype/zeronull"
	"github.com/jackc/pgx/v5/pgxtest"
)

func isExpectedEq(a interface{}) func(interface{}) bool {
	return func(v interface{}) bool {
		return a == v
	}
}

func TestFloat8Transcode(t *testing.T) {
	pgxtest.RunValueRoundTripTests(context.Background(), t, defaultConnTestRunner, nil, "float8", []pgxtest.ValueRoundTripTest{
		{
			(zeronull.Float8)(1),
			new(zeronull.Float8),
			isExpectedEq((zeronull.Float8)(1)),
		},
		{
			nil,
			new(zeronull.Float8),
			isExpectedEq((zeronull.Float8)(0)),
		},
		{
			(zeronull.Float8)(0),
			new(interface{}),
			isExpectedEq(nil),
		},
	})
}
