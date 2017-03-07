#!/usr/bin/env ruby

pwd = `pwd`.strip

dist = './dist'
`mkdir -p #{dist}`

`git submodule update vendor`

gopkg = 'github.com/CenturyLinkCloud/clc-sdk'

ospairs = [
  ['windows', 'amd64', '.exe'],
  ['linux', 'amd64', '.linux'],
  ['darwin', 'amd64', '.osx'],
]

binaries = [
  ['natip', 'bin/natip/main.go'],
  ['baremetal-info', 'bin/baremetal-info/main.go'],
  ['apiv1', 'bin/apiv1/main.go'],
  ['apiv2', 'bin/apiv2/main.go'],
]

binaries.each do |basename, path|
  puts "== building #{basename}"
  ospairs.each do |os, arch, extension|
    cmd = "docker run --rm -it " +
          "-e GOOS=#{os} -e GOARCH=#{arch} " +
          "-v #{pwd}/bin:/go/src/#{gopkg}/bin " +
          "-v #{pwd}/dist:/go/src/#{gopkg}/dist " +
          "-v #{pwd}/../../CenturyLinkCloud/clc-sdk:/go/src/#{gopkg} " +
          "-w /go/src/#{gopkg} " +
          "golang:1.7 go build -o #{dist}/#{basename}#{extension} #{path}"
    puts "BUILD: #{cmd}"
    puts `#{cmd}`
    raise "ERROR BUILDING #{basename}" unless $?.success?
  end
end
