const path = require("path")
const webpack = require('webpack')

const MiniCssExtractPlugin = require('mini-css-extract-plugin')
const HtmlWebpackPlugin = require('html-webpack-plugin')
const { CleanWebpackPlugin } = require('clean-webpack-plugin')
const CopyPlugin = require('copy-webpack-plugin')

const TerserJSPlugin = require('terser-webpack-plugin')
const OptimizeCSSAssetsPlugin = require('optimize-css-assets-webpack-plugin')

const HardSourceWebpackPlugin = require('hard-source-webpack-plugin')

module.exports = {
    mode: 'development',
    entry: path.join(__dirname, "src/main.jsx"),
    output: {
        filename: "bundle.[contenthash:12].js",
        path: path.join(__dirname, "../public/"),
        publicPath: '/'
    },
    //stats: 'errors-only',
    //stats: 'minimal',
    //stats: {
    //    all: false,
    //    assets: true,
    //    errorDetails: true
    //},

    devtool: 'inline-source-map',
    watchOptions: {
        ignored: [ 'node_modules' ],
        aggregateTimeout: 500,
        poll: 1000
    },
    devServer: {
        port: 3000,
        historyApiFallback: true
    },
    module: {
        rules: [
            {
                test: /\.js?x$/,
                exclude: /node_modules/,
                use: ['cache-loader', 'babel-loader'],
                resolve: {
                    extensions: ['.js', '.jsx']
                },
            },
            {
                test: /\.css$/i,
                use: [
                    MiniCssExtractPlugin.loader,
                    {
                        loader: 'css-loader',
                    }
                ],
            },
            {
                test: /\.scss$/i,
                use: [
                    //{
                    //    loader: 'cache-loader',
                    //},

                    MiniCssExtractPlugin.loader,

                    {
                        loader: 'css-loader',

                    },
                    {
                        loader: 'postcss-loader',
                    },
                    {
                        loader: 'sass-loader',
                    },
                ],

            },
            {
                test: /\.(woff(2)?|ttf|eot|svg)(\?v=\d+\.\d+\.\d+)?$/,
                use: [{
                    loader: 'file-loader',
                    options: {
                        name: '[name].[ext]',
                        //outputPath: ''
                    }
                }],
            },

        ]
    },
    plugins: [
        //new HardSourceWebpackPlugin(),
        new webpack.ProgressPlugin(),
        new CleanWebpackPlugin(),
        new HtmlWebpackPlugin({
            template: path.join(__dirname, 'src/index.html'),
            base: "/"
        }),
        new MiniCssExtractPlugin({
            filename: "bundle.[contenthash:12].css",
        }),
        new CopyPlugin([
            { from: 'public', to: '' },
        ]),
    ],
    performance: {
        maxAssetSize: 1 * 1024 * 1024,
        hints: false
    }

}
