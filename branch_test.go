package dbr

import (
	"testing"

	"github.com/embrace-io/dbr/dialect"
	"github.com/stretchr/testify/require"
)

func TestWhen(t *testing.T) {
	for _, test := range []struct {
		WHEN  Builder
		query string
		value []interface{}
	}{
		{
			WHEN:  When(Eq("col", 1), 1),
			query: "WHEN (`col` = ?) THEN ?",
			value: []interface{}{1, 1},
		},
		{
			WHEN:  When(And(Gt("a", 1), Lt("b", 2)), "c"),
			query: "WHEN ((`a` > ?) AND (`b` < ?)) THEN ?",
			value: []interface{}{1, 2, "c"},
		},
		{
			WHEN:  When(Eq("a", 1), Gt("b", 2)),
			query: "WHEN (`a` = ?) THEN `b` > ?",
			value: []interface{}{1, 2},
		},
	} {
		buf := NewBuffer()
		err := test.WHEN.Build(dialect.Clickhouse, buf)
		require.NoError(t, err)
		require.Equal(t, test.query, buf.String())
		require.Equal(t, test.value, buf.Value())
	}
}

func TestCase(t *testing.T) {
	for _, test := range []struct {
		WHEN  Builder
		query string
		value []interface{}
	}{
		{
			WHEN:  Case(When(Eq("col", 1), 1)),
			query: "(CASE WHEN (`col` = ?) THEN ? END)",
			value: []interface{}{1, 1},
		},
		{
			WHEN:  Case(When(Eq("col", 1), 2), Else(3)),
			query: "(CASE WHEN (`col` = ?) THEN ? ELSE ? END)",
			value: []interface{}{1, 2, 3},
		},
		{
			WHEN:  Case(When(Eq("a", 1), 2), Else(Gt("b", 3))),
			query: "(CASE WHEN (`a` = ?) THEN ? ELSE `b` > ? END)",
			value: []interface{}{1, 2, 3},
		},
		{
			WHEN:  Case(When(Eq("col", "a"), 1), When(Eq("col", "b"), 2), Else(3)),
			query: "(CASE WHEN (`col` = ?) THEN ? WHEN (`col` = ?) THEN ? ELSE ? END)",
			value: []interface{}{"a", 1, "b", 2, 3},
		},
		{
			WHEN:  Case(When(Eq("colA", "a"), Lt("colB", 5)), When(Eq("colA", "b"), Lt("colB", 10)), Else(Lt("colB", 15))),
			query: "(CASE WHEN (`colA` = ?) THEN `colB` < ? WHEN (`colA` = ?) THEN `colB` < ? ELSE `colB` < ? END)",
			value: []interface{}{"a", 5, "b", 10, 15},
		},
	} {
		buf := NewBuffer()
		err := test.WHEN.Build(dialect.Clickhouse, buf)
		require.NoError(t, err)
		require.Equal(t, test.query, buf.String())
		require.Equal(t, test.value, buf.Value())
	}
}
