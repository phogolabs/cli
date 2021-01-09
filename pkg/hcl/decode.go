package hcl

import (
	"fmt"
	"strconv"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/zclconf/go-cty/cty"
)

// SpecDescriptor represents the type spec
type SpecDescriptor interface {
	Spec() hcldec.Spec
}

// Unmarshal unmarshal the data for given target
func Unmarshal(data []byte, target interface{}) error {
	var (
		ctx   *hcl.EvalContext
		start = hcl.Pos{Line: 1, Column: 1}
	)

	file, derr := hclsyntax.ParseConfig(data, "config.hcl", start)
	if derr != nil {
		return derr
	}

	if descriptor, ok := target.(SpecDescriptor); ok {
		root, _, _ := hcldec.PartialDecode(file.Body, descriptor.Spec(), nil)
		tree := make(map[string]interface{})

		err := cty.Walk(root, func(path cty.Path, value cty.Value) (bool, error) {
			switch {
			case value.Type().IsMapType():
				fallthrough
			case value.Type().IsPrimitiveType():
				vc := vector(path, root)
				name, kv := leaf(vc, tree)
				kv[name] = value
			}
			return true, nil
		})

		if err != nil {
			return err
		}

		ctx = &hcl.EvalContext{
			Variables: compile(tree).AsValueMap(),
		}
	}

	if err := gohcl.DecodeBody(file.Body, ctx, target); err != nil {
		return err
	}

	return nil
}

func compile(kv map[string]interface{}) cty.Value {
	props := make(map[string]cty.Value, len(kv))
	items := make([]cty.Value, len(kv))

	for k, v := range kv {
		index, err := strconv.Atoi(k)
		yep := err == nil

		if next, ok := v.(map[string]interface{}); ok {
			value := compile(next)

			if yep {
				items[index] = value
			} else {
				props[k] = value
			}
		}

		if value, ok := v.(cty.Value); ok {
			if yep {
				items[index] = value
			} else {
				props[k] = value
			}
		}
	}

	if len(props) > 0 {
		return cty.ObjectVal(props)
	}

	return cty.ListVal(items)

}

func leaf(path []string, tree map[string]interface{}) (string, map[string]interface{}) {
	kv := tree

	for index, key := range path {
		if index == len(path)-1 {
			return key, kv
		}

		next, ok := kv[key].(map[string]interface{})
		if !ok {
			next = make(map[string]interface{})
			kv[key] = next
		}

		kv = next
	}

	return "", kv
}

func vector(path cty.Path, root cty.Value) []string {
	vc := []string{}

	for index, step := range path {
		switch ptr := step.(type) {
		case cty.GetAttrStep:
			vc = append(vc, ptr.Name)
		case cty.IndexStep:
			key := ptr.Key

			switch key.Type() {
			case cty.String:
				vc = append(vc, key.AsString())
			case cty.Number:
				var (
					head      = path[:index+1]
					parent, _ = head.Apply(root)
				)

				if kind := parent.Type(); kind.IsObjectType() {
					if kind.HasAttribute("name") {
						vc = append(vc, parent.GetAttr("name").AsString())
						continue
					}
				}

				kindex, _ := key.AsBigFloat().Int64()
				vc = append(vc, fmt.Sprintf("%v", kindex))
			}
		}
	}

	return vc
}
