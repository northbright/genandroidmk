package main

import (
	"archive/zip"
	"flag"
	"fmt"
	"github.com/northbright/goutil"
	"path"
	"regexp"
	"sort"
)

var (
	DEBUG   bool                = false
	AppName string              = ""
	ApkPath string              = ""
	Arch    string              = "" // armeabi, armeabi-v7a, arm64-v8a, x86, x86_64
	ApkDir  string              = ""
	pattern string              = `^lib/(.*)/(.*\.so)$`     // find arch and .so name in zip file: FindStringSubmatch() should return [3]string, arrs[1] = arch, arrs[2] = .so name
	libMap  map[string][]string = make(map[string][]string) // map to store libs. Key = arch, Value = lib slice
	hasLibs bool                = false
)

func main() {
	flag.StringVar(&ApkPath, "i", "", "input APK file. Ex: -i ./WeChat.apk")
	flag.StringVar(&AppName, "n", "", "LOCAL_MODULE in Android.mk. If not set, it'll use APK's name. Ex: -n WeChat")
	flag.StringVar(&Arch, "a", "", "Arch of libraries(/lib/xx of apk) to be used.\nIt can be these values: armeabi, armeabi-v7a, arm64-v8a, x86, x86_64. If not set, it'll use the first arch in /lib. Ex: -a armeabi-v7a")

	flag.Parse()

	fmt.Printf("AppName = %s\n", AppName)
	fmt.Printf("ApkPath = %s\n", ApkPath)
	fmt.Printf("Arch = %s\n", Arch)

	r, err := zip.OpenReader(ApkPath)
	if err != nil {
		fmt.Printf("zip.OpenReader(%s) error: %s\n", err)
		return
	}
	defer r.Close()

	if AppName == "" {
		AppName = goutil.GetFileNameWithoutExt(ApkPath)
		fmt.Printf("AppName is set to %s\n", AppName)
	}

	ApkDir = path.Dir(ApkPath)
	fmt.Printf("ApkDir = %s\n", ApkDir)

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
			libMap[arch] = append(libMap[arch], so)
		}
	}

	// Check if it contains native libs(.so)
	if len(libMap) != 0 {
		hasLibs = true
		arrs := []string{}

		for k, _ := range libMap {
			arrs = append(arrs, k)
		}

		if _, ok := libMap[Arch]; !ok { // input Arch argument is incorrect, get 1st arch
			// sort arches by string
			sort.Strings(arrs)
			Arch = arrs[0]
		}

		fmt.Printf("hasLibs = %v, Arch = %s, libs:\n", hasLibs, Arch)
		for _, v := range libMap[Arch] {
			fmt.Printf("%s\n", v)
		}
	} else { // no native libs
		fmt.Printf("hasLibs = %v\n", hasLibs)
	}
}
