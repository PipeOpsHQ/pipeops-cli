package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/PipeOpsHQ/pipeops-cli/utils"
	sdk "github.com/PipeOpsHQ/pipeops-go-sdk/pipeops"
	"github.com/spf13/cobra"
)

const workspaceFlagHelp = "Workspace UUID (or set PIPEOPS_WORKSPACE_UUID / pipeops workspace select)"

var groupsCmd = &cobra.Command{
	Use:     "groups",
	Aliases: []string{"group", "project-groups", "project-group"},
	Short:   "Manage project groups (unified project plane)",
	Long: `Manage project groups, members, shared env, topology, and connections.

Examples:
  pipeops groups list
  pipeops groups get <uuid>
  pipeops groups create --name my-plane
  pipeops groups update <uuid> --name new-name
  pipeops groups delete <uuid> --yes
  pipeops groups topology <uuid>
  pipeops groups members attach <uuid> --type project --member-uuid <uuid>
  pipeops groups members detach <uuid> --type addon --member-uuid <uuid>
  pipeops groups env get <uuid>
  pipeops groups env put <uuid> --set KEY=VAL
  pipeops groups env inject <uuid>
  pipeops groups connect <uuid> --consumer-uuid <uuid> --provider-uuid <uuid>
  pipeops groups redeploy <uuid>
  pipeops groups resolve --type project --member-uuid <uuid>
  pipeops groups candidates`,
}

func groupsWorkspaceOpts(cmd *cobra.Command) *sdk.ProjectGroupWorkspaceOptions {
	opts := &sdk.ProjectGroupWorkspaceOptions{}
	if workspace, _ := cmd.Flags().GetString("workspace"); workspace != "" {
		opts.WorkspaceUUID = workspace
	}
	return opts
}

func groupsListOpts(cmd *cobra.Command) *sdk.ProjectGroupListOptions {
	opts := &sdk.ProjectGroupListOptions{}
	if workspace, _ := cmd.Flags().GetString("workspace"); workspace != "" {
		opts.WorkspaceUUID = workspace
	}
	if limit, _ := cmd.Flags().GetInt("limit"); limit > 0 {
		opts.Limit = limit
	}
	if offset, _ := cmd.Flags().GetInt("offset"); offset > 0 {
		opts.Offset = offset
	}
	return opts
}

var groupsListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List project groups",
	RunE: func(cmd *cobra.Command, args []string) error {
		opts := utils.GetOutputOptions(cmd)
		client, err := rootClient(opts)
		if err != nil || client == nil {
			return err
		}
		resp, err := client.ListProjectGroups(context.Background(), groupsListOpts(cmd))
		if err != nil {
			return fmt.Errorf("list project groups: %w", err)
		}
		if opts.Format == utils.OutputFormatJSON {
			return utils.PrintJSON(resp)
		}
		groups := resp.Data.Groups
		if len(groups) == 0 {
			utils.PrintWarning("No project groups found", opts)
			return nil
		}
		rows := make([][]string, 0, len(groups))
		for _, g := range groups {
			rows = append(rows, []string{
				g.UUID,
				g.Name,
				strconv.Itoa(g.MemberCount),
				g.WorkspaceUUID,
				g.CreatedAt,
			})
		}
		utils.PrintTable([]string{"UUID", "NAME", "MEMBERS", "WORKSPACE", "CREATED"}, rows, opts)
		if !opts.Quiet {
			utils.PrintSuccess(fmt.Sprintf("Found %d project groups (total: %d)", len(groups), resp.Data.Total), opts)
		}
		return nil
	},
	Args: cobra.NoArgs,
}

var groupsGetCmd = &cobra.Command{
	Use:   "get <uuid>",
	Short: "Get project group details",
	RunE: func(cmd *cobra.Command, args []string) error {
		opts := utils.GetOutputOptions(cmd)
		client, err := rootClient(opts)
		if err != nil || client == nil {
			return err
		}
		group, err := client.GetProjectGroup(context.Background(), args[0], groupsWorkspaceOpts(cmd))
		if err != nil {
			return fmt.Errorf("get project group: %w", err)
		}
		return printProjectGroup(group, opts)
	},
	Args: cobra.ExactArgs(1),
}

var groupsCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a project group",
	RunE: func(cmd *cobra.Command, args []string) error {
		opts := utils.GetOutputOptions(cmd)
		client, err := rootClient(opts)
		if err != nil || client == nil {
			return err
		}
		name, _ := cmd.Flags().GetString("name")
		body := &sdk.CreateProjectGroupRequest{Name: name}
		if cluster, _ := cmd.Flags().GetString("cluster-uuid"); cluster != "" {
			body.DefaultClusterUUID = &cluster
		}
		if env, _ := cmd.Flags().GetString("environment-uuid"); env != "" {
			body.DefaultEnvironmentUUID = &env
		}
		group, err := client.CreateProjectGroup(context.Background(), body, groupsWorkspaceOpts(cmd))
		if err != nil {
			return fmt.Errorf("create project group: %w", err)
		}
		if opts.Format != utils.OutputFormatJSON {
			utils.PrintSuccess("Project group created", opts)
		}
		return printProjectGroup(group, opts)
	},
	Args: cobra.NoArgs,
}

var groupsUpdateCmd = &cobra.Command{
	Use:   "update <uuid>",
	Short: "Update a project group",
	RunE: func(cmd *cobra.Command, args []string) error {
		opts := utils.GetOutputOptions(cmd)
		client, err := rootClient(opts)
		if err != nil || client == nil {
			return err
		}
		body := &sdk.UpdateProjectGroupRequest{}
		if cmd.Flags().Changed("name") {
			name, _ := cmd.Flags().GetString("name")
			body.Name = &name
		}
		if cmd.Flags().Changed("cluster-uuid") {
			cluster, _ := cmd.Flags().GetString("cluster-uuid")
			body.DefaultClusterUUID = &cluster
		}
		if cmd.Flags().Changed("environment-uuid") {
			env, _ := cmd.Flags().GetString("environment-uuid")
			body.DefaultEnvironmentUUID = &env
		}
		group, err := client.UpdateProjectGroup(context.Background(), args[0], body, groupsWorkspaceOpts(cmd))
		if err != nil {
			return fmt.Errorf("update project group: %w", err)
		}
		if opts.Format != utils.OutputFormatJSON {
			utils.PrintSuccess("Project group updated", opts)
		}
		return printProjectGroup(group, opts)
	},
	Args: cobra.ExactArgs(1),
}

var groupsDeleteCmd = &cobra.Command{
	Use:   "delete <uuid>",
	Short: "Delete a project group",
	RunE: func(cmd *cobra.Command, args []string) error {
		opts := utils.GetOutputOptions(cmd)
		yes, _ := cmd.Flags().GetBool("yes")
		if !yes {
			return fmt.Errorf("--yes is required to delete a project group")
		}
		client, err := rootClient(opts)
		if err != nil || client == nil {
			return err
		}
		if err := client.DeleteProjectGroup(context.Background(), args[0], groupsWorkspaceOpts(cmd)); err != nil {
			return fmt.Errorf("delete project group: %w", err)
		}
		if opts.Format == utils.OutputFormatJSON {
			return utils.PrintJSON(map[string]string{"status": "deleted", "uuid": args[0]})
		}
		utils.PrintSuccess("Project group deleted", opts)
		return nil
	},
	Args: cobra.ExactArgs(1),
}

var groupsTopologyCmd = &cobra.Command{
	Use:   "topology <uuid>",
	Short: "Show project group topology",
	RunE: func(cmd *cobra.Command, args []string) error {
		opts := utils.GetOutputOptions(cmd)
		client, err := rootClient(opts)
		if err != nil || client == nil {
			return err
		}
		resp, err := client.GetProjectGroupTopology(context.Background(), args[0], groupsWorkspaceOpts(cmd))
		if err != nil {
			return fmt.Errorf("get project group topology: %w", err)
		}
		if opts.Format == utils.OutputFormatJSON {
			return utils.PrintJSON(resp)
		}
		topo := resp.Data
		utils.PrintTable([]string{"ATTRIBUTE", "VALUE"}, [][]string{
			{"Group", topo.Group.Name},
			{"Group UUID", topo.Group.UUID},
			{"Members", strconv.Itoa(topo.TotalMemberCount)},
			{"Visible Members", strconv.Itoa(topo.VisibleMemberCount)},
			{"Active Environment", topo.ActiveEnvironment},
		}, opts)
		if len(topo.Nodes) > 0 {
			utils.PrintInfo("Nodes", opts)
			rows := make([][]string, 0, len(topo.Nodes))
			for _, n := range topo.Nodes {
				rows = append(rows, []string{n.MemberType, n.MemberUUID, n.Name, n.ServiceKind, n.Status})
			}
			utils.PrintTable([]string{"TYPE", "UUID", "NAME", "KIND", "STATUS"}, rows, opts)
		}
		if len(topo.Edges) > 0 {
			utils.PrintInfo("Edges", opts)
			rows := make([][]string, 0, len(topo.Edges))
			for _, e := range topo.Edges {
				rows = append(rows, []string{e.Type, e.FromUUID, e.ToUUID, e.Label, e.Confidence})
			}
			utils.PrintTable([]string{"TYPE", "FROM", "TO", "LABEL", "CONFIDENCE"}, rows, opts)
		}
		if len(topo.Warnings) > 0 {
			for _, w := range topo.Warnings {
				utils.PrintWarning(w, opts)
			}
		}
		return nil
	},
	Args: cobra.ExactArgs(1),
}

