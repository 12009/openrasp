//Copyright 2017-2019 Baidu Inc.
//
//Licensed under the Apache License, Version 2.0 (the "License");
//you may not use this file except in compliance with the License.
//You may obtain a copy of the License at
//
//http: //www.apache.org/licenses/LICENSE-2.0
//
//Unless required by applicable law or agreed to in writing, software
//distributed under the License is distributed on an "AS IS" BASIS,
//WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//See the License for the specific language governing permissions and
//limitations under the License.

package models

import (
	"rasp-cloud/mongo"
	"fmt"
	"strconv"
	"time"
	"math/rand"
	"rasp-cloud/tools"
	"gopkg.in/mgo.v2"
	"crypto/sha1"
	"gopkg.in/mgo.v2/bson"
	"rasp-cloud/models/logs"
	"github.com/astaxie/beego"
	"net/smtp"
	"os"
	"net/mail"
	"strings"
	"html/template"
	"bytes"
	"github.com/astaxie/beego/httplib"
	"errors"
	"crypto/sha256"
	"encoding/base64"
	"io/ioutil"
	"rasp-cloud/es"
	"rasp-cloud/environment"
	"crypto/tls"
	"net"
)

type App struct {
	Id               string                 `json:"id" bson:"_id"`
	Name             string                 `json:"name"  bson:"name"`
	Secret           string                 `json:"secret"  bson:"secret"`
	Language         string                 `json:"language"  bson:"language"`
	Description      string                 `json:"description"  bson:"description"`
	CreateTime       int64                  `json:"create_time"  bson:"create_time"`
	ConfigTime       int64                  `json:"config_time"  bson:"config_time"`
	GeneralConfig    map[string]interface{} `json:"general_config"  bson:"general_config"`
	WhitelistConfig  []WhitelistConfigItem  `json:"whitelist_config"  bson:"whitelist_config"`
	SelectedPluginId string                 `json:"selected_plugin_id" bson:"selected_plugin_id"`
	EmailAlarmConf   EmailAlarmConf         `json:"email_alarm_conf" bson:"email_alarm_conf"`
	DingAlarmConf    DingAlarmConf          `json:"ding_alarm_conf" bson:"ding_alarm_conf"`
	HttpAlarmConf    HttpAlarmConf          `json:"http_alarm_conf" bson:"http_alarm_conf"`
}

type WhitelistConfigItem struct {
	Url  string          `json:"url" bson:"url"`
	Hook map[string]bool `json:"hook" bson:"hook"`
}

type EmailAlarmConf struct {
	Enable     bool     `json:"enable" bson:"enable"`
	ServerAddr string   `json:"server_addr" bson:"server_addr"`
	UserName   string   `json:"username" bson:"username"`
	Password   string   `json:"password" bson:"password"`
	Subject    string   `json:"subject" bson:"subject"`
	RecvAddr   []string `json:"recv_addr" bson:"recv_addr"`
	TlsEnable  bool     `json:"tls_enable" bson:"tls_enable"`
}

type DingAlarmConf struct {
	Enable     bool     `json:"enable" bson:"enable"`
	AgentId    string   `json:"agent_id" bson:"agent_id"`
	CorpId     string   `json:"corp_id" bson:"corp_id"`
	CorpSecret string   `json:"corp_secret" bson:"corp_secret"`
	RecvUser   []string `json:"recv_user" bson:"recv_user"`
	RecvParty  []string `json:"recv_party" bson:"recv_party"`
}

type HttpAlarmConf struct {
	Enable   bool     `json:"enable" bson:"enable"`
	RecvAddr []string `json:"recv_addr" bson:"recv_addr"`
}

type emailTemplateParam struct {
	Total        int64
	Alarms       []map[string]interface{}
	DetailedLink string
	AppName      string
}

type dingResponse struct {
	ErrCode     int64  `json:"errcode"`
	ErrMsg      string `json:"errmsg"`
	AccessToken string `json:"access_token"`
}

const (
	appCollectionName = "app"
	defaultAppName    = "PHP 示例应用"
	SecreteMask       = "************"
)

