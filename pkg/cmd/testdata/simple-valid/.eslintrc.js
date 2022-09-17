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