var groupsMembersCmd = &cobra.Command{
	Use:   "members",
	Short: "Manage project group members",
}

var groupsMembersAttachCmd = &cobra.Command{
	Use:   "attach <group-uuid>",
	Short: "Attach a project or addon to a group",
	RunE: func(cmd *cobra.Command, args []string) error {
		opts := utils.GetOutputOptions(cmd)
		client, err := rootClient(opts)
		if err != nil || client == nil {
			return err
		}
		memberType, err := normalizeGroupMemberType(cmd)
		if err != nil {
			return err
		}
		memberUUID, _ := cmd.Flags().GetString("member-uuid")
		move, _ := cmd.Flags().GetBool("move")
		body := &sdk.AttachProjectGroupMemberRequest{
			MemberType: memberType,
			MemberUUID: memberUUID,
			Move:       move,
		}
		if cmd.Flags().Changed("include-session") {
			include, _ := cmd.Flags().GetBool("include-session")
			body.IncludeSession = &include
		}
		resp, err := client.AttachProjectGroupMember(context.Background(), args[0], body, groupsWorkspaceOpts(cmd))
		if err != nil {
			return fmt.Errorf("attach member: %w", err)
		}
		if opts.Format == utils.OutputFormatJSON {
			return utils.PrintJSON(resp)
		}
		utils.PrintSuccess("Member attached", opts)
		if resp != nil {
			utils.PrintTable([]string{"ATTRIBUTE", "VALUE"}, [][]string{
				{"Group UUID", resp.Data.GroupUUID},
				{"Attached", strings.Join(resp.Data.AttachedMemberUUIDs, ", ")},
				{"Session Applied", boolString(resp.Data.IncludeSessionApplied)},
				{"Message", resp.Message},
			}, opts)
		}
		return nil
	},
	Args: cobra.ExactArgs(1),
}

var groupsMembersDetachCmd = &cobra.Command{
	Use:   "detach <group-uuid>",
	Short: "Detach a project or addon from a group",
	RunE: func(cmd *cobra.Command, args []string) error {
		opts := utils.GetOutputOptions(cmd)
		client, err := rootClient(opts)
		if err != nil || client == nil {
			return err
		}
		memberType, err := normalizeGroupMemberType(cmd)
		if err != nil {
			return err
		}
		memberUUID, _ := cmd.Flags().GetString("member-uuid")
		detachOpts := &sdk.ProjectGroupDetachOptions{}
		if workspace, _ := cmd.Flags().GetString("workspace"); workspace != "" {
			detachOpts.WorkspaceUUID = workspace
		}
		if cmd.Flags().Changed("include-session") {
			include, _ := cmd.Flags().GetBool("include-session")
			detachOpts.IncludeSession = &include
		}
		if err := client.DetachProjectGroupMember(context.Background(), args[0], memberType, memberUUID, detachOpts); err != nil {
			return fmt.Errorf("detach member: %w", err)
		}
		if opts.Format == utils.OutputFormatJSON {
			return utils.PrintJSON(map[string]string{
				"status":      "detached",
				"group_uuid":  args[0],
				"member_type": memberType,
				"member_uuid": memberUUID,
			})
		}
		utils.PrintSuccess("Member detached", opts)
		return nil
	},
	Args: cobra.ExactArgs(1),
}