var (
	panelServerURL string
	lastAlarmTime  = time.Now().UnixNano() / 1000000
	TestAlarmData  = []map[string]interface{}{
		{
			"event_time":      time.Now().Format("2006-01-01 15:04:05"),
			"attack_source":   "220.181.57.191",
			"target":          "localhost",
			"attack_type":     "sql",
			"intercept_state": "block",
		},
	}
	DefaultGeneralConfig = map[string]interface{}{
		"clientip.header":    "ClientIP",
		"block.status_code":  302,
		"block.redirect_url": "https://rasp.baidu.com/blocked/?request_id=%request_id%",
		"block.content_xml": "<?xml version=\"1.0\"?><doc><error>true</error><reason>Request blocked by OpenRASP" +
			"</reason><request_id>%request_id%</request_id></doc>",
		"block.content_html": "</script><script>" +
			"location.href=\"https://rasp.baidu.com/blocked2/?request_id=%request_id%\"</script>",
		"block.content_json":        `{"error":true,"reason": "Request blocked by OpenRASP","request_id": "%request_id%"}`,
		"plugin.timeout.millis":     100,
		"body.maxbytes":             4096,
		"plugin.filter":             true,
		"plugin.maxstack":           100,
		"ognl.expression.minlength": 30,
		"log.maxstack":              50,
		"syslog.tag":                "OpenRASP",
		"syslog.url":                "",
		"syslog.facility":           1,
		"syslog.enable":             false,
	}
)

func init() {
	count, err := mongo.Count(appCollectionName)
	if err != nil {
		tools.Panic(tools.ErrCodeMongoInitFailed, "failed to get app collection count", err)
	}
	if count <= 0 {
		index := &mgo.Index{
			Key:        []string{"name"},
			Unique:     true,
			Background: true,
			Name:       "app_name",
		}
		err = mongo.CreateIndex(appCollectionName, index)
		if err != nil {
			tools.Panic(tools.ErrCodeMongoInitFailed, "failed to create index for app collection", err)
		}
	}
	alarmCheckInterval := beego.AppConfig.DefaultInt64("AlarmCheckInterval", 120)
	if alarmCheckInterval <= 0 {
		tools.Panic(tools.ErrCodeMongoInitFailed, "the 'AlarmCheckInterval' config must be greater than 0", nil)
	} else if alarmCheckInterval < 10 {
		beego.Warning("the value of 'AlarmCheckInterval' config is less than 10, it will be set to 10")
		alarmCheckInterval = 10
	}
	if *environment.StartFlag.StartType == environment.StartTypeDefault ||
		*environment.StartFlag.StartType == environment.StartTypeForeground {
		panelServerURL = beego.AppConfig.String("PanelServerURL")
		if panelServerURL == "" {
			tools.Panic(tools.ErrCodeConfigInitFailed,
				"the 'PanelServerURL' config in the app.conf can not be empty", nil)
		}
		if count <= 0 {
			createDefaultApp()
		}
		go startAlarmTicker(time.Second * time.Duration(alarmCheckInterval))
	}
}

func createDefaultApp() {
	_, err := AddApp(&App{
		Name:        defaultAppName,
		Description: "default app",
		Language:    "php",
	})
	if err != nil {
		tools.Panic(tools.ErrCodeInitDefaultAppFailed, "failed to create default app", err)
	}
}

func startAlarmTicker(interval time.Duration) {
	ticker := time.NewTicker(interval)
	for {
		select {
		case <-ticker.C:
			handleAttackAlarm()
			handleRaspExpiredAlarm()
		}
	}
}

func handleAttackAlarm() {
	defer func() {
		if r := recover(); r != nil {
			beego.Error("failed to handle alarm: ", r)
		}
	}()
	var apps []App
	_, err := mongo.FindAllWithSelect(appCollectionName, nil, &apps, bson.M{"plugin": 0}, 0, 0)
	if err != nil {
		beego.Error("failed to get apps for the alarm: " + err.Error())
		return
	}
	now := time.Now().UnixNano() / 1000000
	for _, app := range apps {
		total, result, err := logs.SearchLogs(lastAlarmTime, now, nil, "event_time",
			1, 10, false, logs.AliasAttackIndexName+"-"+app.Id)
		if err != nil {
			beego.Error("failed to get alarm from es: " + err.Error())
			continue
		}
		if total > 0 {
			PushAttackAlarm(&app, total, result, false)
		}
	}
	lastAlarmTime = now + 1
}

