package bsp

func TurnOnLed(devName string) {
	WriteLedFile(devName, 100)
}

func TurnOffLed(devName string) {
	WriteLedFile(devName, 0)
}