var groupsEnvCmd = &cobra.Command{
	Use:   "env",
	Short: "Manage shared environment variables for a project group",
}

var groupsEnvGetCmd = &cobra.Command{
	Use:   "get <uuid>",
	Short: "Get shared environment variables",
	RunE: func(cmd *cobra.Command, args []string) error {
		opts := utils.GetOutputOptions(cmd)
		client, err := rootClient(opts)
		if err != nil || client == nil {
			return err
		}
		resp, err := client.GetProjectGroupSharedEnv(context.Background(), args[0], groupsWorkspaceOpts(cmd))
		if err != nil {
			return fmt.Errorf("get shared env: %w", err)
		}
		if opts.Format == utils.OutputFormatJSON {
			return utils.PrintJSON(resp)
		}
		vars := resp.Data.Variables
		if len(vars) == 0 {
			utils.PrintWarning("No shared environment variables", opts)
			return nil
		}
		rows := make([][]string, 0, len(vars))
		for _, v := range vars {
			rows = append(rows, []string{v.Key, v.Value})
		}
		utils.PrintTable([]string{"KEY", "VALUE"}, rows, opts)
		return nil
	},
	Args: cobra.ExactArgs(1),
}

var groupsEnvPutCmd = &cobra.Command{
	Use:   "put <uuid>",
	Short: "Replace shared environment variables",
	Long: `Replace the group shared env set.

Use --set KEY=VAL (repeatable), --file <json>, or --json-body <file>.
JSON file may be either:
  {"variables":[{"key":"K","value":"V"}],"inject":true,...}
  or a plain object map: {"KEY":"VAL",...}`,
	RunE: func(cmd *cobra.Command, args []string) error {
		opts := utils.GetOutputOptions(cmd)
		client, err := rootClient(opts)
		if err != nil || client == nil {
			return err
		}
		body, err := buildSharedEnvPutRequest(cmd)
		if err != nil {
			return err
		}
		resp, err := client.PutProjectGroupSharedEnv(context.Background(), args[0], body, groupsWorkspaceOpts(cmd))
		if err != nil {
			return fmt.Errorf("put shared env: %w", err)
		}
		if opts.Format == utils.OutputFormatJSON {
			return utils.PrintJSON(resp)
		}
		utils.PrintSuccess("Shared environment variables updated", opts)
		if resp != nil && resp.Data.Message != "" {
			utils.PrintInfo(resp.Data.Message, opts)
		}
		return nil
	},
	Args: cobra.ExactArgs(1),
}

var groupsEnvInjectCmd = &cobra.Command{
	Use:   "inject <uuid>",
	Short: "Inject shared environment variables into members",
	RunE: func(cmd *cobra.Command, args []string) error {
		opts := utils.GetOutputOptions(cmd)
		client, err := rootClient(opts)
		if err != nil || client == nil {
			return err
		}
		body := &sdk.InjectProjectGroupSharedEnvRequest{}
		if jsonBody, _ := cmd.Flags().GetString("json-body"); jsonBody != "" {
			raw, err := os.ReadFile(jsonBody)
			if err != nil {
				return fmt.Errorf("read --json-body: %w", err)
			}
			if err := json.Unmarshal(raw, body); err != nil {
				return fmt.Errorf("parse --json-body: %w", err)
			}
		} else {
			overwrite, _ := cmd.Flags().GetBool("overwrite")
			redeploy, _ := cmd.Flags().GetBool("redeploy")
			keepRefs, _ := cmd.Flags().GetBool("keep-references")
			members, _ := cmd.Flags().GetStringArray("member-uuid")
			body.Overwrite = overwrite
			body.Redeploy = redeploy
			body.KeepReferences = keepRefs
			body.MemberUUIDs = members
		}
		resp, err := client.InjectProjectGroupSharedEnv(context.Background(), args[0], body, groupsWorkspaceOpts(cmd))
		if err != nil {
			return fmt.Errorf("inject shared env: %w", err)
		}
		if opts.Format == utils.OutputFormatJSON {
			return utils.PrintJSON(resp)
		}
		utils.PrintSuccess("Shared environment inject completed", opts)
		if resp != nil {
			utils.PrintTable([]string{"ATTRIBUTE", "VALUE"}, [][]string{
				{"Written Keys", strings.Join(resp.Data.WrittenKeys, ", ")},
				{"Skipped Keys", strings.Join(resp.Data.SkippedKeys, ", ")},
				{"Projects", strings.Join(resp.Data.ProjectsTouched, ", ")},
				{"Addons", strings.Join(resp.Data.AddonsTouched, ", ")},
				{"Redeploy Queued", strings.Join(resp.Data.RedeployQueued, ", ")},
				{"Message", displayOr(resp.Data.Message, resp.Message)},
			}, opts)
		}
		return nil
	},
	Args: cobra.ExactArgs(1),
}

