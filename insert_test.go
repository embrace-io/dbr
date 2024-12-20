package dbr

import (
	"testing"

	"github.com/gocraft/dbr/v2/dialect"
	"github.com/stretchr/testify/require"
)

type insertTest struct {
	A int
	C string `db:"b"`
}

func TestInsertStmt(t *testing.T) {
	buf := NewBuffer()
	builder := InsertInto("table").Ignore().Columns("a", "b").Values(1, "one").Record(&insertTest{
		A: 2,
		C: "two",
	}).Comment("INSERT TEST")
	err := builder.Build(dialect.MySQL, buf)
	require.NoError(t, err)
	require.Equal(t, "/* INSERT TEST */\nINSERT IGNORE INTO `table` (`a`,`b`) VALUES (?,?), (?,?)", buf.String())
	require.Equal(t, []interface{}{1, "one", 2, "two"}, buf.Value())
}

func TestPostgresReturning(t *testing.T) {
	sess := postgresSession
	reset(t, sess)

	var person dbrPerson
	err := sess.InsertInto("dbr_people").Columns("name").Record(&person).
		Returning("id").Load(&person.Id)
	require.NoError(t, err)
	require.True(t, person.Id > 0)
	require.Len(t, sess.EventReceiver.(*testTraceReceiver).started, 1)
	require.Contains(t, sess.EventReceiver.(*testTraceReceiver).started[0].eventName, "dbr.select")
	require.Contains(t, sess.EventReceiver.(*testTraceReceiver).started[0].query, "INSERT")
	require.Contains(t, sess.EventReceiver.(*testTraceReceiver).started[0].query, "dbr_people")
	require.Contains(t, sess.EventReceiver.(*testTraceReceiver).started[0].query, "name")
	require.Equal(t, 1, sess.EventReceiver.(*testTraceReceiver).finished)
	require.Equal(t, 0, sess.EventReceiver.(*testTraceReceiver).errored)
}

func TestOnConflict(t *testing.T) {
	for _, sess := range testSession {
		if sess.Dialect.OnConflict("") == "" {
			// dialect does not support OnConflict
			continue
		}
		t.Run(testSessionName(sess), func(t *testing.T) {
			reset(t, sess)
			for i := 0; i < 2; i++ {
				b := sess.InsertInto("dbr_people").Columns("id", "name", "email").Values(1, "test", "test@test.com")
				b.OnConflict("dbr_people_pkey").Action("email", Expr("CONCAT(?, 2)", Proposed("email")))
				_, err := b.Exec()
				require.NoError(t, err)
			}
			var value string
			_, err := sess.SelectBySql("SELECT email FROM dbr_people WHERE id=?", "1").Load(&value)
			require.NoError(t, err)
			require.Equal(t, "test@test.com2", value)
		})
	}
}

func BenchmarkInsertValuesSQL(b *testing.B) {
	buf := NewBuffer()
	for i := 0; i < b.N; i++ {
		InsertInto("table").Columns("a", "b").Values(1, "one").Build(dialect.MySQL, buf)
	}
}

func BenchmarkInsertRecordSQL(b *testing.B) {
	buf := NewBuffer()
	for i := 0; i < b.N; i++ {
		InsertInto("table").Columns("a", "b").Record(&insertTest{
			A: 2,
			C: "two",
		}).Build(dialect.MySQL, buf)
	}
}
