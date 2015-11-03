module.exports = {
    entry: './src/index.jsx',
    output: {
        filename: './static/bundle.js'
    },
    module: {
        loaders: [
            {
                //tell webpack to use jsx-loader for all *.jsx files
                test: /\.jsx$/,
                loader: 'jsx-loader?insertPragma=React.DOM&harmony'
            }
        ]
    },
    externals: {

    },
    resolve: {
        extensions: ['', '.js', '.jsx']
    }
}
