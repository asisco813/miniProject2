package main

import (
	"fmt"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/aio"
	"gobot.io/x/gobot/drivers/i2c"
	g "gobot.io/x/gobot/platforms/dexter/gopigo3"
	"gobot.io/x/gobot/platforms/raspi"
	"time"
)

func robotMainLoop(piProcessor *raspi.Adaptor, gopigo3 *g.Driver, lidarSensor *i2c.LIDARLiteDriver,

) {

	err := lidarSensor.Start()
	if err != nil {
		fmt.Println("error starting lidarSensor")
	}
	lidarReading, err := lidarSensor.Distance()
	var secondCount = 0

		for lidarReading > 100 {
			lidarReading, err := lidarSensor.Distance()
			if err != nil {
				fmt.Println("Error reading lidar sensor %+v", err)
			}
			message := fmt.Sprintf("Lidar Reading: %d", lidarReading)

			fmt.Println(lidarReading)
			fmt.Println(message)
			time.Sleep(time.Second)

			err = gopigo3.SetMotorDps(g.MOTOR_LEFT+g.MOTOR_RIGHT, 200)
			if err != nil {
				fmt.Println("Error moving forward %+v", err)
			}
		}
		for lidarReading < 100 {
			lidarReading, err := lidarSensor.Distance()
			if err != nil {
				fmt.Println("Error reading lidar sensor %+v", err)
			}
			message := fmt.Sprintf("Lidar Reading: %d", lidarReading)

			fmt.Println(lidarReading)
			fmt.Println(message)
			time.Sleep(time.Second)

			err = gopigo3.SetMotorDps(g.MOTOR_LEFT+g.MOTOR_RIGHT, 200)
			if err != nil {
				fmt.Println("Error moving forward %+v", err)
			}
			secondCount += 1
		}
	err = gopigo3.SetMotorDps(g.MOTOR_LEFT+g.MOTOR_RIGHT, 0)
	if err != nil {
		fmt.Println("Error stopping %+v", err)
	}
	var lengthOfBox float64 = float64(secondCount) * 200 * .5803
	fmt.Println("The length of the box is: %d", lengthOfBox)
}

func main() {
	raspberryPi := raspi.NewAdaptor()
	gopigo3 := g.NewDriver(raspberryPi)
	lidarSensor := i2c.NewLIDARLiteDriver(raspberryPi)
	lightSensor := aio.NewGroveLightSensorDriver(gopigo3, "AD_2_1")
	workerThread := func() {
		robotMainLoop(raspberryPi, gopigo3, lidarSensor)
	}
	robot := gobot.NewRobot("Gopigo Pi4 Bot",
		[]gobot.Connection{raspberryPi},
		[]gobot.Device{gopigo3, lidarSensor, lightSensor},
		workerThread,
	)

	err := robot.Start()

	if err != nil {
		fmt.Errorf("Error starting Robot #{err}")
	}
}
