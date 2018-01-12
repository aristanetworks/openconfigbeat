/*
Copyright 2017 Google Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Package client implements a gNMI client.
package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	log "github.com/golang/glog"
	"context"
	"github.com/golang/protobuf/proto"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc"
	"github.com/openconfig/ygot/ygot"
	"github.com/openconfig/gnmi/client"
	"github.com/openconfig/gnmi/client/grpcutil"
	"github.com/openconfig/gnmi/value"

	gpb "github.com/openconfig/gnmi/proto/gnmi"
)

// Type defines the name resolution for this client type.
const Type = "gnmi"

// Client handles execution of the query and caching of its results.
type Client struct {
	conn      *grpc.ClientConn
	client    gpb.GNMIClient
	sub       gpb.GNMI_SubscribeClient
	query     client.Query
	recv      client.ProtoHandler
	handler   client.NotificationHandler
	connected bool
}

// New returns a new initialized client. If error is nil, returned Client has
// established a connection to d. Close needs to be called for cleanup.
func New(ctx context.Context, d client.Destination) (client.Impl, error) {
	if len(d.Addrs) != 1 {
		return nil, fmt.Errorf("d.Addrs must only contain one entry: %v", d.Addrs)
	}
	opts := []grpc.DialOption{
		grpc.WithTimeout(d.Timeout),
		grpc.WithBlock(),
	}
	if d.TLS != nil {
		opts = append(opts, grpc.WithTransportCredentials(credentials.NewTLS(d.TLS)))
	}
	if d.Credentials != nil {
		pc := newPassCred(d.Credentials.Username, d.Credentials.Password, true)
		opts = append(opts, grpc.WithPerRPCCredentials(pc))
	}
	conn, err := grpc.DialContext(ctx, d.Addrs[0], opts...)
	if err != nil {
		return nil, fmt.Errorf("Dialer(%s, %v): %v", d.Addrs[0], d.Timeout, err)
	}
	return NewFromConn(ctx, conn, d)
}

// NewFromConn creates and returns the client based on the provided transport.
func NewFromConn(ctx context.Context, conn *grpc.ClientConn, d client.Destination) (*Client, error) {
	ok, err := grpcutil.Lookup(ctx, conn, "gnmi.gNMI")
	if err != nil {
		log.V(1).Infof("gRPC reflection lookup on %q for service gnmi.gNMI failed: %v", d.Addrs, err)
		// This check is disabled for now. Reflection will become part of gNMI
		// specification in the near future, so we can't enforce it yet.
	}
	if !ok {
		// This check is disabled for now. Reflection will become part of gNMI
		// specification in the near future, so we can't enforce it yet.
	}

	cl := gpb.NewGNMIClient(conn)
	return &Client{
		conn:   conn,
		client: cl,
	}, nil
}

// Subscribe sends the gNMI Subscribe RPC to the server.
func (c *Client) Subscribe(ctx context.Context, q client.Query) error {
	sub, err := c.client.Subscribe(ctx)
	if err != nil {
		return fmt.Errorf("gpb.GNMIClient.Subscribe(%v) failed to initialize Subscribe RPC: %v", q, err)
	}
	qq, err := subscribe(q)
	if err != nil {
		return fmt.Errorf("generating SubscribeRequest proto: %v", err)
	}
	if err := sub.Send(qq); err != nil {
		return fmt.Errorf("client.Send(%+v): %v", qq, err)
	}

	c.sub = sub
	c.query = q
	if q.ProtoHandler == nil {
		c.recv = c.defaultRecv
		c.handler = q.NotificationHandler
	} else {
		c.recv = q.ProtoHandler
	}
	return nil
}

// Poll will send a single gNMI poll request to the server.
func (c *Client) Poll() error {
	if err := c.sub.Send(&gpb.SubscribeRequest{Request: &gpb.SubscribeRequest_Poll{Poll: &gpb.Poll{}}}); err != nil {
		return fmt.Errorf("client.Poll(): %v", err)
	}
	return nil
}

// Peer returns the peer of the current stream. If the client is not created or
// if the peer is not valid nil is returned.
func (c *Client) Peer() string {
	return c.query.Addrs[0]
}

// Close forcefully closes the underlying connection, terminating the query
// right away. It's safe to call Close multiple times.
func (c *Client) Close() error {
	return c.conn.Close()
}

// Recv will recieve a single message from the server and process it based on
// the provided handlers (Proto or Notification).
func (c *Client) Recv() error {
	n, err := c.sub.Recv()
	if err != nil {
		return err
	}
	return c.recv(n)
}

// defaultRecv is the default implementation of recv provided by the client.
// This function will be replaced by the ProtoHandler member of the Query
// struct passed to New(), if it is set.
func (c *Client) defaultRecv(msg proto.Message) error {
	if !c.connected {
		c.handler(client.Connected{})
		c.connected = true
	}

	resp, ok := msg.(*gpb.SubscribeResponse)
	if !ok {
		return fmt.Errorf("failed to type assert message %#v", msg)
	}
	log.V(1).Info(resp)
	switch v := resp.Response.(type) {
	default:
		return fmt.Errorf("unknown response %T: %s", v, v)
	case *gpb.SubscribeResponse_Error:
		return fmt.Errorf("error in response: %s", v)
	case *gpb.SubscribeResponse_SyncResponse:
		c.handler(client.Sync{})
		if c.query.Type == client.Poll || c.query.Type == client.Once {
			return client.ErrStopReading
		}
	case *gpb.SubscribeResponse_Update:
		n := v.Update
		var p []string
		if n.Prefix != nil {
			var err error
			p, err = ygot.PathToStrings(n.Prefix)
			if err != nil {
				return err
			}
		}
		ts := time.Unix(0, n.Timestamp)
		for _, u := range n.Update {
			if u.Path == nil {
				return fmt.Errorf("invalid nil path in update: %v", u)
			}
			u, err := noti(p, u.Path, ts, u)
			if err != nil {
				return err
			}
			c.handler(u)
		}
		for _, d := range n.Delete {
			u, err := noti(p, d, ts, nil)
			if err != nil {
				return err
			}
			c.handler(u)
		}
	}
	return nil
}

// Set calls the Set RPC, converting request/response to appropriate protos.
func (c *Client) Set(ctx context.Context, sr client.SetRequest) (client.SetResponse, error) {
	req, err := convertSetRequest(sr)
	if err != nil {
		return client.SetResponse{}, err
	}

	resp, err := c.client.Set(ctx, req)
	if err != nil {
		return client.SetResponse{}, err
	}

	return convertSetResponse(resp)
}

func convertSetRequest(sr client.SetRequest) (*gpb.SetRequest, error) {
	req := &gpb.SetRequest{}
	for _, d := range sr.Delete {
		pp, err := ygot.StringToPath(pathToString(d), ygot.StructuredPath, ygot.StringSlicePath)
		if err != nil {
			return nil, fmt.Errorf("invalid delete path %q: %v", d, err)
		}
		req.Delete = append(req.Delete, pp)
	}

	genUpdate := func(v client.Leaf) (*gpb.Update, error) {
		buf, err := json.Marshal(v.Val)
		if err != nil {
			return nil, err
		}
		pp, err := ygot.StringToPath(pathToString(v.Path), ygot.StructuredPath, ygot.StringSlicePath)
		if err != nil {
			return nil, fmt.Errorf("invalid update path %q: %v", v.Path, err)
		}
		return &gpb.Update{
			Path: pp,
			Val:  &gpb.TypedValue{Value: &gpb.TypedValue_JsonVal{buf}},
			// Value is deprecated, remove it at some point.
			Value: &gpb.Value{Type: gpb.Encoding_JSON, Value: buf},
		}, nil
	}

	for _, u := range sr.Update {
		uu, err := genUpdate(u)
		if err != nil {
			return nil, err
		}
		req.Update = append(req.Update, uu)
	}
	for _, u := range sr.Replace {
		uu, err := genUpdate(u)
		if err != nil {
			return nil, err
		}
		req.Replace = append(req.Replace, uu)
	}

	return req, nil
}

func convertSetResponse(sr *gpb.SetResponse) (client.SetResponse, error) {
	resp := client.SetResponse{
		TS: time.Unix(0, sr.GetTimestamp()),
	}
	var errs []string
	for _, r := range sr.GetResponse() {
		if r.Message != nil {
			errs = append(errs, r.GetMessage().String())
		}
	}
	if len(errs) > 0 {
		return resp, errors.New(strings.Join(errs, "; "))
	}

	return resp, nil
}

func getType(t client.Type) gpb.SubscriptionList_Mode {
	switch t {
	case client.Once:
		return gpb.SubscriptionList_ONCE
	case client.Stream:
		return gpb.SubscriptionList_STREAM
	case client.Poll:
		return gpb.SubscriptionList_POLL
	}
	return gpb.SubscriptionList_ONCE
}

func subscribe(q client.Query) (*gpb.SubscribeRequest, error) {
	s := &gpb.SubscribeRequest_Subscribe{
		Subscribe: &gpb.SubscriptionList{
			Mode:   getType(q.Type),
			Prefix: &gpb.Path{Target: q.Target},
		},
	}
	for _, qq := range q.Queries {
		pp, err := ygot.StringToPath(pathToString(qq), ygot.StructuredPath, ygot.StringSlicePath)
		if err != nil {
			return nil, fmt.Errorf("invalid query path %q: %v", qq, err)
		}
		s.Subscribe.Subscription = append(s.Subscribe.Subscription, &gpb.Subscription{Path: pp})
	}
	return &gpb.SubscribeRequest{Request: s}, nil
}

func noti(prefix []string, pp *gpb.Path, ts time.Time, u *gpb.Update) (client.Notification, error) {
	sp, err := ygot.PathToStrings(pp)
	if err != nil {
		return nil, fmt.Errorf("converting path %v to []string: %v", u.GetPath(), err)
	}
	// Make a full new copy of prefix + u.Path to avoid any reuse of underlying
	// slice arrays.
	p := make([]string, 0, len(prefix)+len(sp))
	p = append(p, prefix...)
	p = append(p, sp...)

	if u == nil {
		return client.Delete{Path: p, TS: ts}, nil
	}
	if u.Val != nil {
		val, err := value.ToScalar(u.Val)
		if err != nil {
			return nil, err
		}
		return client.Update{Path: p, TS: ts, Val: val}, nil
	}
	switch v := u.Value; v.Type {
	case gpb.Encoding_BYTES:
		return client.Update{Path: p, TS: ts, Val: v.Value}, nil
	case gpb.Encoding_JSON, gpb.Encoding_JSON_IETF:
		var val interface{}
		if err := json.Unmarshal(v.Value, &val); err != nil {
			return nil, fmt.Errorf("json.Unmarshal(%q, val): %v", v, err)
		}
		return client.Update{Path: p, TS: ts, Val: val}, nil
	default:
		return nil, fmt.Errorf("Unsupported value type: %v", v.Type)
	}
}

func init() {
	client.Register(Type, New)
}

// ProtoResponse converts client library Notification types into gNMI
// SubscribeResponse proto. An error is returned if any notifications have
// invalid paths or if update values can't be converted to gpb.TypedValue.
func ProtoResponse(notifs ...client.Notification) (*gpb.SubscribeResponse, error) {
	n := new(gpb.Notification)

	for _, nn := range notifs {
		switch nn := nn.(type) {
		case client.Update:
			if n.Timestamp == 0 {
				n.Timestamp = nn.TS.UnixNano()
			}

			pp, err := ygot.StringToPath(pathToString(nn.Path), ygot.StructuredPath, ygot.StringSlicePath)
			if err != nil {
				return nil, err
			}
			v, err := value.FromScalar(nn.Val)
			if err != nil {
				return nil, err
			}

			n.Update = append(n.Update, &gpb.Update{
				Path: pp,
				Val:  v,
			})

		case client.Delete:
			if n.Timestamp == 0 {
				n.Timestamp = nn.TS.UnixNano()
			}

			pp, err := ygot.StringToPath(pathToString(nn.Path), ygot.StructuredPath, ygot.StringSlicePath)
			if err != nil {
				return nil, err
			}
			n.Delete = append(n.Delete, pp)

		default:
			return nil, fmt.Errorf("gnmi.ProtoResponse: unsupported type %T", nn)
		}
	}

	resp := &gpb.SubscribeResponse{Response: &gpb.SubscribeResponse_Update{Update: n}}
	return resp, nil
}

func pathToString(q client.Path) string {
	qq := make(client.Path, len(q))
	copy(qq, q)
	// Escape all slashes within a path element. ygot.StringToPath will handle
	// these escapes.
	for i, e := range qq {
		qq[i] = strings.Replace(e, "/", "\\/", -1)
	}
	return strings.Join(qq, "/")
}
