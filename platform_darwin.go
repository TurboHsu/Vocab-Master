package main

import (
	"os"
	"os/exec"
	"path/filepath"
)

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Foundation
#include <Foundation/Foundation.h>

const char* nsstring2cstring(NSString *s) {
    if (s == NULL) { return NULL; }

    const char *cstr = [s UTF8String];
    return cstr;
}

NSString* getFontPath() {
	NSBundle *main = [NSBundle mainBundle];
	return [main pathForResource:@"font" ofType:@"ttf"];
}
*/
import "C"

func CString(s *C.NSString) *C.char {
	return C.nsstring2cstring(s)
}

func GoString(p *C.NSString) string {
	return C.GoString(CString(p))
}

func GetPlatform() (Platform, error) {
	fontPath := GoString(C.getFontPath())
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return Platform{}, err
	}
	dataDir := filepath.Join(homeDir, "Library", "Application Support", "VocabMaster")

	err = os.MkdirAll(dataDir, os.ModePerm)
	if err != nil {
		return Platform{}, err
	}

	//to get cocoa resources dir, where the font is located,
	//a window is required, even if it's not shown

	return Platform{
		DataDir: dataDir,
		CertDir: filepath.Join(dataDir, "cert"),
		Font:    fontPath,
	}, nil
}

func (receiver Platform) OpenCertDir() {
	cmd := exec.Command("/bin/bash", "-c", "open \""+receiver.CertDir+"\"")
	cmd.Start()
}