var groupsConnectCmd = &cobra.Command{
	Use:   "connect <uuid>",
	Short: "Connect provider addon envs into a consumer project",
	RunE: func(cmd *cobra.Command, args []string) error {
		opts := utils.GetOutputOptions(cmd)
		client, err := rootClient(opts)
		if err != nil || client == nil {
			return err
		}
		body := &sdk.ConnectProjectGroupServicesRequest{}
		if jsonBody, _ := cmd.Flags().GetString("json-body"); jsonBody != "" {
			raw, err := os.ReadFile(jsonBody)
			if err != nil {
				return fmt.Errorf("read --json-body: %w", err)
			}
			if err := json.Unmarshal(raw, body); err != nil {
				return fmt.Errorf("parse --json-body: %w", err)
			}
		} else {
			consumerType, _ := cmd.Flags().GetString("consumer-type")
			consumerUUID, _ := cmd.Flags().GetString("consumer-uuid")
			providerType, _ := cmd.Flags().GetString("provider-type")
			providerUUID, _ := cmd.Flags().GetString("provider-uuid")
			variableSet, _ := cmd.Flags().GetString("variable-set")
			overwrite, _ := cmd.Flags().GetBool("overwrite")
			if consumerType == "" {
				consumerType = "project"
			}
			if providerType == "" {
				providerType = "addon_deployment"
			}
			providerType = mapMemberType(providerType)
			body = &sdk.ConnectProjectGroupServicesRequest{
				ConsumerType: consumerType,
				ConsumerUUID: consumerUUID,
				ProviderType: providerType,
				ProviderUUID: providerUUID,
				VariableSet:  variableSet,
				Overwrite:    overwrite,
			}
			if body.ConsumerUUID == "" || body.ProviderUUID == "" {
				return fmt.Errorf("--consumer-uuid and --provider-uuid are required (or use --json-body)")
			}
		}
		resp, err := client.ConnectProjectGroupServices(context.Background(), args[0], body, groupsWorkspaceOpts(cmd))
		if err != nil {
			return fmt.Errorf("connect services: %w", err)
		}
		if opts.Format == utils.OutputFormatJSON {
			return utils.PrintJSON(resp)
		}
		utils.PrintSuccess("Services connected", opts)
		if resp != nil {
			utils.PrintTable([]string{"ATTRIBUTE", "VALUE"}, [][]string{
				{"Written Keys", strings.Join(resp.Data.WrittenKeys, ", ")},
				{"Skipped Keys", strings.Join(resp.Data.SkippedKeys, ", ")},
				{"Restart Triggered", boolString(resp.Data.RestartTriggered)},
				{"Message", displayOr(resp.Data.Message, resp.Message)},
			}, opts)
		}
		return nil
	},
	Args: cobra.ExactArgs(1),
}

var groupsRedeployCmd = &cobra.Command{
	Use:   "redeploy <uuid>",
	Short: "Redeploy application members in a project group",
	RunE: func(cmd *cobra.Command, args []string) error {
		opts := utils.GetOutputOptions(cmd)
		client, err := rootClient(opts)
		if err != nil || client == nil {
			return err
		}
		resp, err := client.RedeployProjectGroupApps(context.Background(), args[0], groupsWorkspaceOpts(cmd))
		if err != nil {
			return fmt.Errorf("redeploy project group apps: %w", err)
		}
		if opts.Format == utils.OutputFormatJSON {
			return utils.PrintJSON(resp)
		}
		utils.PrintSuccess("Redeploy queued", opts)
		if resp != nil {
			utils.PrintTable([]string{"ATTRIBUTE", "VALUE"}, [][]string{
				{"Queued", strings.Join(resp.Data.Queued, ", ")},
				{"Failed", strings.Join(resp.Data.Failed, ", ")},
				{"Message", displayOr(resp.Data.Message, resp.Message)},
			}, opts)
		}
		return nil
	},
	Args: cobra.ExactArgs(1),
}

