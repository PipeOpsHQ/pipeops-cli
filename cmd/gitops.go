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

var gitopsCmd = &cobra.Command{
	Use:     "gitops",
	Aliases: []string{"go", "git-ops"},
	Short:   "Manage GitOps application configurations",
	Long: `Manage GitOps applications (create, sync, status, diff, history).

Examples:
  pipeops gitops list
  pipeops gitops get <uuid>
  pipeops gitops create --name my-app --repo-url https://github.com/org/repo
  pipeops gitops update <uuid> --branch main
  pipeops gitops delete <uuid> --yes
  pipeops gitops sync <uuid>
  pipeops gitops status <uuid>
  pipeops gitops diff <uuid>
  pipeops gitops history <uuid>`,
}

var gitopsListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List GitOps applications",
	RunE: func(cmd *cobra.Command, args []string) error {
		opts := utils.GetOutputOptions(cmd)
		client, err := rootClient(opts)
		if err != nil || client == nil {
			return err
		}
		listOpts := &sdk.GitOpsListOptions{}
		if page, _ := cmd.Flags().GetInt("page"); page > 0 {
			listOpts.Page = page
		}
		if limit, _ := cmd.Flags().GetInt("limit"); limit > 0 {
			listOpts.Limit = limit
		}
		resp, err := client.ListGitOps(context.Background(), listOpts)
		if err != nil {
			return fmt.Errorf("list gitops: %w", err)
		}
		if opts.Format == utils.OutputFormatJSON {
			return utils.PrintJSON(resp)
		}
		items := resp.Data.Items
		if len(items) == 0 {
			utils.PrintWarning("No GitOps applications found", opts)
			return nil
		}
		rows := make([][]string, 0, len(items))
		for _, item := range items {
			rows = append(rows, []string{
				item.UUID,
				item.Name,
				item.RepoURL,
				item.Branch,
				item.SyncStatus,
				item.HealthStatus,
			})
		}
		utils.PrintTable([]string{"UUID", "NAME", "REPO", "BRANCH", "SYNC", "HEALTH"}, rows, opts)
		if !opts.Quiet {
			utils.PrintSuccess(fmt.Sprintf("Found %d GitOps applications (total: %d)", len(items), resp.Data.Total), opts)
		}
		return nil
	},
	Args: cobra.NoArgs,
}

var gitopsGetCmd = &cobra.Command{
	Use:   "get <uuid>",
	Short: "Get GitOps application details",
	RunE: func(cmd *cobra.Command, args []string) error {
		opts := utils.GetOutputOptions(cmd)
		client, err := rootClient(opts)
		if err != nil || client == nil {
			return err
		}
		cfg, err := client.GetGitOps(context.Background(), args[0])
		if err != nil {
			return fmt.Errorf("get gitops: %w", err)
		}
		return printGitOpsConfig(cfg, opts)
	},
	Args: cobra.ExactArgs(1),
}

var gitopsCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a GitOps application",
	RunE: func(cmd *cobra.Command, args []string) error {
		opts := utils.GetOutputOptions(cmd)
		client, err := rootClient(opts)
		if err != nil || client == nil {
			return err
		}
		name, _ := cmd.Flags().GetString("name")
		repoURL, _ := cmd.Flags().GetString("repo-url")
		branch, _ := cmd.Flags().GetString("branch")
		path, _ := cmd.Flags().GetString("path")
		targetRevision, _ := cmd.Flags().GetString("target-revision")
		manifestType, _ := cmd.Flags().GetString("manifest-type")

		body := &sdk.CreateGitOpsConfigRequest{
			Name:           name,
			RepoURL:        repoURL,
			Branch:         branch,
			Path:           path,
			TargetRevision: targetRevision,
			ManifestType:   manifestType,
		}
		if projectID, err := optionalUintFlag(cmd, "project-id"); err != nil {
			return err
		} else if projectID != nil {
			body.ProjectID = projectID
		}
		if envID, err := optionalUintFlag(cmd, "environment-id"); err != nil {
			return err
		} else if envID != nil {
			body.EnvironmentID = envID
		}

		cfg, err := client.CreateGitOps(context.Background(), body)
		if err != nil {
			return fmt.Errorf("create gitops: %w", err)
		}
		if opts.Format != utils.OutputFormatJSON {
			utils.PrintSuccess("GitOps application created", opts)
		}
		return printGitOpsConfig(cfg, opts)
	},
	Args: cobra.NoArgs,
}

