#!/usr/bin/env python
# encoding: utf-8
from __future__ import print_function
import json
import sys

# This script is called for each convolution task
# The convolution key is passed in the start arguments
key = sys.argv[1]
reducerId = sys.argv[2]

# The list of values ​​is passed as an array of JSON objects
# in the standard input stream. Each object in the array is the result
# output some sort of mapping task
values = json.JSONDecoder().decode(sys.stdin.read())

# In this map-reduce, the mapper spits out for each key a number,
# hence as a value in the convolution an array of numbers
# The script should write to its standard output a valid JSON object
# is the result of convolution. The number is also a valid object,
# and it can also be written to the result
print(sum(values))

# The script could also write the result to a file.
# It would be wiser in a distributed case, where each reducer records
# result in your own output shard.
with open(reducerId, "a") as f:
    print("Key:%s Sum:%d" % (key, sum(values)), file=f)