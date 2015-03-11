# genandroidmk

#### Description

###### Before Android 5.0 L  
To integrate prebuilt apps which contains native libraries, you need to:  

1. Write the `Android.mk` and specify the native libraries in the APK.  
2. Extract the libraries from `/lib` in APK and copy them to BSP.  

###### Since Android 5.0 L  
There's no need to extract the libraries in APK and copy them to BSP.  
Android will extract the libs and copy them to `/system/app/APP_NAME/lib` automatically while building the system image.  
We only need to specify the prebuilt libs(`LOCAL_PREBUILT_JNI_LIBS`) in `Android.mk`.  

`genandroidmk` is a tool that help to generate `LOCAL_PREBUILT_JNI_LIBS` in `Android.mk`.  
It will check the `/lib` in APK and let users choose the CPU arch of libraries to fill the `LOCAL_PREBUILT_JNI_LIBS` variable.  
