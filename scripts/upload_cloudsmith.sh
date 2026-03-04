#!/bin/bash
set -e

if [ -z "$CLOUDSMITH_TOKEN" ]; then
    echo "CLOUDSMITH_TOKEN is not set. Skipping Cloudsmith upload."
    exit 0
fi

pip install --quiet cloudsmith-cli

ORG="ksckaan1"
REPO="crtui"

echo "Uploading packages to Cloudsmith..."

for deb in dist/*.deb; do
    if [ -f "$deb" ]; then
        echo "Uploading $deb..."
        cloudsmith push deb "$ORG/$REPO/ubuntu/any-version" "$deb" --no-wait-for-sync || true
    fi
done

for rpm in dist/*.rpm; do
    if [ -f "$rpm" ]; then
        echo "Uploading $rpm..."
        cloudsmith push rpm "$ORG/$REPO/any-distro/any-version" "$rpm" --no-wait-for-sync || true
    fi
done

for apk in dist/*.apk; do
    if [ -f "$apk" ]; then
        echo "Uploading $apk..."
        cloudsmith push alpine "$ORG/$REPO/alpine/edge" "$apk" --no-wait-for-sync || true
    fi
done

echo "Cloudsmith upload completed."
