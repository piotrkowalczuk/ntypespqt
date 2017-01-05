package ntypespqt

import (
	"go/types"
	"strings"
	"fmt"
	"github.com/piotrkowalczuk/pqt"
	"github.com/piotrkowalczuk/pqt/pqtgo"
)

type Plugin struct {
	Formatter  *pqtgo.Formatter
	Visibility pqtgo.Visibility
}

func (*Plugin) PropertyType(c *pqt.Column, m int32) string {
	switch {
	case useString(c, m):
		return "ntypes.String"
	case useInt64(c, m):
		return "ntypes.Int64"
	case useBool(c, m):
		return "ntypes.Bool"
	}
	return ""
}

// WhereClause implements pqtgo Plugin interface.
func (p *Plugin) WhereClause(c *pqt.Column) string {
	txt := `if {{ .selector }}.Valid {
			if {{ .composer }}.Dirty {
				if _, err := {{ .composer }}.WriteString(", "); err != nil {
					return "", nil, err
				}
			}
			if _, err := {{ .composer }}.WriteString({{ .column }}); err != nil {
				return "", nil, err
			}
			if _, err := {{ .composer }}.WriteString("="); err != nil {
				return "", nil, err
			}
			if err := {{ .composer }}.WritePlaceholder(); err != nil {
				return "", nil, err
			}
			{{ .composer }}.Add({{ .selector }})
			{{ .composer }}.Dirty=true
		}`
	switch {
	case useBool(c, 3):
		return txt
	}
	return ""
}

func (p *Plugin) SetClause(c *pqt.Column) string {
	txt := `if {{ .selector }}.Valid {
			if {{ .composer }}.Dirty {
				if _, err := {{ .composer }}.WriteString(", "); err != nil {
					return "", nil, err
				}
			}
			if _, err := {{ .composer }}.WriteString({{ .column }}); err != nil {
				return "", nil, err
			}
			if _, err := {{ .composer }}.WriteString("="); err != nil {
				return "", nil, err
			}
			if err := {{ .composer }}.WritePlaceholder(); err != nil {
				return "", nil, err
			}
			{{ .composer }}.Add({{ .selector }})
			{{ .composer }}.Dirty=true
		}`
	switch {
	case useString(c, 2):
		return txt
	case useInt64(c, 2):
		return txt
	case useBool(c, 2), useBool(c, 3):
		return txt
	}
	return ""
}

func (p *Plugin) Static(s *pqt.Schema) string {
	return ""
}

func useInt64(c *pqt.Column, m int32) (use bool) {
	if !(m == 2 || (m == 0 && !c.NotNull)) {
		return false
	}
	switch t := c.Type.(type) {
	case pqtgo.BuiltinType:
		switch types.BasicKind(t) {
		case types.Int:
			use = true
		case types.Int8:
			use = true
		case types.Int16:
			use = true
		case types.Int32:
			use = true
		case types.Int64:
			use = true
		}
	case pqt.BaseType:
		switch t {
		case pqt.TypeInteger():
			use = true
		case pqt.TypeIntegerBig():
			use = true
		case pqt.TypeIntegerSmall():
			use = true
		case pqt.TypeSerial():
			use = true
		case pqt.TypeSerialBig():
			use = true
		case pqt.TypeSerialSmall():
			use = true
		default:
			switch {
			case strings.HasPrefix(c.Name, "INTEGER["):
				use = true
			case strings.HasPrefix(c.Name, "BIGINT["):
				use = true
			case strings.HasPrefix(c.Name, "SMALLINT["):
				use = true
			}
		}
	}
	return
}

func useString(c *pqt.Column, m int32) (use bool) {
	if !(m == 2 || (m == 0 && !c.NotNull)) {
		if c.Name == "description" {
			fmt.Printf("%d - %#v\n", m, c)
		}
		return false
	}
	switch t := c.Type.(type) {
	case pqtgo.BuiltinType:
		switch types.BasicKind(t) {
		case types.String:
			use = true
		}
	case pqt.BaseType:
		switch t {
		case pqt.TypeText():
			use = true
		case pqt.TypeUUID():
			use = true
		default:
			switch {
			case strings.HasPrefix(c.Name, "TEXT["):
				use = true
			case strings.HasPrefix(c.Name, "VARCHAR"), strings.HasPrefix(c.Name, "CHARACTER"):
				use = true
			}
		}
	}
	return
}

func useBool(c *pqt.Column, m int32) (use bool) {
	if m == 1 && c.NotNull {
		return false
	}
	switch t := c.Type.(type) {
	case pqtgo.BuiltinType:
		switch types.BasicKind(t) {
		case types.Bool:
			use = true
		}
	case pqt.BaseType:
		switch t {
		case pqt.TypeBool():
			use = true
		default:
			switch {
			case strings.HasPrefix(c.Name, "BOOL["):
				use = true
			}
		}
	}
	return
}
