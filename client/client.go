package main

import (
	"bufio"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/jessevdk/go-flags"
	log "github.com/sirupsen/logrus"
	"go-gnet/constant"
	"go-gnet/database/mysql/model"
	"go-gnet/server/handler"
	"go-gnet/server/protocol"
	"go-gnet/util"
	"math/rand"
	"net"
	"os"
	"strings"
	"time"
)

var opts struct {
	Uuid   string `long:"uuid" required:"true" description:"设备登录UUID"`
	Key    string `long:"key"  description:"设备登录密钥"`
	Host   string `long:"host" default:"127.0.0.1" description:"通讯网关IP地址"`
	Port   int    `long:"port" default:"9000" description:"通讯网关端口"`
	Config string `long:"config" default:"json" choice:"xml" choice:"json" choice:"yaml" description:"模拟该设备的配置文件格式"`
	Debug  bool   `long:"debug" description:"打印调试信息"`
}

func main() {
	_, err := flags.Parse(&opts)
	if err != nil {
		panic(err)
	}
	util.InitUtil()
	if !opts.Debug {
		log.SetLevel(log.InfoLevel)
	}
	log.Info("启动边缘计算设备模拟器...")
	// tcp
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", opts.Host, opts.Port))
	if err != nil {
		panic(err)
	}

	/* 不同格式的配置 */

	configJson := `{
    "configuration": {
        "status": "error",
        "name": "JSONConfigDemo",
        "packages": "com.howtodoinjava",
        "ThresholdFilter": {
            "level": "debug"
        },
        "appenders": {
            "Console": {
                "name": "STDOUT",
                "PatternLayout": {
                    "pattern": "%d [%t] %-5p %c - %m%n"
                }
            }
        },
        "loggers": {
            "root": {
                "level": "debug",
                "AppenderRef": {
                    "ref": "STDOUT"
                }
            }
        }
    }
}`
	configXml := `<settings>
  <setting name="LogEntries">200</setting>
  <setting name="LogListViewColumns"></setting>
  <setting name="LogWindowPosition">300,300</setting>
  <setting name="LogWindowSize">@120|450,500</setting>
  <setting name="ColorRemoved">283cff</setting>
  <setting name="UseColorOwnProcesses">1</setting>
  <setting name="ColorOwnProcesses">aaffff</setting>
  <setting name="UseColorSystemProcesses">1</setting>
  <setting name="ColorSystemProcesses">ffccaa</setting>
  <setting name="UseColorServiceProcesses">1</setting>
  <setting name="ProcessHacker.ToolStatus.ToolbarConfig"></setting>
  <setting name="ProcessHacker.ToolStatus.StatusbarConfig"></setting>
  <setting name="ProcessHacker.UpdateChecker.PromptStart">1</setting>
  <setting name="ProcessHacker.UpdateChecker.LastUpdateCheckTime">132951649178685166</setting>
  <setting name="ProcessHacker.UserNotes.DatabasePath">%APPDATA%\Process Hacker 2\usernotesdb.xml</setting>
  <setting name="ProcessHacker.UserNotes.ColorCustomList"></setting>
  <setting name="ProcessHacker.WindowExplorer.ShowDesktopWindows">0</setting>
  <setting name="ProcessHacker.WindowExplorer.WindowTreeListColumns"></setting>
  <setting name="ProcessHacker.WindowExplorer.WindowsWindowPosition">100,100</setting>
  <setting name="ProcessHacker.WindowExplorer.WindowsWindowSize">@120|690,540</setting>
  <setting name="MemoryListViewColumns"></setting>
  <setting name="UseColorPosixProcesses">1</setting>
  <setting name="ColorPosixProcesses">8b3d48</setting>
  <setting name="ProcessHacker.ToolStatus.EnableToolBar">1</setting>
  <setting name="ProcessHacker.ToolStatus.EnableSearchBox">1</setting>
  <setting name="ProcessHacker.ToolStatus.EnableStatusBar">1</setting>
  <setting name="ProcessHacker.ToolStatus.EnableWicImaging">1</setting>
  <setting name="ProcessHacker.ToolStatus.ResolveGhostWindows">1</setting>
  <setting name="ProcessHacker.ToolStatus.StatusMask">d</setting>
  <setting name="ProcessHacker.Updater.PromptStart">1</setting>
</settings>
`
	configYaml := `lol-audio:
    data:
        ambienceMusicEnabled: false
        championSelectionMusicEnabled: false
        loginMusicEnabled: false
        musicEnabled: true
        musicVolume: 82
        sfxEnabled: true
        sfxVolume: 82
        voiceEnabled: true
        voiceVolume: 41
    schemaVersion: 5
lol-clash:
    data:
        simpleStateFlagIds:
        - "registration_open_tournament_1521"
        - "lock_in_open_tournament_1581_phase_1581"
        - "registration_open_tournament_1501"
        - "registration_open_tournament_1581"
        - "lock_in_open_tournament_1661_phase_1661"
        - "lock_in_open_tournament_1621_phase_1621"
        - "registration_open_tournament_1641"
        - "lock_in_open_tournament_1641_phase_1641"
        - "registration_open_tournament_1621"
        - "registration_open_tournament_1681"
        - "lock_in_open_tournament_1681_phase_1681"
        - "registration_open_tournament_1661"
        - "registration_open_tournament_1701"
    schemaVersion: 1
lol-home:
    data:
        useTFTOverride: false
    schemaVersion: 0
lol-premade-voice:
    data:
        currentCaptureDeviceHandle: "{173DAB58-AB46-433A-BEB5-4395FF7ECD5D}"
        inputVolume: 100
        vadSensitivity: 65
    schemaVersion: 1
lol-replays:
    data:
        dontShowEogUsageTip: true
    schemaVersion: 1
lol-user-experience:
    data:
        potatoModeEnabled: true
    schemaVersion: 3
video:
    data:
        ZoomScale: 1
    schemaVersion: 1
`
	config := configJson
	if opts.Config == "xml" {
		config = configXml
	} else if opts.Config == "yaml" {
		config = configYaml
	}

	/* 发包函数 */

	codec := protocol.PacketCodec{}
	send := func(code uint16, d interface{}) {
		context, err := json.Marshal(d)
		if err != nil {
			log.Error(err)
			return
		}
		p := protocol.PacketDataReceived{
			Code: code,
			Data: context,
		}
		bytes, err := json.Marshal(p)
		if err != nil {
			log.Error("Marshal error: ", err)
			return
		}
		encode, err := codec.Encode(nil, bytes)
		if err != nil {
			log.Error("Encode error: ", err)
			return
		}
		_, err = conn.Write(encode)
	}

	/* 自动回包函数 */

	autoRet := func() {
		for {
			_version := make([]byte, 2)
			n, err := conn.Read(_version)
			if err != nil {
				panic(err)
			}
			log.Debugf("read %d", n)
			_dataLength := make([]byte, 4)
			n, err = conn.Read(_dataLength)
			if err != nil {
				panic(err)
			}
			log.Debugf("read %d", n)
			// version := binary.BigEndian.Uint16(_version)
			dataLength := binary.BigEndian.Uint32(_dataLength)
			data := make([]byte, dataLength)
			n, err = conn.Read(data)
			if err != nil {
				panic(err)
			}
			if uint32(n) != dataLength {
				log.Errorf("read %d", n)
				panic(n)
			}
			log.Debugf("read %d", n)
			packet := protocol.PacketDataSent{}
			err = json.Unmarshal(data, &packet)
			if err != nil {
				panic(err)
			}
			switch packet.Code {
			case constant.ServerRequestUpdateDeviceConfig:
				send(constant.ClientUploadDeviceConfig, &model.DeviceConfig{
					Content: config,
					Type:    opts.Config,
				})
				log.Info("设备配置已上传")
				break
			case constant.ServerIssueDeviceConfig:
				log.Info("新的配置已应用")
				c := model.DeviceConfig{}
				err := json.Unmarshal(packet.Data, &c)
				if err != nil {
					panic(err)
				}
				config = c.Content
				// log.Info(c.Content)
				break
			default:
				log.Warnf("Unexpected packet code %d", packet.Code)
				break
			}
		}
	}
	go autoRet()

	// 发送认证信息
	log.Infof("发送认证 uuid=%s key=%s .", opts.Uuid, opts.Key)
	send(constant.ClientConnect, &handler.AuthPackage{
		Uuid:      opts.Uuid,
		SecretKey: opts.Key,
	})
	log.Info("已发送")

	/* 心跳 */

	type HeartbeatControl uint
	const (
		HeartbeatControlStart HeartbeatControl = iota
		HeartbeatControlStop
	)
	heartbeatChannel := make(chan HeartbeatControl, 100)

	/* 心跳协程 */

	heartbeatRoutine := func() {
		flag := true
		for {
			if len(heartbeatChannel) > 0 {
				control := <-heartbeatChannel
				switch control {
				case HeartbeatControlStop:
					flag = false
					log.Info("已关闭心跳")
					break
				case HeartbeatControlStart:
					flag = true
					log.Info("已开启心跳")
					break
				}
			}
			if !flag {
				time.Sleep(1 * time.Second)
				continue
			}
			packet := protocol.PacketDataSent{
				Code: constant.ClientHeartbeat,
				Msg:  "",
				Data: nil,
			}
			bytes, err := json.Marshal(packet)
			if err != nil {
				panic(err)
			}
			bytes, err = codec.Encode(nil, bytes)
			if err != nil {
				panic(err)
			}
			_, err = conn.Write(bytes)
			if err != nil {
				panic(err)
			}
			// log.Infof("Heartbeat %d", time.Now().Unix())
			time.Sleep(1 * time.Second)
		}
	}
	go heartbeatRoutine()

	/* 持续上报温度 */

	keepReportingTemperature := func() {
		base := 20 + float64(int64(1000*rand.Float64()))/100
		for {
			time.Sleep(1 * time.Second)
			change := float64(int64((rand.Float64()*0.2)*100))/100 - 0.1
			log.Debugf("温度涨幅: %.2f", change)
			base += change
			send(constant.ClientDataUploadTemperature, &model.DeviceDataTemperature{
				Time: time.Now(),
				Data: base,
			})
		}
	}
	go keepReportingTemperature()

	/* 进入命令行 */
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		text := scanner.Text()
		switch {
		case strings.HasPrefix(text, "udi"):
			var updateDeviceInfoOpts struct {
				SoftwareInfo string `short:"s"`
				HardwareInfo string `long:"hi"`
				Remark       string `short:"r"`
			}
			_args := strings.TrimPrefix(text, "udi")
			args := strings.Split(_args, " ")
			_, err = flags.ParseArgs(&updateDeviceInfoOpts, args)
			if err != nil {
				log.Error("Parse error: ", err)
				break
			}
			send(constant.ClientUpdateDeviceInfo, model.DeviceInfo{
				SoftwareInfo: updateDeviceInfoOpts.SoftwareInfo,
				HardwareInfo: updateDeviceInfoOpts.HardwareInfo,
				Remark:       updateDeviceInfoOpts.Remark,
			})
			log.Info("新的设备信息已上传")
			break
		case text == "start-heartbeat":
			heartbeatChannel <- HeartbeatControlStart
			break
		case text == "stop-heartbeat":
			heartbeatChannel <- HeartbeatControlStop
			break
		case strings.HasPrefix(text, "login"):
			params := strings.Split(text, " ")
			if len(params) < 3 {
				log.Info("Usage: login <username> <purpose>")
				break
			}
			send(constant.ClientEventUserLogin, model.DeviceEventLogin{
				Username:     params[1],
				LoginPurpose: params[2],
			})
			log.Info("登录信息已上传")
			break
		case strings.HasPrefix(text, "logout"):
			params := strings.Split(text, " ")
			if len(params) < 3 {
				log.Info("Usage: logout <username> <work-contents>")
				break
			}
			send(constant.ClientEventUserLogout, model.DeviceEventLogout{
				Username:     params[1],
				WorkContents: params[2],
			})
			log.Info("注销信息已上传")
			break
		case strings.HasPrefix(text, "operate"):
			params := strings.Split(text, " ")
			if len(params) < 4 {
				log.Info("Usage: operate <username> <type> <detail>")
				break
			}
			send(constant.ClientEventOperation, &model.DeviceEventOperation{
				Username:      params[1],
				OperationType: params[2],
				Detail:        params[3],
			})
			log.Info("操作记录已上传")
			break
		case strings.HasPrefix(text, "upload-log"):
			filename := "C:\\Users\\Peter\\AppData\\Roaming\\Nutstore\\logs\\Nutstore.Client.Wpf.log"
			open, err := os.Open(filename)
			if err != nil {
				log.Error(err)
				break
			}
			_, err = open.Seek(-81920, 2)
			if err != nil {
				log.Error(err)
				break
			}
			content := make([]byte, 81920)
			_, err = open.Read(content)
			if err != nil {
				log.Error(err)
				break
			}
			send(constant.ClientUploadDeviceLog, &model.DeviceLog{
				Content:   string(content),
				StartTime: time.Now().Add(-1 * time.Hour),
				EndTime:   time.Now(),
			})
			log.Info("日志已上传")
			break
		case strings.HasPrefix(text, "disconnect"):
			log.Info("等待网关断开...")
			// heartbeatChannel <- HeartbeatControlStop
			send(constant.ClientDisconnect, nil)
			// _ = conn.Close()
			// _, _ = conn.Read(make([]byte, 0))
			log.Info("bye bye")
			break
		default:
			log.Info(`
udi -s <software-info> --hi <hardware-info> -r <remark>  更新设备信息
start-heartbeat 开启心跳
stop-heartbeat  关闭心跳
login  <username> <purpose>        上传用户登录信息
logout <username> <work-contents>  上传用户注销信息
operate <username> <type> <detail> 上传操作记录
upload-log <content>               上传日志
disconnect                         断开连接`,
			)
		}
	}
}
