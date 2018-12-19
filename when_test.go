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
