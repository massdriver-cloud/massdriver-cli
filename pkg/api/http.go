package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"sync/atomic"
	"time"

	"github.com/hasura/go-graphql-client"
	"github.com/rs/zerolog/log"
	"nhooyr.io/websocket"
)

// TODO remove this once SA token working
// for now this can be populated with gql`query me { me { token } }`
var removeMeToken = "SFMyNTY.g2gDdAAAAA5kAAhfX21ldGFfX3QAAAAGZAAKX19zdHJ1Y3RfX2QAG0VsaXhpci5FY3RvLlNjaGVtYS5NZXRhZGF0YWQAB2NvbnRleHRkAANuaWxkAAZwcmVmaXhkAANuaWxkAAZzY2hlbWFkACNFbGl4aXIuTWFzc2RyaXZlci5NZWF0c3BhY2UuQWNjb3VudGQABnNvdXJjZW0AAAAIYWNjb3VudHNkAAVzdGF0ZWQABWJ1aWx0ZAAKX19zdHJ1Y3RfX2QAI0VsaXhpci5NYXNzZHJpdmVyLk1lYXRzcGFjZS5BY2NvdW50ZAALYXR0cmlidXRpb25kAANuaWxkABNiZXRhX2FjY2Vzc19lbmFibGVkZAADbmlsZAAKY3JlYXRlZF9hdGQAA25pbGQABWVtYWlsZAADbmlsZAAKZmlyc3RfbmFtZWQAA25pbGQAEWdyb3VwX21lbWJlcnNoaXBzdAAAAARkAA9fX2NhcmRpbmFsaXR5X19kAARtYW55ZAAJX19maWVsZF9fZAARZ3JvdXBfbWVtYmVyc2hpcHNkAAlfX293bmVyX19kACNFbGl4aXIuTWFzc2RyaXZlci5NZWF0c3BhY2UuQWNjb3VudGQACl9fc3RydWN0X19kACFFbGl4aXIuRWN0by5Bc3NvY2lhdGlvbi5Ob3RMb2FkZWRkAAZncm91cHN0AAAABGQAD19fY2FyZGluYWxpdHlfX2QABG1hbnlkAAlfX2ZpZWxkX19kAAZncm91cHNkAAlfX293bmVyX19kACNFbGl4aXIuTWFzc2RyaXZlci5NZWF0c3BhY2UuQWNjb3VudGQACl9fc3RydWN0X19kACFFbGl4aXIuRWN0by5Bc3NvY2lhdGlvbi5Ob3RMb2FkZWRkAAJpZG0AAAAkNmI0ZjRjMjktMDFkMi00YjYyLWIzNzEtMjg3NTM2Y2Y2MDhlZAATaWRlbnRpdHlfc2VydmljZV9pZGQAA25pbGQACWxhc3RfbmFtZWQAA25pbGQADW9yZ2FuaXphdGlvbnN0AAAABGQAD19fY2FyZGluYWxpdHlfX2QABG1hbnlkAAlfX2ZpZWxkX19kAA1vcmdhbml6YXRpb25zZAAJX19vd25lcl9fZAAjRWxpeGlyLk1hc3Nkcml2ZXIuTWVhdHNwYWNlLkFjY291bnRkAApfX3N0cnVjdF9fZAAhRWxpeGlyLkVjdG8uQXNzb2NpYXRpb24uTm90TG9hZGVkZAAKdXBkYXRlZF9hdGQAA25pbG4GAAShM0eDAWIAAVGA.qdmtujqw9gw9wyM7TPMG27GMI8WU5nSIlbq6OC8pX_c" //nolint: gosec
func NewClient() *graphql.Client {
	c := http.Client{Transport: &transport{underlyingTransport: http.DefaultTransport}}
	client := graphql.NewClient("https://api.massdriver.cloud/api/", &c)
	return client
}

type PhoenixWebsocket struct {
	*websocket.Conn
	// holds the join ref for the channel
	JoinRefByTopic map[string]*int64
}

type PhoenixMessage struct {
	JoinRef *int64      `json:"join_ref,omitempty"`
	Ref     *int64      `json:"ref,omitempty"`
	Topic   string      `json:"topic,omitempty"`
	Event   string      `json:"event,omitempty"`
	Payload interface{} `json:"payload"`
}

