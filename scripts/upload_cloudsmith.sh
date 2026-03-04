#!/bin/bash
set -e

if [ -z "$CLOUDSMITH_TOKEN" ]; then
    echo "CLOUDSMITH_TOKEN is not set. Skipping Cloudsmith upload."
    exit 0
fi

ORG="ksckaan1"
REPO="crtui"

echo "Uploading packages to Cloudsmith..."

# DEB
for deb in dist/*.deb; do
    if [ -f "$deb" ]; then
        echo "Uploading $deb..."
        cloudsmith push deb "$ORG/$REPO/any-distro/any-version" "$deb" --no-wait-for-sync || true
    fi
done

# RPM
for rpm in dist/*.rpm; do
    if [ -f "$rpm" ]; then
        echo "Uploading $rpm..."
        cloudsmith push rpm "$ORG/$REPO/any-distro/any-version" "$rpm" --no-wait-for-sync || true
    fi
done

# APK (Alpine)
for apk in dist/*.apk; do
    if [ -f "$apk" ]; then
        echo "Uploading $apk..."
        cloudsmith push alpine "$ORG/$REPO/alpine/any-version" "$apk" --no-wait-for-sync || true
    fi
done

echo "Cloudsmith upload completed."