func handleRaspExpiredAlarm() {
	//defer func() {
	//	if r := recover(); r != nil {
	//		beego.Error("failed to handle alarm: ", r)
	//	}
	//}()
	//Rasp
}

func AddApp(app *App) (result *App, err error) {
	app.Id = generateAppId(app)
	app.Secret = generateSecret(app)
	app.CreateTime = time.Now().Unix()
	if mongo.FindOne(appCollectionName, bson.M{"name": app.Name}, &App{}) != mgo.ErrNotFound {
		return nil, errors.New("duplicate app name")
	}
	HandleApp(app, true)
	err = es.CreateEsIndex(logs.PolicyIndexName+"-"+app.Id, logs.AliasPolicyIndexName+"-"+app.Id, logs.PolicyEsMapping)
	if err != nil {
		return
	}
	err = es.CreateEsIndex(logs.AttackIndexName+"-"+app.Id, logs.AliasAttackIndexName+"-"+app.Id, logs.AttackEsMapping)
	if err != nil {
		return
	}
	err = es.CreateEsIndex(ReportIndexName+"-"+app.Id, AliasReportIndexName+"-"+app.Id, ReportEsMapping)
	if err != nil {
		return
	}
	// ES must be created before mongo
	err = mongo.Insert(appCollectionName, app)
	if err != nil {
		return nil, errors.New("failed to insert app to db: " + err.Error())
	}
	result = app
	beego.Info("Succeed to create app, name: " + app.Name)
	selectDefaultPlugin(app)
	return
}

func selectDefaultPlugin(app *App) {
	// if setting default plugin fails, continue to initialize
	currentPath, err := tools.GetCurrentPath()
	if err != nil {
		beego.Warn("failed to create default plugin", err)
		return
	}
	content, err := ioutil.ReadFile(currentPath + "/resources/plugin.js")
	if err != nil {
		beego.Warn(tools.ErrCodeInitDefaultAppFailed, "failed to get default plugin: "+err.Error())
		return
	}
	plugin, err := AddPlugin(content, app.Id)
	if err != nil {
		beego.Warn(tools.ErrCodeInitDefaultAppFailed, "failed to insert default plugin: "+err.Error())
		return
	}
	err = SetSelectedPlugin(app.Id, plugin.Id)
	if err != nil {
		beego.Warn(tools.ErrCodeInitDefaultAppFailed, "failed to select default plugin for app: " + err.Error()+
			", app_id: "+ app.Id+ ", plugin_id: "+ plugin.Id)
		return
	}
	beego.Info("Succeed to set up default plugin for app, version: " + plugin.Version)
}

func generateAppId(app *App) string {
	random := "openrasp_app" + app.Name + strconv.FormatInt(time.Now().UnixNano(), 10) + strconv.Itoa(rand.Intn(10000))
	return fmt.Sprintf("%x", sha1.Sum([]byte(random)))
}

func generateSecret(app *App) string {
	random := "openrasp_app" + app.Name + app.Id +
		strconv.FormatInt(time.Now().UnixNano(), 10) + strconv.Itoa(rand.Intn(10000))
	sha256Data := sha256.Sum256([]byte(random))
	base64Data := base64.NewEncoding("OPQRSTYZabcdefgABCDEFGHIJKLMNhijklmnopqrUVWXstuvwxyz01234567891q").
		EncodeToString(sha256Data[0:])
	return base64Data[0 : len(base64Data)-1]
}

func GetAllApp(page int, perpage int, mask bool) (count int, result []*App, err error) {
	count, err = mongo.FindAll(appCollectionName, nil, &result, perpage*(page-1), perpage, "name")
	if err == nil && result != nil {
		for _, app := range result {
			if mask {
				HandleApp(app, false)
			}
		}
	}
	return
}

func GetAppByIdWithoutMask(id string) (app *App, err error) {
	err = mongo.FindId(appCollectionName, id, &app)
	return
}

func GetAppById(id string) (app *App, err error) {
	err = mongo.FindId(appCollectionName, id, &app)
	if err == nil && app != nil {
		HandleApp(app, false)
	}
	return
}

func GetSecretByAppId(appId string) (secret string, err error) {
	newSession := mongo.NewSession()
	defer newSession.Close()
	var result *App
	err = newSession.DB(mongo.DbName).C(appCollectionName).FindId(appId).Select(bson.M{"secret": 1}).One(&result)
	if err != nil {
		return
	}
	if result != nil {
		secret = result.Secret
	}
	return
}