func (p PhoenixMessage) String() string {
	return fmt.Sprintf("{joinRef: %d, ref: %d, topic: \"%v\", event: \"%v\", payload: %+v}", p.JoinRef, p.Ref, p.Topic, p.Event, p.Payload)
}

// replace the value if i, with that of v if types are compatible
func replace(i, v interface{}) error {
	val := reflect.ValueOf(i)
	if val.Kind() != reflect.Ptr {
		return fmt.Errorf("expected pointer, got %T", i)
	}

	val = val.Elem()

	newVal := reflect.Indirect(reflect.ValueOf(v))

	if !val.Type().AssignableTo(newVal.Type()) {
		return fmt.Errorf("cannot assign %T to %T", v, i)
	}

	val.Set(newVal)
	return nil
}

func (p PhoenixWebsocket) ReadJSON(v interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	_, b, err := p.Conn.Read(ctx)
	if err != nil {
		return err
	}
	log.Debug().Msgf("received raw message from phoenix socket: %s", string(b))

	phxMsg, alreadyPhoenixMsg := v.(PhoenixMessage)
	if alreadyPhoenixMsg {
		return json.Unmarshal(b, v)
	}

	// if we reached here, unwrap the phoenix implementation details just grabbing the payload the user cares about
	if marshalErr := json.Unmarshal(b, &phxMsg); marshalErr != nil {
		p.Conn.Close(websocket.StatusInvalidFramePayloadData, "failed to unmarshal JSON")
		return fmt.Errorf("failed to unmarshal JSON: %w", marshalErr)
	}
	return replace(v, phxMsg.Payload)
}

var msgRef = new(int64)
var joinRef = new(int64)

func nextMsgRef() int64 {
	return atomic.AddInt64(msgRef, 1)
}

func nextJoinRef() int64 {
	return atomic.AddInt64(joinRef, 1)
}

func (p PhoenixWebsocket) WriteJSON(v interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	ref := nextMsgRef()
	phxMsg, alreadyPhxMsg := v.(PhoenixMessage)
	if !alreadyPhxMsg {
		phxMsg = PhoenixMessage{
			JoinRef: joinRef,
			Ref:     &ref,
			Topic:   "__absinthe__:control",
			Event:   "doc",
			Payload: v,
		}
	}
	// hack to transform the connection_init into a phx_join event,
	// wrap the value in the phoenix implementation details
	msg, isGQLMsg := v.(graphql.OperationMessage)
	if isGQLMsg {
		if msg.Type == "connection_init" {
			jr := nextJoinRef()
			phxMsg.JoinRef = &jr
			phxMsg.Event = "phx_join"
			phxMsg.Payload = map[string]interface{}{}
		}
		if msg.Type == "start" {
			var p map[string]interface{}
			if err := json.Unmarshal(msg.Payload, &p); err != nil {
				return err
			}
			phxMsg.Payload = p
		}
		if msg.Type == "stop" {
			phxMsg.Event = "phx_leave"
			phxMsg.Payload = map[string]interface{}{}
		}
		if msg.Type == "connection_terminate" {
			phxMsg.Event = "phx_leave"
			phxMsg.Payload = map[string]interface{}{}
		}
	}

	w, err := p.Conn.Writer(ctx, websocket.MessageText)
	if err != nil {
		return err
	}
	defer w.Close()

	// json.Marshal cannot reuse buffers between calls as it has to return
	// a copy of the byte slice but Encoder does as it directly writes to w.
	b, err := json.Marshal(phxMsg)
	if err != nil {
		return err
	}
	log.Debug().Msgf("writing message to phoenix socket: %v", string(b))
	return json.NewEncoder(w).Encode(phxMsg)
}

func (p PhoenixWebsocket) Close() error {
	return p.Conn.Close(websocket.StatusAbnormalClosure, "client closed")
}

func (p PhoenixWebsocket) SetReadLimit(limit int64) {
	p.Conn.SetReadLimit(limit)
}

