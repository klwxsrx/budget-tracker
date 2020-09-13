module.exports = {
    extends: [
        "eslint:recommended",
        "plugin:vue/base",
        "plugin:vue/essential",
        "plugin:vue/strongly-recommended",
        "plugin:vue/recommended"
    ],
    rules: {
        indent: ["error", 2],
        semi: ["error", "never"],
        quotes: ["error", "single"],
        "comma-dangle": ["error", "always-multiline"]
    }
}