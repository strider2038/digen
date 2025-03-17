package di

import (
	"go/ast"
	"reflect"
	"strings"
)

type ServiceDefinitionsOptions struct {
	Flags           []string
	PublicName      string
	FactoryPackage  string
	FactoryFilename string
}

func (o ServiceDefinitionsOptions) merge(o2 ServiceDefinitionsOptions) ServiceDefinitionsOptions {
	o.Flags = append(o.Flags, o2.Flags...)

	if o2.PublicName != "" {
		o.PublicName = o2.PublicName
	}
	if o2.FactoryPackage != "" {
		o.FactoryPackage = o2.FactoryPackage
	}
	if o2.FactoryFilename != "" {
		o.FactoryFilename = o2.FactoryFilename
	}

	return o
}

type OptionsParser struct {
	Logger Logger
}

func (p OptionsParser) ParseServiceDefinitionOptions(field *ast.Field) ServiceDefinitionsOptions {
	options := p.parseServiceDefinitionsOptionsFromComments(field)
	options = options.merge(p.parseServiceDefinitionsOptionsFromTags(field))

	return options
}

func (p OptionsParser) parseServiceDefinitionsOptionsFromTags(field *ast.Field) ServiceDefinitionsOptions {
	if field.Tag == nil || len(field.Tag.Value) == 0 {
		return ServiceDefinitionsOptions{}
	}

	tag := reflect.StructTag(field.Tag.Value[1 : len(field.Tag.Value)-1])

	return ServiceDefinitionsOptions{
		Flags:           split(tag.Get("di"), ","),
		PublicName:      tag.Get("public_name"),
		FactoryPackage:  tag.Get("factory_pkg"),
		FactoryFilename: tag.Get("factory_file"),
	}
}

func (p OptionsParser) parseServiceDefinitionsOptionsFromComments(field *ast.Field) ServiceDefinitionsOptions {
	options := ServiceDefinitionsOptions{}
	if field.Doc == nil {
		return options
	}

	for _, comment := range field.Doc.List {
		text := strings.TrimSpace(strings.TrimPrefix(comment.Text, "//"))
		if strings.HasPrefix(text, "di:") {
			parts := split(text, ":")
			if len(parts) == 2 {
				options.Flags = append(options.Flags, split(parts[1], ",")...)
			} else if len(parts) == 3 {
				switch parts[1] {
				case "public_name":
					options.PublicName = parts[2]
				case "factory_pkg":
					options.FactoryPackage = parts[2]
				case "factory_file":
					options.FactoryFilename = parts[2]
				default:
					p.Logger.Warning("unknown comment service definition option:", parts[1])
				}
			} else {
				p.Logger.Warning("cannot parse comment service definition option:", text)
			}
		}
	}

	return options
}

func split(s, sep string) []string {
	if s == "" {
		return nil
	}

	parts := strings.Split(s, sep)
	parsed := make([]string, 0, len(parts))

	for _, part := range parts {
		p := strings.TrimSpace(part)
		if p != "" {
			parsed = append(parsed, p)
		}
	}

	return parsed
}
