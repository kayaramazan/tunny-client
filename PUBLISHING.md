# Publishing Tunny Client

This guide explains how to publish the Tunny client to GitHub and NPM.

## Prerequisites

1. GitHub account with `tunny` repository created
2. NPM account (npmjs.com)
3. GoReleaser installed (`brew install goreleaser`)

## Step 1: Push to GitHub

```bash
# Initialize repository (if not done)
git init
git add .
git commit -m "feat: initial release of tunny client"

# Add remote
git remote add origin https://github.com/kayaramazan/tunny.git

# Push to GitHub
git branch -M main
git push -u origin main
```

## Step 2: Create GitHub Release

```bash
# Tag the release
git tag -a v1.0.0 -m "Release v1.0.0"
git push origin v1.0.0

# Set GitHub token
export GITHUB_TOKEN="your_github_personal_access_token"

# Run GoReleaser (builds binaries and creates GitHub release)
goreleaser release --clean
```

This will:
- Build binaries for all platforms
- Create GitHub Release with binaries
- Generate checksums

## Step 3: Publish to NPM

```bash
cd npm

# Login to NPM (first time only)
npm login

# Publish
npm publish

# Or publish with tag
npm publish --tag latest
```

## Step 4: Verify

```bash
# Test NPM installation
npm install -g tunny-tunnel

# Test the CLI
tunny --help

# Test connection (requires running server)
tunny connect localhost:3000 --token dev --subdomain test
```

## Updating Versions

### For patch updates (1.0.0 → 1.0.1)

```bash
# Update version in npm/package.json
cd npm
npm version patch

# Commit and tag
git add .
git commit -m "chore: bump version to 1.0.1"
git tag -a v1.0.1 -m "Release v1.0.1"
git push origin main --tags

# Release
goreleaser release --clean

# Publish to NPM
cd npm
npm publish
```

### For minor updates (1.0.0 → 1.1.0)

```bash
cd npm
npm version minor
# ... same process as above
```

### For major updates (1.0.0 → 2.0.0)

```bash
cd npm
npm version major
# ... same process as above
```

## Troubleshooting

### GoReleaser fails

```bash
# Check config
goreleaser check

# Test without publishing
goreleaser release --snapshot --skip=publish --clean
```

### NPM publish fails

```bash
# Check if logged in
npm whoami

# Check package name availability
npm view tunny-tunnel

# Login again
npm login
```

## Automation with GitHub Actions

Create `.github/workflows/release.yml` for automated releases:

```yaml
name: Release

on:
  push:
    tags:
      - 'v*'

jobs:
  release:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - uses: actions/setup-go@v4
        with:
          go-version: '1.22'
      
      - uses: goreleaser/goreleaser-action@v5
        with:
          version: latest
          args: release --clean
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
      
      - uses: actions/setup-node@v3
        with:
          node-version: '18'
          registry-url: 'https://registry.npmjs.org'
      
      - name: Publish to NPM
        run: |
          cd npm
          npm publish
        env:
          NODE_AUTH_TOKEN: ${{ secrets.NPM_TOKEN }}
```

## Notes

- Never commit `GITHUB_TOKEN` or `NPM_TOKEN` to Git
- Store secrets in GitHub repository settings
- Test releases with `--snapshot` flag first
- Always update CHANGELOG.md before releases
