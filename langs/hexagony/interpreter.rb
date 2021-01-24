#!/usr/bin/env ruby
# coding: utf-8

require_relative 'hexagony'
require_relative 'grid'
require 'stringio'

input = StringIO.new ARGV*"\0"
ARGV.clear

Hexagony.run(ARGF.read, 0, input)
