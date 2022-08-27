# Goplicate

Goplicate is a CLI tool that helps define common code or configuration snippets once, and sync it to multiple projects.

## Why and how

In cases where we have many snippets that are repeated between different repositories in the same project or across projects, it becomes a real hassle to keep them up-to-date.

We want to stay [DRY](https://en.wikipedia.org/wiki/Don%27t_repeat_yourself).

Goplicate achieves that by defining "blocks" around such shared snippets and automates their update via a shared source that contains the most up-to-date version of those snippets.

## Design principles

- Keep it simple - Treat snippets as text, not assuming anything about structure or correctness.

## Example use case

Let's say that we have a common configuration that we need to maintain for multiple projects. In this example, we'll use an imaginary `.pre-commit-config.yaml` ([https://pre-commit.com](pre-commit.com)):

Project 1:

```yaml
repos:
  - repo: https://github.com/some/repo
    rev: v1.2.3
    hooks:
      - id: my-common-pre-commit-hook
  - repo: local
    hooks:
      - id: my-project-1-pre-commit-hook
```

Project 2:

```yaml
repos:
  - repo: https://github.com/some/repo
    rev: v1.2.3
    hooks:
      - id: my-common-pre-commit-hook
  - repo: local
    hooks:
      - id: my-project-2-pre-commit-hook
```

If we have many such projects that we have to maintain the common pre-commit hook for, it starts getting messy. Goplicate comes to the rescue!

For each project, we can add a `goplicate` comment that will denote a section as managed by Goplicate:

```yaml
repos:
  # goplicate(name=common,pos=start)
  - repo: https://github.com/some/repo
    rev: v1.2.3
    hooks:
      - id: my-common-pre-commit-hook
  # goplicate(name=common,pos=end)
  - repo: local
    hooks:
      - id: my-project-1-pre-commit-hook
```

Then, initialize a new `.goplicate.yaml` file in each of the projects:

```yaml
source: ../goplicate
targets:
  - path: .pre-commit-hooks.yaml
    source: pre-commit-common.yaml
```

Where:

- `source` is the path to the Goplicate repository containing the common configuration that we're going to define in the next step.
- `targets` is a list of configurations to apply to `path` from `source` inside the common Goplicate repository.

Finally, define the Goplicate repository in the `source` path (in our example `../goplicate`) with the following file inside:

.pre-commit-common.yaml:

```yaml
  # goplicate(name=common,pos=start)
  - repo: https://github.com/some/repo
    rev: v1.3.0
    hooks:
      - id: my-common-pre-commit-hook
  # goplicate(name=common,pos=end)
```

Now, if we run `goplicate run` from one of our defined projects, we'll see that `rev` was changed from `v1.2.4` to `v1.3.0`
