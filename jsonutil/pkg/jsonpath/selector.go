package jsonpath

import (
	"context"

	"github.com/PaesslerAG/gval"
)

type selector func(context.Context, any) (values, error)

type currentElementContextKey struct{}

// $
func rootElement(c context.Context, v any) (any, error) {
	return v, nil
}

// @
func currentElement(c context.Context, _ any) (any, error) {
	return c.Value(currentElementContextKey{}), nil
}

func currentContext(c context.Context, v any) context.Context {
	return context.WithValue(c, currentElementContextKey{}, v)
}

func varSelector(variable gval.Evaluable) selector {
	return func(c context.Context, v any) (values, error) {
		r, err := variable(currentContext(c, v), v)
		if err != nil {
			return nil, err
		} else if r == nil {
			return empty{}, nil
		}

		vs, ok := r.(values)
		if !ok {
			return nil, &VarSelectorTypeError{Variable: r}
		}

		return vs, nil
	}
}

// .
func directSelector(sel selector) selector {
	// This selector just drops the wildcard from the resulting values since
	// it's known to be a direct path.
	return func(c context.Context, v any) (values, error) {
		vs, err := sel(c, v)
		if err != nil {
			return nil, err
		}

		return vs.flatMap(func(v value) (values, error) {
			v.wildcards = nil
			return v, nil
		})
	}
}

// [x, ...]
func multiSelector(sels []selector) selector {
	return func(c context.Context, v any) (values, error) {
		var vs values = valueSlice{}
		for _, sel := range sels {
			r, err := sel(c, v)
			if err != nil {
				// An individual selector failing in a multi-select should not
				// propagate (only if the number of selectors > 1).
				continue
			}

			r.concat(&vs)
		}
		return vs, nil
	}
}

// [? ]
func filterSelector(sel selector, filter gval.Evaluable) selector {
	return func(c context.Context, v any) (values, error) {
		r, err := sel(c, v)
		if err != nil {
			return nil, err
		}

		return r.flatMap(func(v value) (values, error) {
			cond, err := filter.EvalBool(currentContext(c, v.value), v)
			if err != nil || !cond {
				// Again, errors don't propagate here.
				return nil, nil
			}

			return v, nil
		})
	}
}

// (
func scriptSelector(script gval.Evaluable) selector {
	return func(c context.Context, v any) (values, error) {
		v, err := script(currentContext(c, v), v)
		if err != nil {
			return nil, err
		}

		return value{value: v}, nil
	}
}
