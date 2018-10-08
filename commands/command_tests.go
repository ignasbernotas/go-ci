/*
 * Copyright (C) 2018 The "MysteriumNetwork/goci" Authors.
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package commands

import (
	"os"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"github.com/mysteriumnetwork/goci/util"
)

// Runs the test suite against the repo
func Test() error {
	mg.Deps(Deps)
	path := os.Getenv(util.MagePathOverrideEnvVar)
	if path == "" {
		path = "../..."
	}
	return sh.RunV("go", "test", "-race", "-cover", path)
}
