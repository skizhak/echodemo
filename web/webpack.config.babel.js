import {join} from 'path'
import ExtractTextPlugin from 'extract-text-webpack-plugin'
import UglifyJSPlugin from 'uglifyjs-webpack-plugin'

let fileName = 'echodemo'
const libraryName = 'echodemo'
const paths = {}
function absolute (...args) {
  return join(__dirname, ...args)
}
const defaultEnv = {'dev': true}
/**
 * Following allows publishing compiled files to multiple paths.
 * add path relative to directory of this config file.
 * For eg: if you need to output the compiled output in lib dir of parent dir echodemo-demo,
 * use '../echodemo-demo/lib'
 * @type {Array}
 */
const publishPaths = []

export default (env = defaultEnv) => {
  const plugins = []
  const rules = [{
    test: /\.scss$/,
    loader: ExtractTextPlugin.extract({
      fallback: 'style-loader',
      use: ['css-loader', 'sass-loader']
    }),
  }, {
    test: /\.html/,
    loader: 'handlebars-loader',
  }, {
    loader: 'babel-loader',
    test: /\.js$/,
    include: /(src)/,
    query: {
      presets: ['es2015'],
    }
  }]
  
  const externals = {
    'jquery': {amd: 'jquery', root: 'jQuery'},
    'lodash': {amd: 'lodash', root: '_'},
    'backbone': {amd: 'backbone', root: 'Backbone'},
  }

  if (env.prod) {
    plugins.push(new UglifyJSPlugin({
      compress: {
        warnings: false
      },
      mangle: {
        keep_fnames: true,
      },
      sourceMap: true,
      include: /\.min\.js$/,
    }))
  }

  // For every library added in the include env, we will remove from externals.
  if (env.include) {
    env.include.split(',').forEach((lib) => {
      delete externals[lib.trim()]
    })
    // Will append .bundle to the output file name.
    fileName = `${fileName}.bundle`
  }
  // Let's put css under css directory.
  plugins.push(new ExtractTextPlugin(fileName + '.css'))

  const configList = []

  const config = {
    entry: {
      [fileName]: absolute('src/main.js'),
      [`${fileName}.min`]: absolute('src/main.js')
    },
    devtool: 'source-map',
    module: {rules},
    externals: externals,
    resolve: {
      modules: [absolute('src'), 'node_modules'],
      alias: {},
      extensions: ['.js'],
    },
    plugins: plugins,
    stats: { children: false }
  }

  const defaultConfig = Object.assign({}, config, {
    output: {
      path: absolute('dist'),
      filename: 'js/[name].js',
      library: libraryName,
      libraryTarget: 'umd',
      umdNamedDefine: false,
    }
  })
  configList.push(defaultConfig)

  if (publishPaths) {
    publishPaths.forEach((outPath) => {
      configList.push(Object.assign({}, config, {
        output: {
          path: absolute(outPath),
          filename: 'js/[name].js',
          library: libraryName,
          libraryTarget: 'umd',
          umdNamedDefine: false,
        }
      }))
    })
  }

  return configList
}
