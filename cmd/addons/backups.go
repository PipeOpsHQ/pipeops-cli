package addons

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/PipeOpsHQ/pipeops-cli/utils"
	sdk "github.com/PipeOpsHQ/pipeops-go-sdk/pipeops"
	"github.com/spf13/cobra"
)

var backupsCmd = &cobra.Command{
	Use:     "backups",
	Aliases: []string{"backup"},
	Short:   "Manage addon deployment backups",
	Long: `List and export backup snapshots for addon deployments.

Examples:
  pipeops addons backups list <deployment-uid>
  pipeops addons backups export <deployment-uid> --snapshot-id <id>
  pipeops addons backups export-status <deployment-uid> <export-id>`,
}

var backupsListCmd = &cobra.Command{
	Use:     "list <deployment-uid>",
	Aliases: []string{"ls"},
	Short:   "List backup snapshots for an addon deployment",
	RunE: func(cmd *cobra.Command, args []string) error {
		opts := utils.GetOutputOptions(cmd)
		client, err := addonsClient(cmd, opts)
		if err != nil || client == nil {
			return err
		}
		resp, err := client.ListAddonBackups(context.Background(), args[0])
		if err != nil {
			return fmt.Errorf("list addon backups: %w", err)
		}
		if opts.Format == utils.OutputFormatJSON {
			return utils.PrintJSON(resp)
		}
		snapshots := resp.Data.Snapshots
		if len(snapshots) == 0 {
			utils.PrintWarning("No backup snapshots found", opts)
			return nil
		}
		rows := make([][]string, 0, len(snapshots))
		for _, s := range snapshots {
			size := ""
			if s.TotalSizeBytes > 0 {
				size = formatBytes(s.TotalSizeBytes)
			} else if s.SizeUnknown {
				size = "unknown"
			}
			rows = append(rows, []string{
				s.ID,
				s.Name,
				s.Time,
				size,
				s.TypeChip,
				boolYesNo(s.Useful),
			})
		}
		utils.PrintTable([]string{"SNAPSHOT ID", "NAME", "TIME", "SIZE", "TYPE", "USEFUL"}, rows, opts)
		if !opts.Quiet {
			utils.PrintSuccess(fmt.Sprintf("Found %d snapshots", len(snapshots)), opts)
		}
		return nil
	},
	Args: cobra.ExactArgs(1),
}

var backupsExportCmd = &cobra.Command{
	Use:   "export <deployment-uid>",
	Short: "Start an async export of an addon backup snapshot",
	RunE: func(cmd *cobra.Command, args []string) error {
		opts := utils.GetOutputOptions(cmd)
		client, err := addonsClient(cmd, opts)
		if err != nil || client == nil {
			return err
		}
		snapshotID, _ := cmd.Flags().GetString("snapshot-id")
		path, _ := cmd.Flags().GetString("path")
		format, _ := cmd.Flags().GetString("format")
		resp, err := client.StartAddonBackupExport(context.Background(), args[0], &sdk.AddonBackupExportRequest{
			SnapshotID: snapshotID,
			Path:       path,
			Format:     format,
		})
		if err != nil {
			return fmt.Errorf("start addon backup export: %w", err)
		}
		return printAddonBackupExport(resp, opts, "Addon backup export started")
	},
	Args: cobra.ExactArgs(1),
}

var backupsExportStatusCmd = &cobra.Command{
	Use:   "export-status <deployment-uid> <export-id>",
	Short: "Get addon backup export status",
	RunE: func(cmd *cobra.Command, args []string) error {
		opts := utils.GetOutputOptions(cmd)
		client, err := addonsClient(cmd, opts)
		if err != nil || client == nil {
			return err
		}
		resp, err := client.GetAddonBackupExport(context.Background(), args[0], args[1])
		if err != nil {
			return fmt.Errorf("get addon backup export status: %w", err)
		}
		return printAddonBackupExport(resp, opts, "")
	},
	Args: cobra.ExactArgs(2),
}

func printAddonBackupExport(resp *sdk.AddonBackupExportResponse, opts utils.OutputOptions, successMsg string) error {
	if resp == nil {
		return fmt.Errorf("empty export response")
	}
	if opts.Format == utils.OutputFormatJSON {
		return utils.PrintJSON(resp)
	}
	if successMsg != "" {
		utils.PrintSuccess(successMsg, opts)
	}
	size := ""
	if resp.Data.SizeBytes > 0 {
		size = formatBytes(resp.Data.SizeBytes)
	}
	utils.PrintTable([]string{"ATTRIBUTE", "VALUE"}, [][]string{
		{"Export ID", resp.Data.ExportID},
		{"Status", resp.Data.Status},
		{"Snapshot ID", resp.Data.SnapshotID},
		{"Download URL", resp.Data.DownloadURL},
		{"Filename", resp.Data.Filename},
		{"Content Type", resp.Data.ContentType},
		{"Size", size},
		{"Path", resp.Data.Path},
		{"Error", resp.Data.ErrorMessage},
		{"Created", resp.Data.CreatedAt},
	}, opts)
	return nil
}

func formatBytes(n int64) string {
	const unit = 1024
	if n < unit {
		return strconv.FormatInt(n, 10) + " B"
	}
	div, exp := int64(unit), 0
	for v := n / unit; v >= unit; v /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(n)/float64(div), "KMGTPE"[exp])
}

func boolYesNo(v bool) string {
	if v {
		return "yes"
	}
	return "no"
}

func init() {
	backupsExportCmd.Flags().String("snapshot-id", "", "Snapshot ID to export")
	backupsExportCmd.Flags().String("path", "", "Optional path within the snapshot")
	backupsExportCmd.Flags().String("format", "", "Export format: auto, sql, rdb, or archive")
	_ = backupsExportCmd.MarkFlagRequired("snapshot-id")

	// Normalize format values if provided.
	backupsExportCmd.PreRunE = func(cmd *cobra.Command, args []string) error {
		format, _ := cmd.Flags().GetString("format")
		if format == "" {
			return nil
		}
		format = strings.ToLower(strings.TrimSpace(format))
		switch format {
		case "auto", "sql", "rdb", "archive":
			return cmd.Flags().Set("format", format)
		default:
			return fmt.Errorf("--format must be one of: auto, sql, rdb, archive")
		}
	}

	backupsCmd.AddCommand(backupsListCmd, backupsExportCmd, backupsExportStatusCmd)
	AddonsCmd.AddCommand(backupsCmd)
}
