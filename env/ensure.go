/*
 * Copyright (C) 2019 The "MysteriumNetwork/go-ci" Authors.
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

package env

import (
	"os"

	"github.com/pkg/errors"
)

// EnsureEnvVars ensures that specified environment variables have a value.
// Purpose is to reduce boilerplate in mage target when reading multiple env vars.
func EnsureEnvVars(vars ...BuildVar) error {
	missingVars := make([]BuildVar, 0)
	for _, v := range vars {
		if os.Getenv(string(v)) == "" {
			missingVars = append(missingVars, v)
		}
	}
	if len(missingVars) > 0 {
		return errors.Errorf("the following environment variables are missing: %v", missingVars)
	}
	return nil
}
