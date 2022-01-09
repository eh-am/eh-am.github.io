+++
date = "2022-01-09"
title = "Setting up evdi after updating kernel to 5.15.12"
slug = "setting-up-evdi-after-updating-kernel-to-5-15-12"
categories = ["blog", "technical", "linux"]
+++

# intro
[as you may know, I use an old ipad as a second screen.](./ipad-second-monitor-setup).
in a nutshell, it relies on the `evdi` kernel module to create a virtual display which is then accessed by the ipad via a VNC Viewer.

pne quirk is that almost every time I update the kernel, `evdi` (and consequently my ipad setup) also breaks. I am not super familiar with kernel shenanigans so I don't know if it's a fundamental issue or not. [^1]
[^1]: based on the github issues I would say most people also have that issue. (https://github.com/DisplayLink/evdi/issues?q=is%3Aissue+is%3Aclosed+kernel+issues)

that being said, fixing it is relative fine as long as:
* a) there's been a fix released
* b) I remember the steps on how to reinstall the module


for a), well, running fedora helps a bit with this (I reckon it would be much worse with a rolling release distro, like arch). It's still possible to mess it up though, so I usually wait a few **weeks** before updating the OS version. Also watching [new releases](https://github.com/DisplayLink/evdi/releases) help.


now for b), that's what this blog post is for. I usually maintain notes for almost everything technical I do, but it's designed to be a write-intensive workflow (which means I suffer between `grep`ing for the right keywords and sequentially going over my daily notes).


# the meat
i am going from kernel `5.15.8-200.fc35.x86_64` to `5.15.12-200.fc35.x86_64`


last time I updated was 2021-12-19.


let's clone the `evdi` repo and check out the `devel` branch, or in my case I already had it cloned

```shell
git@github.com:DisplayLink/evdi.git
git checkout devel
```

pull latest changes

```
git pull origin devel
```

`make` it
```shell
❯ make
CFLAGS="-isystem./include -isystem./include/uapi -Werror -Wextra -Wall -Wmissing-prototypes -Wstrict-prototypes -Wno-error=missing-field-initializers -Werror=sign-compare " make -C module 
make[1]: Entering directory '/home/eduardo/projects/evdi/module'
make -C /lib/modules/5.15.12-200.fc35.x86_64/build M=$PWD
make[2]: Entering directory '/usr/src/kernels/5.15.12-200.fc35.x86_64'
  CC [M]  /home/eduardo/projects/evdi/module/evdi_platform_drv.o
  CC [M]  /home/eduardo/projects/evdi/module/evdi_platform_dev.o
  CC [M]  /home/eduardo/projects/evdi/module/evdi_sysfs.o
  CC [M]  /home/eduardo/projects/evdi/module/evdi_modeset.o
  CC [M]  /home/eduardo/projects/evdi/module/evdi_connector.o
  CC [M]  /home/eduardo/projects/evdi/module/evdi_encoder.o
  CC [M]  /home/eduardo/projects/evdi/module/evdi_drm_drv.o
  CC [M]  /home/eduardo/projects/evdi/module/evdi_fb.o
  CC [M]  /home/eduardo/projects/evdi/module/evdi_gem.o
  CC [M]  /home/eduardo/projects/evdi/module/evdi_painter.o
  CC [M]  /home/eduardo/projects/evdi/module/evdi_params.o
  CC [M]  /home/eduardo/projects/evdi/module/evdi_cursor.o
  CC [M]  /home/eduardo/projects/evdi/module/evdi_debug.o
  CC [M]  /home/eduardo/projects/evdi/module/evdi_i2c.o
  CC [M]  /home/eduardo/projects/evdi/module/evdi_ioc32.o
  LD [M]  /home/eduardo/projects/evdi/module/evdi.o
  MODPOST /home/eduardo/projects/evdi/module/Module.symvers
  CC [M]  /home/eduardo/projects/evdi/module/evdi.mod.o
  LD [M]  /home/eduardo/projects/evdi/module/evdi.ko
  BTF [M] /home/eduardo/projects/evdi/module/evdi.ko
Skipping BTF generation for /home/eduardo/projects/evdi/module/evdi.ko due to unavailability of vmlinux
make[2]: Leaving directory '/usr/src/kernels/5.15.12-200.fc35.x86_64'
make[1]: Leaving directory '/home/eduardo/projects/evdi/module'
CFLAGS="-I../module -Werror -Wextra -Wall -Wmissing-prototypes -Wstrict-prototypes -Wno-error=missing-field-initializers -Werror=sign-compare " make -C library 
make[1]: Entering directory '/home/eduardo/projects/evdi/library'
cc evdi_lib.o -shared -Wl,-soname,libevdi.so.0 -o libevdi.so.1.10.0 -lc -lgcc 
cp libevdi.so.1.10.0 libevdi.so
make[1]: Leaving directory '/home/eduardo/projects/evdi/library'
```



