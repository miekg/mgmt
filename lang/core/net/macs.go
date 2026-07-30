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

package corenet

import (
	"context"
	"net"

	"github.com/purpleidea/mgmt/lang/funcs"
	"github.com/purpleidea/mgmt/lang/interfaces"
	"github.com/purpleidea/mgmt/lang/types"
)

const (
	// MacsFuncName is the name this function is registered as.
	MacsFuncName = "macs"
)

func init() {
	funcs.ModuleRegister(ModuleName, MacsFuncName, func() interfaces.Func {
		return &MacsFunc{}
	})
}

// MacsFunc returns the visible MAC addresses and streams events when network
// links or addresses change.
type MacsFunc struct {
	interfaces.Textarea

	init *interfaces.Init
}

// String returns a simple name for this function.
func (obj *MacsFunc) String() string {
	return MacsFuncName
}

// Validate makes sure the function was built correctly.
func (obj *MacsFunc) Validate() error {
	return nil
}

// Info returns static information about this function.
func (obj *MacsFunc) Info() *interfaces.Info {
	return &interfaces.Info{
		Pure: false,
		Memo: false,
		Fast: false,
		Spec: false,
		Sig:  types.NewType("func() []str"),
	}
}

// Init initializes this function.
func (obj *MacsFunc) Init(init *interfaces.Init) error {
	obj.init = init
	return nil
}

// Stream emits an initial event and subsequent network change events.
func (obj *MacsFunc) Stream(ctx context.Context) error {
	return networkEventStream(ctx, obj.init.Event)
}

// Call returns the current visible MAC addresses.
func (obj *MacsFunc) Call(ctx context.Context, args []types.Value) (types.Value, error) {
	return Macs(ctx, args)
}

// Macs returns the list of MAC addresses that are seen on the machine.
func Macs(ctx context.Context, input []types.Value) (types.Value, error) {
	ifs, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	values := []types.Value{}
	for _, iface := range ifs {
		// Check if the interface has the loopback flag set
		if iface.Flags&net.FlagLoopback != 0 {
			continue
		}

		mac := iface.HardwareAddr.String()

		// Check if the MAC address is valid.
		if len(mac) != len("00:00:00:00:00:00") {
			continue // maybe something weird we want to ignore...
		}

		v := &types.StrValue{
			V: mac,
		}
		values = append(values, v)
	}

	return &types.ListValue{
		T: types.TypeListStr,
		V: values,
	}, nil
}
