package cmd

import (
	"fmt"

	"github.com/aymenhmaidiwastaken/devboot/internal/state"
	"github.com/aymenhmaidiwastaken/devboot/internal/ui"
	"github.com/spf13/cobra"
)

var historyLimit int

var historyCmd = &cobra.Command{
	Use:   "history",
	Short: "Show history of devboot actions",
	RunE:  runHistory,
}

func init() {
	historyCmd.Flags().IntVarP(&historyLimit, "limit", "n", 50, "number of entries to show")
	rootCmd.AddCommand(historyCmd)
}

func runHistory(cmd *cobra.Command, args []string) error {
	st, err := state.Load()
	if err != nil {
		return fmt.Errorf("loading state: %w", err)
	}

	if len(st.Actions) == 0 {
		ui.Info("no history yet — run: devboot apply")
		return nil
	}

	fmt.Printf("\n  devboot history\n")

	if !st.LastApply.IsZero() {
		ui.Info(fmt.Sprintf("last apply: %s", st.LastApply.Format("2006-01-02 15:04:05")))
	}
	if len(st.ManagedTools) > 0 {
		ui.Info(fmt.Sprintf("managed tools: %d", len(st.ManagedTools)))
	}

	actions := st.History(historyLimit)

	currentSection := ""
	for _, a := range actions {
		sectionLabel := fmt.Sprintf("%s/%s", a.Section, a.Type)
		if sectionLabel != currentSection {
			ui.Section(sectionLabel)
			currentSection = sectionLabel
		}

		timestamp := a.Timestamp.Format("01/02 15:04")
		switch a.Status {
		case "ok":
			ui.Success(fmt.Sprintf("[%s] %s", timestamp, a.Target))
		case "skipped":
			ui.Skip(fmt.Sprintf("[%s] %s", timestamp, a.Target))
		case "failed":
			ui.Fail(fmt.Sprintf("[%s] %s — %s", timestamp, a.Target, a.Detail))
		default:
			ui.Info(fmt.Sprintf("[%s] %s", timestamp, a.Target))
		}
	}

	fmt.Println()
	return nil
}
