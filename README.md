# Welcome
Welcome to yet another image to terminal project.
Target of this project is to create a tool that will be able to display images in Linux terminal (should work on all posix compliant systems though).

## Usage
Just run that program with file name as argument.
Additionally there are some flags you can use if you need:
- `-dql` - Instead using background color only, program will use U+2584 (Lower half block) and text color. By default it's disabled.
- `-sampling [fast/average/uv]` - `fast` is default because it's faster. `average` will make things look smoother but it will be much slower on bigger images.`uv` is for debugging only.
- `-width [size]` and `-height [size]` - Allows you to set custom output size (setting this to 0 will make output size same as terminal size).
