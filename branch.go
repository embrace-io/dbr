package dbr

// When creates a WHEN statement given a condition and a value that's evaluated
// if the condition is true.
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

// Case creates a CASE statement from a list of conditions.
// If there are more than 1 conditions, the last one will be an else statement.
func Case(conds ...Builder) Builder {
	return BuildFunc(func(d Dialect, buf Buffer) error {
		buf.WriteString("(case ")
		l := len(conds)
		for i, cond := range conds {
			if l > 1 && i == l-1 {
				buf.WriteString("else ")
			}
			cond.Build(d, buf)
			if i < l-1 {
				buf.WriteString(" ")
			}
		}
		buf.WriteString(" end)")
		return nil
	})
}