func RegenerateSecret(appId string) (secret string, err error) {
	var app *App
	err = mongo.FindId(appCollectionName, appId, &app)
	if err == nil {
		return
	}
	newSecret := generateSecret(app)
	err = mongo.UpdateId(appCollectionName, appId, bson.M{"secret": newSecret})
	return
}

func HandleApp(app *App, isCreate bool) {
	if app.EmailAlarmConf.RecvAddr == nil {
		app.EmailAlarmConf.RecvAddr = make([]string, 0)
	}
	if app.DingAlarmConf.RecvParty == nil {
		app.DingAlarmConf.RecvParty = make([]string, 0)
	}
	if app.DingAlarmConf.RecvUser == nil {
		app.DingAlarmConf.RecvUser = make([]string, 0)
	}
	if app.HttpAlarmConf.RecvAddr == nil {
		app.HttpAlarmConf.RecvAddr = make([]string, 0)
	}
	if !isCreate {
		if app.EmailAlarmConf.Password != "" {
			app.EmailAlarmConf.Password = SecreteMask
		}
		if app.DingAlarmConf.CorpSecret != "" {
			app.DingAlarmConf.CorpSecret = SecreteMask
		}
	} else {
		if app.GeneralConfig == nil {
			app.GeneralConfig = DefaultGeneralConfig
		}
	}
	if app.WhitelistConfig == nil {
		app.WhitelistConfig = make([]WhitelistConfigItem, 0)
	}
	if app.GeneralConfig == nil {
		app.GeneralConfig = make(map[string]interface{})
	}
}

func UpdateAppById(id string, doc interface{}) (app *App, err error) {
	err = mongo.UpdateId(appCollectionName, id, doc)
	if err != nil {
		return
	}
	return GetAppById(id)
}

func UpdateGeneralConfig(appId string, config map[string]interface{}) (*App, error) {
	return UpdateAppById(appId, bson.M{"general_config": config, "config_time": time.Now().UnixNano()})
}

func UpdateWhiteListConfig(appId string, config []WhitelistConfigItem) (app *App, err error) {
	return UpdateAppById(appId, bson.M{"whitelist_config": config, "config_time": time.Now().UnixNano()})
}

func RemoveAppById(id string) (app *App, err error) {
	err = mongo.FindId(appCollectionName, id, &app)
	if err != nil {
		return
	}
	return app, mongo.RemoveId(appCollectionName, id)
}

func GetAppCount() (count int, err error) {
	return mongo.Count(appCollectionName)
}

func PushAttackAlarm(app *App, total int64, alarms []map[string]interface{}, isTest bool) {
	if app != nil {
		if app.DingAlarmConf.Enable {
			PushDingAttackAlarm(app, total, alarms, isTest)
		}
		if app.EmailAlarmConf.Enable {
			PushEmailAttackAlarm(app, total, alarms, isTest)
		}
		if app.HttpAlarmConf.Enable {
			PushHttpAttackAlarm(app, total, alarms, isTest)
		}
	}
}

