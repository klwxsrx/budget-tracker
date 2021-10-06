module.exports = {
  extends: [
    'eslint:recommended',
    'plugin:vue/base',
    'plugin:vue/essential',
    'plugin:vue/strongly-recommended',
    'plugin:vue/recommended',
  ],
  rules: {
    'comma-dangle': ['warn', 'always-multiline'],
    'vue/html-closing-bracket-spacing': ['warn', {
      'startTag': 'never',
      'endTag': 'never',
      'selfClosingTag': 'never',
    }],
    'vue/max-attributes-per-line': ['warn', {
      'singleline': {
        'max': 4,
        'allowFirstLine': true,
      },
      'multiline': {
        'max': 1,
        'allowFirstLine': false,
      },
    }],
    'vue/html-indent': ['warn', 2, {
      'attribute': 2,
      'baseIndent': 1,
      'closeBracket': 0,
      'alignAttributesVertically': true,
      'ignores': [],
    }],
  },
  env: {
    'browser': true,
    'node': true,
  },
}