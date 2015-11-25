#!/bin/sh

while : ; do
    echo "starting scraper"
    (node scrape.js > /dev/null 2>&1)
done
