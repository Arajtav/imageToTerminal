# Welcome
Welcome to yet another image to terminal project.
Target of this project is to create a tool that will be able to display images in linux terminal (should work on all posix compliant systems though).

## Usage
Just run that program with file name as argument.
Additionaly there are some flags you can use if you need:
- `-sqpx` - Pixel will be 2 chars wide, which should make it look like squares.
- `-sampling [fast/average]` - `fast` is default because it's faster. `average` will make things look smoother but it will be much slower on bigger images.
- `-width [size]` and `-height [size]` - Allows you to set custom output size (setting this to 0 will make output size same as terminal size).
