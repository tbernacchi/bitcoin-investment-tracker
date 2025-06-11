# ArgoCD Configuration

> This directory contains the ArgoCD application configuration for the Bitcoin Investment Tracker.

## Prerequisites

Follow the installation guides in [home-lab](https://github.com/tbernacchi/home-lab) repository:
1. [ArgoCD Installation](https://github.com/tbernacchi/home-lab?tab=readme-ov-file#argo-cd)
2. [ArgoCD Image Updater Setup](https://github.com/tbernacchi/home-lab?tab=readme-ov-file#argo-image-updater)
3. Docker Hub registry access
4. Git repository access with write permissions

## Configuration Steps

### 1. Docker Hub Secret

Create a secret for Docker Hub authentication:
```bash
kubectl create secret docker-registry dockerhub-secret \
  --namespace bitcoin-investment-tracker \
  --docker-server=docker.io \
  --docker-username=ambrosiaaaaa \
  --docker-password=<your-dockerhub-token>
```

### 2. Git Credentials

Create a secret for Git repository access:
```bash
kubectl create secret generic git-creds \
  --namespace argocd \
  --from-literal=username=<your-git-username> \
  --from-literal=password=<your-git-pat>
```

## Application Configuration

The [`application.yaml`](./application.yaml) contains several important annotations:

### Image Update Configuration
```yaml
argocd-image-updater.argoproj.io/image-list: "myalias=ambrosiaaaaa/bitcoin-investment-tracker"
```
- Configures which image to track for updates
- `myalias` is a reference name used in other annotations

### Update Strategy
```yaml
argocd-image-updater.argoproj.io/myalias.update-strategy: "semver"
```
- Uses semantic versioning for image updates
- Ensures proper version ordering

### Docker Hub Authentication
```yaml
argocd-image-updater.argoproj.io/myalias.pull-secret: "dockerhub-secret"
```
- References the Docker Hub secret created earlier
- Required for private registry access

### Git Write-back Configuration
```yaml
argocd-image-updater.argoproj.io/write-back-method: "git:secret:argocd/git-creds"
```
- Configures how updated manifests are written back to Git
- Uses the git-creds secret for authentication

```yaml
argocd-image-updater.argoproj.io/write-back-target: "kustomization"
```
- Updates will be written to the kustomization.yaml file
- **Note**: Remember to `git pull` regularly as ArgoCD will commit changes directly to the repository

## Sync Policy

The application is configured with automated sync:
```yaml
syncPolicy:
  automated:
    prune: true    # Remove resources that are no longer defined
    selfHeal: true # Revert manual changes in the cluster
```

## Monitoring

To check the status of image updates:
```bash
# Get image updater logs
kubectl logs -n argocd -l app.kubernetes.io/name=argocd-image-updater

# Check application sync status
argocd app get bitcoin-investment-tracker

# View application events
kubectl get events -n argocd --field-selector involvedObject.name=bitcoin-investment-tracker
```

## Troubleshooting

Common issues and solutions:

1. **Image not updating:**
   - Check Docker Hub credentials
   - Verify image tag format matches update strategy
   - Check image updater logs

2. **Sync failures:**
   - Verify Git credentials
   - Check if target branch is protected
   - Review ArgoCD application events

3. **Write-back failures:**
   - Ensure git-creds secret is correctly configured
   - Check repository permissions
   - Verify write-back target path exists

## References

### ArgoCD
- [ArgoCD Official Documentation](https://argo-cd.readthedocs.io/en/stable/)
- [Declarative Setup](https://argo-cd.readthedocs.io/en/stable/operator-manual/declarative-setup/)
- [Application Specification](https://argo-cd.readthedocs.io/en/stable/operator-manual/application.yaml)
- [Sync Options](https://argo-cd.readthedocs.io/en/stable/user-guide/sync-options/)

### ArgoCD Image Updater
- [Image Updater Documentation](https://argocd-image-updater.readthedocs.io/en/stable/)
- [Installation Guide](https://argocd-image-updater.readthedocs.io/en/stable/install/installation/)
- [Update Strategies](https://argocd-image-updater.readthedocs.io/en/stable/configuration/update-strategies/)
- [Write-Back Configuration](https://argocd-image-updater.readthedocs.io/en/stable/configuration/write-back/)
- [Application Configuration Examples](https://argocd-image-updater.readthedocs.io/en/stable/configuration/applications/)

### Specific Features Used
- [Image Update Strategy - Semver](https://argocd-image-updater.readthedocs.io/en/stable/configuration/update-strategies/#semantic-version-update-strategy)
- [Git Write-Back Method](https://argocd-image-updater.readthedocs.io/en/stable/configuration/write-back/#git)
- [Pull Secrets Configuration](https://argocd-image-updater.readthedocs.io/en/stable/configuration/registries/#using-pull-secrets)
- [Kustomize Integration](https://argocd-image-updater.readthedocs.io/en/stable/configuration/applications/#kustomize) 
