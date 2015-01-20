package controller

import (
    "log"
    "testing"
    "time"

    "th4t/vboxgo" //only for the test constant
    "th4t/desktopcontroller"
)

// just call a few functions, and see if they crash and burn
func TestSmoke(t *testing.T) {
    vmController := NewVBoxGoController()
    // this can be performed after requesting a screenshot, but it's no issue
    go vmController.Start(vboxgo.TEST_VM_NAME)

    log.Println("Fetching a screenshot")
    img := vmController.RequestScreenshot()
    bounds := img.Bounds()
    log.Printf("Got a screenshot with dimensions [%d,%d]\n", bounds.Max.X,bounds.Max.Y)

    log.Println("Performing a click activity")
    clickX, clickY := 0, 0
    duration := time.Duration(50 * time.Millisecond)
    activity := &desktopcontroller.SimpleActivity {
        []desktopcontroller.Action{
        desktopcontroller.Action{desktopcontroller.ACTION_MOUSE_DOWN, clickX, clickY, 0},
        desktopcontroller.Action{desktopcontroller.ACTION_SLEEP, clickX, clickY, duration},
        desktopcontroller.Action{desktopcontroller.ACTION_MOUSE_UP, clickX, clickY, 0},
        },
    }
    vmController.RequestActivity(activity)
    log.Println("Click executed")
}
