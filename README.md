## Setting Up Git Hooks

To ensure consistent code formatting, linting, and testing, please set up the Git hooks by running the following command:

**NB**: Make sure that Golint & Gosec are installed

````bash
go install github.com/securego/gosec/v2/cmd/gosec@latest
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
````
**NB**: Make sure that go/bin path is added to your profile

```bash
chmod +x .githooks/pre-commit

chmod +x install-hooks.sh

./install-hooks.sh
```

