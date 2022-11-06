#!/bin/bash

if [ -d ./VocabMaster.app ]; then
  rm -r ./VocabMaster.app
fi

tmpdir=$(mktemp -d)
bundleDir="$tmpdir/VocabMaster.app"
contentDir="$bundleDir/Contents"
resourcesDir="$contentDir/Resources"
exeDir="$contentDir/MacOS"

mkdir "$bundleDir"
mkdir "$contentDir"
mkdir "$exeDir"
mkdir "$resourcesDir"

cp ./icon/VocabMaster.icns "$resourcesDir/icon.icns"
cp ./font/red_bean.ttf "$resourcesDir/font.ttf"

CGO_ENABLED=1 go build -o vocabmaster
mv ./vocabmaster "$exeDir"

plistContent=$(
  cat <<-END
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
	<key>CFBundleExecutable</key>
	<string>vocabmaster</string>
	<key>CFBundleIconFile</key>
	<string>icon.icns</string>
	<key>CFBundleIdentifier</key>
	<string>com.turbohsu.vocabmaster</string>
	<key>NSHighResolutionCapable</key>
	<true/>
</dict>
</plist>
END
)
echo "$plistContent" > "$contentDir/Info.plist"

cp -r "$bundleDir" .
rm -r "$tmpdir"