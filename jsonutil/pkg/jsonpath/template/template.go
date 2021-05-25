// Package template implements the JSONPath template format used by kubectl.
//
// See https://kubernetes.io/docs/reference/kubectl/jsonpath/ for more
// information.
//
// In some cases, it deviates slightly from the syntax accepted by kubectl's
// JSONPath expressions:
//
// - The use of \ to escape the next character in identifiers is not supported.
// - The use of @['x.y'] (equivalent to @.x.y) inside brackets is not supported,
//   as it could conflict with an actual key in a JSON object.
package template

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"text/scanner"

	"github.com/PaesslerAG/gval"
	"github.com/puppetlabs/leg/gvalutil/pkg/template"
	"github.com/puppetlabs/leg/jsonutil/pkg/jsonpath"
)

func ExpressionLanguage() gval.Language {
	return expressionLanguage
}

type TemplateOption func(tl *templateLanguage)

func WithExpressionLanguageVariableVisitor(visitor jsonpath.VariableVisitor) TemplateOption {
	return func(tl *templateLanguage) {
		tl.tpl = gval.NewLanguage(tl.tpl, gval.VariableSelector(jsonpath.VariableSelector(visitor)))
	}
}

type FormatterFunc func(value interface{}) (string, error)

func DefaultFormatter(v interface{}) (string, error) {
	switch vt := v.(type) {
	case nil:
		return "", nil
	case []interface{}:
		vs := make([]string, len(vt))
		for i, vi := range vt {
			vs[i] = fmt.Sprintf("%v", vi)
		}

		return strings.Join(vs, " "), nil
	default:
		return fmt.Sprintf("%v", vt), nil
	}
}

func WithFormatter(formatter FormatterFunc) TemplateOption {
	return func(tl *templateLanguage) {
		tl.fmt = formatter
	}
}

func TemplateLanguage(opts ...TemplateOption) gval.Language {
	return newTemplateLanguage(opts).generate()
}

func parseRange(ctx context.Context, p *gval.Parser, lang gval.Language) (gval.Evaluable, error) {
	query, err := p.ParseSublanguage(ctx, expressionLanguage)
	if err != nil {
		return nil, err
	}

	switch p.Scan() {
	case '}':
	default:
		return nil, p.Expected("JSONPath template range", '}')
	}

	sub, err := p.ParseSublanguage(ctx, lang)
	if err != nil {
		return nil, err
	}

	switch p.Scan() {
	case scanner.Ident:
		if p.TokenText() == "end" {
			break
		}

		fallthrough
	default:
		return nil, p.Expected("JSONPath template range end")
	}

	return func(ctx context.Context, parameter interface{}) (interface{}, error) {
		candidate, err := query(ctx, parameter)
		if err != nil {
			return nil, err
		}

		var s string
		if els, ok := candidate.([]interface{}); ok {
			for _, el := range els {
				v, err := sub.EvalString(ctx, el)
				if err != nil {
					return nil, err
				}

				s += v
			}
		}

		return s, nil
	}, nil
}

func eq(a, b interface{}) bool {
	// Support matrix-y == against scalar values.
	if as, ok := a.([]interface{}); ok {
		for _, av := range as {
			if reflect.DeepEqual(av, b) {
				return true
			}
		}
	} else if bs, ok := b.([]interface{}); ok {
		for _, bv := range bs {
			if reflect.DeepEqual(a, bv) {
				return true
			}
		}
	}

	return reflect.DeepEqual(a, b)
}

// expressionLanguage is the language of JSONPath expressions adjusted to
// support matrix equality comparison.
var expressionLanguage = gval.NewLanguage(
	jsonpath.Language(jsonpath.WithInitialPath{}, jsonpath.WithMissingKeysAllowed{}),
	gval.InfixOperator("==", func(a, b interface{}) (interface{}, error) { return eq(a, b), nil }),
	gval.InfixOperator("!=", func(a, b interface{}) (interface{}, error) { return !eq(a, b), nil }),
)

// templateLanguage is the total language, which includes literal handling outside of curly braces
type templateLanguage struct {
	tpl gval.Language
	fmt FormatterFunc
}

func (tl *templateLanguage) generate() gval.Language {
	return template.Language(
		template.WithJoiner{
			Joiner: template.NewStringJoiner(
				template.WithStringFormatter{
					StringFormatter: template.StringFormatterFunc(tl.fmt),
				},
			),
		},
		template.WithDelimitedLanguageFactory{
			DelimitedLanguageFactory: template.DelimitedLanguageFactoryFunc(func(outer gval.Language) *template.DelimitedLanguage {
				tpl := gval.NewLanguage(
					tl.tpl,
					gval.PrefixMetaPrefix(scanner.Ident, func(ctx context.Context, p *gval.Parser) (call string, alternative func() (gval.Evaluable, error), err error) {
						token := p.TokenText()
						return token,
							func() (gval.Evaluable, error) {
								if token == "range" {
									return parseRange(ctx, p, outer)
								}

								p.Camouflage("JSONPath template")
								return p.Const(""), nil
							},
							nil
					}),
				)

				return &template.DelimitedLanguage{
					Start: "{",
					Language: gval.Init(func(ctx context.Context, p *gval.Parser) (gval.Evaluable, error) {
						eval, err := p.ParseSublanguage(ctx, tpl)
						if err != nil {
							return nil, err
						}

						switch p.Scan() {
						case '}':
							return eval, nil
						case scanner.Ident:
							if p.TokenText() == "end" {
								template.StopFlag.MustSet(ctx)

								p.Camouflage("JSONPath template", '}')
								return eval, nil
							}

							fallthrough
						default:
							return nil, p.Expected("JSONPath template", '}')
						}
					}),
				}
			}),
		},
	)
}

func newTemplateLanguage(opts []TemplateOption) *templateLanguage {
	tl := &templateLanguage{
		tpl: expressionLanguage,
		fmt: DefaultFormatter,
	}
	for _, opt := range opts {
		opt(tl)
	}
	return tl
}
