package json

import (
	"fmt"
	"path/filepath"
	"strings"
	"text/template"

	pgs "github.com/lyft/protoc-gen-star"
)

const ImportsInsertionPoint = "plugin_imports"

var _ pgs.Module = (*module)(nil)

type module struct {
	*pgs.ModuleBase
	messageTemplate *template.Template
	importsTemplate *template.Template
	clientTemplate  *template.Template
	mapMessages     map[string]mapMessages
	fileImports     map[string]typeImports
}

type mapMessages map[string]*mapMessage

type mapMessage struct {
	Name      pgs.Name
	Type      pgs.Name
	IsMessage bool
}

type typeImports map[string]*typeImport

type typeImport struct {
	Name pgs.Name
	Path pgs.FilePath
}

func Module() *module {
	return &module{
		ModuleBase:  &pgs.ModuleBase{},
		mapMessages: make(map[string]mapMessages),
		fileImports: make(map[string]typeImports),
	}
}

func (m *module) InitContext(c pgs.BuildContext) {
	m.ModuleBase.InitContext(c)

	m.messageTemplate = template.Must(
		template.New("message").Funcs(map[string]interface{}{
			"EntityNamer":     m.entityNamer,
			"FieldNamer":      fieldNamer,
			"FieldTyper":      m.fieldTyper,
			"JSONFieldNamer":  jsonFieldNamer,
			"JSONFieldTyper":  m.jsonFieldTyper,
			"Caster":          m.caster,
			"JSONCaster":      m.jsonCaster,
			"IsMessage":       isMessage,
			"LowerCamelCaser": lowerCamelCaser,
		}).Parse(messageTemplate),
	)

	m.importsTemplate = template.Must(
		template.New("imports").Funcs(map[string]interface{}{}).Parse(importsTemplate),
	)

	m.clientTemplate = template.Must(
		template.New("client").Funcs(map[string]interface{}{
			"EntityNamer":     m.entityNamer,
			"LowerCamelCaser": lowerCamelCaser,
		}).Parse(clientTemplate),
	)
}

func (m *module) Name() string {
	return "json"
}

func (m *module) addMapMessage(fName pgs.Name, f pgs.Field) {
	var (
		ttype     pgs.Name
		isMessage bool
	)

	wktT, wellKnown := isWellKnownType(f)
	switch {
	case wellKnown:
		ttype = pgs.Name(m.wkTyper(wktT))
	case isEnum(f):
		ttype = pgs.Name(m.entityNamer(f.Type().Element().Embed()))
	case isNumeric(f):
		ttype = pgs.Name("number")
	case isBool(f):
		ttype = pgs.Name("boolean")
	case isString(f):
		ttype = pgs.Name("string")
	default:
		ttype = pgs.Name(m.entityNamer(f.Type().Element().Embed()))
		isMessage = true
	}

	inputPath := f.File().InputPath().String()

	if m.mapMessages[inputPath] == nil {
		m.mapMessages[inputPath] = make(mapMessages)
	}

	m.mapMessages[inputPath][fName.String()] = &mapMessage{
		Name:      fName,
		Type:      ttype,
		IsMessage: isMessage,
	}
}

func (m *module) addTypeImport(f pgs.Field) {
	if len(f.Imports()) != 1 {
		m.Failf("field %s should have 1 import (%d)", f.Name(), len(f.Imports()))
	}

	inputPath := f.File().InputPath()

	rel, err := filepath.Rel(inputPath.Dir().String(), f.Imports()[0].File().InputPath().String())
	if err != nil {
		m.Failf("type import relative path: %s", err)
	}

	fp := pgs.FilePath(rel)
	fp = pgs.FilePath(fmt.Sprintf("%s/%s", fp.Dir(), fp.SetExt(".pb").Base()))

	if m.fileImports[inputPath.String()] == nil {
		m.fileImports[inputPath.String()] = make(typeImports)
	}

	var ttype pgs.Name

	_, wellKnown := isWellKnownType(f)
	switch {
	case wellKnown:
		return
	case isRepeatedMessage(f), f.Type().IsMap():
		ttype = pgs.Name(m.entityNamer(f.Type().Element().Embed()))
	default:
		ttype = pgs.Name(m.entityNamer(f.Type().Embed()))
	}

	if m.fileImports[inputPath.String()] == nil {
		m.fileImports[inputPath.String()] = make(typeImports)
	}

	m.fileImports[inputPath.String()][ttype.String()] = &typeImport{
		Name: ttype,
		Path: fp,
	}
}

