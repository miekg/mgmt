// Mgmt
// Copyright (C) James Shubin and the project contributors
// Written by James Shubin <james@shubin.ca> and the project contributors
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.
//
// Additional permission under GNU GPL version 3 section 7
//
// If you modify this program, or any covered work, by linking or combining it
// with embedded mcl code and modules (and that the embedded mcl code and
// modules which link with this program, contain a copy of their source code in
// the authoritative form) containing parts covered by the terms of any other
// license, the licensors of this program grant you additional permission to
// convey the resulting work. Furthermore, the licensors of this program grant
// the original author, James Shubin, additional permission to update this
// additional permission if he deems it necessary to achieve the goals of this
// additional permission.

package coreiter

import (
	"context"
	"fmt"

	"github.com/purpleidea/mgmt/lang/funcs"
	"github.com/purpleidea/mgmt/lang/interfaces"
	"github.com/purpleidea/mgmt/lang/types"
	"github.com/purpleidea/mgmt/util/errwrap"
)

const (
	// FlattenFuncName is the name this function is registered as.
	FlattenFuncName = "flatten"

	// arg names...
	flattenArgNameInputs = "inputs"
)

func init() {
	funcs.ModuleRegister(ModuleName, FlattenFuncName, func() interfaces.Func { return &FlattenFunc{} }) // must register the func and name
}

var _ interfaces.BuildableFunc = &FlattenFunc{} // ensure it meets this expectation

// FlattenFunc is a function that takes a list of lists and concatenates them
// into a single list, in order. Unlike list.concat, which takes a fixed number
// of list arguments specified at compile time, this works on a dynamic outer
// list whose length can vary at runtime. It composes naturally with iter.map to
// produce a flatMap-style operation.
type FlattenFunc struct {
	interfaces.Textarea

	Type *types.Type // element type of the inner list (and of the result)

	init *interfaces.Init
}

// String returns a simple name for this function. This is needed so this struct
// can satisfy the pgraph.Vertex interface.
func (obj *FlattenFunc) String() string {
	return FlattenFuncName
}

// ArgGen returns the Nth arg name for this function.
func (obj *FlattenFunc) ArgGen(index int) (string, error) {
	seq := []string{flattenArgNameInputs}
	if l := len(seq); index >= l {
		return "", fmt.Errorf("index %d exceeds arg length of %d", index, l)
	}
	return seq[index], nil
}

// helper
func (obj *FlattenFunc) sig() *types.Type {
	// func(inputs [][]?1) []?1
	t := "?1"
	if obj.Type != nil {
		t = obj.Type.String()
	}
	return types.NewType(fmt.Sprintf("func(%s [][]%s) []%s", flattenArgNameInputs, t, t))
}

// Build is run to turn the polymorphic, undetermined function, into the
// specific statically typed version. It is usually run after Unify completes,
// and must be run before Info() and any of the other Func interface methods are
// used. This function is idempotent, as long as the arg isn't changed between
// runs.
func (obj *FlattenFunc) Build(typ *types.Type) (*types.Type, error) {
	if typ.Kind != types.KindFunc {
		return nil, fmt.Errorf("input type must be of kind func")
	}

	if len(typ.Ord) != 1 {
		return nil, fmt.Errorf("the flatten function needs exactly one arg")
	}
	if typ.Map == nil {
		return nil, fmt.Errorf("the map is nil")
	}

	tInputs, exists := typ.Map[typ.Ord[0]]
	if !exists || tInputs == nil {
		return nil, fmt.Errorf("first argument was missing")
	}
	if tInputs.Kind != types.KindList {
		return nil, fmt.Errorf("first argument must be of kind list")
	}
	if tInputs.Val == nil || tInputs.Val.Kind != types.KindList {
		return nil, fmt.Errorf("first argument must be a list of lists")
	}

	if typ.Out == nil {
		return nil, fmt.Errorf("return type must be specified")
	}
	if typ.Out.Kind != types.KindList {
		return nil, fmt.Errorf("return type must be a list")
	}
	if err := tInputs.Val.Val.Cmp(typ.Out.Val); err != nil {
		return nil, errwrap.Wrapf(err, "the inner list element type must match the returned list element type")
	}

	obj.Type = typ.Out.Val // element type

	return obj.sig(), nil
}

// Validate tells us if the input struct takes a valid form.
func (obj *FlattenFunc) Validate() error {
	if obj.Type == nil {
		return fmt.Errorf("type is not yet known")
	}
	return nil
}

// Info returns some static info about itself. Build must be called before this
// will return correct data.
func (obj *FlattenFunc) Info() *interfaces.Info {
	return &interfaces.Info{
		Pure: true,
		Memo: true,
		Fast: true,
		Spec: true,
		Sig:  obj.sig(),
		Err:  obj.Validate(),
	}
}

// Init runs some startup code for this function.
func (obj *FlattenFunc) Init(init *interfaces.Init) error {
	obj.init = init
	return nil
}

// Call returns the result of this function.
func (obj *FlattenFunc) Call(ctx context.Context, args []types.Value) (types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("expected exactly one arg")
	}

	outer, ok := args[0].(*types.ListValue)
	if !ok {
		return nil, fmt.Errorf("expected a ListValue argument")
	}

	out := []types.Value{}
	for _, v := range outer.List() {
		inner, ok := v.(*types.ListValue)
		if !ok {
			return nil, fmt.Errorf("expected each inner element to be a ListValue")
		}
		out = append(out, inner.List()...)
	}

	return &types.ListValue{
		T: types.NewType(fmt.Sprintf("[]%s", obj.Type)),
		V: out,
	}, nil
}

// Copy is implemented so that the type value is not lost if we copy this
// function.
func (obj *FlattenFunc) Copy() interfaces.Func {
	return &FlattenFunc{
		Textarea: obj.Textarea,

		Type: obj.Type, // don't copy because we use this after unification

		init: obj.init, // likely gets overwritten anyways
	}
}
