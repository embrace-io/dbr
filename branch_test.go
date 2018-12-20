package dbr

import (
	"testing"

	"github.com/embrace-io/dbr/dialect"
	"github.com/stretchr/testify/require"
)

func TestWhen(t *testing.T) {
	for _, test := range []struct {
		when  Builder
		query string
		value []interface{}
	}{
		{
			when:  When(Eq("col", 1), 1),
			query: "when (`col` = ?) then ?",
			value: []interface{}{1, 1},
		},
		{
			when:  When(And(Gt("a", 1), Lt("b", 2)), "c"),
			query: "when ((`a` > ?) AND (`b` < ?)) then ?",
			value: []interface{}{1, 2, "c"},
		},
		{
			when:  When(Eq("a", 1), Gt("b", 2)),
			query: "when (`a` = ?) then `b` > ?",
			value: []interface{}{1, 2},
		},
	} {
		buf := NewBuffer()
		err := test.when.Build(dialect.Clickhouse, buf)
		require.NoError(t, err)
		require.Equal(t, test.query, buf.String())
		require.Equal(t, test.value, buf.Value())
	}
}

func TestCase(t *testing.T) {
	for _, test := range []struct {
		when  Builder
		query string
		value []interface{}
	}{
		{
			when:  Case(When(Eq("col", 1), Expr("?", 1))),
			query: "case (when (`col` = ?) then ?)",
			value: []interface{}{1, 1},
		},
		{
			when:  Case(When(Eq("col", 1), 2), Expr("?", 3)),
			query: "case (when (`col` = ?) then ? else ?)",
			value: []interface{}{1, 2, 3},
		},
		{
			when:  Case(When(Eq("a", 1), 2), Gt("b", 3)),
			query: "case (when (`a` = ?) then ? else `b` > ?)",
			value: []interface{}{1, 2, 3},
		},
		{
			when:  Case(When(Eq("col", "a"), 1), When(Eq("col", "b"), 2), Expr("?", 3)),
			query: "case (when (`col` = ?) then ? when (`col` = ?) then ? else ?)",
			value: []interface{}{"a", 1, "b", 2, 3},
		},
		{
			when:  Case(When(Eq("colA", "a"), Lt("colB", 5)), When(Eq("colA", "b"), Lt("colB", 10)), Lt("colB", 15)),
			query: "case (when (`colA` = ?) then `colB` < ? when (`colA` = ?) then `colB` < ? else `colB` < ?)",
			value: []interface{}{"a", 5, "b", 10, 15},
		},
	} {
		buf := NewBuffer()
		err := test.when.Build(dialect.Clickhouse, buf)
		require.NoError(t, err)
		require.Equal(t, test.query, buf.String())
		require.Equal(t, test.value, buf.Value())
	}
}
