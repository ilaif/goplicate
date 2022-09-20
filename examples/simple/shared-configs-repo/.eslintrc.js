// A sample .eslintrc.js "source" (shared) configuration,
// where only the sections that are surrounded by
// goplicate-start/end comments are synced.
//
// Note the '{{.indent}}' template variable, which will be
// replaced with the value from params.yaml when running goplicate.

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
