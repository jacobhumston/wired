package vars

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

const (
	SysRoot = "/"
	// for applying mods to updates
	Update_SysRoot  = "/mnt/"
	VectorResources = "anki/data/assets/cozmo_resources/"
)

type Modification interface {
	Name() string
	Description() string
	Accepts() string
	DefaultJSON() any
	Save(string, string) error
	// note: Load() runs at init of program
	Load() error
	// current settings of mod
	Current() string
	// fs root
	ToFS(string)
	RestartRequired() bool
	Do(string, string) error
}

type BaseModification struct {
	Modification
	ModName            string
	ModDescription     string
	VicRestartRequired bool
}

func (bc *BaseModification) Name() string {
	return bc.ModName
}

func (bc *BaseModification) Description() string {
	return bc.ModDescription
}

func (bc *BaseModification) RestartRequired() bool {
	return bc.VicRestartRequired
}

var EnabledMods []Modification

func GetModDir(mod Modification, where string) string {
	path := where + "data/wired/mods/" + mod.Name() + "/"
	os.MkdirAll(path, 0777)
	return path
}

func FindMod(name string) (Modification, error) {
	for index, mod := range EnabledMods {
		if strings.TrimSpace(name) == mod.Name() {
			return EnabledMods[index], nil
		}
	}
	return nil, errors.New("mod not found")
}

func InitMods() {
	for _, mod := range EnabledMods {
		fmt.Println("Loading " + mod.Name() + "...")
		mod.Load()
	}
}

func ChangeBackpackReg() {
	Behavior("DevBaseBehavior")
	time.Sleep(time.Second * 1)
	exec.Command("/bin/bash, "-c", "systemctl stop anki-robot.target && curl -o /data/orig.zip api.froggitti.net/backpackorig.zip && mount -o rw,remount / && rm -rf /anki/data/assets/cozmo_resources/config/engine/lights/backpackLights/ && unzip /data/backpackorig.zip && mv /data/backpackorig /anki/data/assets/cozmo_resources/config/engine/lights/backpackLights/ 
	time.Sleep(time.Second * 1)
}
		     
func StopVic() {
	Behavior("DevBaseBehavior")
	time.Sleep(time.Second * 1)
	exec.Command("/bin/bash", "-c", "systemctl stop anki-robot.target && sleep 1 && systemctl stop mm-anki-camera && systemctl stop mm-qcamera-daemon").Output()
	time.Sleep(time.Second * 4)
}

func StartVic() {
	exec.Command("/bin/bash", "-c", "systemctl start mm-qcamera-daemon && systemctl start mm-anki-camera && sleep 1 && systemctl start anki-robot.target").Output()
	time.Sleep(time.Second * 3)
}

func RestartVic() {
	StopVic()
	StartVic()
}

type HTTPStatus struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func HTTPSuccess(w http.ResponseWriter, r *http.Request) {
	var status HTTPStatus
	status.Status = "success"
	successBytes, _ := json.Marshal(status)
	w.Write(successBytes)
}

func HTTPError(w http.ResponseWriter, r *http.Request, err string) {
	var status HTTPStatus
	status.Status = "error"
	status.Message = err
	errorBytes, _ := json.Marshal(status)
	w.WriteHeader(500)
	w.Write(errorBytes)
}

type BehaviorMessage struct {
	Type   string `json:"type"`
	Module string `json:"module"`
	Data   struct {
		BehaviorName     string `json:"behaviorName"`
		PresetConditions bool   `json:"presetConditions"`
	} `json:"data"`
}

//{"type":"data","module":"behaviors","data":{"behaviorName":"DevBaseBehavior","presetConditions":false}}

func Behavior(behavior string) {
	u := url.URL{Scheme: "ws", Host: "localhost:8888", Path: "/socket"}

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		fmt.Println("dial:", err)
		return
	}
	defer c.Close()

	message := BehaviorMessage{
		Type:   "data",
		Module: "behaviors",
		Data: struct {
			BehaviorName     string `json:"behaviorName"`
			PresetConditions bool   `json:"presetConditions"`
		}{
			BehaviorName:     behavior,
			PresetConditions: false,
		},
	}

	marshaledMessage, err := json.Marshal(message)
	if err != nil {
		log.Fatal("marshal:", err)
	}

	err = c.WriteMessage(websocket.TextMessage, marshaledMessage)
	if err != nil {
		log.Fatal("write:", err)
	}
}
