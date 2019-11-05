package json

import pgs "github.com/lyft/protoc-gen-star"

func isWellKnown(f pgs.Field) (pgs.WellKnownType, bool) {
	if f.Type().IsEmbed() {
		wkt := pgs.LookupWKT(f.Type().Embed().Name())
		return wkt, wkt.Valid()
	}

	if (f.Type().IsRepeated() || f.Type().IsMap()) && f.Type().Element().IsEmbed() {
		wkt := pgs.LookupWKT(f.Type().Element().Embed().Name())
		return wkt, wkt.Valid()
	}

	return pgs.UnknownWKT, false
}

func isEnum(f pgs.Field) bool {
	if f.Type().IsRepeated() || f.Type().IsMap() {
		return f.Type().Element().IsEnum()
	}

	return f.Type().IsEnum()
}

func isRepeatedEnum(f pgs.Field) bool {
	return f.Type().IsRepeated() && f.Type().Element().IsEnum()
}

func isRepeatedMessage(f pgs.Field) bool {
	return f.Type().IsRepeated() && f.Type().Element().ProtoType() == pgs.MessageT
}

func isNumeric(f pgs.Field) bool {
	pType := f.Type().ProtoType()

	if f.Type().IsRepeated() || f.Type().IsMap() {
		pType = f.Type().Element().ProtoType()
	}

	switch pType {
	case pgs.DoubleT, pgs.FloatT, pgs.Int64T, pgs.UInt64T, pgs.Int32T, pgs.Fixed64T, pgs.Fixed32T:
		return true
	}

	return false
}

func isBool(f pgs.Field) bool {
	if f.Type().ProtoType() == pgs.BoolT {
		return true
	}

	if f.Type().IsRepeated() || f.Type().IsMap() {
		return f.Type().Element().ProtoType() == pgs.BoolT
	}

	return false
}

func isString(f pgs.Field) bool {
	if f.Type().ProtoType() == pgs.StringT {
		return true
	}

	if f.Type().IsRepeated() || f.Type().IsMap() {
		return f.Type().Element().ProtoType() == pgs.StringT
	}

	return false
}

func isRepeatedPrimitive(f pgs.Field) bool {
	if f.Type().IsRepeated() {
		return !f.Type().Element().IsEmbed()
	}

	return false
}

func isMapPrimitive(f pgs.Field) bool {
	if f.Type().IsMap() {
		return !f.Type().Element().IsEmbed()
	}

	return false
}

func isMessage(f pgs.Field) bool {
	return f.Type().ProtoType() == pgs.MessageT
}

func lowerCamelCaser(n pgs.Name) pgs.Name {
	return n.LowerCamelCase()
}
