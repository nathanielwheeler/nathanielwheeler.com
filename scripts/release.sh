#!/bin/bash

cd "$GOPATH/src/nathanielwheeler.com"
echo "==== Releasing nathanielwheeler.com ===="

# Merge dev into master
echo "	Merging dev to master..."
git checkout master
git merge dev

# Pull current code from master
echo "	Pulling from master in production..."
ssh root@128.199.11.200 "cd /root/go/src/nathanielwheeler.com; \
	git pull"

# Update dependencies
echo "	Updating dependencies..."
ssh root@128.199.11.200 "cd /root/go/src/nathanielwheeler.com; \
	go get -u; \
	go mod tidy"

# Building binaries
echo "	Building nathanielwheeler.com..."
ssh root@128.199.11.200 "cd /root/app; \
	go build -o ./server /root/go/src/nathanielwheeler.com/*.go"

# Transferring files
echo "	Transferring files..."
ssh root@128.199.11.200 "cd /root/app; \
	cp -R /root/go/src/nathanielwheeler.com/public .; \
	cp -R /root/go/src/nathanielwheeler.com/views .; \
	cp /root/go/src/nathanielwheeler.com/Caddyfile ."

# Restarting server
echo "	Restarting server..."
ssh root@128.199.11.200 "service nathanielwheeler.com restart"

# Restarting server
echo "	Restarting Caddy..."
ssh root@128.199.11.200 "service caddy restart"

echo "==== nathanielwheeler.com Released! ===="

# Checkout back to dev
git checkout dev