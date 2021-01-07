package hcl

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/zclconf/go-cty/cty"
)

type SpecDescriptor interface {
	Spec() hcldec.Spec
}

func Unmarshal(data []byte, target interface{}) error {
	start := hcl.Pos{Line: 1, Column: 1}

	file, err := hclsyntax.ParseConfig(data, "config.hcl", start)
	if err != nil {
		return err
	}

	var ctx *hcl.EvalContext

	if descriptor, ok := target.(SpecDescriptor); ok {
		value, _, _ := hcldec.PartialDecode(file.Body, descriptor.Spec(), nil)

		ctx = &hcl.EvalContext{
			Variables: NewVariables(value),
		}
	}

	if err := gohcl.DecodeBody(file.Body, ctx, target); err != nil {
		return err
	}

	return nil
}

func NewVariables(value cty.Value) map[string]cty.Value {
	var (
		kv   = make(map[string]cty.Value)
		kind = value.Type()
	)

	switch {
	case kind.IsObjectType():
		for it := value.ElementIterator(); it.Next(); {
			k, v := it.Element()
			props := NewVariables(v)
			kv[k.AsString()] = cty.ObjectVal(props)
		}
	case kind.IsListType():
		for it := value.ElementIterator(); it.Next(); {
			_, v := it.Element()

			if !v.Type().HasAttribute("name") {
				continue
			}

			k := v.GetAttr("name")
			kv[k.AsString()] = v
		}
	}

	return kv
}

