# PipeOps CLI TODO

## CLI Features

### Currently Working
- [x] **Authentication**: `pipeops auth login/logout/status` - OAuth2 with PKCE flow
- [x] **User Management**: `pipeops auth me/debug` - User info and debugging
- [x] **Basic Project Listing**: `pipeops project list` - List projects (basic implementation)
- [x] **Project Linking**: `pipeops link <project-id>` - Link project to directory
- [x] **Addon Deployment**: `pipeops deploy --addon <addon-id>` - Deploy addons to projects
- [x] **Server Management**: `pipeops server list` - List servers
- [x] **Proxy Management**: `pipeops proxy start/stop/status` - Local proxy connections
- [x] **Status Monitoring**: `pipeops status` - Show current status
- [x] **Version Management**: `pipeops version` - Show CLI version

### 🚧 Needs Implementation

1. **Project Management**
   - [ ] `pipeops list` - List all projects (enhance current `project list`)
   - [ ] `pipeops status` - Show current project info (enhance current status)
   - [ ] `pipeops unlink` - Disconnect from project

2. **Service/Addon Management**
   - [ ] `pipeops addons list` - List all available addons
   - [ ] `pipeops addons <addon-name>` - Link addon to project
   - [ ] `pipeops addons status` - Show addon status and details

3. **Logs & Variables**
   - [ ] `pipeops logs` - Stream logs (enhance current logs)
   - [ ] `pipeops logs --project <project-name>` - Project logs
   - [ ] `pipeops logs --addon <addon-name>` - Addon logs
   - [ ] `pipeops variables` - Show environment variables
   - [ ] `pipeops variables --set "KEY=value"` - Set variables

4. **Database & Local Development**
   - [ ] `pipeops connect [ADDON_NAME]` - Connect to databases (PostgreSQL, MySQL, Redis, MongoDB)
   - [ ] `pipeops run [COMMAND]` - Run commands with PipeOps variables
   - [ ] `pipeops shell` - Interactive shell with variables

5. **Project Operations**
   - [ ] `pipeops redeploy --project <project-name>` - Redeploy project
   - [ ] `pipeops down` - Rollback deployment
   - [ ] `pipeops open` - Open dashboard in browser

---

## 🗺️ Implementation Plan

**Phase 1 (Week 1): Enhance Existing Features**
- **Enhance project listing**: Improve `pipeops project list`
- **Enhance status command**: Improve `pipeops status` to show current project info
- **Add unlink command**: `pipeops unlink` to disconnect from project
- **Enhance logs**: Improve `pipeops logs` with better filtering and real-time streaming

**Phase 2 (Week 2): Add Missing Features**
- **Addon management**: `pipeops addons list` and `pipeops addons <name>`
- **Variables system**: `pipeops variables` and `pipeops variables --set`
- **Database connections**: `pipeops connect [ADDON_NAME]` for database access
- **Local development**: `pipeops run [COMMAND]` and `pipeops shell`

**Phase 3 (Week 3): Polish & Testing**
- **Project operations**: `pipeops redeploy`, `pipeops down`, `pipeops open`
- **Auto-completion**: Command and addon name completion
- **Better error handling**: User-friendly error messages
- **Unit tests**: Test all commands thoroughly

---

> This plan is strictly focused on providing the basic developer experience that [Railway CLI](https://docs.railway.com/reference/cli-api) offers
