const path = require('path');
//const UglifyJSPlugin = require('uglifyjs-webpack-plugin');

module.exports = {
  devtool: 'source-map',
  entry: path.resolve(__dirname, 'webroot', 'js', 'index.js'),
  output: {
    filename: 'bundle.js',
    path: path.resolve(__dirname, 'webroot', 'js', 'dist')
  }//,
  //plugins: [
  //  new UglifyJSPlugin({
  //    sourceMap: true
  //  })
  //]
};