func (m *module) Execute(targets map[string]pgs.File, packages map[string]pgs.Package) []pgs.Artifact {
	for _, p := range packages {
		if isWellKnownPackage(p.ProtoName()) {
			continue
		}

		for _, f := range p.Files() {
			m.AddGeneratorTemplateFile(
				f.InputPath().SetExt(".pb.ts").String(),
				m.messageTemplate,
				map[string]interface{}{
					"InputPath":   f.InputPath().String(),
					"Messages":    f.AllMessages(),
					"MapMessages": m.mapMessages,
					"Enums":       f.AllEnums(),
				},
			)

			// generate twirp related things only if services are present
			if f.BuildTarget() && len(f.Services()) > 0 {
				m.AddGeneratorTemplateFile(
					f.InputPath().SetExt(".client.ts").String(),
					m.clientTemplate,
					map[string]interface{}{
						"Services": f.Services(),
					},
				)
				m.AddGeneratorFile(f.InputPath().Dir().Push("twirp.ts").String(), twirpTemplate)
			}
		}
	}

	for _, p := range packages {
		if isWellKnownPackage(p.ProtoName()) {
			continue
		}

		for _, f := range p.Files() {
			m.AddGeneratorTemplateInjection(
				f.InputPath().SetExt(".pb.ts").String(),
				ImportsInsertionPoint, m.importsTemplate,
				map[string]interface{}{
					"InputPath":   f.InputPath().String(),
					"FileImports": m.fileImports,
				},
			)
		}
	}

	return m.Artifacts()
}

func (m *module) entityNamer(e pgs.Entity) pgs.Name {
	name := strings.TrimPrefix(e.FullyQualifiedName(), fmt.Sprintf(".%s", e.Package().ProtoName()))
	return pgs.Name(name).UpperCamelCase()
}

func fieldNamer(f pgs.Field) pgs.Name {
	return pgs.Name(f.Name().LowerCamelCase())
}

var wellKnownTypeMap = map[pgs.WellKnownType]pgs.Name{
	pgs.TimestampWKT:   pgs.Name("string"),
	pgs.DoubleValueWKT: pgs.Name("number"),
	pgs.FloatValueWKT:  pgs.Name("number"),
	pgs.Int64ValueWKT:  pgs.Name("number"),
	pgs.UInt64ValueWKT: pgs.Name("number"),
	pgs.Int32ValueWKT:  pgs.Name("number"),
	pgs.UInt32ValueWKT: pgs.Name("number"),
	pgs.BoolValueWKT:   pgs.Name("bool"),
	pgs.StringValueWKT: pgs.Name("string"),
	pgs.BytesValueWKT:  pgs.Name("string"),
}

func (m *module) wkTyper(wkType pgs.WellKnownType) pgs.Name {
	if m, ok := wellKnownTypeMap[wkType]; ok {
		return m
	}

	m.Failf("unknown well known type: %s", wkType.Name())
	return pgs.Name("")
}

