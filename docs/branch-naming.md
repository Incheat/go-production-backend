# Git Branch Naming & Merge Strategy (Go Microservices)

## Principles

-   Clear intent from branch name
-   Machine-friendly (CI/CD, automation)
-   Consistency over personal preference
-   Short but descriptive
-   Traceable to issues/tickets

------------------------------------------------------------------------

## Main Branches

-   `main`\
    Production-ready code.\
    Protected branch, PRs only.

-   `develop` (optional)\
    Integration branch for larger teams.

------------------------------------------------------------------------

## Branch Types & Naming

### Feature

New functionality, APIs, or services.

    feature/<service>-<short-description>

Examples: - `feature/user-add-jwt-auth` -
`feature/order-create-endpoint`

------------------------------------------------------------------------

### Bugfix

Non-critical bug fixes.

    bugfix/<service>-<issue>

Examples: - `bugfix/user-fix-nil-pointer`

------------------------------------------------------------------------

### Hotfix

Critical production fixes (branched from `main`).

    hotfix/<service>-<issue>

Examples: - `hotfix/payment-double-charge`

------------------------------------------------------------------------

### Docs

Documentation only (Markdown, guides, ADRs).

    docs/<scope>-<description>

Examples: - `docs/readme-update` - `docs/payment-runbook`

------------------------------------------------------------------------

### Chore

Tooling, dependencies, repo maintenance.

    chore/<scope>-<description>

Examples: - `chore/add-golangci-lint` - `chore/bump-go-1.22`

------------------------------------------------------------------------

### CI / CD

Pipeline, automation, deployment config.

    ci/<description>

Examples: - `ci/add-go-test-matrix` - `ci/optimize-docker-build`

------------------------------------------------------------------------

### Release

Release preparation and stabilization.

    release/v1.4.2
    release/v2.0.0-rc1

------------------------------------------------------------------------

### Experiment

Proof of concepts or spikes.

    experiment/<description>

------------------------------------------------------------------------

## Recommended Naming Template

    <branch-type>/<scope>-<short-description>

------------------------------------------------------------------------

## Merge Strategy (Recommended)

### Merge Commit (Default)

-   Best for `feature/*`, `bugfix/*`, `hotfix/*`
-   Preserves history and context
-   Preferred for microservices teams

### Squash Merge

-   Recommended for `docs/*`, `chore/*`, `ci/*`
-   Keeps `main` clean and readable

### Rebase

-   Allowed only on personal feature branches
-   Use before opening PR
-   Never rebase shared branches

------------------------------------------------------------------------

## Summary Rules

-   Rebase to clean **your own branch**
-   Merge commit for core functionality
-   Squash non-functional changes
-   Never force-push shared branches

------------------------------------------------------------------------

## Example Convention Block

    main              # production
    develop           # integration (optional)

    feature/*
    bugfix/*
    hotfix/*
    docs/*
    chore/*
    ci/*
    release/*
    experiment/*
