package ntypespqt

import (
	"go/types"
	"strings"

	"github.com/piotrkowalczuk/pqt"
	"github.com/piotrkowalczuk/pqt/pqtgo"
)

// Plugin ...
type Plugin struct{}

// PropertyType implements pqtgo Plugin interface.
func (*Plugin) PropertyType(c *pqt.Column, m int32) string {
	switch {
	case useString(c, m):
		return "ntypes.String"
	case useStringArray(c, m):
		return "ntypes.StringArray"
	case useInt64(c, m):
		return "ntypes.Int64"
	case useInt64Array(c, m):
		return "ntypes.Int64Array"
	case useFloat64(c, m):
		return "ntypes.Float64"
	case useFloat64Array(c, m):
		return "ntypes.Float64Array"
	case useBool(c, m):
		return "ntypes.Bool"
	case useBoolArray(c, m):
		return "ntypes.BoolArray"
	}
	return ""
}

// WhereClause implements pqtgo Plugin interface.
func (p *Plugin) WhereClause(c *pqt.Column) string {
	txt := `
		if {{ .selector }}.Valid {
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

// ScanClause implements pqtgo Plugin interface.
func (p *Plugin) ScanClause(c *pqt.Column) string {
	return ""
}

// SetClause implements pqtgo Plugin interface.
func (p *Plugin) SetClause(c *pqt.Column) string {
	txt := func(t string) string {
		r := `
		if {{ .selector }}.Valid {
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
			}`
		r += t
		r += `
			{{ .composer }}.Dirty=true
		}`
		return r
	}
	switch {
	case useString(c, 2), useFloat64(c, 2), useInt64(c, 2), useBool(c, 2), useBool(c, 3):
		return txt(`
		{{ .composer }}.Add({{ .selector }})`)
	case useStringArray(c, 2):
		switch c.NotNull {
		case true:
			return txt(`
			if {{ .selector }}.StringArray == nil {
				{{ .selector }}.StringArray = []string{}
			}
			{{ .composer }}.Add({{ .selector }})`)
		case false:
			return txt(`
			{{ .composer }}.Add({{ .selector }})`)
		}
	case useFloat64Array(c, 2):
		switch c.NotNull {
		case true:
			return txt(`
			if {{ .selector }}.Float64Array == nil {
				{{ .selector }}.Float64Array = []float64{}
			}
			{{ .composer }}.Add({{ .selector }})`)
		case false:
			return txt(`
			{{ .composer }}.Add({{ .selector }})`)
		}
	case useInt64Array(c, 2):
		switch c.NotNull {
		case true:
			return txt(`
			if {{ .selector }}.Int64Array == nil {
				{{ .selector }}.Int64Array = []int64{}
			}
			{{ .composer }}.Add({{ .selector }})`)
		case false:
			return txt(`
			{{ .composer }}.Add({{ .selector }})`)
		}
	case useBoolArray(c, 2):
		switch c.NotNull {
		case true:
			return txt(`
			if {{ .selector }}.BoolArray == nil {
				{{ .selector }}.BoolArray = []bool{}
			}
			{{ .composer }}.Add({{ .selector }})`)
		case false:
			return txt(`
			{{ .composer }}.Add({{ .selector }})`)
		}
	}
	return ""
}

// Static implements pqtgo Plugin interface.
func (p *Plugin) Static(s *pqt.Schema) string {
	return ""
}

func useInt64Array(c *pqt.Column, m int32) (use bool) {
	if ignore(c, m) {
		return false
	}
	switch {
	case strings.HasPrefix(c.Type.String(), "INTEGER["):
		use = true
	case strings.HasPrefix(c.Type.String(), "BIGINT["):
		use = true
	case strings.HasPrefix(c.Type.String(), "SMALLINT["):
		use = true
	}
	return
}

func useInt64(c *pqt.Column, m int32) (use bool) {
	if ignore(c, m) {
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
		}
	}
	return
}

func useFloat64Array(c *pqt.Column, m int32) (use bool) {
	if ignore(c, m) {
		return false
	}

	switch {
	case strings.HasPrefix(c.Type.String(), "DOUBLE PRECISION["):
		use = true
	}
	return
}

func useFloat64(c *pqt.Column, m int32) (use bool) {
	if ignore(c, m) {
		return false
	}
	switch t := c.Type.(type) {
	case pqtgo.BuiltinType:
		switch types.BasicKind(t) {
		case types.Float32:
			use = true
		case types.Float64:
			use = true
		}
	case pqt.BaseType:
		switch t {
		case pqt.TypeDoublePrecision():
			use = true
		default:
			switch {
			case strings.HasPrefix(t.String(), "DECIMAL"):
				use = true
			case strings.HasPrefix(t.String(), "NUMERIC"):
				use = true
			}
		}
	}
	return
}

func useStringArray(c *pqt.Column, m int32) (use bool) {
	if ignore(c, m) {
		return false
	}
	switch {
	case strings.HasPrefix(c.Type.String(), "TEXT["):
		use = true
	case strings.HasPrefix(c.Type.String(), "VARCHAR["):
		use = true
	case strings.HasPrefix(c.Type.String(), "CHARACTER["):
		use = true
	}
	return
}

func useString(c *pqt.Column, m int32) (use bool) {
	if ignore(c, m) {
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
			switch c.Type.String() {
			case "TEXT", "VARCHAR", "CHARACTER":
				use = true
			}
		}
	}
	return
}

func useBoolArray(c *pqt.Column, m int32) (use bool) {
	if ignore(c, m) {
		return false
	}
	switch {
	case strings.HasPrefix(c.Type.String(), "BOOL["):
		use = true
	}
	return
}

func useBool(c *pqt.Column, m int32) (use bool) {
	if m == 1 && (c.NotNull || c.PrimaryKey) {
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
		}
	}
	return
}

func ignore(c *pqt.Column, m int32) bool {
	return !(m == 2 || (m < 2 && !c.PrimaryKey && !c.NotNull))
}
