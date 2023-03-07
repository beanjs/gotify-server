package apns

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"net/http"
	"runtime"
	"strings"
	"time"

	"github.com/sideshow/apns2"
	"github.com/sideshow/apns2/payload"
	"github.com/sideshow/apns2/token"
	"golang.org/x/net/http2"
)

type PushMessage struct {
	DeviceToken string
	Title       string
	Body        string
	Category    string
	// ios notification sound(system sound please refer to http://iphonedevwiki.net/index.php/AudioServices)
	ExtParams map[string]interface{}
}

const (
	topic          = "me.fin.bark"
	keyID          = "LH4T9V5U4R"
	teamID         = "5U8LBRXG3A"
	PayloadMaximum = 4096
)

var cli *apns2.Client

func init() {
	authKey, err := token.AuthKeyFromBytes([]byte(apnsPrivateKey))
	if err != nil {
		log.Fatalf("failed to create APNS auth key: %v", err)
	}

	var rootCAs *x509.CertPool
	if runtime.GOOS == "windows" {
		rootCAs = x509.NewCertPool()
	} else {
		rootCAs, err = x509.SystemCertPool()
		if err != nil {
			log.Fatalf("failed to get rootCAs: %v", err)
		}
	}

	for _, ca := range apnsCAs {
		rootCAs.AppendCertsFromPEM([]byte(ca))
	}

	cli = &apns2.Client{
		Token: &token.Token{
			AuthKey: authKey,
			KeyID:   keyID,
			TeamID:  teamID,
		},
		HTTPClient: &http.Client{
			Transport: &http2.Transport{
				DialTLS: apns2.DialTLS,
				TLSClientConfig: &tls.Config{
					RootCAs: rootCAs,
				},
			},
			Timeout: apns2.HTTPClientTimeout,
		},
		Host: apns2.HostProduction,
	}
	log.Println("init apns client success...")
}

func Push(msg *PushMessage) error {
	pl := payload.NewPayload().
		AlertTitle(msg.Title).
		AlertBody(msg.Body).
		Category(msg.Category)

	group, exist := msg.ExtParams["group"]
	if exist {
		pl = pl.ThreadID(group.(string))
	}

	for k, v := range msg.ExtParams {
		// Change all parameter names to lowercase to prevent inconsistent capitalization
		pl.Custom(strings.ToLower(k), fmt.Sprintf("%v", v))
	}

	// JSON payload maximum size of 4 KB (4096 bytes)
	// https://developer.apple.com/documentation/usernotifications/setting_up_a_remote_notification_server/sending_notification_requests_to_apns#2947607
	plContentForJson, _ := pl.MarshalJSON()
	if len(plContentForJson) > PayloadMaximum {
		return fmt.Errorf("APNS Push Msg Payload too Large %d > 4096 bytes", len(plContentForJson))
	}

	resp, err := cli.Push(&apns2.Notification{
		DeviceToken: msg.DeviceToken,
		Topic:       topic,
		Payload:     pl.MutableContent(),
		Expiration:  time.Now().Add(24 * time.Hour),
	})
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("APNS push failed: %s", resp.Reason)
	}
	return nil
}
