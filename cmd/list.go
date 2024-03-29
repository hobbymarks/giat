/*
Copyright © 2022 hobbymarks ihobbymarks@gmail.com
*/
package cmd

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all git managed directory",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if results, err := AllGitDirs(args); err != nil {
			log.Error(err)
		} else {
			for _, repo := range results {
				fmt.Println(repo)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}

func AllGitDirs(rootPaths []string) ([]string, error) {
	var roots []string
	var results []string

	if len(rootPaths) >= 1 {
		for _, arg := range rootPaths {
			if _, err := os.Stat(arg); err != nil {
				log.Warning(err)
			} else {
				roots = append(roots, arg)
			}
		}
		//FIXME:should merge paths
	} else {
		roots = []string{"."}
	}

	for _, root := range roots {
		if dirs, err := Dirs(root); err != nil {
			log.Error(err)
		} else {
			if gitDirs, err := GitDirs(dirs); err != nil {
				log.Error(err)
			} else {
				log.Trace(gitDirs)
				results = append(results, gitDirs...)
			}
		}
	}
	//TODO:support colorfull output
	//TODO:supoort output grouped,such as by directory or by arg ...
	sort.Slice(results, func(i, j int) bool {
		return results[i] > results[j]
	})
	return results, nil
}

// filter gitdir from rootDirs
func GitDirs(rootDirs []string) ([]string, error) {
	var gitDirs []string

	for _, d := range rootDirs {
		if ok, err := IsGitDir(d); err != nil {
			log.Error(err)
		} else {
			if ok {
				if abs, err := filepath.Abs(d); err != nil {
					log.Error(err)
				} else {
					gitDirs = append(gitDirs, abs)
				}
			}
		}
	}

	return gitDirs, nil
}

// return all dirs path in rootPath
func Dirs(rootPath string) ([]string, error) {
	var dirs []string

	rootPath = filepath.Clean(rootPath)
	log.Trace(rootPath)
	err := filepath.WalkDir(rootPath, func(path string, info fs.DirEntry, err error) error {
		if err != nil {
			log.Error(err)
			return err
			//FIXME:should check error type,some error should ignore,such as permission ...
		}
		if info.IsDir() {
			log.Trace("IsDir:", path)
			rel := path
			dirs = append(dirs, filepath.ToSlash(rel))
			return nil
		} else {
			log.Trace("Skipped:", path)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return dirs, nil
}
