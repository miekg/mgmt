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

package traits

import (
	"github.com/purpleidea/mgmt/engine"
)

// Edgeable contains a general implementation with some of the properties and
// methods needed to support autoedges on resources. It may be used as a start
// point to avoid re-implementing the straightforward methods.
type Edgeable struct {
	// Xmeta is the stored meta. It should be called `meta` but it must be
	// public so that the `encoding/gob` package can encode it properly.
	Xmeta *engine.AutoEdgeMeta

	// Bug5819 works around issue https://github.com/golang/go/issues/5819
	Bug5819 interface{} // XXX: workaround
}

// AutoEdgeMeta lets you get or set meta params for the automatic edges trait.
func (obj *Edgeable) AutoEdgeMeta() *engine.AutoEdgeMeta {
	if obj.Xmeta == nil { // set the defaults if previously empty
		obj.Xmeta = &engine.AutoEdgeMeta{
			Disabled: false,
		}
	}
	return obj.Xmeta
}

// SetAutoEdgeMeta lets you set all of the meta params for the automatic edges
// trait in a single call.
func (obj *Edgeable) SetAutoEdgeMeta(meta *engine.AutoEdgeMeta) {
	obj.Xmeta = meta
}
