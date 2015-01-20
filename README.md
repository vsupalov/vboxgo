#Vboxgo
Get user-like access to VirtualBox VMs from Go code.
This library wraps some define-tainted **VirtualBox SDK** functions, making it possible to get screenshots from a virtual machine and generate user input events (in particular mouse clicks). Constants and interfaces from [th4t/desktopcontroller](https://github.com/th4t/desktopcontroller) are used. Also a desktopcontroller interface implementation in *./controller* is provided, which makes the usage simpler.

At this point I would like to note, that the particular details are quite inelegant. I can't help but suspect that there must is a better way to interact with this SDK from Golang. Nevertheless it works reliably for me, and is superiour to a previous approach of mine, which involved an application written in Python and communication through sockets.

I am sure there are many things which can be improved upon, pull requests or suggestions are more than welcome.

##Prerequisites
The code has been developed and verified with the following moving parts:

* Known to work on **Ubuntu 14.04** and **Arch Linux**
* **VirtualBox** version: *4.3.20*
* **VirtualBox SDK** version: *4.3.20-96996*

You can download VirtualBox and the SDK [here](https://www.virtualbox.org/wiki/Downloads), but should probably use the VirtualBox version provided by the package manager for simplicity. Unpack the SDK to a directory of you choice.

##Setup
A few manual actions, adjustments and fiddling is needed to make this code ready to run.

Either, you need copy the SDK content into a new folder *./vbox_sdk* inside this repository directory, or create a symlink from here to where the SDK is unpacked/installed on your system using *ln -s* with the name *vbox_sdk*.

The VirtualBox SDK install path on Arch Linux is */usr/lib/virtualbox/sdk*, thus the command would be the following in the *clone*-d/*go get*-ted repository directory
```
$ ln -s /usr/lib/virtualbox/sdk vbox_sdk
```

One single *.c* file *VBoxCAPIGlue.c*  needs to be copied to from: *vbox_sdk/bindings/c/glue/* to the current directory. This particular bit bothers me quite a lot, how to help cgo to find this file in its original location?

#Tests
For the tests to work, a virtual machine named appropriately has to be up and running.

As soon as a virtual machine, with the name and operating system of your choice (Windows has been tested), has been created, the *TEST_VM_NAME* constant in *const.go* needs to be adjusted accordingly. This constant is used for tests of the controller as well. Now the tests  of *vboxgo* and *vboxgo/controller* should pass, using the following command in the respective directory.

```
$ go test
```

##Usage
To see simple usage examples, take a look at the tests. They should be sufficiently self-explaining and cover all relevant usecases. For simplicity, it is highly recommended to stick to the controller interface, as it is defined in [*th4t/desktopcontroller*](https://github.com/th4t/desktopcontroller) and implemented for VirtualBox in *th4t/vboxgo/controller*.

##License: MIT
Copyright (c) 2015 Vladislav Supalov

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
