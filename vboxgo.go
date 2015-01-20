package vboxgo

// the following comment block is interpreted by cgo
// http://golang.org/cmd/cgo/

// ------------------------------------------------------ cgo

/*
#cgo CFLAGS: -I vbox_sdk/bindings/c/include -I vbox_sdk/bindings/c/glue
#cgo LDFLAGS: -ldl -lpthread

#include "VBoxCAPIGlue.h"
#include <stdlib.h>

HRESULT MY_FAILED(HRESULT rc) {
    return FAILED(rc);
}

// formerly : wrappers around pointer-functions, as go can not handle those
// currently: wrappers around define macros
// oh the humanity
void ClientInitialize(IVirtualBoxClient** vboxclient) {
    g_pVBoxFuncs->pfnClientInitialize(NULL, vboxclient);
}

void ClientUninitialize() {
    g_pVBoxFuncs->pfnClientUninitialize();
}

void ReleaseEverything(ISession* session, IVirtualBox* vbox, IVirtualBoxClient* vboxclient) {
    if (session) {
        ISession_Release(session);
    }
    if (vbox) {
        IVirtualBox_Release(vbox);
    }
    if (vboxclient) {
        IVirtualBoxClient_Release(vboxclient);
    }
}

HRESULT GetVirtualBox(IVirtualBoxClient* vboxclient, IVirtualBox** vbox) {
    return IVirtualBoxClient_get_VirtualBox(vboxclient, vbox);
}

HRESULT GetSession(IVirtualBoxClient* vboxclient, ISession** session) {
    return IVirtualBoxClient_get_Session(vboxclient, session);
}

HRESULT FindMachine(IVirtualBox* vbox, char* nameOrId, IMachine** machine) {
    BSTR BSTR_nameOrId = NULL;

    HRESULT res;
    g_pVBoxFuncs->pfnUtf8ToUtf16(nameOrId, &BSTR_nameOrId);

    res = IVirtualBox_FindMachine(vbox, BSTR_nameOrId, machine);

    g_pVBoxFuncs->pfnUtf16Free(BSTR_nameOrId);
    return res;
}

HRESULT LockMachine(IMachine* machine, ISession* session) {
    return IMachine_LockMachine(machine, session, LockType_Shared);
}

HRESULT UnlockMachine(ISession* session) {
    return ISession_UnlockMachine(session);
}

HRESULT GetConsole(ISession* session, IConsole** console) {
    return ISession_GetConsole(session, console);
}

HRESULT GetKeyboard(IConsole* console, IKeyboard** keyboard) {
    return IConsole_GetKeyboard(console, keyboard);
}

HRESULT GetMouse(IConsole* console, IMouse** mouse) {
    return IConsole_GetMouse(console, mouse);
}

HRESULT GetDisplay(IConsole* console, IDisplay** display) {
    return IConsole_GetDisplay(console, display);
}

const int MOUSE_1_DOWN = 0x1; //actually, see https://www.virtualbox.org/sdkref/interface_i_mouse.html
const int MOUSE_2_DOWN = 0x2;
const int MOUSE_ALL_UP = 0x0;
HRESULT MouseClickEvent(IMouse* mouse, int x, int y, int upOrDown) {
    return IMouse_PutMouseEventAbsolute(mouse, x, y, 0, 0, upOrDown);
}

HRESULT GetScreenResolution(IDisplay* display, uint* out_width, uint* out_height, uint* out_bitsPerPixel) {
    PRUint32 screenId = 0;
    PRUint32 width, height, bitsPerPixel;
    PRInt32 xOrigin, yOrigin;

    HRESULT res;
    res = IDisplay_GetScreenResolution(display, screenId, &width, &height, &bitsPerPixel, &xOrigin, &yOrigin);

    *out_width = (int)width;
    *out_height = (int)height;
    *out_bitsPerPixel = (int)bitsPerPixel;

    // assume, xOrigin and yOrigin are somewhere at 0? Do we even care?
    return res;
}

//TODO: supposedly, this api call is slow?
HRESULT TakeScreenShotPNGToArray(IDisplay* display, uint width_in, uint height_in, PRUint32* screenDataSize, PRUint8** screenData) {
    PRUint32 screenId = 0;
    PRUint32 width = width_in;
    PRUint32 height = height_in;

    HRESULT res;
    //res = IDisplay_TakeScreenShotPNGToArray(display, screenId, width, height, imageData);
    res = display->lpVtbl->TakeScreenShotPNGToArray(display, screenId, width, height, screenDataSize, screenData);

    return res;
}

*/
import "C" //this triggers the comment to be interpreted