var groupsResolveCmd = &cobra.Command{
	Use:   "resolve",
	Short: "Resolve which group a member belongs to",
	RunE: func(cmd *cobra.Command, args []string) error {
		opts := utils.GetOutputOptions(cmd)
		client, err := rootClient(opts)
		if err != nil || client == nil {
			return err
		}
		memberType, err := normalizeGroupMemberType(cmd)
		if err != nil {
			return err
		}
		memberUUID, _ := cmd.Flags().GetString("member-uuid")
		resolveOpts := &sdk.ProjectGroupResolveOptions{
			MemberType: memberType,
			MemberUUID: memberUUID,
		}
		if workspace, _ := cmd.Flags().GetString("workspace"); workspace != "" {
			resolveOpts.WorkspaceUUID = workspace
		}
		resp, err := client.ResolveProjectGroupMember(context.Background(), resolveOpts)
		if err != nil {
			return fmt.Errorf("resolve member: %w", err)
		}
		if opts.Format == utils.OutputFormatJSON {
			return utils.PrintJSON(resp)
		}
		utils.PrintTable([]string{"ATTRIBUTE", "VALUE"}, [][]string{
			{"Group UUID", resp.Data.GroupUUID},
			{"Member Type", resp.Data.MemberType},
			{"Member UUID", resp.Data.MemberUUID},
		}, opts)
		return nil
	},
	Args: cobra.NoArgs,
}

var groupsCandidatesCmd = &cobra.Command{
	Use:   "candidates",
	Short: "List attachable projects and addons",
	RunE: func(cmd *cobra.Command, args []string) error {
		opts := utils.GetOutputOptions(cmd)
		client, err := rootClient(opts)
		if err != nil || client == nil {
			return err
		}
		candOpts := &sdk.ProjectGroupCandidatesOptions{}
		if workspace, _ := cmd.Flags().GetString("workspace"); workspace != "" {
			candOpts.WorkspaceUUID = workspace
		}
		if groupUUID, _ := cmd.Flags().GetString("group-uuid"); groupUUID != "" {
			candOpts.GroupUUID = groupUUID
		}
		resp, err := client.ListProjectGroupCandidates(context.Background(), candOpts)
		if err != nil {
			return fmt.Errorf("list candidates: %w", err)
		}
		if opts.Format == utils.OutputFormatJSON {
			return utils.PrintJSON(resp)
		}
		printCandidates := func(title string, items []sdk.ProjectGroupAttachCandidate) {
			if len(items) == 0 {
				return
			}
			utils.PrintInfo(title, opts)
			rows := make([][]string, 0, len(items))
			for _, item := range items {
				rows = append(rows, []string{
					item.MemberType,
					item.MemberUUID,
					item.Name,
					item.Status,
					item.CurrentGroupName,
					boolString(item.InTargetGroup),
				})
			}
			utils.PrintTable([]string{"TYPE", "UUID", "NAME", "STATUS", "CURRENT GROUP", "IN TARGET"}, rows, opts)
		}
		printCandidates("Projects", resp.Data.Projects)
		printCandidates("Addons", resp.Data.Addons)
		if len(resp.Data.Projects) == 0 && len(resp.Data.Addons) == 0 {
			utils.PrintWarning("No candidates found", opts)
		}
		return nil
	},
	Args: cobra.NoArgs,
}

func printProjectGroup(group *sdk.ProjectGroup, opts utils.OutputOptions) error {
	if group == nil {
		return fmt.Errorf("project group not found")
	}
	if opts.Format == utils.OutputFormatJSON {
		return utils.PrintJSON(group)
	}
	utils.PrintTable([]string{"ATTRIBUTE", "VALUE"}, [][]string{
		{"UUID", group.UUID},
		{"Name", group.Name},
		{"Slug", group.NameSlug},
		{"Workspace", group.WorkspaceUUID},
		{"Members", strconv.Itoa(group.MemberCount)},
		{"Default Cluster", group.DefaultClusterUUID},
		{"Default Environment", group.DefaultEnvironmentUUID},
		{"Created", group.CreatedAt},
		{"Updated", group.UpdatedAt},
	}, opts)
	if len(group.Members) > 0 {
		utils.PrintInfo("Members", opts)
		rows := make([][]string, 0, len(group.Members))
		for _, m := range group.Members {
			rows = append(rows, []string{m.MemberType, m.MemberUUID, m.Name, m.ServiceKind, m.Status})
		}
		utils.PrintTable([]string{"TYPE", "UUID", "NAME", "KIND", "STATUS"}, rows, opts)
	}
	return nil
}

