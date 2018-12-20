package dbr

// When creates a WHEN statement given a condition and a value that's evaluated
// if the condition is true.
func When(cond Builder, value interface{}) Builder {
	return BuildFunc(func(d Dialect, buf Buffer) error {
		buf.WriteString("WHEN (")
		err := cond.Build(d, buf)
		if err != nil {
			return err
		}
		buf.WriteString(") THEN ")

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

// Else creates an ELSE statement given a value or a condition.
func Else(value interface{}) Builder {
	return BuildFunc(func(d Dialect, buf Buffer) error {
		buf.WriteString("ELSE ")
		if builder, ok := value.(Builder); ok {
			if err := builder.Build(d, buf); err != nil {
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
		buf.WriteString("(CASE ")
		l := len(conds)
		for i, cond := range conds {
			cond.Build(d, buf)
			if i < l-1 {
				buf.WriteString(" ")
			}
		}
		buf.WriteString(" END)")
		return nil
	})
}