// ------------------------------------------------------ /cgo

import (
    "bytes"
    "fmt"
    "log"
    "errors"
    "unsafe"
    "image"
    "image/png"
    "image/draw"
)

type VirtualBox struct {
    vboxclient      *C.IVirtualBoxClient
    vbox            *C.IVirtualBox
    session         *C.ISession

    // those variables could be gathered, but are not needed
    //versionApi      uint
    //versionVbox     uint

    //revision        C.ULONG
    //versionUtf16    C.BSTR
    //homefolderUtf16 C.BSTR
}

func (vb *VirtualBox) Init() error {
    if C.VBoxCGlueInit() != 0 {
        message := fmt.Sprintf("VBoxCGlueInit failed: %v", C.g_szVBoxErrMsg)
        log.Println(message)
        return errors.New(message)
    }

    //unsigned ver = g_pVBoxFuncs->pfnGetVersion();
    //printf("VirtualBox version: %u.%u.%u\n", ver / 1000000, ver / 1000 % 1000, ver % 1000);
    //ver = g_pVBoxFuncs->pfnGetAPIVersion();
    //printf("VirtualBox API version: %u.%u\n", ver / 1000, ver % 1000);

    C.ClientInitialize(&vb.vboxclient)
    if vb.vboxclient == nil {
        message := fmt.Sprintf("Could not get vboxclient reference")
        log.Println(message)
        return errors.New(message)
    }

    var rc C.HRESULT = C.GetVirtualBox(vb.vboxclient, &vb.vbox)
    if C.MY_FAILED(rc) != 0 || vb.vbox == nil {
        message := fmt.Sprintf("Could not get vbox reference")
        log.Println(message)
        return errors.New(message)
    }

    rc = C.GetSession(vb.vboxclient, &vb.session)
    if C.MY_FAILED(rc) != 0 || vb.session == nil {
        message := fmt.Sprintf("Could not get session reference")
        log.Println(message)
        return errors.New(message)
    }

    return nil
}

func (vb *VirtualBox) Uninit() {
    C.ReleaseEverything(vb.session, vb.vbox, vb.vboxclient)
    C.ClientUninitialize()
    C.VBoxCGlueTerm()
}

type VirtualMachine struct {
    Machine  *C.IMachine
    Name string
    Session  *C.ISession
    locked   bool

    Console  *C.IConsole

    Display  *C.IDisplay
    Mouse    *C.IMouse
    Keyboard *C.IKeyboard
}

// TODO: name is actually nameOrId
func (vb *VirtualBox) ConnectToMachine(name string) (*VirtualMachine, error) {
    var machine *C.IMachine

    cs := C.CString(name)
    defer C.free(unsafe.Pointer(cs))
    var rc C.HRESULT = C.FindMachine(vb.vbox, cs, &machine)
    if C.MY_FAILED(rc) != 0 || machine == nil {
        if machine != nil {
            defer C.free(unsafe.Pointer(machine))
        }

        message := fmt.Sprintf("Could not get machine reference")
        log.Println(message)
        return nil, errors.New(message)
    }

    vm := new(VirtualMachine)
    vm.Machine = machine
    vm.Name = name
    vm.Session = vb.session
    vm.locked = false

    return vm, nil
}

