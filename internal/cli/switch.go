package cli

import (
	"fmt"

	"github.com/FeLvi-zzz/tentez"
	"github.com/spf13/cobra"
)

var switchCmd = &cobra.Command{
	Use:   "switch",
	Short: "Switch targets weights",
	Long: `
# show plan and switch
$ tentez -f ./examples/example.yaml switch --weight 30,70`,
	RunE: func(cmd *cobra.Command, args []string) error {
		t, err := tentez.NewFromYaml(filename)
		if err != nil {
			return err
		}

		if len(weights) != 2 {
			return fmt.Errorf("the length of weights must be 2, got %v, len: %d", weights, len(weights))
		}

		return t.Switch(weights, !noPause)
	},
}

var weights []int

func init() {
	switchCmd.Flags().StringVarP(&filename, "filename", "f", "", "config file for tentez")
	switchCmd.Flags().IntSliceVarP(&weights, "weights", "w", []int{}, "weights for switch, 'old,new'")
	switchCmd.Flags().BoolVar(&noPause, "no-pause", false, "skip pause")

	if err := switchCmd.MarkFlagRequired("filename"); err != nil {
		panic(err)
	}

	rootCmd.AddCommand(switchCmd)
}
