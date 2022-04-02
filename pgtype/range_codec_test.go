package pgtype_test

import (
	"context"
	"testing"

	pgx "github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxtest"
	"github.com/stretchr/testify/require"
)

func TestRangeCodecTranscode(t *testing.T) {
	skipCockroachDB(t, "Server does not support range types (see https://github.com/cockroachdb/cockroach/issues/27791)")

	pgxtest.RunValueRoundTripTests(context.Background(), t, defaultConnTestRunner, nil, "int4range", []pgxtest.ValueRoundTripTest{
		{
			pgtype.Int4range{LowerType: pgtype.Empty, UpperType: pgtype.Empty, Valid: true},
			new(pgtype.Int4range),
			isExpectedEq(pgtype.Int4range{LowerType: pgtype.Empty, UpperType: pgtype.Empty, Valid: true}),
		},
		{
			pgtype.Int4range{
				LowerType: pgtype.Inclusive,
				Lower:     pgtype.Int4{Int32: 1, Valid: true},
				Upper:     pgtype.Int4{Int32: 5, Valid: true},
				UpperType: pgtype.Exclusive, Valid: true,
			},
			new(pgtype.Int4range),
			isExpectedEq(pgtype.Int4range{
				LowerType: pgtype.Inclusive,
				Lower:     pgtype.Int4{Int32: 1, Valid: true},
				Upper:     pgtype.Int4{Int32: 5, Valid: true},
				UpperType: pgtype.Exclusive, Valid: true,
			}),
		},
		{pgtype.Int4range{}, new(pgtype.Int4range), isExpectedEq(pgtype.Int4range{})},
		{nil, new(pgtype.Int4range), isExpectedEq(pgtype.Int4range{})},
	})
}

func TestRangeCodecTranscodeCompatibleRangeElementTypes(t *testing.T) {
	ctr := defaultConnTestRunner
	ctr.AfterConnect = func(ctx context.Context, t testing.TB, conn *pgx.Conn) {
		pgxtest.SkipCockroachDB(t, conn, "Server does not support range types (see https://github.com/cockroachdb/cockroach/issues/27791)")
	}

	pgxtest.RunValueRoundTripTests(context.Background(), t, ctr, nil, "numrange", []pgxtest.ValueRoundTripTest{
		{
			pgtype.Float8range{LowerType: pgtype.Empty, UpperType: pgtype.Empty, Valid: true},
			new(pgtype.Float8range),
			isExpectedEq(pgtype.Float8range{LowerType: pgtype.Empty, UpperType: pgtype.Empty, Valid: true}),
		},
		{
			pgtype.Float8range{
				LowerType: pgtype.Inclusive,
				Lower:     pgtype.Float8{Float64: 1, Valid: true},
				Upper:     pgtype.Float8{Float64: 5, Valid: true},
				UpperType: pgtype.Exclusive, Valid: true,
			},
			new(pgtype.Float8range),
			isExpectedEq(pgtype.Float8range{
				LowerType: pgtype.Inclusive,
				Lower:     pgtype.Float8{Float64: 1, Valid: true},
				Upper:     pgtype.Float8{Float64: 5, Valid: true},
				UpperType: pgtype.Exclusive, Valid: true,
			}),
		},
		{pgtype.Float8range{}, new(pgtype.Float8range), isExpectedEq(pgtype.Float8range{})},
		{nil, new(pgtype.Float8range), isExpectedEq(pgtype.Float8range{})},
	})
}

func TestRangeCodecScanRangeTwiceWithUnbounded(t *testing.T) {
	skipCockroachDB(t, "Server does not support range types (see https://github.com/cockroachdb/cockroach/issues/27791)")

	defaultConnTestRunner.RunTest(context.Background(), t, func(ctx context.Context, t testing.TB, conn *pgx.Conn) {

		var r pgtype.Int4range

		err := conn.QueryRow(context.Background(), `select '[1,5)'::int4range`).Scan(&r)
		require.NoError(t, err)

		require.Equal(
			t,
			pgtype.Int4range{
				Lower:     pgtype.Int4{Int32: 1, Valid: true},
				Upper:     pgtype.Int4{Int32: 5, Valid: true},
				LowerType: pgtype.Inclusive,
				UpperType: pgtype.Exclusive,
				Valid:     true,
			},
			r,
		)

		err = conn.QueryRow(ctx, `select '[1,)'::int4range`).Scan(&r)
		require.NoError(t, err)

		require.Equal(
			t,
			pgtype.Int4range{
				Lower:     pgtype.Int4{Int32: 1, Valid: true},
				Upper:     pgtype.Int4{},
				LowerType: pgtype.Inclusive,
				UpperType: pgtype.Unbounded,
				Valid:     true,
			},
			r,
		)

		err = conn.QueryRow(ctx, `select 'empty'::int4range`).Scan(&r)
		require.NoError(t, err)

		require.Equal(
			t,
			pgtype.Int4range{
				Lower:     pgtype.Int4{},
				Upper:     pgtype.Int4{},
				LowerType: pgtype.Empty,
				UpperType: pgtype.Empty,
				Valid:     true,
			},
			r,
		)
	})
}

func TestRangeCodecDecodeValue(t *testing.T) {
	skipCockroachDB(t, "Server does not support range types (see https://github.com/cockroachdb/cockroach/issues/27791)")

	defaultConnTestRunner.RunTest(context.Background(), t, func(ctx context.Context, _ testing.TB, conn *pgx.Conn) {

		for _, tt := range []struct {
			sql      string
			expected interface{}
		}{
			{
				sql: `select '[1,5)'::int4range`,
				expected: pgtype.GenericRange{
					Lower:     int32(1),
					Upper:     int32(5),
					LowerType: pgtype.Inclusive,
					UpperType: pgtype.Exclusive,
					Valid:     true,
				},
			},
		} {
			t.Run(tt.sql, func(t *testing.T) {
				rows, err := conn.Query(ctx, tt.sql)
				require.NoError(t, err)

				for rows.Next() {
					values, err := rows.Values()
					require.NoError(t, err)
					require.Len(t, values, 1)
					require.Equal(t, tt.expected, values[0])
				}

				require.NoError(t, rows.Err())
			})
		}
	})
}
