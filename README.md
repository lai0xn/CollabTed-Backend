Here is a basic contribution guide for your startup code that includes all the steps you mentioned:

## Contribution Guide
### 1. Clone the Repository
Clone the repository to your local machine using SSH or HTTPS:
```bash
git clone https://github.com/isif00/CollabTed-Backend.git
cd CollabTed-Backend
```

### 2. Create a Feature Branch
Create a new branch for your feature or bug fix:
```bash
git checkout -b feature/your-feature-name
```

### 3. Work on Your Feature
Make your changes and commit them:
```bash
git add --all
git commit -m "[feat] Describe your feature or fix"
```
**NB:** make sure you follow this commiting messages form:

- [feat]: A new feature
- [fix]: A bug fix
- [docs]: Documentation changes
- [style]: Code style changes (formatting, missing semi colons, etc.)
- [refactor]: Code refactoring without changing functionality
- [perf]: Performance improvements
- [test]: Adding missing tests or correcting existing tests
- [chore]: Other changes that don't modify src or test files

### 4. Pull Latest Changes from Main
Ensure your feature branch is up to date with the latest changes from the `main` branch:
```bash
git checkout main
git pull origin main
```

### 5. Rebase Your Feature Branch
Rebase your feature branch on top of the latest `main` changes:
```bash
git checkout feature/your-feature-name
git rebase main
```

### 6. Resolve Conflicts
If there are any conflicts, resolve them and continue the rebase:
```bash
# Resolve conflicts in your editor
git add --all
git rebase --continue
```

### 7. Rebase Main Branch
Rebase the `main` branch on top of your feature branch:
```bash
git checkout main
git rebase feature/your-feature-name
```

### 8. Open a Pull Request
Push your changes and open a pull request:
```bash
git push origin feature/your-feature-name
```
Then go to the repository on GitHub and open a pull request against the `main` branch.

### 9. Wait for Approval
Wait for your pull request to be reviewed and approved by a maintainer.

## Setting Up Git Hooks

To ensure consistent code formatting, linting, and testing, please set up the Git hooks by running the following command:

**NB**: Make sure that Golint & Gosec are installed

```bash
go install github.com/securego/gosec/v2/cmd/gosec@latest
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```

**NB**: Make sure that `$HOME/go/bin` path is added to your profile

```bash
chmod +x .githooks/pre-commit
chmod +x install-hooks.sh
./install-hooks.sh
```

By following these steps, we ensure a smooth development process.