func PushEmailAttackAlarm(app *App, total int64, alarms []map[string]interface{}, isTest bool) error {
	var emailConf = app.EmailAlarmConf
	if len(emailConf.RecvAddr) > 0 && emailConf.ServerAddr != "" {
		var (
			subject   string
			msg       string
			emailAddr = &mail.Address{Address: emailConf.UserName}
		)
		hostName, err := os.Hostname()
		if err == nil {
			emailAddr.Name = hostName
		} else {
			emailAddr.Name = "OpenRASP"
		}
		if emailConf.Subject == "" {
			subject = "OpenRASP alarm"
		} else {
			subject = emailConf.Subject
		}
		if isTest {
			subject = "【测试邮件】" + subject
			alarms = TestAlarmData
			total = int64(len(TestAlarmData))
		}
		head := map[string]string{
			"from":         emailAddr.String(),
			"To":           strings.Join(emailConf.RecvAddr, ","),
			"Content-Type": "text/html; charset=UTF-8",
			"Subject":      subject,
		}
		t, err := template.ParseFiles("views/email.tpl")
		if err != nil {
			beego.Error("failed to render email template: " + err.Error())
			return err
		}
		alarmData := new(bytes.Buffer)
		err = t.Execute(alarmData, &emailTemplateParam{
			Total:        total - int64(len(alarms)),
			Alarms:       alarms,
			AppName:      app.Name,
			DetailedLink: panelServerURL + "/#/events/" + app.Id,
		})
		if err != nil {
			beego.Error("failed to execute email template: " + err.Error())
			return err
		}
		for k, v := range head {
			msg += fmt.Sprintf("%s: %s\r\n", k, v)
		}
		msg += "\r\n" + alarmData.String()
		host, _, err := net.SplitHostPort(emailConf.ServerAddr)
		if err != nil {
			errMsg := "failed to get email serve host: " + err.Error()
			beego.Error(errMsg)
			return errors.New(errMsg)
		}
		auth := smtp.PlainAuth("", emailConf.UserName, emailConf.Password, host)
		if emailConf.Password == "" {
			auth = nil
		}

		if emailConf.TlsEnable {
			return sendEmailWithTls(emailConf, auth, msg)
		} else {
			return sendNormalEmail(emailConf, auth, msg)
		}
	} else {
		beego.Error(
			"failed to send email alarm: the email receiving address and email server address can not be empty", emailConf)
		return errors.New("the email receiving address and email server address can not be empty")
	}
	beego.Debug("succeed in pushing email alarm for app: " + app.Name)
	return nil
}

func sendNormalEmail(emailConf EmailAlarmConf, auth smtp.Auth, msg string) (err error) {
	err = smtp.SendMail(emailConf.ServerAddr, auth, emailConf.UserName, emailConf.RecvAddr, []byte(msg))
	if err != nil {
		beego.Error("failed to push email alarms: " + err.Error())
		return
	}
	return
}

func sendEmailWithTls(emailConf EmailAlarmConf, auth smtp.Auth, msg string) error {
	client, err := smtpTlsDial(emailConf.ServerAddr)
	if err != nil {
		errMsg := "failed to start tls: " + err.Error()
		beego.Error(errMsg)
		return errors.New(errMsg)
	}
	defer client.Close()
	if auth != nil {
		if ok, _ := client.Extension("AUTH"); ok {
			if err = client.Auth(auth); err != nil {
				errMsg := "failed to auth with tls: " + err.Error()
				beego.Error(errMsg)
				return errors.New(errMsg)
			}
		}
	}
	if err = client.Mail(emailConf.UserName); err != nil {
		errMsg := "failed to mail from 'emailConf.UserName': " + err.Error()
		beego.Error(errMsg)
		return errors.New(errMsg)
	}

	for _, addr := range emailConf.RecvAddr {
		if err = client.Rcpt(addr); err != nil {
			errMsg := "failed to push email to " + addr + " with tls: " + err.Error()
			beego.Error(errMsg)
			return errors.New(errMsg)
		}
	}

	writer, err := client.Data()
	if err != nil {
		errMsg := "failed to get writer for email with tls: " + err.Error()
		beego.Error(errMsg)
		return errors.New(errMsg)
	}
	defer writer.Close()

	_, err = writer.Write([]byte(msg))
	if err != nil {
		errMsg := "failed to write msg with tls: " + err.Error()
		beego.Error(errMsg)
		return errors.New(errMsg)
	}

	client.Quit()
	return nil
}

func smtpTlsDial(addr string) (*smtp.Client, error) {
	host, _, err := net.SplitHostPort(addr)
	if err != nil {
		errMsg := "failed to get email serve host: " + err.Error()
		beego.Error(errMsg)
		return nil, errors.New(errMsg)
	}
	conn, err := tls.Dial("tcp", addr, nil)
	if err != nil {
		errMsg := "smtp dialing error: " + err.Error()
		beego.Error(errMsg)
		return nil, errors.New(errMsg)
	}
	return smtp.NewClient(conn, host)
}

