#!/bin/sh

while : ; do
    echo "starting scraper"
    (node scrape.js > scrape.log 2>&1)
done
