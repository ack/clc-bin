# clc-bin

Collected standalone binaries for use on CLC. 

# Adding binaries

1. add under bin, test, then run: `godep save ./...` to pick up any add'l dependencies
2. add to buildbin.rb so it's automatically rebuilt in docker builds
3. publish as a release

