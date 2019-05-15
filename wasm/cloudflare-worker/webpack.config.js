const path = require('path')

module.exports = {
  entry: {
    bundle: path.join(__dirname, './index.js'),
  },

  output: {
    filename: 'bundle.js',
    path: path.join(__dirname, 'dist'),
  },

  target: 'webworker',

  // Switch to production for minified javascript
  mode: 'development',
  //mode: 'production', 

  watchOptions: {
    ignored: /dist/g,
  },

  devtool: "source-map",

  resolve: {
    extensions: ['.ts', '.tsx', '.js', '.json', '.wasm'],
    plugins: [],
  },

  module: {
    rules: [
    ],
  },
}