/*
locking machine or starting it up:
see startVM in tstCAPIGlue.c
https://www.virtualbox.org/sdkref/interface_i_console.html#_details
https://www.virtualbox.org/sdkref/interface_i_mouse.html#_details
*/
// I guess this one expects the vm to be active already, the other possibility would be "LaunchVMProcess"
func (vm *VirtualMachine) LockMachine() error {
    var rc C.HRESULT = C.LockMachine(vm.Machine, vm.Session)
    if C.MY_FAILED(rc) != 0 {
        vm.locked = false
        message := fmt.Sprintf("Could not lock machine")
        log.Println(message)
        return errors.New(message)
    }
    vm.locked = true

    rc = C.GetConsole(vm.Session, &vm.Console)
    if C.MY_FAILED(rc) != 0 || vm.Console == nil {
        vm.UnlockMachine()
        message := fmt.Sprintf("Could not get console reference")
        log.Println(message)
        return errors.New(message)
    }

    rc = C.GetMouse(vm.Console, &vm.Mouse)
    if C.MY_FAILED(rc) != 0 || vm.Mouse == nil {
        vm.UnlockMachine()
        message := fmt.Sprintf("Could not get mouse reference")
        log.Println(message)
        return errors.New(message)
    }

    rc = C.GetKeyboard(vm.Console, &vm.Keyboard)
    if C.MY_FAILED(rc) != 0 || vm.Keyboard == nil {
        vm.UnlockMachine()
        message := fmt.Sprintf("Could not get keyboard reference")
        log.Println(message)
        return errors.New(message)
    }

    rc = C.GetDisplay(vm.Console, &vm.Display)
    if C.MY_FAILED(rc) != 0 || vm.Display == nil {
        vm.UnlockMachine()
        message := fmt.Sprintf("Could not get display reference")
        log.Println(message)
        return errors.New(message)
    }

    return nil
}

func (vm *VirtualMachine) UnlockMachine() (error) {
    var rc C.HRESULT = C.UnlockMachine(vm.Session)

    if C.MY_FAILED(rc) != 0 {
        message := fmt.Sprintf("Could not unlock machine")
        log.Println(message)
        return errors.New(message)
    }

    vm.locked = false
    return nil
}

const (
    MOUSE_EVENT_1_DOWN = 0x1
    MOUSE_EVENT_2_DOWN = 0x2
    MOUSE_EVENT_ALL_UP = 0x0
)

//state should be one of MOUSE_EVENT_...
func (vm *VirtualMachine) MouseEvent(x,y, state int) (error) {
    if !vm.locked {
        message := fmt.Sprintf("Machine not locked, can't perform mouse down")
        log.Println(message)
        return errors.New(message)
    }

    var rc C.HRESULT = C.MouseClickEvent(vm.Mouse, C.int(x),C.int(y), C.int(state));
    if C.MY_FAILED(rc) != 0 {
        message := fmt.Sprintf("Failed to perform mouse down")
        log.Println(message)
        return errors.New(message)
    }

    return nil
}

func (vm *VirtualMachine) GetScreenshot() (*image.NRGBA, error) {
    if !vm.locked {
        message := fmt.Sprintf("Machine not locked, can't take screenshot")
        log.Println(message)
        return nil, errors.New(message)
    }

    width := C.uint(0)
    height := C.uint(0)
    bitsPerPixel := C.uint(0)

    //TODO: is vm.Display threadsafe?
    var rc C.HRESULT = C.GetScreenResolution(vm.Display, &width, &height, &bitsPerPixel);
    if C.MY_FAILED(rc) != 0 {
        message := fmt.Sprintf("Failed to get screen resolution")
        log.Println(message)
        return nil, errors.New(message)
    }

    //var imageData [int(width)*(height)*4]char
    var imageDataSize    C.PRUint32
    var imageData       *C.PRUint8
    defer C.free(unsafe.Pointer(imageData))

    rc = C.TakeScreenShotPNGToArray(vm.Display, width, height, &imageDataSize, &imageData)
    if C.MY_FAILED(rc) != 0 {
        message := fmt.Sprintf("Failed to get screen shot data")
        log.Println(message)
        return nil, errors.New(message)
    }

    // convert data to a go byte array
    dataSize := int(imageDataSize)
    dataBytes := C.GoBytes(unsafe.Pointer(imageData), C.int(dataSize))

    reader := bytes.NewReader(dataBytes)
    imgUnknown, err := png.Decode(reader) //this guy is probably RGBA
    if err != nil {
        message := fmt.Sprintf("Failed to decode image from png raw data", err)
        log.Println(message)
        return nil, err
    }

    // copy the image into NRGBA format. This is wasteful (TODO?)
    b := imgUnknown.Bounds()
    imgNRGBA := image.NewNRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
    draw.Draw(imgNRGBA, imgNRGBA.Bounds(), imgUnknown, b.Min, draw.Src)

    return imgNRGBA, nil
}
