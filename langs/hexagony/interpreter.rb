#!/usr/bin/env ruby

require_relative 'hexagony'
require 'stringio'

input = StringIO.new ARGV*"\0"
ARGV.clear

Hexagony.run(ARGF.read, 0, input)
