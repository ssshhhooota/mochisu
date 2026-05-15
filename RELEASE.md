# Release

Release procedure and branch naming convention for this project.

## Branch naming

```
release-v<MAJOR>_<MINOR>_<PATCH>
```

Examples:

```
release-v0_1_0
release-v1_2_3
```

- Use `_` as separator (follows the git-branch rule)
- One branch = one release. Do not reuse
- Delete the branch after the release completes

The tag is derived by simple `_` → `.` substitution:

```
release-v0_1_0  →  v0.1.0
```

### Pre-releases

When needed, add a suffix such as `_rc_N`:

```
release-v0_1_0_rc_1   →  v0.1.0-rc.1
release-v0_1_0_beta_2 →  v0.1.0-beta.2
```

## Tags

Follows semver. The `v` prefix is required.

```
v0.1.0
v1.2.3
v0.1.0-rc.1
```

## Release procedure

```bash
# Update main
git switch main
git pull --rebase origin main

# Create and push the release branch
git switch -c release-v0_1_0 origin/main --no-track
git push origin release-v0_1_0
```

The push triggers GitHub Actions, which:

1. Extracts the version from the branch name (`release-v0_1_0` → `v0.1.0`)
2. Creates and pushes the tag (`v0.1.0`)
3. Runs goreleaser to build binaries, create the GitHub Release, and update the Homebrew tap
4. Deletes the release branch

## Rollback

If you need to abort before the tag is pushed, delete the remote branch:

```bash
git push origin --delete release-v0_1_0
```

After the tag is pushed, manually delete the tag and the GitHub Release, then bump the version and re-release. Reusing the same version is not allowed.
