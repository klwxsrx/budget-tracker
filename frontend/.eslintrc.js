module.exports = {
    extends: [
        "eslint:recommended",
        "plugin:vue/base",
        "plugin:vue/essential",
        "plugin:vue/strongly-recommended",
        "plugin:vue/recommended"
    ],
    rules: {
        indent: ["warn", 4],
        semi: ["warn", "never"],
        quotes: ["warn", "single"],
        "comma-dangle": ["warn", "always-multiline"],
    }
}