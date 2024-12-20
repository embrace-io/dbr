package dbr

type joinType uint8

const (
	inner joinType = iota
	left
	right
	full
)

func join(t joinType, table interface{}, on interface{}, indexHints []Builder) Builder {
	return BuildFunc(func(d Dialect, buf Buffer) error {
		buf.WriteString(" ")
		switch t {
		case left:
			buf.WriteString("LEFT ")
		case right:
			buf.WriteString("RIGHT ")
		case full:
			buf.WriteString("FULL ")
		}
		buf.WriteString("JOIN ")
		switch table := table.(type) {
		case string:
			buf.WriteString(d.QuoteIdent(table))
		default:
			buf.WriteString(placeholder)
			buf.WriteValue(table)
		}
		for _, hint := range indexHints {
			buf.WriteString(" ")
			if err := hint.Build(d, buf); err != nil {
				return err
			}
		}
		if d.SupportsOn() {
			buf.WriteString(" ON ")
		} else {
			buf.WriteString(" USING ")
		}
		switch on := on.(type) {
		case string:
			buf.WriteString(on)
		case Builder:
			buf.WriteString(placeholder)
			buf.WriteValue(on)
		}
		return nil
	})
}
