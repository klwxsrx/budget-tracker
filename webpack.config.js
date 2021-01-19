const HtmlWebpackPlugin = require('html-webpack-plugin')
const MiniCssExtractPlugin = require('mini-css-extract-plugin')
const VueLoaderPlugin = require('vue-loader/lib/plugin')

module.exports = {
  entry: './web/main.js',
  output: {
    filename: 'bundle.js',
  },
  devtool: 'source-map',
  devServer: {
    contentBase: __dirname,
  },
  module: {
    rules: [
      {test: /\.js$/, use: 'babel-loader'},
      {test: /\.vue$/, use: 'vue-loader'},
      {test: /\.css$/, use: ['vue-style-loader', MiniCssExtractPlugin.loader, 'css-loader']},
      {test: /\.scss$/, use: ['vue-style-loader', MiniCssExtractPlugin.loader, 'css-loader', 'sass-loader']},
      {test: /\.svg$/, use: 'file-loader'},
      {test: /\.(woff|woff2|eot|ttf|otf)$/, use: 'file-loader'},
      {
        resourceQuery: /blockType=i18n/,
        type: 'javascript/auto',
        loader: '@kazupon/vue-i18n-loader',
      },
      {
        enforce: 'pre',
        test: /\.(js|vue)$/,
        loader: 'eslint-loader',
        exclude: /node_modules/,
      },
    ],
  },
  plugins: [
    new HtmlWebpackPlugin({
      template: './web/index.html',
    }),
    new MiniCssExtractPlugin({
      filename: 'styles.css',
    }),
    new VueLoaderPlugin(),
  ],
}