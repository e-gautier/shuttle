const path = require('path');
const MinifyPlugin = require("babel-minify-webpack-plugin");
const CopyWebpackPlugin = require('copy-webpack-plugin');

module.exports = {
    entry: './public/js/main.js',
    output: {
        filename: 'main.js',
        path: path.resolve(__dirname, 'dist'),
        publicPath: '/dist/'
    },
    module: {
        rules: [
            {
                test: /\.(scss|css)$/,
                use: [
                    'style-loader',
                    'css-loader',
                    'sass-loader'
                ]
            },
            {
                test: /\.(png|svg|jpg|gif)$/,
                use: [
                    'file-loader'
                ]
            },
            {
                test: /\.json$/,
                use: [
                    'json-loader'
                ]
            },
            {
                test: /\.html$/,
                loader: 'raw-loader!html-minifier-loader'
            },
            {
                test: /\.(woff|woff2|eot|ttf|otf)$/,
                use: [
                    'file-loader'
                ]
            }
        ]
    },
    plugins: [
        new MinifyPlugin(),
        new CopyWebpackPlugin([
            { from: 'public/static' }
        ])
    ]
};
