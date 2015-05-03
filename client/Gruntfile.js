(function() {
  'use strict';

var cfg = {
  build: {
    path: {
      root: "build/",
      js: "build/js/",
      css: "build/css/",
      fonts: "build/fonts/"
    }
  },
  files: {
    grunt: ['Gruntfile.js'],
    app: {
      js: ['src/js/**/*.js', '!src/node_modules/bower_components/**/*.js'],
      css: ['src/js/**/*.css'],
      html: [ 'src/js/**/*.html' ]
    },
    vendor: {
      root: 'src/node_modules/bower_components/',
      js: [
        'src/node_modules/bower_components/jquery/dist/jquery.min.js',
        'src/node_modules/bower_components/bootstrap/dist/js/bootstrap.min.js',
        'src/node_modules/bower_components/angular/angular.min.js'
      ],
      css: ['src/node_modules/bower_components/bootstrap/dist/css/*.min.css'],
      fonts: ['src/node_modules/bower_components/bootstrap/dist/fonts/**']
    }
  }
};

module.exports = function(grunt) {

  grunt.initConfig({
    jshint: {
      files: cfg.files.app.js
    },
    watch: {
      build: {
        files: [cfg.files.grunt, 'src/**', '!src/node_modules/bower_components/**'],
        tasks: ['build']
      },
      js: {
        files: ['<%= jshint.files %>'],
        tasks: ['jshint']
      }
    },
    clean: [ cfg.build.path.root ],
    copy: {
      build: {
        files: [
        {
          src: cfg.files.vendor.fonts,
          dest: cfg.build.path.fonts,
          flatten: true,
          expand: true
        }]
      },
      vendorjs: {
        files: [{
            src: [ cfg.files.vendor.root + 'modernizr/modernizr.js' ],
            dest: cfg.build.path.js,
            expand: true,
            flatten: true
          }
        ]
      },
      staticContent: {
        files: [{
          src: ['*'],
          dest: cfg.build.path.root,
          cwd: 'src/static/',
          expand: true
        }]
      }
    },
    html2js: {
      main: {
        src: cfg.files.app.html,
        dest: cfg.build.path.js + 'templates.js',
        options: {
          base: "src/js"
        }
      }
    },
    concat: {
      maincss: {
        src: cfg.files.app.css,
        dest: cfg.build.path.css + 'app.css'
      },
      vendorcss: {
        src: cfg.files.vendor.css, 
        dest: cfg.build.path.css + 'vendor.css'
      }
    },
    browserify: {
      app: {
        src: 'src/js/index.js',
        dest: cfg.build.path.js + 'app.js'
      }
    }
  });

  grunt.loadNpmTasks('grunt-contrib-jshint');
  grunt.loadNpmTasks('grunt-contrib-watch');
  grunt.loadNpmTasks('grunt-contrib-copy');
  grunt.loadNpmTasks('grunt-contrib-clean');
  grunt.loadNpmTasks('grunt-contrib-concat');
  grunt.loadNpmTasks('grunt-browserify');
  grunt.loadNpmTasks('grunt-html2js');


  grunt.registerTask('default', ['build', 'watch:build']);
  grunt.registerTask('build', ['jshint', 'clean', 'browserify', 'html2js', 'concat', 'copy']);


};
})();
