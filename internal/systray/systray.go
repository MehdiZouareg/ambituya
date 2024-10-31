package systray

import (
	_ "embed"
	"fmt"
	"os"

	"github.com/MehdiZouareg/ambituya/internal/tuya"
	"github.com/getlantern/systray"
)

func Systray() {
	systray.Run(onReady, onExit)
}

var menuItemsPtr []*systray.MenuItem
var menuOptions []MenuOption
var menuItems []MenuItem
var programPath string

type MenuOption struct {
	label    string
	command  func() error
	deviceID string
}

type MenuItem struct {
	menuItemType MenuItemType
	label        string
	command      func() error
}

type MenuItemType int64

const (
	Choice    MenuItemType = 0
	Separator              = 1
)

func onReady() {
	mainIcon := loadIcon()
	systray.SetIcon(mainIcon)
	systray.SetTitle("LightControl")
	systray.SetTooltip("Lumières de la maison")

	menuItemsPtr = make([]*systray.MenuItem, 0)
	menuItems := make([]MenuItem, 0)
	options := make([]MenuOption, 0)

	///////////////////////////////////////////////////////////////////////////////
	// AMBILIGHT
	///////////////////////////////////////////////////////////////////////////////

	///////////////////////////////////////////////////////////////////////////////
	// SWITCH LIGHTS
	///////////////////////////////////////////////////////////////////////////////

	for _, dev := range tuya.RegisteredDevices {
		menuItem := MenuItem{
			label:        fmt.Sprintf("ON/OFF %v", dev.Name),
			command:      dev.Switch,
			menuItemType: Choice,
		}
		menuItems = append(menuItems, menuItem)

		option := MenuOption{
			label: dev.Name,
			command: func(device tuya.Device) func() error {
				return func() error {
					return device.Switch()
				}
			}(dev),
			deviceID: dev.ID,
		}

		options = append(options, option)
	}

	menuItemsPtr = make([]*systray.MenuItem, 0)

	indexOption := 0
	for _, v := range menuItems {
		if v.menuItemType == Separator {
			systray.AddSeparator()
			continue
		}
		menuItemPtr := systray.AddMenuItem(options[indexOption].label, options[indexOption].label)

		// var icon []byte
		// status, err := v.device.GetDeviceStatus("switch_led")
		// if err != nil {
		// 	tuyalogger.Log.Error("ERROR ON GETTING DEVICE STATUS")
		// 	continue
		// }
		// if status != nil {
		// 	fmt.Printf("%v is %v", v.device.ID, status.(bool))
		// }

		// menuItemPtr.SetIcon(icon)
		menuItemsPtr = append(menuItemsPtr, menuItemPtr)

		indexOption++
	}

	cmdChan := make(chan func() error)

	///////////////////////////////////////////////////////////////////////////////
	systray.AddSeparator()

	///////////////////////////////////////////////////////////////////////////////
	// QUIT AND SETTINGS
	///////////////////////////////////////////////////////////////////////////////

	mSettings := systray.AddMenuItem("Settings", "Settings")
	mQuit := systray.AddMenuItem("Quit", "Quit the whole app")

	for i, menuItenPtr := range menuItemsPtr {
		go func(c chan struct{}, cmd func() error) {
			for range c {
				cmdChan <- cmd
			}
		}(menuItenPtr.ClickedCh, options[i].command)
	}

	go func() {
		for {
			select {
			case cmd := <-cmdChan:
				execute(cmd) // Handle command
			case <-mSettings.ClickedCh:
				systray.Quit()
				os.Exit(1)
				return
			case <-mQuit.ClickedCh:
				systray.Quit()
				os.Exit(1)
				return
			}
		}
	}()
}

func execute(commands func() error) {
	// command_array := strings.Split(commands, " ")
	// command := ""
	// command, command_array = command_array[0], command_array[1:]
	// cmd := exec.Command(command, command_array...)
	// var out bytes.Buffer
	// cmd.Stdout = &out
	// err := cmd.Run()
	// if err != nil {
	// 	errMsg := fmt.Sprintf("Error executing command: %s", err)
	// 	fmt.Println(errMsg)
	// 	dialog.Error(errMsg)
	// 	// log.Fatal(err)
	// }
	// fmt.Printf("Output: %s\n", commands)

	commands()
}

func onExit() {
	// clean up here
}

//go:embed icons/house.ico
var icon []byte

func loadIcon() []byte {
	return icon
}
