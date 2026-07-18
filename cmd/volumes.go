package cmd

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/PipeOpsHQ/pipeops-cli/utils"
	sdk "github.com/PipeOpsHQ/pipeops-go-sdk/pipeops"
	"github.com/spf13/cobra"
)

var volumesCmd = &cobra.Command{
	Use:     "volumes",
	Aliases: []string{"volume", "vol"},
	Short:   "Manage workspace volumes",
	Long: `Manage workspace volumes (PVC inventory, remount, delete, and export).

Examples:
  pipeops volumes list
  pipeops volumes list --workspace <workspace-uuid>
  pipeops volumes get <volume-uuid>
  pipeops volumes remount <volume-uuid> --target-type project --target-uuid <uuid>
  pipeops volumes delete <volume-uuid> --yes
  pipeops volumes export <volume-uuid>
  pipeops volumes export-status <volume-uuid>`,
}

func volumeListOpts(cmd *cobra.Command) *sdk.VolumeListOptions {
	opts := &sdk.VolumeListOptions{}
	if workspace, _ := cmd.Flags().GetString("workspace"); workspace != "" {
		opts.WorkspaceUUID = workspace
	}
	if status, err := cmd.Flags().GetString("status"); err == nil && status != "" {
		opts.Status = status
	}
	return opts
}

var volumesListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List workspace volumes",
	RunE: func(cmd *cobra.Command, args []string) error {
		opts := utils.GetOutputOptions(cmd)
		client, err := rootClient(opts)
		if err != nil || client == nil {
			return err
		}
		resp, err := client.ListVolumes(context.Background(), volumeListOpts(cmd))
		if err != nil {
			return fmt.Errorf("list volumes: %w", err)
		}
		if opts.Format == utils.OutputFormatJSON {
			return utils.PrintJSON(resp)
		}
		volumes := resp.Data.Volumes
		if len(volumes) == 0 {
			utils.PrintWarning("No volumes found", opts)
			return nil
		}
		rows := make([][]string, 0, len(volumes))
		for _, v := range volumes {
			size := ""
			if v.SizeGB > 0 {
				size = strconv.FormatFloat(float64(v.SizeGB), 'f', 1, 32) + " GB"
			}
			rows = append(rows, []string{
				v.UUID,
				displayOr(v.DisplayName, v.PVCName),
				v.Status,
				v.OwnerType,
				displayOr(v.OwnerName, v.OwnerUUID),
				size,
				v.MountPath,
			})
		}
		utils.PrintTable([]string{"UUID", "NAME", "STATUS", "OWNER TYPE", "OWNER", "SIZE", "MOUNT"}, rows, opts)
		if !opts.Quiet {
			utils.PrintSuccess(fmt.Sprintf("Found %d volumes (mounted: %d, unattached: %d)",
				resp.Data.Total, resp.Data.Summary.Mounted, resp.Data.Summary.Unattached), opts)
		}
		return nil
	},
	Args: cobra.NoArgs,
}

var volumesGetCmd = &cobra.Command{
	Use:   "get <volume-uuid>",
	Short: "Get volume details",
	RunE: func(cmd *cobra.Command, args []string) error {
		opts := utils.GetOutputOptions(cmd)
		client, err := rootClient(opts)
		if err != nil || client == nil {
			return err
		}
		volume, err := client.GetVolume(context.Background(), args[0], volumeListOpts(cmd))
		if err != nil {
			return fmt.Errorf("get volume: %w", err)
		}
		return printVolume(volume, opts)
	},
	Args: cobra.ExactArgs(1),
}

var volumesRemountCmd = &cobra.Command{
	Use:   "remount <volume-uuid>",
	Short: "Remount an unattached volume onto a project or addon",
	RunE: func(cmd *cobra.Command, args []string) error {
		opts := utils.GetOutputOptions(cmd)
		client, err := rootClient(opts)
		if err != nil || client == nil {
			return err
		}
		targetType, _ := cmd.Flags().GetString("target-type")
		targetUUID, _ := cmd.Flags().GetString("target-uuid")
		mountPath, _ := cmd.Flags().GetString("mount-path")
		targetType = strings.ToLower(strings.TrimSpace(targetType))
		if targetType != "project" && targetType != "addon" {
			return fmt.Errorf("--target-type must be project or addon")
		}
		resp, err := client.RemountVolume(context.Background(), args[0], &sdk.RemountVolumeRequest{
			TargetType: targetType,
			TargetUUID: targetUUID,
			MountPath:  mountPath,
		}, volumeListOpts(cmd))
		if err != nil {
			return fmt.Errorf("remount volume: %w", err)
		}
		if opts.Format == utils.OutputFormatJSON {
			return utils.PrintJSON(resp)
		}
		utils.PrintSuccess("Volume remount scheduled", opts)
		if resp != nil {
			if msg := resp.Data.Message; msg != "" {
				utils.PrintInfo(msg, opts)
			}
			return printVolume(&resp.Data.Volume, opts)
		}
		return nil
	},
	Args: cobra.ExactArgs(1),
}

