#!/usr/bin/env python
# encoding: utf-8
from __future__ import print_function
import json
import sys


# The file name - the shard is passed in the first argument
shard = sys.argv[1]

# Mapping function
def mapfn(key, value):
    for w in value.split():
        yield w, 1


# MAPPER reads JSON objects from an array. According to the protocol,
# each object always consists of two fields:
# Key - a string that is an object identifier
# Value - in general, a JSON object, but this mapper
# expects it to be a string.
#
# According to the protocol, the mapper spits out a JSON object that is an array
# The array consists of JSON objects, each of which has fields
# Key - Convolution key
# Value - in general, an arbitrary JSON object that is passed to convolution tasks
# This mapper spits numbers as Value
with open(shard, "r") as f:
    input_values = json.JSONDecoder().decode(f.read())
    output_values = []

    for v in input_values:
        key = v["Key"]
        value = v["Value"]
        output_values += [{"Key": k, "Value": v} for (k, v) in mapfn(key, value)]

    print(json.JSONEncoder().encode(output_values))