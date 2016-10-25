package main

import (
	"archive/zip"
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"path"
	"regexp"
	"sort"

	"github.com/northbright/pathhelper"
)

var (
	// DEBUG is debug mode.
	DEBUG         = false
	appName       = ""
	apkPath       = ""
	arch          = "" // armeabi, armeabi-v7a, arm64-v8a, x86, x86_64
	apkDir        = ""
	androidmkPath = ""
	pattern       = `^lib/(.*)/(.*\.so)$`     // find arch and .so name in zip file: FindStringSubmatch() should return [3]string, arrs[1] = arch, arrs[2] = .so name
	libMap        = make(map[string][]string) // map to store libs. Key = arch, Value = lib slice
)

func genAndroidmk(file, apk, moduleName, arch string, libs []string) (err error) {
	var buffer bytes.Buffer

	buffer.WriteString("LOCAL_PATH := $(call my-dir)\n")
	buffer.WriteString("\nmy_archs := arm x86 arm64\n")
	buffer.WriteString("my_src_arch := $(call get-prebuilt-src-arch, $(my_archs))\n")
	buffer.WriteString("\ninclude $(CLEAR_VARS)\n")
	buffer.WriteString(fmt.Sprintf("LOCAL_MODULE := %s\n", moduleName))
	buffer.WriteString("LOCAL_MODULE_CLASS := APPS\n")
	buffer.WriteString("LOCAL_MODULE_TAGS := optional\n")
	buffer.WriteString("LOCAL_BUILT_MODULE_STEM := package.apk\n")
	buffer.WriteString("LOCAL_MODULE_SUFFIX := $(COMMON_ANDROID_PACKAGE_SUFFIX)\n")
	buffer.WriteString("LOCAL_CERTIFICATE := PRESIGNED\n")
	buffer.WriteString(fmt.Sprintf("LOCAL_SRC_FILES := %s\n", apk))

	if len(libs) > 0 {
		buffer.WriteString("\n")
		buffer.WriteString("LOCAL_PREBUILT_JNI_LIBS := \\\n")
		for i, v := range libs {
			buffer.WriteString("  @")
			buffer.WriteString(v)
			if i != len(libs)-1 {
				buffer.WriteString(" \\\n")
			} else {
				buffer.WriteString("\n")
			}
		}
	}

	buffer.WriteString("\nLOCAL_MODULE_TARGET_ARCH := $(my_src_arch)\n")
	buffer.WriteString("\ninclude $(BUILD_PREBUILT)\n")

	if err = ioutil.WriteFile(file, buffer.Bytes(), 0777); err != nil {
		fmt.Printf("ioutil.WriteFile(%s) error: %s\n", file, err)
		return err
	}

	return nil
}

func main() {
	flag.StringVar(&apkPath, "i", "", "input APK file. Ex: -i ./WeChat.apk")
	flag.StringVar(&appName, "n", "", "LOCAL_MODULE in Android.mk. If not set, it'll use APK's name. Ex: -n WeChat")

	flag.Parse()

	fmt.Printf("appName = %s\n", appName)
	fmt.Printf("apkPath = %s\n", apkPath)

	if apkPath == "" {
		flag.PrintDefaults()
		return
	}

	r, err := zip.OpenReader(apkPath)
	if err != nil {
		fmt.Printf("zip.OpenReader(%s) error: %s\n", apkPath, err)
		return
	}
	defer r.Close()

	if appName == "" {
		appName = pathhelper.GetFileNameWithoutExt(apkPath)
		fmt.Printf("appName is set to %s\n", appName)
	}

	apkDir = path.Dir(apkPath)
	androidmkPath = path.Join(apkDir, "Android.mk")
	fmt.Printf("Output Android.mk = %s\n", androidmkPath)

	re := regexp.MustCompile(pattern)

	for _, f := range r.File {
		arrs := re.FindStringSubmatch(f.Name)
		if len(arrs) == 3 {
			arch := arrs[1]
			so := arrs[2]
			if DEBUG {
				fmt.Printf("%s\n", f.Name)
				fmt.Printf("arch: %s\n", arch)
				fmt.Printf("so: %s\n", so)
			}
			libMap[arch] = append(libMap[arch], f.Name)
		}
	}

	// Check if it contains native libs(.so)
	if len(libMap) != 0 {
		arrs := []string{}

		for k, _ := range libMap {
			arrs = append(arrs, k)
		}

		if _, ok := libMap[arch]; !ok { // input arch argument is incorrect, get 1st arch
			// sort arches by string
			sort.Strings(arrs)
		}

		for {
			index := 1
			fmt.Printf("\nPlease select one of avaialbe arches in current APK:\n====================================\n")
			for i, v := range arrs {
				fmt.Printf("%d: %s\n", i+1, v)
			}

			// Wait user imput
			if _, err := fmt.Scanf("%d", &index); err != nil {
				fmt.Printf("fmt.Scanf() error: %s\n", err)
				return
			}

			index -= 1
			if 0 <= index && index < len(arrs) {
				arch = arrs[index]
				fmt.Printf("You choice: %s, libs:\n", arch)
				for _, v := range libMap[arch] {
					fmt.Printf("%s\n", v)
				}
				break
			} else {
				fmt.Printf("Please input the right number: %d -- %d\n", 1, len(arrs))
			}
		}
	} else { // no native libs
		fmt.Printf("no native libs in APK\n")
	}

	fmt.Printf("Start gnerating %s\n", androidmkPath)
	if err := genAndroidmk(androidmkPath, apkPath, appName, arch, libMap[arch]); err != nil {
		fmt.Printf("genAndroidmk() error: %s\n", err)
		return
	}
	fmt.Printf("Done\n")
}