`make install` it
```
sudo make install

make -C module install
make[1]: Entering directory '/home/eduardo/projects/evdi/module'
make -C /lib/modules/5.15.12-200.fc35.x86_64/build M=$PWD INSTALL_MOD_PATH= INSTALL_MOD_DIR=/kernel/drivers/gpu/drm/evdi modules_install
make[2]: Entering directory '/usr/src/kernels/5.15.12-200.fc35.x86_64'
  INSTALL /lib/modules/5.15.12-200.fc35.x86_64//kernel/drivers/gpu/drm/evdi/evdi.ko
  SIGN    /lib/modules/5.15.12-200.fc35.x86_64//kernel/drivers/gpu/drm/evdi/evdi.ko
At main.c:160:
- SSL error:02001002:system library:fopen:No such file or directory: crypto/bio/bss_file.c:69
- SSL error:2006D080:BIO routines:BIO_new_file:no such file: crypto/bio/bss_file.c:76
sign-file: certs/signing_key.pem: No such file or directory
  DEPMOD  /lib/modules/5.15.12-200.fc35.x86_64
make[2]: Leaving directory '/usr/src/kernels/5.15.12-200.fc35.x86_64'
true
make[1]: Leaving directory '/home/eduardo/projects/evdi/module'
make -C library install
make[1]: Entering directory '/home/eduardo/projects/evdi/library'
install -d /usr/local/lib
install -m 755 libevdi.so.1.10.0 /usr/local/lib/libevdi.so.1.10.0
ln -sf libevdi.so.1.10.0 /usr/local/lib/libevdi.so.0
ln -sf libevdi.so.0 /usr/local/lib/libevdi.so
make[1]: Leaving directory '/home/eduardo/projects/evdi/library'
```

ok it failed due to a missing `certs/signing_key.pem` to sign the module. [You need to generate one](https://askubuntu.com/questions/820883/how-to-resolve-ssl-error-during-make-modules-install-command/1178467#1178467).

in my case I already had one, so I just copied from the previous kernel.

```shell
find /usr/src/kernels/ -name "*.pem"

/usr/src/kernels/5.15.8-200.fc35.x86_64/certs/signing_key.pem
```


my kernel is
```shell
uname -r

5.15.12-200.fc35.x86_64
```


copying
```shell
sudo cp -r /usr/src/kernels/5.13.14-200.fc34.x86_64/certs/ /usr/src/kernels/5.15.12-200.fc35.x86_64/
```


you can use something like `$(uname -r)`, but since I run fish the syntax is different and I don't want to publish something with a syntax that will break to people. How considerate I am.


alright now `sudo make install` works

```shell
sudo make install

make -C module install
make[1]: Entering directory '/home/eduardo/projects/evdi/module'
make -C /lib/modules/5.15.12-200.fc35.x86_64/build M=$PWD INSTALL_MOD_PATH= INSTALL_MOD_DIR=/kernel/drivers/gpu/drm/evdi modules_install
make[2]: Entering directory '/usr/src/kernels/5.15.12-200.fc35.x86_64'
  INSTALL /lib/modules/5.15.12-200.fc35.x86_64//kernel/drivers/gpu/drm/evdi/evdi.ko
  SIGN    /lib/modules/5.15.12-200.fc35.x86_64//kernel/drivers/gpu/drm/evdi/evdi.ko
  DEPMOD  /lib/modules/5.15.12-200.fc35.x86_64
make[2]: Leaving directory '/usr/src/kernels/5.15.12-200.fc35.x86_64'
true
make[1]: Leaving directory '/home/eduardo/projects/evdi/module'
make -C library install
make[1]: Entering directory '/home/eduardo/projects/evdi/library'
install -d /usr/local/lib
install -m 755 libevdi.so.1.10.0 /usr/local/lib/libevdi.so.1.10.0
ln -sf libevdi.so.1.10.0 /usr/local/lib/libevdi.so.0
ln -sf libevdi.so.0 /usr/local/lib/libevdi.so
make[1]: Leaving directory '/home/eduardo/projects/evdi/library'
```


as you can see the module is there:
```shell
find /lib/modules/(uname -r) -name "*.ko" | grep -i evdi
/lib/modules/5.15.12-200.fc35.x86_64/kernel/drivers/gpu/drm/evdi/evdi.ko
```

so it's just a matter of loading that module with
```shell
sudo modprobe evdi
```

and done!

# conclusion
as you can see the process is pretty manual process, which most of the times is similar. let's see if i can remember posting the process every time it breaks again.