var volumesDeleteCmd = &cobra.Command{
	Use:   "delete <volume-uuid>",
	Short: "Delete a volume permanently",
	RunE: func(cmd *cobra.Command, args []string) error {
		opts := utils.GetOutputOptions(cmd)
		yes, _ := cmd.Flags().GetBool("yes")
		if !yes {
			return fmt.Errorf("--yes is required to delete a volume")
		}
		client, err := rootClient(opts)
		if err != nil || client == nil {
			return err
		}
		if err := client.DeleteVolume(context.Background(), args[0], volumeListOpts(cmd)); err != nil {
			return fmt.Errorf("delete volume: %w", err)
		}
		if opts.Format == utils.OutputFormatJSON {
			return utils.PrintJSON(map[string]string{"status": "deleted", "volume_uuid": args[0]})
		}
		utils.PrintSuccess("Volume deleted", opts)
		return nil
	},
	Args: cobra.ExactArgs(1),
}

var volumesExportCmd = &cobra.Command{
	Use:   "export <volume-uuid>",
	Short: "Start an async volume export",
	RunE: func(cmd *cobra.Command, args []string) error {
		opts := utils.GetOutputOptions(cmd)
		client, err := rootClient(opts)
		if err != nil || client == nil {
			return err
		}
		resp, err := client.StartVolumeExport(context.Background(), args[0], volumeListOpts(cmd))
		if err != nil {
			return fmt.Errorf("start volume export: %w", err)
		}
		return printVolumeExport(resp, opts, "Volume export started")
	},
	Args: cobra.ExactArgs(1),
}

var volumesExportStatusCmd = &cobra.Command{
	Use:   "export-status <volume-uuid>",
	Short: "Get volume export status",
	RunE: func(cmd *cobra.Command, args []string) error {
		opts := utils.GetOutputOptions(cmd)
		client, err := rootClient(opts)
		if err != nil || client == nil {
			return err
		}
		resp, err := client.GetVolumeExport(context.Background(), args[0], volumeListOpts(cmd))
		if err != nil {
			return fmt.Errorf("get volume export status: %w", err)
		}
		return printVolumeExport(resp, opts, "")
	},
	Args: cobra.ExactArgs(1),
}

func printVolume(volume *sdk.Volume, opts utils.OutputOptions) error {
	if volume == nil {
		return fmt.Errorf("volume not found")
	}
	if opts.Format == utils.OutputFormatJSON {
		return utils.PrintJSON(volume)
	}
	size := ""
	if volume.SizeGB > 0 {
		size = strconv.FormatFloat(float64(volume.SizeGB), 'f', 1, 32) + " GB"
	}
	utils.PrintTable([]string{"ATTRIBUTE", "VALUE"}, [][]string{
		{"UUID", volume.UUID},
		{"Name", displayOr(volume.DisplayName, volume.PVCName)},
		{"PVC", volume.PVCName},
		{"Status", volume.Status},
		{"Size", size},
		{"Mount Path", volume.MountPath},
		{"Owner Type", volume.OwnerType},
		{"Owner", displayOr(volume.OwnerName, volume.OwnerUUID)},
		{"Cluster", displayOr(volume.ClusterName, volume.ClusterUUID)},
		{"Namespace", volume.Namespace},
		{"Export Status", volume.ExportStatus},
		{"Export URL", volume.ExportURL},
		{"Created", volume.CreatedAt},
		{"Updated", volume.UpdatedAt},
	}, opts)
	return nil
}

func printVolumeExport(resp *sdk.VolumeExportResponse, opts utils.OutputOptions, successMsg string) error {
	if resp == nil {
		return fmt.Errorf("empty export response")
	}
	if opts.Format == utils.OutputFormatJSON {
		return utils.PrintJSON(resp)
	}
	if successMsg != "" {
		utils.PrintSuccess(successMsg, opts)
	}
	utils.PrintTable([]string{"ATTRIBUTE", "VALUE"}, [][]string{
		{"UUID", resp.Data.UUID},
		{"Status", resp.Data.Status},
		{"Download URL", resp.Data.DownloadURL},
		{"Filename", resp.Data.Filename},
		{"Error", resp.Data.Error},
		{"Message", displayOr(resp.Data.Message, resp.Message)},
	}, opts)
	return nil
}

func displayOr(primary, fallback string) string {
	if strings.TrimSpace(primary) != "" {
		return primary
	}
	return fallback
}

func init() {
	workspaceFlag := "Workspace UUID (or set PIPEOPS_WORKSPACE_UUID / pipeops workspace select)"

	volumesListCmd.Flags().String("workspace", "", workspaceFlag)
	volumesListCmd.Flags().String("status", "", "Filter by status (mounted, unattached)")
	volumesGetCmd.Flags().String("workspace", "", workspaceFlag)

	volumesRemountCmd.Flags().String("workspace", "", workspaceFlag)
	volumesRemountCmd.Flags().String("target-type", "", "Target type: project or addon")
	volumesRemountCmd.Flags().String("target-uuid", "", "Target project or addon UUID")
	volumesRemountCmd.Flags().String("mount-path", "", "Optional mount path")
	_ = volumesRemountCmd.MarkFlagRequired("target-type")
	_ = volumesRemountCmd.MarkFlagRequired("target-uuid")

	volumesDeleteCmd.Flags().String("workspace", "", workspaceFlag)
	volumesDeleteCmd.Flags().Bool("yes", false, "Confirm volume deletion")

	volumesExportCmd.Flags().String("workspace", "", workspaceFlag)
	volumesExportStatusCmd.Flags().String("workspace", "", workspaceFlag)

	volumesCmd.AddCommand(
		volumesListCmd,
		volumesGetCmd,
		volumesRemountCmd,
		volumesDeleteCmd,
		volumesExportCmd,
		volumesExportStatusCmd,
	)
	rootCmd.AddCommand(volumesCmd)
}
