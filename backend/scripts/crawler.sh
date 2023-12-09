#!/bin/bash

echo "  Starting Crawler at $(date)"

for word in $(shuf -n 5 /usr/share/dict/words)
do
    echo "  Fetching word $word"
    curl "0.0.0.0:8080/paper?type=QUERY&id=$word&limit=5&maxDepth=2" > /dev/null
done

echo "  Ending Crawler at $(date)"
