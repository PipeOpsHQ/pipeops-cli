# PipeOps CLI - Railway-Style Developer Experience TODO

## CLI Features

### Core Commands

1. **Project Management**
   - [ ] `pipeops list` - List all projects
   - [ ] `pipeops link <project-id>` - Associate project with directory
   - [ ] `pipeops status` - Show current project info
   - [ ] `pipeops unlink` - Disconnect from project

2. **Service Management**
   - [ ] `pipeops addons <service-name>` - Link service to project
   - [ ] `pipeops addons list` - List all services

<!-- 3. **Deployment & Environment**
   - [ ] `pipeops up [PATH]` - Deploy from directory
   - [ ] `pipeops environment <env-name>` - Switch environments
   - [ ] `pipeops environment new <name>` - Create environments
   - [ ] `pipeops up --detach` - Deploy without logs -->

4. **Logs & Variables**
   - [ ] `pipeops logs` - Stream logs
   - [ ] `pipeops logs --project <project-name>` - Project logs
   - [ ] `pipeops logs --service <deployed-addon-name>` - Addons logs
   - [ ] `pipeops variables` - Show environment variables
   - [ ] `pipeops variables --set "KEY=value"` - Set variables

5. **Database & Local Development**
   - [ ] `pipeops connect [SERVICE_NAME]` - Connect to databases (PostgreSQL, MySQL, Redis, MongoDB)
   - [ ] `pipeops run [COMMAND]` - Run commands with Railway variables
   - [ ] `pipeops shell` - Interactive shell with variables
   <!-- - [ ] `pipeops ssh [COMMAND]` - SSH to services -->

6. **Domain & Operations**
   <!-- - [ ] `pipeops domain [DOMAIN]` - Manage domains -->
   - [ ] `pipeops redeploy --project <project-name>` - Redeploy
   - [ ] `pipeops down` - Rollback deployment
   - [ ] `pipeops open` - Open dashboard

---

## 🗺️ 6-Week Implementation Plan

**Phase 1 (Weeks 1): Core commands**
- **Week 1:** Project & service management
- **Week 2:** Deployment & environments
- **Week 3:** Logs & variables
- **Week 4:** Database connections & local development

**Phase 2 (Weeks 2): Advanced features**
- **Week 5:** Domain management & operations
- **Week 6:** Polish & testing

---

> This plan is strictly focused on providing the basic developer experience that [Railway CLI](https://docs.railway.com/reference/cli-api) offers. All features and commands should be implemented to match or closely mirror the Railway CLI workflow and user expectations.
