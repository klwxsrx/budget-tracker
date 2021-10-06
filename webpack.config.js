const HtmlWebpackPlugin = require('html-webpack-plugin')
const MiniCssExtractPlugin = require('mini-css-extract-plugin')
const {VueLoaderPlugin} = require('vue-loader')

module.exports = (env, options) => {
  const isProduction = options.mode === 'production'
  const styleLoader = isProduction ? MiniCssExtractPlugin.loader : 'vue-style-loader'
  const devtool = isProduction ? 'source-map' : 'eval'

  return {
    mode: options.mode,
    entry: './web/main.js',
    output: {
      filename: 'bundle.js',
    },
    devtool: devtool,
    devServer: {
      static: __dirname,
      historyApiFallback: true,
    },
    module: {
      rules: [
        {test: /\.js$/, use: 'babel-loader'},
        {test: /\.vue$/, use: 'vue-loader'},
        {test: /\.css$/, use: [styleLoader, 'css-loader']},
        {test: /\.scss$/, use: [styleLoader, 'css-loader', 'sass-loader']},
        {test: /\.svg$/, use: 'file-loader'},
        {test: /\.(woff|woff2|eot|ttf|otf)$/, use: 'file-loader'},
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
        filename: '[contenthash:8].min.css',
      }),
      new VueLoaderPlugin(),
    ],
  }
}