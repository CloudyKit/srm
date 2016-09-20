package scheme

import (
	"unicode"
	"bytes"
)

func genName(text string) string {
	buf := bytes.NewBuffer(make([]byte, 0, len(text)))
	for k, r := range text {
		if unicode.IsUpper(r) {
			buf.WriteRune(unicode.ToLower(r))
		} else {
			buf.WriteRune(r)
			if unicode.IsLower(r) {
				if k + 1 < len(text) && unicode.IsUpper(text[k + 1]) {
					buf.WriteByte('_')
				}
			}
		}
	}
	return buf.String()
}

// HasOneEmbed creates a relation between the current field and his target field
func (f FieldDef) HasOneEmbed(args ...string) FieldDef {
	f.f.RelKind = HasOneEmbed
	if len(args) > 0 {
		f.f.RelDst = args[0]
	}
	return f
}

func (f FieldDef) HasManyEmbed(args ...string) FieldDef {
	f.f.RelKind = HasManyEmbed
	if len(args) > 0 {
		f.f.RelDst = args[0]
	}
	return f
}

func (f FieldDef) HasOne(args ...string) FieldDef {
	f.f.RelKind = HasOne
	if len(args) > 0 {
		f.f.RelDst = args[0]
	}
	return f
}

func (f FieldDef) HasMany(args ...string) FieldDef {
	f.f.RelKind = HasMany
	if len(args) > 0 {
		f.f.RelDst = args[0]
	}
	return f
}

func (f FieldDef) Belongs(args ...string) FieldDef {
	f.f.RelKind = Belongs
	if len(args) > 0 {
		f.f.RelDst = args[0]
	}
	return f
}

func (f FieldDef) BelongsEmbed(args ...string) FieldDef {
	f.f.RelKind = BelongsEmbed
	if len(args) > 0 {
		f.f.RelDst = args[0]
	}
	return f
}