var gitopsUpdateCmd = &cobra.Command{
	Use:   "update <uuid>",
	Short: "Update a GitOps application",
	RunE: func(cmd *cobra.Command, args []string) error {
		opts := utils.GetOutputOptions(cmd)
		client, err := rootClient(opts)
		if err != nil || client == nil {
			return err
		}
		body := &sdk.UpdateGitOpsConfigRequest{}
		if cmd.Flags().Changed("name") {
			body.Name, _ = cmd.Flags().GetString("name")
		}
		if cmd.Flags().Changed("branch") {
			body.Branch, _ = cmd.Flags().GetString("branch")
		}
		if cmd.Flags().Changed("path") {
			body.Path, _ = cmd.Flags().GetString("path")
		}
		if cmd.Flags().Changed("target-revision") {
			body.TargetRevision, _ = cmd.Flags().GetString("target-revision")
		}
		cfg, err := client.UpdateGitOps(context.Background(), args[0], body)
		if err != nil {
			return fmt.Errorf("update gitops: %w", err)
		}
		if opts.Format != utils.OutputFormatJSON {
			utils.PrintSuccess("GitOps application updated", opts)
		}
		return printGitOpsConfig(cfg, opts)
	},
	Args: cobra.ExactArgs(1),
}

var gitopsDeleteCmd = &cobra.Command{
	Use:   "delete <uuid>",
	Short: "Delete a GitOps application",
	RunE: func(cmd *cobra.Command, args []string) error {
		opts := utils.GetOutputOptions(cmd)
		yes, _ := cmd.Flags().GetBool("yes")
		if !yes {
			return fmt.Errorf("--yes is required to delete a GitOps application")
		}
		client, err := rootClient(opts)
		if err != nil || client == nil {
			return err
		}
		if err := client.DeleteGitOps(context.Background(), args[0]); err != nil {
			return fmt.Errorf("delete gitops: %w", err)
		}
		if opts.Format == utils.OutputFormatJSON {
			return utils.PrintJSON(map[string]string{"status": "deleted", "uuid": args[0]})
		}
		utils.PrintSuccess("GitOps application deleted", opts)
		return nil
	},
	Args: cobra.ExactArgs(1),
}

var gitopsSyncCmd = &cobra.Command{
	Use:   "sync <uuid>",
	Short: "Trigger a GitOps sync",
	RunE: func(cmd *cobra.Command, args []string) error {
		opts := utils.GetOutputOptions(cmd)
		client, err := rootClient(opts)
		if err != nil || client == nil {
			return err
		}
		revision, _ := cmd.Flags().GetString("revision")
		prune, _ := cmd.Flags().GetBool("prune")
		dryRun, _ := cmd.Flags().GetBool("dry-run")
		resp, err := client.TriggerGitOpsSync(context.Background(), args[0], &sdk.TriggerGitOpsSyncRequest{
			Revision: revision,
			Prune:    prune,
			DryRun:   dryRun,
		})
		if err != nil {
			return fmt.Errorf("sync gitops: %w", err)
		}
		if opts.Format == utils.OutputFormatJSON {
			return utils.PrintJSON(resp)
		}
		utils.PrintSuccess("GitOps sync triggered", opts)
		if resp != nil {
			utils.PrintTable([]string{"ATTRIBUTE", "VALUE"}, [][]string{
				{"Status", resp.Data.Status},
				{"Revision", resp.Data.Revision},
				{"Dry Run", boolString(resp.Data.DryRun)},
				{"Message", resp.Message},
			}, opts)
		}
		return nil
	},
	Args: cobra.ExactArgs(1),
}

