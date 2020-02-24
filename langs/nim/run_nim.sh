#!/bin/sh -e
/usr/bin/nim --cc:tcc --hints:off --nimcache:/tmp --verbosity:0 -o:/tmp/code c - && /tmp/code "$@"
