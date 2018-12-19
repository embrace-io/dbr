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