var gitopsStatusCmd = &cobra.Command{
	Use:   "status <uuid>",
	Short: "Get GitOps sync status",
	RunE: func(cmd *cobra.Command, args []string) error {
		opts := utils.GetOutputOptions(cmd)
		client, err := rootClient(opts)
		if err != nil || client == nil {
			return err
		}
		resp, err := client.GetGitOpsSyncStatus(context.Background(), args[0])
		if err != nil {
			return fmt.Errorf("gitops status: %w", err)
		}
		if opts.Format == utils.OutputFormatJSON {
			return utils.PrintJSON(resp)
		}
		lastSynced := ""
		if resp.Data.LastSyncedAt != nil {
			lastSynced = *resp.Data.LastSyncedAt
		}
		utils.PrintTable([]string{"ATTRIBUTE", "VALUE"}, [][]string{
			{"Sync Status", resp.Data.SyncStatus},
			{"Sync Message", resp.Data.SyncMessage},
			{"Last Synced Commit", resp.Data.LastSyncedCommit},
			{"Last Synced At", lastSynced},
			{"Health Status", resp.Data.HealthStatus},
			{"Health Message", resp.Data.HealthMessage},
		}, opts)
		return nil
	},
	Args: cobra.ExactArgs(1),
}

var gitopsDiffCmd = &cobra.Command{
	Use:   "diff <uuid>",
	Short: "Show GitOps diff (git vs live)",
	RunE: func(cmd *cobra.Command, args []string) error {
		opts := utils.GetOutputOptions(cmd)
		client, err := rootClient(opts)
		if err != nil || client == nil {
			return err
		}
		resp, err := client.GetGitOpsDiff(context.Background(), args[0])
		if err != nil {
			return fmt.Errorf("gitops diff: %w", err)
		}
		if opts.Format == utils.OutputFormatJSON {
			return utils.PrintJSON(resp)
		}
		utils.PrintTable([]string{"ATTRIBUTE", "VALUE"}, [][]string{
			{"Current Commit", resp.Data.CurrentCommit},
			{"Target Commit", resp.Data.TargetCommit},
			{"Sync Required", boolString(resp.Data.SyncRequired)},
		}, opts)
		if resp.Data.Diff != nil {
			printGitOpsDiffSnapshot(resp.Data.Diff, opts)
		}
		return nil
	},
	Args: cobra.ExactArgs(1),
}

var gitopsHistoryCmd = &cobra.Command{
	Use:   "history <uuid>",
	Short: "Show GitOps sync history",
	RunE: func(cmd *cobra.Command, args []string) error {
		opts := utils.GetOutputOptions(cmd)
		client, err := rootClient(opts)
		if err != nil || client == nil {
			return err
		}
		listOpts := &sdk.GitOpsListOptions{}
		if page, _ := cmd.Flags().GetInt("page"); page > 0 {
			listOpts.Page = page
		}
		if limit, _ := cmd.Flags().GetInt("limit"); limit > 0 {
			listOpts.Limit = limit
		}
		resp, err := client.GetGitOpsHistory(context.Background(), args[0], listOpts)
		if err != nil {
			return fmt.Errorf("gitops history: %w", err)
		}
		if opts.Format == utils.OutputFormatJSON {
			return utils.PrintJSON(resp)
		}
		items := resp.Data.Items
		if len(items) == 0 {
			utils.PrintWarning("No sync history found", opts)
			return nil
		}
		rows := make([][]string, 0, len(items))
		for _, item := range items {
			rows = append(rows, []string{
				strconv.FormatUint(uint64(item.ID), 10),
				shortSHA(item.CommitSHA),
				item.SyncStatus,
				item.TriggeredBy,
				item.StartedAt,
				item.FinishedAt,
			})
		}
		utils.PrintTable([]string{"ID", "COMMIT", "STATUS", "TRIGGERED BY", "STARTED", "FINISHED"}, rows, opts)
		if !opts.Quiet {
			utils.PrintSuccess(fmt.Sprintf("Found %d history entries (total: %d)", len(items), resp.Data.Total), opts)
		}
		return nil
	},
	Args: cobra.ExactArgs(1),
}

