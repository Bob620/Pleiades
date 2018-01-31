require 'compass/import-once/activate'
require 'ceaser-easing'

http_path = "/"
css_dir = "static/dist"
sass_dir = "static/sass"
images_dir = "static/img"
javascripts_dir = "static/js"

output_style = :compressed
sourcemap = false

# To enable relative paths to assets via compass helper functions. Uncomment:
# relative_assets = true
relative_assets = false
sass_options = {:debug_info=>false}

# To disable debugging comments that display the original location of your selectors. Uncomment:
# line_comments = false
line_comments = true

preferred_syntax = :scss