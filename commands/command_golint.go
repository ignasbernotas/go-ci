/*
 * Copyright (C) 2018 The "MysteriumNetwork/go-ci" Authors.
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
	"fmt"
	"path"
	"regexp"
	"strings"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"

	"github.com/zolia/go-ci/shell"
	"github.com/zolia/go-ci/util"
)

// GetLint Checks if revive exists, if not installs it
func GetLint() error {
	path, _ := util.GetGoBinaryPath("revive")
	if path != "" {
		fmt.Println("Tool 'revive' already installed")
		return nil
	}
	err := sh.RunV("go", "install", "github.com/mgechev/revive@latest")
	if err != nil {
		fmt.Printf("Could not go install revive: %s\n", err)
		return err
	}
	return nil
}

var packageRegexp = regexp.MustCompile(`\.\./(.*)\/.*\.go`)

func getPackageFromGoLintOutput(line string) string {
	results := packageRegexp.FindAllStringSubmatch(line, -1)
	for i := range results {
		return results[i][1]
	}
	return ""
}

func formatAndPrintGoLintOutput(rawGolint string) {
	packageErrorMap := make(map[string][]string)
	separateLines := strings.Split(rawGolint, "\n")

	for i := range separateLines {
		pkg := getPackageFromGoLintOutput(separateLines[i])
		if val, ok := packageErrorMap[pkg]; ok {
			packageErrorMap[pkg] = append(val, separateLines[i])
		} else {
			lines := []string{separateLines[i]}
			packageErrorMap[pkg] = lines
		}
	}

	fmt.Println()
	for k := range packageErrorMap {
		fmt.Println("PACKAGE: ", k)
		fmt.Println()
		for _, v := range packageErrorMap[k] {
			fmt.Println(v)
		}
		fmt.Println()
	}
}

// GoLint checks for linting errors in the solution
func GoLint(pathToCheck string, excludes ...string) error {
	mg.Deps(GetLint)
	golintPath, err := util.GetGoBinaryPath("revive")
	if err != nil {
		return err
	}

	gopath := util.GetGoPath()
	dirs, err := util.GetPackagePathsWithExcludes(pathToCheck, excludes...)
	if err != nil {
		fmt.Println("go list crashed")
		return err
	}

	var files []string

	for _, dir := range dirs {
		absolutePath := path.Join(gopath, "src", dir)
		files = append(files, absolutePath)
	}

	args := []string{"--set_exit_status"}
	args = append(args, files...)
	output, err := sh.Output(golintPath, args...)
	exitStatus := sh.ExitStatus(err)
	if exitStatus == 0 {
		fmt.Println("No linting errors")
		return nil
	}

	formatAndPrintGoLintOutput(output)
	fmt.Printf("Linting failed: %s\n", err)
	return err
}

// GoLintD checks for linting errors in the solution
//
// Instead of packages, it operates on directories, thus it is compatible with gomodules outside GOPATH.
//
// Example:
//     commands.GoLintD(".", "docs")
func GoLintD(dir string, excludes ...string) error {
	mg.Deps(GetLint)
	reviveBin, err := util.GetGoBinaryPath("revive")
	if err != nil {
		return err
	}

	var allExcludes []string
	allExcludes = append(allExcludes, excludes...)
	allExcludes = append(allExcludes, util.GoLintExcludes()...)
	dirs, err := util.GetProjectFileDirectories(allExcludes)
	if err != nil {
		fmt.Printf("revive: go list crashed: %s\n", err)
		return err
	}

	output, err := shell.NewCmd(reviveBin + " --set_exit_status " + strings.Join(dirs, " ")).Output()
	exitStatus := sh.ExitStatus(err)
	if exitStatus != 0 {
		formatAndPrintGoLintOutput(output)
		fmt.Printf("revive: linting failed: %s\n", err)
		return err
	}

	fmt.Println("revive: no linting errors")
	return nil
}