func normalizeGroupMemberType(cmd *cobra.Command) (string, error) {
	raw, _ := cmd.Flags().GetString("type")
	if raw == "" {
		return "", fmt.Errorf("--type is required (project or addon)")
	}
	return mapMemberType(raw), nil
}

func mapMemberType(raw string) string {
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "addon", "addon_deployment", "addons":
		return "addon_deployment"
	case "project", "projects":
		return "project"
	default:
		return strings.ToLower(strings.TrimSpace(raw))
	}
}

func buildSharedEnvPutRequest(cmd *cobra.Command) (*sdk.UpsertProjectGroupSharedEnvRequest, error) {
	body := &sdk.UpsertProjectGroupSharedEnvRequest{}
	jsonBody, _ := cmd.Flags().GetString("json-body")
	filePath, _ := cmd.Flags().GetString("file")
	sets, _ := cmd.Flags().GetStringArray("set")

	path := jsonBody
	if path == "" {
		path = filePath
	}
	if path != "" {
		raw, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("read env file: %w", err)
		}
		// Prefer full request shape.
		if err := json.Unmarshal(raw, body); err == nil && (len(body.Variables) > 0 || strings.Contains(string(raw), "variables")) {
			// ok
		} else {
			// Fallback: map[string]string
			var m map[string]string
			if err := json.Unmarshal(raw, &m); err != nil {
				return nil, fmt.Errorf("parse env file: expected variables array or key/value object: %w", err)
			}
			body.Variables = make([]sdk.ProjectGroupSharedEnvVar, 0, len(m))
			for k, v := range m {
				body.Variables = append(body.Variables, sdk.ProjectGroupSharedEnvVar{Key: k, Value: v})
			}
		}
	}

	if len(sets) > 0 {
		for _, pair := range sets {
			key, value, ok := strings.Cut(pair, "=")
			key = strings.TrimSpace(key)
			if !ok || key == "" {
				return nil, fmt.Errorf("invalid --set %q; expected KEY=value", pair)
			}
			body.Variables = append(body.Variables, sdk.ProjectGroupSharedEnvVar{Key: key, Value: value})
		}
	}

	if len(body.Variables) == 0 && path == "" {
		return nil, fmt.Errorf("provide --set KEY=VAL, --file <json>, or --json-body <file>")
	}

	if cmd.Flags().Changed("inject") {
		body.Inject, _ = cmd.Flags().GetBool("inject")
	}
	if cmd.Flags().Changed("overwrite") {
		body.Overwrite, _ = cmd.Flags().GetBool("overwrite")
	}
	if cmd.Flags().Changed("redeploy") {
		body.Redeploy, _ = cmd.Flags().GetBool("redeploy")
	}
	if cmd.Flags().Changed("keep-references") {
		body.KeepReferences, _ = cmd.Flags().GetBool("keep-references")
	}
	return body, nil
}

func displayOr(primary, fallback string) string {
	if strings.TrimSpace(primary) != "" {
		return primary
	}
	return fallback
}

func addWorkspaceFlag(cmds ...*cobra.Command) {
	for _, c := range cmds {
		c.Flags().String("workspace", "", workspaceFlagHelp)
	}
}

