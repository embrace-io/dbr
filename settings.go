package dbr

import (
	"errors"
	"strings"

	"github.com/embrace-io/dbr/dialect"
)

var (
	ErrUnsupportedDialectForSettings = errors.New("only the Clickhouse dialect supports Settings")
)

type QuerySettings []string

func (qs QuerySettings) Append(setting, value string) QuerySettings {
	setting = strings.TrimSpace(setting)
	value = strings.TrimSpace(value)
	qs = append(qs, setting+" = "+value)
	return qs
}

// Build writes each setting in the form of "SETTINGS setting_key=value[,setting_key=value] \n"
func (qs QuerySettings) Build(d Dialect, buf Buffer) error {
	if d != dialect.Clickhouse {
		return ErrUnsupportedDialectForSettings
	}
	if len(qs) == 0 {
		return nil
	}
	_, err := buf.WriteString("\nSETTINGS " + strings.Join(qs, ", "))
	return err
}
