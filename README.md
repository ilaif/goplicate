# Goplicate

<img src="https://github.com/ilaif/goplicate/raw/main/assets/logo.png" width="700">

---

Goplicate is a CLI tool that helps define common code or configuration snippets once and sync them to multiple projects.

## Why and how

In cases where we have many snippets that are repeated between different repositories or projects, it becomes a real hassle to keep them up-to-date.

We want to stay [DRY](https://en.wikipedia.org/wiki/Don%27t_repeat_yourself).

Goplicate achieves that by defining "blocks" around such shared snippets and automates their update via a shared source that contains the most up-to-date version of those snippets.

## Installation

### MacOS

```sh
brew install ilaif/tap/goplicate
brew upgrade ilaif/tap/goplicate
```

### Install from source

```sh
go install github.com/ilaif/goplicate/cmd/goplicate@latest
```

## Usage

`goplicate --help`

## Design principles

* 🌵 Stay [DRY](https://en.wikipedia.org/wiki/Don%27t_repeat_yourself) - Write a configuration once, and have it synced across many projects.
* 🤤 [Keep It Stupid Simple (KISS)](https://en.wikipedia.org/wiki/KISS_principle) - Treat configuration snippets as simple text, not assuming anything about structure.
* 🙆🏻‍♀️ Allow flexibility, but not too much - Allow syncing whole files, or parts of them (currently, line-based).
* 😎 Automate all the things - After an initial configuration, automates the rest.

## Features

* Configure line-based blocks that should be synced across multiple projects and files.
* See comfortable diffs while updating config files.
* Template support using [Go Templates](https://pkg.go.dev/text/template) with dynamic parameters or conditions.
* Sync multiple repositories with a single command.
* Automatically run post hooks to validate that the updates worked well before opening a pull request.
* Open a GitHub Pull Request (assuming [GitHub CLI](https://cli.github.com/) is installed and configured).

## Quick start example

In the following simplified example, we'll sync an [eslint](https://eslint.org) configuration.

We'll end up having the following folder structure:

```diff
+ shared-configs-repo/
+   .eslintrc.js.tpl
  repo-1/
    .eslintrc.js
+   .goplicate.yaml
  repo-2/
    .eslintrc.js
+   .goplicate.yaml
  ...
```

1️⃣ Choose a config file that some of its contents are copied across multiple projects, and add goplicate block comments for the `common-rules` section of your desire:

`repo-1/.eslintrc.js`:

```diff
module.exports = {
    "extends": "eslint:recommended",
    "rules": {
+       // goplicate-start:common-rules
        // enable additional rules
        "indent": ["error", 2],
        "linebreak-style": ["error", "unix"],
        "quotes": ["error", "double"],
        "semi": ["error", "always"],
+       // goplicate-end:common-rules

        // override configuration set by extending "eslint:recommended"
        "no-empty": "warn",
        "no-cond-assign": ["error", "always"],
    }
}
```

2️⃣ Create a separate, centralized repository to manage all of the shared config files. We'll name it `shared-configs-repo`. Then, add an `.eslintrc.js.tpl` file with the `common-rules` snippet that we want to sync:

`shared-configs-repo/.eslint.js.tpl`:

```txt
module.exports = {
     "rules": {
          // goplicate-start:common-rules
          // enable additional rules
          "indent": ["error", 4],
          "linebreak-style": ["error", "unix"],
          "quotes": ["error", "double"],
          "semi": ["error", "always"],
          // goplicate-end:common-rules
    }
}
```

> Goplicate snippets are simply the sections of the config file that we'd like to sync. In this example, we've also added the surrounding configuration to make it more readable, but it's not really needed.

3️⃣ Go back to the original project, and create a `.goplicate.yaml` file in your project root folder:

`repo-1/.goplicate.yaml`:

```yaml
targets:
  - path: .eslintrc.js
    source: ../shared-configs-repo/.eslintrc.js.tpl
```

4️⃣ Finally, run goplicate on the repository to sync any updates:

<img src="https://github.com/ilaif/goplicate/raw/main/assets/goplicate-run.gif" width="700">

## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change. See [CONTRIBUTING.md](CONTRIBUTING.md) for more information.

## License

Goplicate is licensed under the [MIT](https://choosealicense.com/licenses/mit/) license. For more information, please see the [LICENSE](LICENSE) file.
