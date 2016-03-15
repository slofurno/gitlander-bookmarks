#!/bin/sh

curl -X POST -H "Authorization: 2799535" 127.0.0.1:555/api/follow?follow=$1
