package mqtt_driver

import (
	"ByzoroAC/aclog"
	"ByzoroAC/conf"
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"io/ioutil"
	"sync"
	"time"
)

type mqttClient struct {
	addr string
	ssl  bool
	client  mqtt.Client
}

var (
	// ConnectRetries says how many times the client should retry a failed connection
   ConnectRetries = 10
   // ConnectRetryDelay says how long the client should wait between retries
   ConnectRetryDelay = time.Second
   )

var MqttHandle * mqttClient


func newTLSConfig() *tls.Config{
	certpool := x509.NewCertPool()
	pemCerts,err := ioutil.ReadFile("crt/ca.crt")
	if err == nil{
		certpool.AppendCertsFromPEM(pemCerts)
	}else{
		fmt.Println("read ca file error")
	}
	cert, err := tls.LoadX509KeyPair("crt/client.crt","crt/client.key")
	if err != nil{
		fmt.Println("load client cert fail")
		panic(err)
	}
	cert.Leaf,err = x509.ParseCertificate(cert.Certificate[0])
	return &tls.Config{
		RootCAs: certpool,
		ClientAuth: tls.NoClientCert,
		ClientCAs: nil,
		InsecureSkipVerify: true,
		Certificates: []tls.Certificate{cert},
	}
}

func (handle *mqttClient)newClient() {
	opts := mqtt.NewClientOptions().AddBroker(handle.addr)
	opts.SetClientID("AC-MQTT-CLIENT")
	opts.SetCleanSession(false)
	opts.SetAutoReconnect(true)
	opts.SetMaxReconnectInterval(300*time.Second)
	if handle.ssl{
		tlsconf := newTLSConfig()
		opts.SetTLSConfig(tlsconf)
	}
	handle.client = mqtt.NewClient(opts)
}

func (handle *mqttClient)connect() error{
	if handle.client.IsConnected(){
		return nil
	}
	var err error
	for retries := 0; retries < ConnectRetries; retries++{
		token := handle.client.Connect()
		token.Wait()
		err = token.Error()
		if err == nil{
			break
		}
		<- time.After(ConnectRetryDelay)
	}
	if err != nil{
		return  err
	}
	return  nil
}

func(handle *mqttClient)SendMsgToAP(mac []string,data string){
	var wg sync.WaitGroup
	if handle.client.IsConnected() == false{
		aclog.Warning("Send message to ap process disconnect to server will reconnect to server")
		InitMqtt()
	}
	for _,v := range mac{
		var tmp bytes.Buffer
		tmp.WriteString("v1/ap/cfg/")
		tmp.WriteString(v)
		topic := tmp.String()
		wg.Add(1)
		go func(topic,data string) {
			token := handle.client.Publish(topic,0,false,data)
			if token.Wait() && token.Error() != nil{
				aclog.Error("topic:[%s] data:[%s] error: [%s] fail",
					topic,data,token.Error())
			}
			wg.Done()
		}(topic,data)
	}
	wg.Wait()
}

func(handle *mqttClient)SubscribeTopic(topic string, qos byte, callback mqtt.MessageHandler )error{
	if handle.client.IsConnected() == false{
		aclog.Warning("Subscribe process disconnect to server will reconnect to server")
		InitMqtt()
	}
	token := handle.client.Subscribe(topic, qos, callback)
	if token.Wait() && token.Error() != nil{
		return token.Error()
	}
	return  nil
}

func InitMqtt(){
	var addr bytes.Buffer
	MqttHandle = new(mqttClient)
	if conf.GlobalConf["SSL"] != "" &&  conf.GlobalConf["SSL"] =="enable"{
		addr.WriteString("ssl://")
		MqttHandle.ssl = true
	}else{
		addr.WriteString("tcp://")
	}
	if conf.GlobalConf["MqttHost"] != ""{
		addr.WriteString(conf.GlobalConf["MqttHost"])
	}else {
		addr.WriteString("127.0.0.1")
	}
	addr.WriteString(":")
	if conf.GlobalConf["MqttPort"] != ""{
		addr.WriteString(conf.GlobalConf["MqttPort"])
	}else{
		if MqttHandle.ssl{
			addr.WriteString("8883")
		}else{
			addr.WriteString("6379")
		}
	}
	MqttHandle.addr = addr.String()
	MqttHandle.newClient()
	err := MqttHandle.connect()
	if err != nil{
		aclog.Error("Mqtt server init error can not connect server")
	}
}

