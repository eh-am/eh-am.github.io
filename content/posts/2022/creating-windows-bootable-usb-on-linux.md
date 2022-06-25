+++
date = "2022-06-22"
title = "Creating Windows bootable USB on Linux"
categories = ["blog", "technical", "bootable", "usb", "linux"]
+++

This is more for self documentation than anything else.


Had positive results using WoeUSB (https://github.com/woeusb/WoeUSB)


```shell
fdisk -l
```

Then
```shell
sudo ./woeusb-5.2.4.bash --device $ISO_FILENAME $DEVICE
```

Eg.
```shell
sudo ./woeusb-5.2.4.bash --device Win10_21H2_English_x64.iso /dev/sdd
```


