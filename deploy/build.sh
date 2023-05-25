#!/bin/sh
#This is a script for building and indexing the chart.

helm lint simply-dns-webhook
#helm package simply-dns-webhook
helm repo index . --url https://runnerm.github.io/simply-dns-webhook/ --merge index.yaml