package tplinator

import "io"

func Tplinate(tplReader io.Reader, parserOptions ...ParserOptionFunc) (*Template, error) {
	defaultParserOptions := []ParserOptionFunc{
		NodeProcessorsParserOption(
			ConditionalExtensionNodeProcessor,
			ConditionalClassExtensionNodeProcessor,
			RangeExtensionNodeProcessor,
			StringInterpolationNodeProcessor,
		),
	}
	defaultParserOptions = append(defaultParserOptions, parserOptions...)

	template, err := CreateTemplateFromReader(tplReader, defaultParserOptions...)
	if err != nil {
		return nil, err
	}
	template.extDeps = compoundExtensionDependencies{
		defaultExtDep: NewDefaultExtensionDependencies(),
	}

	return template, nil
}
