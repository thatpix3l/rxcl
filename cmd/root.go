package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
    shortAbout = "Random xkcd comic link generator"
    longAbout = `A needlessly overcomplicated program for
generating an xkcd comic or image link between either the default
range or a custom range if provided.`

    rootCmd = &cobra.Command{
        Use: "rxcl",
        Short: shortAbout,
        Long: longAbout,
    }
)

func Execute() {

    if err := rootCmd.Execute(); err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }

}