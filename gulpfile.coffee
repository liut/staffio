gulp       = require 'gulp'
uglify     = require 'gulp-uglify'
sourcemaps = require 'gulp-sourcemaps'
stylus     = require 'gulp-stylus'
rename     = require 'gulp-rename'
gutil      = require 'gulp-util'
source     = require 'vinyl-source-stream'
buffer     = require 'vinyl-buffer'
browserify = require 'gulp-browserify'
del        = require 'del'
concat     = require 'gulp-concat'
cleanCSS   = require 'gulp-clean-css'


paths =
  scripts: [
    './fe/*/*.js',
  ],
  pwstrength: [
    './node_modules/pwstrength-bootstrap/src/i18n.js',
    './node_modules/pwstrength-bootstrap/src/rules.js',
    './node_modules/pwstrength-bootstrap/src/options.js',
    './node_modules/pwstrength-bootstrap/src/ui.js',
    './node_modules/pwstrength-bootstrap/src/methods.js'
  ],
  stylesheets: [
    './fe/*/*.styl',
  ],
  stylus: [
    './node_modules/bootstrap/dist/css',
    './node_modules/bootstrapValidator/dist/css',
    './node_modules/bootstrap-select/dist/css',
    './node_modules/x-editable/dist/bootstrap3-editable/css',
    './node_modules/datatables.net-bs/css',
    './node_modules/font-awesome-stylus/lib',
  ],
  fonts: [
    './node_modules/bootstrap/fonts/*',
    './node_modules/font-awesome-stylus/fonts/fontawesome-webfont.*',
  ],
  images: [
    './node_modules/x-editable/dist/bootstrap3-editable/img/*.*'
  ],
  dest: './fe/build/static',
  dest_maps: './fe/build/static_maps',
  release: './htdocs/static'


# Fonts
gulp.task 'fonts', () ->
  gulp.src(paths.fonts)
    .pipe gulp.dest(paths.release + '/fonts/')


gulp.task 'build:pwstrength', () ->
  gulp.src(paths.pwstrength)
    .pipe concat('pwstrength.js')
    .pipe gulp.dest('./fe/scripts/plugin')


gulp.task 'build:scripts', gulp.series ['build:pwstrength'], () ->
  gulp.src(paths.scripts, { sourcemaps: true })
    .pipe browserify({transform: 'babelify'})
    .pipe buffer()
    .pipe sourcemaps.init(loadMaps: false)
    .pipe gulp.dest(paths.dest)
    .pipe uglify().on('error', gutil.log)
    # .pipe rename({ extname: '.min.js' })
    .pipe gulp.dest(paths.release)
    .pipe sourcemaps.write('./')
    .pipe gulp.dest(paths.dest_maps)


gulp.task 'build:stylesheets', () ->
  gulp.src(paths.stylesheets)
    .pipe sourcemaps.init(loadMaps: false)
    .pipe stylus(compress: true, paths: paths.stylus, 'include css': true)
    .pipe rename('css/main.css')
    .pipe gulp.dest(paths.dest)
    .pipe cleanCSS {specialComments: '*', format:'breaks.afterComment'}, (details) ->
      console.log(details.name + ': original ' + details.stats.originalSize + ' bytes')
      console.log(details.name + ': minified ' + details.stats.minifiedSize + ' bytes')
    .pipe gulp.dest(paths.release)
    .pipe sourcemaps.write('./')
    .pipe gulp.dest(paths.dest_maps)


gulp.task 'build:images', () ->
  gulp.src(paths.images)
    .pipe gulp.dest(paths.release + '/img')


gulp.task 'build', gulp.parallel ['fonts', 'build:scripts', 'build:stylesheets', 'build:images']

gulp.task 'watch', gulp.series ['build'], () ->
  gulp.watch paths.scripts, gulp.series ['build:scripts']
  gulp.watch paths.stylesheets, gulp.series ['build:stylesheets']


gulp.task 'clean', (cb) ->
  del([paths.dest + '/**/*.{js,css,map,png,gif}', paths.release + '/**/*.{js,css,map,png,gif}'], cb)

