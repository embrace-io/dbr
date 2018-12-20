package dbr

func When(cond Builder, value interface{}) Builder {
	return BuildFunc(func(d Dialect, buf Buffer) error {
		buf.WriteString("when (")
		err := cond.Build(d, buf)
		if err != nil {
			return err
		}
		buf.WriteString(") then ")

		builder, ok := value.(Builder)
		if ok {
			err = builder.Build(d, buf)
			if err != nil {
				return err
			}
			return nil
		}

		buf.WriteString(placeholder)
		buf.WriteValue(value)
		return nil
	})
}

func Case(values ...interface{}) Builder {
	return BuildFunc(func(d Dialect, buf Buffer) error {
		buf.WriteString("case (")
		l := len(values)
		for i, value := range values {
			if l > 1 && i == l-1 {
				buf.WriteString("else ")
			}
			switch v := value.(type) {
			case Builder:
				v.Build(d, buf)
			default:
				buf.WriteString(placeholder)
				buf.WriteValue(v)
			}
			if i < l-1 {
				buf.WriteString(" ")
			}
		}
		buf.WriteString(")")
		return nil
	})
}
