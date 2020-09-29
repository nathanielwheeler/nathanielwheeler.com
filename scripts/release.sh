#!/bin/bash

cd "$GOPATH/src/nathanielwheeler.com"
echo "==== Releasing nathanielwheeler.com ===="

# Merge dev into master
echo "	Merging dev to master..."
git checkout master
git merge dev
git push

# Pull current code from master
echo "	Pulling from master in production..."
ssh nathanielwheeler.com "cd /root/go/src/nathanielwheeler.com; \
	git pull"

# Update dependencies
echo "	Updating dependencies..."
ssh nathanielwheeler.com "cd /root/go/src/nathanielwheeler.com; \
	go get -u ./...; \
	go mod tidy; \
	yarn"

# Building binaries
echo "	Building nathanielwheeler.com..."
ssh nathanielwheeler.com "cd /root/app; \
	go build -o ./server /root/go/src/nathanielwheeler.com/main.go"

# Preprocess Sass
echo "	Preprocessing Sass..."
ssh nathanielwheeler.com "cd /root/go/src/nathanielwheeler.com; \
	dart-sass app/sass/main.sass public/stylesheets/main.css"

# Transferring files
echo "	Transferring files..."
ssh nathanielwheeler.com "cd /root/app; \
	cp -R /root/go/src/nathanielwheeler.com/public .; \
	cp -R /root/go/src/nathanielwheeler.com/views .; \
	cp /root/go/src/nathanielwheeler.com/Caddyfile ."

# Restarting server
echo "	Restarting server..."
ssh nathanielwheeler.com "service nathanielwheeler.com restart"

# Restarting server
echo "	Restarting Caddy..."
ssh nathanielwheeler.com "service caddy restart"

echo "==== nathanielwheeler.com Released! ===="

# Checkout back to dev
git checkout dev