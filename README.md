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
* Open a GitHub Pull Request (requires [GitHub CLI](https://cli.github.com/) to be installed and configured).

## Examples

### Quick start

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
    source:
      path: ../shared-configs-repo/.eslintrc.js.tpl
```

4️⃣ Finally, run goplicate on the repository to sync any updates:

<img src="https://github.com/ilaif/goplicate/raw/main/assets/goplicate-run.gif" width="700">

### Using a remote git repository

In this example, we'll use a remote git repository as the source of the shared snippets instead of a local folder.

1️⃣ Fork [goplicate-example-repo-1](https://github.com/ilaif/goplicate-example-repo-1) and clone it.

Looking inside, we see 2 files:

`.eslintrc.js`:

```js
module.exports = {
  extends: 'eslint:recommended',
  rules: {
    // goplicate-start:common-rules
    // enable additional rules
    indent: ['error', 4],
    'linebreak-style': ['error', 'unix'],
    quotes: ['error', 'double'],
    semi: ['error', 'always'],
    // goplicate-end:common-rules

    // override configuration set by extending "eslint:recommended"
    'no-empty': 'warn',
    'no-cond-assign': ['error', 'always'],
  },
}
```

> Notice the `// goplicate-start:common-rules` and `// goplicate-end:common-rules` block annotations that will be synced by goplicate.

`.goplicate.yaml`:

```yaml
targets:
  - path: .eslintrc.js
    source:
      repository: https://github.com/ilaif/goplicate-example-shared-configs
      path: .eslintrc.js
    params:
      - repository: https://github.com/ilaif/goplicate-example-shared-configs
        path: params.yaml
```

If we go to [goplicate-example-shared-configs](https://github.com/ilaif/goplicate-example-shared-configs), we'll see that `.eslintrc.js` contains the `common-rules` source of truth with the `params.yaml` containing a parameter as well:

`.eslintrc.js` in [goplicate-example-shared-configs](https://github.com/ilaif/goplicate-example-shared-configs):

```js
module.exports = {
  rules: {
    // goplicate-start:common-rules
    // enable additional rules
    indent: ['error', {{.indent}}],
    'linebreak-style': ['error', 'unix'],
    quotes: ['error', 'double'],
    semi: ['error', 'always'],
    // goplicate-end:common-rules
  },
}
```

`params.yaml` in [goplicate-example-shared-configs](https://github.com/ilaif/goplicate-example-shared-configs):

```yaml
indent: 2
```

2️⃣ From the cloned repository, run goplicate to create a new PR with synced changes:

```sh
~/git/oss/goplicate-example-repo-1 (main ✔) ᐅ goplicate run --publish
• Cloning 'https://github.com/ilaif/goplicate-example-shared-configs'
• Target '.eslintrc.js': Block 'common-rules' needs to be updated. Diff:
     // goplicate-start:common-rules
     // enable additional rules
-    indent: ['error', 4],
+    indent: ['error', 2],
     'linebreak-style': ['error', 'unix'],
     quotes: ['error', 'double'],
     semi: ['error', 'always'],
...

? Do you want to apply the above changes? Yes
• Target '.eslintrc.js': Updated
? Do you want to publish the above changes? Yes
• Publishing changes...
• Created PR: https://github.com/ilaif/goplicate-example-repo-1/pull/4
```

3️⃣ Open the PR and review it! We'll see that the indentation indeed changed:

```diff
  rules: {
    // goplicate-start:common-rules
    // enable additional rules
-   indent: ['error', 4],
+   indent: ['error', 2],
    'linebreak-style': ['error', 'unix'],
    quotes: ['error', 'double'],
    semi: ['error', 'always'],
```

4️⃣ Finally, merge it and maintain consistency and standardization across your configuration files!

### More examples

See the [Examples](https://github.com/ilaif/goplicate/tree/main/examples) folder for usage examples.

## Questions, bug reporting and feature requests

You're more than welcome to [Create a new issue](https://github.com/ilaif/goplicate/issues/new) or contribute.

## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change. See [CONTRIBUTING.md](CONTRIBUTING.md) for more information.

## License

Goplicate is licensed under the [MIT](https://choosealicense.com/licenses/mit/) license. For more information, please see the [LICENSE](LICENSE) file.
