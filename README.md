# genandroidmk

[![Build Status](https://travis-ci.org/northbright/genandroidmk.svg?branch=master)](https://travis-ci.org/northbright/genandroidmk)
[![Go Report Card](https://goreportcard.com/badge/github.com/northbright/genandroidmk)](https://goreportcard.com/report/github.com/northbright/genandroidmk)

`genandroidmk` is a tool written in [Go](http://golang.org) that help to integrate prebuilt apps on Android(5.0 and later).  

#### Details of Integration of Prebuilt Apps on Android.

* Before Android 5.0 L  
To integrate prebuilt apps which contains native libraries, you need to:  

  1. Write the `Android.mk` and specify the native libraries in the APK.  
  2. Extract the libraries from `/lib` in APK and copy them to BSP.  

* Since Android 5.0 L  
There's no need to extract the libraries in APK and copy them to BSP.  
Android will extract the libs and copy them to `/system/app/APP_NAME/lib` automatically while building the system image.  
We only need to specify the prebuilt libs(`LOCAL_PREBUILT_JNI_LIBS`) in `Android.mk`.  

#### How it Works
* It'll generate `LOCAL_PREBUILT_JNI_LIBS` and the whole `Android.mk` automatically.
* It will check the `/lib` in APK and use a default CPU arch or let users choose the CPU arch of libraries to fill the `LOCAL_PREBUILT_JNI_LIBS` variable.

#### How to Use

    ./genandroidmk -i <input APK file> -n <app name> -d <default arch>
    Ex:
    ./genandroidmk -i WeChat.apk -n WeChat

* Arguments:

  * `<input APK file>`  
    
    abs path of APK and `Android.mk` will be outputed in the same folder
  * `<app name>`  
    
    It will be used to set `LOCAL_MODULE` in `Android.mk`  
    if it's not set, `LOCAL_MODULE` will be set to the APK name(without ".apk") by default.
  * `<default arch>`

    Specify a default architecture.
    When a default architecture is specified, it will not require users to choose the CPU architecure that's useful for CI / CD integration.

* Choose CPU Arch of native libraries:  
If the APK contains multi-arch native libraries, 
and there's **NO** default architecture specified(`-d` argument),
you will be asked to select one arch:

        Please select one of available arches in current APK:
        ====================================
        1: armeabi
        2: armeabi-v7a
        3: x86

* Output
`Android.mk` will be put under the same folder as input APK file.

#### License
* [MIT License](./LICENSE)

#### Thanks
* [Kushtrim Pacaj](https://github.com/KushtrimPacaj)
  * Contributed the new argument: [`-d <default arch>`](https://github.com/northbright/genandroidmk/pull/1) that's useful for CI / CD integration
