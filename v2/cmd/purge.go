package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var purgeCmd = &cobra.Command{
	Use:     "purge",
	Args:    cobra.NoArgs,
	Short:   "Purge all of chezmoi's configuration and data",
	Long:    mustGetLongHelp("purge"),
	Example: getExample("purge"),
	RunE:    config.runPurgeCmd,
}

func init() {
	rootCmd.AddCommand(purgeCmd)
}

func (c *Config) runPurgeCmd(cmd *cobra.Command, args []string) error {
	// Build a list of chezmoi-related paths.
	var paths []string
	for _, dirs := range [][]string{
		c.bds.ConfigDirs,
		c.bds.DataDirs,
	} {
		for _, dir := range dirs {
			paths = append(paths, filepath.Join(dir, "chezmoi"))
		}
	}
	paths = append(paths, c.configFile, c.getPersistentStateFile())

	// Remove all paths that exist.
PATH:
	for _, path := range paths {
		_, err := c.fs.Stat(path)
		switch {
		case os.IsNotExist(err):
			continue PATH
		case err != nil:
			return err
		}
		if !c.force {
			choice, err := c.prompt(fmt.Sprintf("Remove %s", path), "ynqa")
			if err != nil {
				return err
			}
			switch choice {
			case 'a':
				c.force = true
			case 'n':
				continue PATH
			case 'q':
				return nil
			}
		}
		if err := c.system.RemoveAll(path); err != nil {
			return err
		}
	}

	return nil
}