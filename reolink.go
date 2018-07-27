/* Reolink Camera Control.
 *
 * This is intended to control at least some aspects of the Reolink IP
 * cameras.  The functionality contained herein has been tested against the
 * RLC-411WS camera.
 *
 * Not all methods are implemented.  You also have the misfortune of reading
 * code that is my attempt at learning Go; proceed at your own risk!
 *
 * Copyright (c) 2018 Jeff Cody <jeff@codyprime.org>
 *
 * This program is free software; you can redistribute it and/or modify it under
 * the terms of the GNU General Public License as published by the Free Software
 * Foundation; under version 2 of the license only.
 *
 * This program is distributed in the hope that it will be useful, but WITHOUT
 * ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS
 * FOR A PARTICULAR PURPOSE.  See the GNU General Public License for more
 * details.
 *
 * You should have received a copy of the GNU General Public License along with
 * this program; if not, see <http://www.gnu.org/licenses/gpl-2.0.html>.
 */

package main

import (
    "net/http"
    "crypto/tls"
    "bytes"
    "fmt"
    "io/ioutil"
    "encoding/json"
    "flag"
    "log"
)


// json API structs
type Token struct {
    LeaseTime int   `json:"leaseTime"`
    Name string     `json:"name"`
}

type MaxMin struct {
    max   int      `json:"max"`
    min   int      `json:"min"`
}

type Isp struct {
    AntiFlicker     string  `json:"antiFlicker"`
    Backlight       string  `json:"backLight"`
    BLC             int     `json:"blc"`
    BlueGain        int     `json:"blueGain"`
    Channel         int     `json:"channel"`
    DayNight        string  `json:"dayNight"`
    DRC             int     `json:"drc"`
    Exposure        string  `json:"exposure"`
    Gain            MaxMin  `json:"gain"`
    Mirroring       int     `json:"mirroring"`
    NR3D            int     `json:"nr3d"`
    RedGain         int     `json:"redGain"`
    Rotation        int     `json:"rotation"`
    Shutter         MaxMin  `json:"shutter"`
    WhiteBalance    string  `json:"whiteBalance"`
}

type Value struct {
    Token Token     `json:"Token"`
    Isp Isp         `json:"Isp"`
}

type CameraResponse struct {
    Cmd     string  `json:"cmd"`
    Code    int     `json:"code"`
    Initial Value   `json:"initial"`
    //Range   Value   `json:"range"`
    Value   Value   `json:"value"`
}

// Send a command to the camera, and returned the response in the appropriate struct
// after unmarshaling the json response.
func sendCmd(ip string, token string, cmd string, jsonStr string) []CameraResponse{

    url := "https://" + ip + "/cgi-bin/api.cgi?cmd=" + cmd + "&token=" + token

    req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(jsonStr)))
    if err != nil {
        panic(err)
    }
    req.Header.Set("Content-Type", "application/json")

    cfg := &http.Transport{ TLSClientConfig: &tls.Config{InsecureSkipVerify: true} }

    client := &http.Client{Transport: cfg}
    resp, err := client.Do(req)
    defer resp.Body.Close()
    if err != nil {
        panic(err)
    }

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        panic(err)
    }

    var camResp []CameraResponse

    err = json.Unmarshal([]byte(body), &camResp)
    if err != nil {
        panic(err)
    }

    return camResp
}

// All APIs (except this one) require a token.  A token is obtained via
// the 'Login' command, with the appropriate username/password combo.
func getToken(ip string, user string, passwd string) string {
    jsonStr := fmt.Sprintf(`[{"cmd":"Login","action":0,"param":{"User":{"userName":"%s","password":"%s"}}}]`, user, passwd)

    camResp := sendCmd(ip, "null", "Login", jsonStr)

    return camResp[0].Value.Token.Name
}

// TODO: Have a single 'get' that grabs all supported structures / parameters

// See the 'Isp' struct above for the fields contained; this structure is
// used for both setting and getting the parameters.
//
// Notable parameters are day/night, exposure, shutter control, etc..
func getIsp(ip string, token string) *Isp {
    jsonStr := `[{"cmd":"GetIsp","action":1,"param":{"channel":0}}]`

    camResp := sendCmd(ip, token, "GetIsp", jsonStr)
    return &camResp[0].Value.Isp
}

func setIsp(ip string, token string, isp *Isp) {
    b, err := json.Marshal(isp)
    if err != nil {
        panic(err)
    }

    jsonStr := fmt.Sprintf(`[{"cmd":"SetIsp","action":0,"param":{"Isp":%s}}]`, string(b))
    sendCmd(ip, token, "SetIsp", jsonStr)
}


func main() {

    ip     := flag.String("ip", "192.168.15.124", "ip address of camera")
    user   := flag.String("user", "admin", "username")
    //TODO: add option to read username/password from file
    passwd := flag.String("password", "password", "password")
    token  := flag.String("token", "", "token")
    cmd    := flag.String("cmd", "nil", "command to issue")
    data   := flag.String("data", "", "data associated with 'cmd'")

    flag.Parse()

    if *token == "" {
        *token = getToken(*ip, *user, *passwd)
        fmt.Println(*token)
    }


    // TODO: Factor the cmds into something more scaleable (e.g. table lookup.
    //       Right now, just initial PoC for talking to the camera.
    switch *cmd {
    // no-op, currently happens by default if no token specified
    case "get-token":

    // daynight expects 'data' to be one of: 'day', 'night', 'auto'
    case "set-daynight":
        var isp *Isp
        isp = getIsp(*ip, *token)
        fmt.Println("previous setting: ", isp.DayNight)
        if *data == "day" {
            isp.DayNight = "Color"
        } else if *data == "night" {
            isp.DayNight = "Black&White"
        } else if *data == "auto" {
            isp.DayNight = "Auto"
        } else {
            log.Fatal("data for 'set-daynight' should be 'day', 'night', or 'auto'")
        }
        setIsp(*ip, *token, isp)
        isp = getIsp(*ip, *token)
        fmt.Println("current setting: ", isp.DayNight)
    case "get-daynight":
        var isp *Isp
        isp = getIsp(*ip, *token)
        fmt.Println("current Day/Night setting: ", isp.DayNight)
    default:
        // no-op
    }

}