func (m *module) fieldTyper(f pgs.Field) pgs.Name {
	wktT, wellKnown := isWellKnownType(f)

	if len(f.Imports()) == 1 {
		m.addTypeImport(f)
	}

	// NB. the ordering of these cases matter
	switch {
	case wellKnown:
		return pgs.Name(m.wkTyper(wktT))
	case isRepeatedEnum(f):
		return pgs.Name(m.entityNamer(f.Type().Element().Enum()))
	case isEnum(f):
		return pgs.Name(m.entityNamer(f.Type().Enum()))
	case f.Type().IsMap():
		typeName := pgs.Name(fmt.Sprintf("%s%sEntry", m.entityNamer(f.Message()), f.Name().UpperCamelCase()))
		m.addMapMessage(typeName, f)
		return typeName
	case isRepeatedMessage(f):
		return pgs.Name(m.entityNamer(f.Type().Element().Embed()))
	case isNumeric(f):
		return pgs.Name("number")
	case isBool(f):
		return pgs.Name("boolean")
	case isString(f):
		return pgs.Name("string")
	default: // it has to be a message
		if f.Type().IsEmbed() { // but we're still checking just in case
			return pgs.Name(m.entityNamer(f.Type().Embed()))
		}
	}

	m.Failf("unknown field type: %s", f.FullyQualifiedName())
	return pgs.Name("")
}

func (m *module) caster(f pgs.Field) pgs.Name {
	jsonFieldName := jsonFieldNamer(f)
	_, wellKnown := isWellKnownType(f)

	switch {
	case isRepeatedPrimitive(f):
		return pgs.Name(fmt.Sprintf("m.%s as %s[]", jsonFieldName, m.fieldTyper(f)))
	case wellKnown, isMapPrimitive(f):
		return pgs.Name(fmt.Sprintf("m.%s", jsonFieldName))
	case isRepeatedMessage(f):
		return pgs.Name(fmt.Sprintf("m.%s && m.%s.map(JSONTo%s)", jsonFieldName, jsonFieldName, m.fieldTyper(f)))
	case isEnum(f):
		return pgs.Name(fmt.Sprintf("m.%s as %s", jsonFieldName, m.fieldTyper(f)))
	case isNumeric(f):
		return pgs.Name(fmt.Sprintf("m.%s || 0", jsonFieldName))
	case isBool(f):
		return pgs.Name(fmt.Sprintf("m.%s || false", jsonFieldName))
	case isString(f):
		return pgs.Name(fmt.Sprintf("m.%s || \"\"", jsonFieldName))
	}

	return pgs.Name(fmt.Sprintf("m.%s && JSONTo%s(m.%s)", jsonFieldName, m.fieldTyper(f), jsonFieldName))
}

func jsonFieldNamer(f pgs.Field) pgs.Name {
	return pgs.Name(f.Name().LowerSnakeCase())
}

func (m *module) jsonFieldTyper(f pgs.Field) pgs.Name {
	fieldType := m.fieldTyper(f)

	switch {
	case fieldType == "number", fieldType == "string", fieldType == "boolean":
		return fieldType
	case isEnum(f):
		return "string"
	}

	return pgs.Name(fmt.Sprintf("%sJSON", fieldType))
}

func (m *module) jsonCaster(f pgs.Field) pgs.Name {
	fieldName := fieldNamer(f)
	_, wellKnown := isWellKnownType(f)

	switch {
	case wellKnown, isRepeatedPrimitive(f), isMapPrimitive(f):
		return pgs.Name(fmt.Sprintf("m.%s", fieldName))
	case isRepeatedMessage(f):
		return pgs.Name(fmt.Sprintf("m.%s && m.%s.map(%sToJSON)", fieldName, fieldName, m.fieldTyper(f)))
	case isEnum(f):
		return pgs.Name(fmt.Sprintf("m.%s as %s", fieldName, m.fieldTyper(f)))
	case isNumeric(f), isBool(f), isString(f):
		return pgs.Name(fmt.Sprintf("m.%s", fieldName))
	}

	return pgs.Name(fmt.Sprintf("m.%s && %sToJSON(m.%s)", fieldName, m.fieldTyper(f), fieldName))
}
