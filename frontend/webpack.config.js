const HtmlWebpackPlugin = require('html-webpack-plugin');
const MiniCssExtractPlugin = require('mini-css-extract-plugin')
const VueLoaderPlugin = require('vue-loader/lib/plugin');

module.exports = {
    entry: './src/main.js',
    output: {
        filename: 'bundle.js'
    },
    module: {
        rules: [
            {test: /\.js$/, use: 'babel-loader'},
            {test: /\.vue$/, use: 'vue-loader'},
            {test: /\.css$/, use: ['vue-style-loader', MiniCssExtractPlugin.loader, 'css-loader']},
            {
                test: /\.scss$/, use: ['vue-style-loader', MiniCssExtractPlugin.loader, 'css-loader', {
                    loader: 'sass-loader',
                    options: {
                        additionalData: '@import "~bulma/sass/utilities/initial-variables";'
                    }
                }]
            },
            {test: /\.svg$/, use: 'file-loader'},
            {test: /\.(woff|woff2|eot|ttf|otf)$/, use: 'file-loader'}
        ]
    },
    plugins: [
        new HtmlWebpackPlugin({
            template: './src/index.html',
        }),
        new MiniCssExtractPlugin({
            filename: "styles.css"
        }),
        new VueLoaderPlugin(),
    ]
};