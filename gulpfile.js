var gulp       = require('gulp')
var uglify     = require('gulp-uglify')
var sourcemaps = require('gulp-sourcemaps')
var stylus     = require('gulp-stylus')
var rename     = require('gulp-rename')
var gutil      = require('gulp-util')
var source     = require('vinyl-source-stream')
var buffer     = require('vinyl-buffer')
var browserify = require('gulp-browserify')
var del        = require('del')
var concat     = require('gulp-concat');
var cleanCSS   = require('gulp-clean-css');

var paths = {
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
  ],
  fonts: [
    './node_modules/bootstrap/fonts/*',
  ],
  images: [
    './node_modules/x-editable/dist/bootstrap3-editable/img/*.*'
  ],
  dest: './fe/build/static',
  dest_maps: './fe/build/static_maps',
  release: './htdocs/static'
};

gulp.task('build', ['fonts', 'build:scripts', 'build:stylesheets', 'build:images'])

// Fonts
gulp.task('fonts', function() {
  gulp.src(paths.fonts)
    .pipe(gulp.dest(paths.release + '/fonts/'));
});

gulp.task('build:pwstrength', function() {
  gulp.src(paths.pwstrength)
    .pipe(concat('pwstrength.js'))
    .pipe(gulp.dest('./fe/scripts/plugin'));
})

gulp.task('build:scripts', ['build:pwstrength'], function() {
  gulp.src(paths.scripts, { sourcemaps: true })
    .pipe(browserify({transform: 'babelify'}))
    .pipe(buffer())
    .pipe(sourcemaps.init({loadMaps: false}))
    .pipe(gulp.dest(paths.dest))
    .pipe(uglify().on('error', gutil.log))
    // .pipe(rename({ extname: '.min.js' }))
    .pipe(gulp.dest(paths.release))
    .pipe(sourcemaps.write('./'))
    .pipe(gulp.dest(paths.dest_maps))
});

gulp.task('build:stylesheets', function() {
  gulp.src(paths.stylesheets)
    .pipe(sourcemaps.init({loadMaps: false}))
    .pipe(stylus({compress: true, paths: paths.stylus, 'include css': true}))
    .pipe(rename('css/main.css'))
    .pipe(gulp.dest(paths.dest))
    .pipe(cleanCSS({specialComments: '*', format:'breaks.afterComment'}, function(details) {
            console.log(details.name + ': ' + details.stats.originalSize);
            console.log(details.name + ': ' + details.stats.minifiedSize);
        }))
    .pipe(gulp.dest(paths.release))
    .pipe(sourcemaps.write('./'))
    .pipe(gulp.dest(paths.dest_maps))
});

gulp.task('build:images', function() {
  gulp.src(paths.images)
    .pipe(gulp.dest(paths.release + '/img'))
});

gulp.task('watch', ['build'], function() {
  gulp.watch(paths.scripts, ['build:scripts'])
  gulp.watch(paths.stylesheets, ['build:stylesheets'])
});

gulp.task('clean', function(cb) {
  del([paths.dest + '/**/*.{js,css,map,png,gif}', paths.release + '/**/*.{js,css,map,png,gif}'], cb)
});