func printGitOpsConfig(cfg *sdk.GitOpsConfig, opts utils.OutputOptions) error {
	if cfg == nil {
		return fmt.Errorf("gitops application not found")
	}
	if opts.Format == utils.OutputFormatJSON {
		return utils.PrintJSON(cfg)
	}
	projectID := ""
	if cfg.ProjectID != nil {
		projectID = strconv.FormatUint(uint64(*cfg.ProjectID), 10)
	}
	envID := ""
	if cfg.EnvironmentID != nil {
		envID = strconv.FormatUint(uint64(*cfg.EnvironmentID), 10)
	}
	utils.PrintTable([]string{"ATTRIBUTE", "VALUE"}, [][]string{
		{"UUID", cfg.UUID},
		{"Name", cfg.Name},
		{"Repo URL", cfg.RepoURL},
		{"Branch", cfg.Branch},
		{"Path", cfg.Path},
		{"Target Revision", cfg.TargetRevision},
		{"Project ID", projectID},
		{"Project Name", cfg.ProjectName},
		{"Environment ID", envID},
		{"Environment Name", cfg.EnvironmentName},
		{"Sync Status", cfg.SyncStatus},
		{"Sync Message", cfg.SyncMessage},
		{"Health Status", cfg.HealthStatus},
		{"Health Message", cfg.HealthMessage},
		{"Last Synced Commit", cfg.LastSyncedCommit},
		{"Last Synced At", cfg.LastSyncedAt},
		{"Created", cfg.CreatedAt},
		{"Updated", cfg.UpdatedAt},
	}, opts)
	return nil
}

func printGitOpsDiffSnapshot(diff *sdk.GitOpsDiffSnapshot, opts utils.OutputOptions) {
	printChangeGroup := func(title string, changes []sdk.GitOpsResourceChange) {
		if len(changes) == 0 {
			return
		}
		utils.PrintInfo(title, opts)
		rows := make([][]string, 0, len(changes))
		for _, c := range changes {
			rows = append(rows, []string{c.Kind, c.Name, c.Field, fmt.Sprint(c.OldValue), fmt.Sprint(c.NewValue)})
		}
		utils.PrintTable([]string{"KIND", "NAME", "FIELD", "OLD", "NEW"}, rows, opts)
	}
	printChangeGroup("Added", diff.Added)
	printChangeGroup("Modified", diff.Modified)
	printChangeGroup("Removed", diff.Removed)
}

func shortSHA(sha string) string {
	if len(sha) > 8 {
		return sha[:8]
	}
	return sha
}

func optionalUintFlag(cmd *cobra.Command, name string) (*uint, error) {
	raw, _ := cmd.Flags().GetString(name)
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil, nil
	}
	parsed, err := strconv.ParseUint(raw, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid --%s: %w", name, err)
	}
	v := uint(parsed)
	return &v, nil
}

func init() {
	gitopsListCmd.Flags().Int("page", 0, "Page number")
	gitopsListCmd.Flags().Int("limit", 0, "Page size")

	gitopsCreateCmd.Flags().String("name", "", "Application name")
	gitopsCreateCmd.Flags().String("repo-url", "", "Git repository URL")
	gitopsCreateCmd.Flags().String("branch", "", "Git branch (default: main)")
	gitopsCreateCmd.Flags().String("path", "", "Path within the repository")
	gitopsCreateCmd.Flags().String("target-revision", "", "Target revision (branch/tag/commit)")
	gitopsCreateCmd.Flags().String("project-id", "", "Optional numeric project ID")
	gitopsCreateCmd.Flags().String("environment-id", "", "Optional numeric environment ID")
	gitopsCreateCmd.Flags().String("manifest-type", "", "Manifest type: pipeops | raw")
	_ = gitopsCreateCmd.MarkFlagRequired("name")
	_ = gitopsCreateCmd.MarkFlagRequired("repo-url")

	gitopsUpdateCmd.Flags().String("name", "", "Application name")
	gitopsUpdateCmd.Flags().String("branch", "", "Git branch")
	gitopsUpdateCmd.Flags().String("path", "", "Path within the repository")
	gitopsUpdateCmd.Flags().String("target-revision", "", "Target revision")

	gitopsDeleteCmd.Flags().Bool("yes", false, "Confirm deletion")

	gitopsSyncCmd.Flags().String("revision", "", "Revision to sync")
	gitopsSyncCmd.Flags().Bool("prune", false, "Prune resources")
	gitopsSyncCmd.Flags().Bool("dry-run", false, "Dry-run sync without applying")

	gitopsHistoryCmd.Flags().Int("page", 0, "Page number")
	gitopsHistoryCmd.Flags().Int("limit", 0, "Page size")

	gitopsCmd.AddCommand(
		gitopsListCmd,
		gitopsGetCmd,
		gitopsCreateCmd,
		gitopsUpdateCmd,
		gitopsDeleteCmd,
		gitopsSyncCmd,
		gitopsStatusCmd,
		gitopsDiffCmd,
		gitopsHistoryCmd,
	)
	rootCmd.AddCommand(gitopsCmd)
}
