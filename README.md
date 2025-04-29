# GitLab LLM Agent Monorepo (Go)

This monorepo includes lightweight agents for GitLab automation using local LLMs and label-driven workflows:

- `approval-agent`: Approves merge requests based on labels
- `reasoning-labeler`: Uses an LLM to suggest auto-approve or manual-review
- `security-labeler`: Compares Docker image versions and applies security labels

Includes:
- Dockerfiles
- Helm charts
- Kubernetes manifests
- Shared GitLab client