func authURL(u *url.URL) *url.URL {
	q := u.Query()
	// TODO better pluggable auth here.
	q.Add("token", removeMeToken)
	u.RawQuery = q.Encode()
	return u
}

func newPhoenixWebsocketConn(sc *graphql.SubscriptionClient) (graphql.WebsocketConn, error) {
	options := &websocket.DialOptions{
		Subprotocols: []string{"graphql-ws"},
	}
	url, err := url.Parse(sc.GetURL())
	if err != nil {
		return nil, err
	}
	authURL(url)

	c, resp, err := websocket.Dial(sc.GetContext(), url.String(), options)
	if resp.Body != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return nil, err
	}
	phxSocket := PhoenixWebsocket{Conn: c}

	// lazy heartbeater probably a better way to do this
	go func(phxSocket PhoenixWebsocket) {
		for {
			ref := nextMsgRef()
			hbErr := phxSocket.WriteJSON(PhoenixMessage{
				JoinRef: nil,
				Ref:     &ref,
				Topic:   "phoenix",
				Event:   "heartbeat",
				Payload: map[string]interface{}{},
			})
			if hbErr != nil {
				log.Err(hbErr).Msg("failed to send heartbeat")
			}
			time.Sleep(time.Second)
		}
	}(phxSocket)

	return &phxSocket, nil
}

func NewSubscriptionClient() *graphql.SubscriptionClient {
	c := http.Client{Transport: &wstransport{underlyingTransport: http.DefaultTransport}}
	client := graphql.NewSubscriptionClient("wss://api.massdriver.cloud/socket/websocket?vsn=2.0.0")
	client.WithWebSocketOptions(graphql.WebsocketOptions{HTTPClient: &c})
	client.WithWebSocket(newPhoenixWebsocketConn)
	client.OnConnected(func() {
		log.Debug().Msg("connected to massdriver websocket")
	})
	client.OnDisconnected(func() {
		log.Debug().Msg("disconnected to massdriver websocket")
	})
	// TODO the UI does _not_ do this
	// client.WithConnectionParams(map[string]interface{}{
	// 		"token": removeMeToken,
	// })
	// TODO not sure if this is necessary for GQL_CONNECTION_INIT
	// client.WithConnectionParams(map[string]interface{}{
	// 	"token": removeMeToken,
	// })
	client.WithLog(func(args ...interface{}) {
		log.Debug().Msgf("%#v", args)
	})
	return client
}

type transport struct {
	underlyingTransport http.RoundTripper
}

func (t *transport) RoundTrip(req *http.Request) (*http.Response, error) {
	apiKey := os.Getenv("MASSDRIVER_API_KEY")

	if apiKey == "" {
		log.Fatal().Msg("MASSDRIVER_API_KEY must be set")
	}

	bearer := fmt.Sprintf("Bearer %s", apiKey)
	req.Header.Add("authorization", bearer)
	return t.underlyingTransport.RoundTrip(req)
}

// wstransport is a transport for the graphql subscription client because it needs the token set as a url param
// as websockets do not support headers.
type wstransport struct {
	underlyingTransport http.RoundTripper
}

func (t *wstransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// TODO eventually we will use the api key in the token query
	// apiKey := os.Getenv("MASSDRIVER_API_KEY")

	// if apiKey == "" {
	// 	log.Fatal().Msg("MASSDRIVER_API_KEY must be set")
	// }

	// req.URL.Query().Add("token", apiKey)

	// TODO remove this hard coded token
	q := req.URL.Query()
	q.Add("token", removeMeToken)
	req.URL.RawQuery = q.Encode()
	// TODO these did not seem to help but were missing when compared to what is set when inspecting what the UI is doing in chrome
	req.Header.Add("Cache-Control", "no-cache")
	req.Header.Add("Pragma", "no-cache")
	req.Header.Add("Accept-Encoding", "gzip, deflate, br")
	req.Header.Add("Accept-Language", "en-US,en;q=0.9")
	req.Header.Add("Host", req.Host)
	req.Header.Add("Origin", "https://app.massdriver.cloud")
	log.Debug().Msgf("wstransport req: %#v", req)
	return t.underlyingTransport.RoundTrip(req)
}
