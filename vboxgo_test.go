package vboxgo

import (
    "testing"
    "image"
)

func TestLifecycle(t *testing.T) {
    vbox := VirtualBox{}

    if err := vbox.Init(); err != nil {
        t.Errorf("Init failed with: %v", err)
    }

    vbox.Uninit()
}

func TestConnect(t *testing.T) {
    vbox := VirtualBox{}

    if err := vbox.Init(); err != nil {
        t.Errorf("Init failed with: %v", err)
    }

    _, err := vbox.ConnectToMachine(TEST_VM_NAME)

    if err != nil {
        t.Errorf("Failed to connect to VM: %v", err)
    }

    vbox.Uninit()
}

func TestLockMachine(t *testing.T) {
    vbox := VirtualBox{}

    if err := vbox.Init(); err != nil {
        t.Errorf("Init failed with: %v", err)
    }

    machine, err := vbox.ConnectToMachine(TEST_VM_NAME)

    if err != nil {
        t.Errorf("Failed to connect to VM: %v", err)
    }

    err = machine.LockMachine()

    if err != nil {
        t.Errorf("Failed to lock VM: %v", err)
    }

    err = machine.UnlockMachine()
    if err != nil {
        t.Errorf("Failed to unlock VM: %v", err)
    }

    vbox.Uninit()
}

func TestMouse(t *testing.T) {
    vbox := VirtualBox{}

    if err := vbox.Init(); err != nil {
        t.Errorf("Init failed with: %v", err)
    }

    machine, err := vbox.ConnectToMachine(TEST_VM_NAME)

    if err != nil {
        t.Errorf("Failed to connect to VM: %v", err)
    }

    err = machine.LockMachine()

    if err != nil {
        t.Errorf("Failed to lock VM: %v", err)
    }

    clickX, clickY := 50, 50
    err = machine.MouseEvent(clickX, clickY, MOUSE_EVENT_1_DOWN) //click somewhere
    err = machine.MouseEvent(clickX, clickY, MOUSE_EVENT_ALL_UP)
    if err != nil {
        t.Errorf("Failed to press down mouse: %v", err)
    }

    err = machine.UnlockMachine()
    if err != nil {
        t.Errorf("Failed to unlock VM: %v", err)
    }

    vbox.Uninit()
}

func TestScreenshot(t *testing.T) {
    vbox := VirtualBox{}

    if err := vbox.Init(); err != nil {
        t.Errorf("Init failed with: %v", err)
    }

    machine, err := vbox.ConnectToMachine(TEST_VM_NAME)

    if err != nil {
        t.Errorf("Failed to connect to VM: %v", err)
    }

    err = machine.LockMachine()

    if err != nil {
        t.Errorf("Failed to lock VM: %v", err)
    }

    var img image.Image
    img, err = machine.GetScreenshot()
    if err != nil {
        t.Errorf("Failed to get image: %v", err)
    }

    if img == nil {
        t.Errorf("Image pointer is nil: %v", err)
    }

    /*
    // save the image, just for laughs and giggles
    toimg, _ := os.Create("/tmp/vbox_test.png")
    defer toimg.Close()
    png.Encode(toimg, *img)
    */

    //TODO: check that image is not 'boring'

    err = machine.UnlockMachine()
    if err != nil {
        t.Errorf("Failed to unlock VM: %v", err)
    }

    vbox.Uninit()
}