func init() {
	addWorkspaceFlag(
		groupsListCmd, groupsGetCmd, groupsCreateCmd, groupsUpdateCmd, groupsDeleteCmd,
		groupsTopologyCmd, groupsMembersAttachCmd, groupsMembersDetachCmd,
		groupsEnvGetCmd, groupsEnvPutCmd, groupsEnvInjectCmd,
		groupsConnectCmd, groupsRedeployCmd, groupsResolveCmd, groupsCandidatesCmd,
	)

	groupsListCmd.Flags().Int("limit", 0, "Maximum number of groups to return")
	groupsListCmd.Flags().Int("offset", 0, "Offset for pagination")

	groupsCreateCmd.Flags().String("name", "", "Project group name")
	groupsCreateCmd.Flags().String("cluster-uuid", "", "Default cluster UUID")
	groupsCreateCmd.Flags().String("environment-uuid", "", "Default environment UUID")
	_ = groupsCreateCmd.MarkFlagRequired("name")

	groupsUpdateCmd.Flags().String("name", "", "Project group name")
	groupsUpdateCmd.Flags().String("cluster-uuid", "", "Default cluster UUID")
	groupsUpdateCmd.Flags().String("environment-uuid", "", "Default environment UUID")

	groupsDeleteCmd.Flags().Bool("yes", false, "Confirm project group deletion")

	for _, c := range []*cobra.Command{groupsMembersAttachCmd, groupsMembersDetachCmd, groupsResolveCmd} {
		c.Flags().String("type", "", "Member type: project or addon")
		c.Flags().String("member-uuid", "", "Member UUID")
		_ = c.MarkFlagRequired("type")
		_ = c.MarkFlagRequired("member-uuid")
	}
	groupsMembersAttachCmd.Flags().Bool("move", false, "Move member from another group if already attached")
	groupsMembersAttachCmd.Flags().Bool("include-session", false, "Include related session members")
	groupsMembersDetachCmd.Flags().Bool("include-session", false, "Include related session members")

	groupsEnvPutCmd.Flags().StringArray("set", nil, "Environment variable KEY=VAL; repeatable")
	groupsEnvPutCmd.Flags().String("file", "", "JSON file of variables (map or full request body)")
	groupsEnvPutCmd.Flags().String("json-body", "", "JSON file for full UpsertProjectGroupSharedEnvRequest")
	groupsEnvPutCmd.Flags().Bool("inject", false, "Inject after upsert")
	groupsEnvPutCmd.Flags().Bool("overwrite", false, "Overwrite existing keys on inject")
	groupsEnvPutCmd.Flags().Bool("redeploy", false, "Queue redeploy after inject")
	groupsEnvPutCmd.Flags().Bool("keep-references", false, "Keep references when upserting")

	groupsEnvInjectCmd.Flags().Bool("overwrite", false, "Overwrite existing keys")
	groupsEnvInjectCmd.Flags().Bool("redeploy", false, "Queue redeploy after inject")
	groupsEnvInjectCmd.Flags().Bool("keep-references", false, "Keep references")
	groupsEnvInjectCmd.Flags().StringArray("member-uuid", nil, "Limit inject to specific member UUIDs")
	groupsEnvInjectCmd.Flags().String("json-body", "", "JSON file for InjectProjectGroupSharedEnvRequest")

	groupsConnectCmd.Flags().String("consumer-type", "project", "Consumer type (default: project)")
	groupsConnectCmd.Flags().String("consumer-uuid", "", "Consumer project UUID")
	groupsConnectCmd.Flags().String("provider-type", "addon_deployment", "Provider type (addon or addon_deployment)")
	groupsConnectCmd.Flags().String("provider-uuid", "", "Provider addon deployment UUID")
	groupsConnectCmd.Flags().String("variable-set", "", "Optional variable set name")
	groupsConnectCmd.Flags().Bool("overwrite", false, "Overwrite existing connection env keys")
	groupsConnectCmd.Flags().String("json-body", "", "JSON file for ConnectProjectGroupServicesRequest")

	groupsCandidatesCmd.Flags().String("group-uuid", "", "Target group UUID for in-target markers")

	groupsEnvCmd.AddCommand(groupsEnvGetCmd, groupsEnvPutCmd, groupsEnvInjectCmd)
	groupsMembersCmd.AddCommand(groupsMembersAttachCmd, groupsMembersDetachCmd)

	groupsCmd.AddCommand(
		groupsListCmd,
		groupsGetCmd,
		groupsCreateCmd,
		groupsUpdateCmd,
		groupsDeleteCmd,
		groupsTopologyCmd,
		groupsMembersCmd,
		groupsEnvCmd,
		groupsConnectCmd,
		groupsRedeployCmd,
		groupsResolveCmd,
		groupsCandidatesCmd,
	)
	rootCmd.AddCommand(groupsCmd)
}
