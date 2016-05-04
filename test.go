//
// @Author: Geoffrey Bauduin <bauduin.geo@gmail.com>
//

package main

import (
    "fmt"
    "./pushnotifications"
)

func main () {
    accessKey := ""
    secretKey := ""
    sender := pushnotifications.NewPushNotification(accessKey, secretKey, "eu-west-1")
    deviceArn := ""
    text := "Hello world!"
    badge := 5
    err := sender.Send(deviceArn, &pushnotifications.Push{
        Alert: &pushnotifications.Alert{
            Body: &text,
        },
        Badge: &badge, 
    })
    if err != nil {
        fmt.Printf("%v", err)
    }
}