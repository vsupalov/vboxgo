package controller

import (
    "time"
    "log"
    "image"

    "th4t/desktopcontroller"
    "th4t/vboxgo"
)

type ActivityResponseStruct struct {
    //an ack is sent on this channel, as soon as the activity has been executed
    ResponseChan chan bool
    Activity desktopcontroller.Activity
}

type VBoxGoController struct {
    vbox *vboxgo.VirtualBox
    machine *vboxgo.VirtualMachine

    inChanRequestScreenshot chan chan *image.NRGBA
    inChanRequestActivity chan *ActivityResponseStruct
}

func NewVBoxGoController() (*VBoxGoController) {
    c := new(VBoxGoController)

    c.inChanRequestScreenshot = make(chan chan *image.NRGBA, 0)
    c.inChanRequestActivity = make(chan *ActivityResponseStruct, 0)

    return c
}

func (c *VBoxGoController) RequestScreenshot() *image.NRGBA {
    responseChan := make(chan *image.NRGBA, 0)
    //TODO: maybe this does not need to be blocking? Potential speedup.
    c.inChanRequestScreenshot <- responseChan
    img := <-responseChan
    return img
}

func (c *VBoxGoController) RequestActivity(activity desktopcontroller.Activity) {
    responseChan := make(chan bool, 0)
    ar := ActivityResponseStruct{responseChan, activity}
    c.inChanRequestActivity <- &ar
    <-responseChan //wait for the command to finish
}

// starts the listening loop, should be called with "go"
// TODO: way to shut it down gracefully?
func (c *VBoxGoController) Start(machineName string) {
    log.Println("Initializing vbox")

    // establish connection with virtualbox
    vbox := vboxgo.VirtualBox{}
    if err := vbox.Init(); err != nil {
        log.Fatalf("Vboxgo init failed with: %v\n", err)
    }
    // TODO: grab errors?
    defer c.vbox.Uninit()

    machine, err := vbox.ConnectToMachine(machineName)
    if err != nil {
        log.Fatalf("Failed to connect to VM: %v\n", err)
    }

    // TODO: grab errors?
    err = machine.LockMachine()
    if err != nil {
        log.Fatalf("Failed to lock VM: %v\n", err)
    }
    defer c.machine.UnlockMachine()

    c.vbox = &vbox
    c.machine = machine

    for {
        select {
            case outChan := <-c.inChanRequestScreenshot:
                img, err := c.machine.GetScreenshot()
                if err != nil {
                    log.Fatalf("Failed to get image: %v\n", err)
                }
                outChan <- img
            case activityResponse := <-c.inChanRequestActivity:
                for a := range(activityResponse.Activity.GetActions()) {
                    // perform actions one after another
                    switch(a.Type) {
                        case desktopcontroller.ACTION_MOUSE_DOWN:
                            err := c.machine.MouseEvent(a.X, a.Y, vboxgo.MOUSE_EVENT_1_DOWN)
                            if err != nil {
                                log.Fatalf("Failed to pass mouse event: %v\n", err)
                            }
                        case desktopcontroller.ACTION_MOUSE_UP:
                            err := c.machine.MouseEvent(a.X, a.Y, vboxgo.MOUSE_EVENT_ALL_UP)
                            if err != nil {
                                log.Fatalf("Failed to pass mouse event: %v\n", err)
                            }
                        case desktopcontroller.ACTION_SLEEP:
                            time.Sleep(a.Duration)
                        default:
                            log.Fatalln("Unknown action type in", a)
                    }
                }
                activityResponse.ResponseChan <- true
        }
    }
}