func PushHttpAttackAlarm(app *App, total int64, alarms []map[string]interface{}, isTest bool) error {
	var httpConf = app.HttpAlarmConf
	if len(httpConf.RecvAddr) != 0 {
		body := make(map[string]interface{})
		body["app_id"] = app.Id
		if isTest {
			body["data"] = TestAlarmData
		} else {
			body["data"] = alarms
		}
		for _, addr := range httpConf.RecvAddr {
			request := httplib.Post(addr)
			request.JSONBody(body)
			request.SetTimeout(10*time.Second, 10*time.Second)
			response, err := request.Response()
			if err != nil {
				beego.Error("failed to push http alarms to: " + addr + ", with error: " + err.Error())
				return err
			}
			if response.StatusCode > 299 || response.StatusCode < 200 {
				err := errors.New("failed to push http alarms to: " + addr + ", with status code: " +
					strconv.Itoa(response.StatusCode))
				beego.Error(err.Error())
				return err
			}
		}
	} else {
		beego.Error("failed to send http alarm: the http receiving address can not be empty", httpConf)
		return errors.New("the http receiving address can not be empty")
	}
	beego.Debug("succeed in pushing http alarm for app: " + app.Name + " ,with urls: " +
		fmt.Sprintf("%v", httpConf.RecvAddr))
	return nil
}

func PushDingAttackAlarm(app *App, total int64, alarms []map[string]interface{}, isTest bool) error {
	var dingCong = app.DingAlarmConf
	if dingCong.CorpId != "" && dingCong.CorpSecret != "" && dingCong.AgentId != "" &&
		!(len(dingCong.RecvParty) == 0 && len(dingCong.RecvUser) == 0) {

		request := httplib.Get("https://oapi.dingtalk.com/gettoken")
		request.SetTimeout(10*time.Second, 10*time.Second)
		request.Param("corpid", dingCong.CorpId)
		request.Param("corpsecret", dingCong.CorpSecret)
		response, err := request.Response()
		errMsg := "failed to get ding ding token with corp id: " + dingCong.CorpId
		if err != nil {
			beego.Error(errMsg + ", with error: " + err.Error())
			return err
		}
		if response.StatusCode != 200 {
			err := errors.New(errMsg + ", with status code: " + strconv.Itoa(response.StatusCode))
			beego.Error(err.Error())
			return err
		}
		var result dingResponse
		err = request.ToJSON(&result)
		if err != nil {
			beego.Error(errMsg + ", with error: " + err.Error())
			return err
		}
		if result.ErrCode != 0 {
			err := errors.New(errMsg + ", with errmsg: " + result.ErrMsg)
			beego.Error(err.Error())
			return err
		}
		token := result.AccessToken
		body := make(map[string]interface{})
		dingText := ""
		if isTest {
			dingText = "OpenRASP test message from app: " + app.Name + ", time: " + time.Now().Format(time.RFC3339)
		} else {
			dingText = "时间：" + time.Now().Format(time.RFC3339) + "， 来自 OpenRAS 的报警\n共有 " +
				strconv.FormatInt(total, 10) + " 条报警信息来自 APP：" + app.Name + "，详细信息：" + panelServerURL + "/#/events/" + app.Id
		}
		if len(dingCong.RecvUser) > 0 {
			body["touser"] = strings.Join(dingCong.RecvUser, "|")
		}
		if len(dingCong.RecvParty) > 0 {
			body["toparty"] = strings.Join(dingCong.RecvParty, "|")
		}
		body["agentid"] = dingCong.AgentId
		body["msgtype"] = "text"
		body["text"] = map[string]string{"content": dingText}
		request = httplib.Post("https://oapi.dingtalk.com/message/send?access_token=" + token)
		request.JSONBody(body)
		request.SetTimeout(10*time.Second, 10*time.Second)
		response, err = request.Response()
		errMsg = "failed to push ding ding alarms with corp id: " + dingCong.CorpId
		if err != nil {
			beego.Error(errMsg + ", with error: " + err.Error())
			return err
		}
		if response.StatusCode != 200 {
			err := errors.New(errMsg + ", with status code: " + strconv.Itoa(response.StatusCode))
			beego.Error(err.Error())
			return err
		}
		err = request.ToJSON(&result)
		if err != nil {
			beego.Error(errMsg + ", with error: " + err.Error())
			return err
		}
		if result.ErrCode != 0 {
			err := errors.New(errMsg + ", with errmsg: " + result.ErrMsg)
			beego.Error(err.Error())
			return err
		}
	} else {
		beego.Error("failed to send ding ding alarm: invalid ding ding alarm conf", dingCong)
		return errors.New("invalid ding ding alarm conf")
	}
	beego.Debug("succeed in pushing ding ding alarm for app: " + app.Name + " ,with corp id: " + dingCong.CorpId)
	return nil
}
