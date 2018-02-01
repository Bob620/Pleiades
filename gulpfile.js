const gulp = require('gulp'),
      concat = require('gulp-concat'),
      sourcemaps = require('gulp-sourcemaps'),
      gutil = require('gulp-util'),
      uglify = require('gulp-uglify-es').default;

process.chdir('static');
/*
gulp.task('concat-vendor-js', () => {
    return gulp.src(['js/vendor/jquery.min.js', 'js/vendor/*.js'])
        .pipe(concat('minnehack.vendor.js'))
        .pipe(sourcemaps.write())
        .pipe(gulp.dest('dist'));
});
*/
gulp.task('minify-app-js', () => {
    return gulp.src(['js/*.js'])
        .pipe(sourcemaps.init({ loadMaps: true }))
        .pipe(uglify().on('error', e => {
            gutil.log("Failed to minify javascript:");
gutil.log(e.message);
}))
.pipe(concat('pleiades.min.js'))
    .pipe(sourcemaps.write("."))
    .pipe(gulp.dest('dist'));
});

//const taskTree = ['concat-vendor-js', 'minify-app-js'];
const taskTree = ['minify-app-js'];
gulp.task('default', taskTree);
gulp.task('stop', taskTree, () => {
    process.exit();
});

gulp.watch(['js/**/*.js'], ['minify-app-js']).on('change', (event) => {
    console.log(`File ${event.path} was ${event.type}, running tasks...`);
});
/*
gulp.watch(['js/vendor/*.js'], ['concat-vendor-js']).on('change', (event) => {
    console.log(`File ${event.path} was ${event.type}, running tasks...`);
});
*/