/*
 * Copyright (c) Huawei Technologies Co., Ltd. 2020-2021. All rights reserved.
 * secPaver is licensed under the Mulan PSL v2.
 * You can use this software according to the terms and conditions of the Mulan PSL v2.
 * You may obtain a copy of Mulan PSL v2 at:
 *     http://license.coscl.org.cn/MulanPSL2
 * THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR
 * PURPOSE.
 * See the Mulan PSL v2 for more details.
 */

/*
Package client is a tool of the grpc client creation.
*/
package client

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
	"google.golang.org/grpc"
	"net"
	"secpaver/common/global"
	"time"
)

const maxDelay = 5

// Client is the grpc client struct
type Client struct {
	conn *grpc.ClientConn
}

// Connection returns the grpc connection of the client
func (c *Client) Connection() *grpc.ClientConn {
	return c.conn
}

// Close method close the grpc connection
func (c *Client) Close() {
	_ = c.conn.Close()
}

// NewClient creates a grpc client with the grpc connection
func NewClient(conn *grpc.ClientConn) (*Client, error) {
	return &Client{
		conn: conn,
	}, nil
}

// NewClientFromContext method create a grpc client with the cmd context
// now secPaver can only allow local connections by unix socket
func NewClientFromContext(ctx *cli.Context, opts ...grpc.DialOption) (*Client, error) {
	addr := ctx.GlobalString("socket")
	return NewClientWithConnectParam("unix", addr, opts...)
}

// NewClientWithConnectParam method create a grpc with the specified connect parameters
func NewClientWithConnectParam(protocol, address string, opts ...grpc.DialOption) (*Client, error) {
	if protocol != "unix" {
		return nil, fmt.Errorf("now secPaver can only allow local connections by unix socket")
	}

	dialer := func(ctx context.Context, addr string) (net.Conn, error) {
		unixAddr, err := net.ResolveUnixAddr("unix", addr)
		if err != nil {
			return nil, errors.Wrap(err, "fail to resolve unix socket address")
		}

		return net.DialUnix("unix", nil, unixAddr)
	}

	opts = append(opts, grpc.WithContextDialer(dialer))
	opts = append(opts, []grpc.DialOption{
		grpc.WithBlock(),
		grpc.WithInsecure(),
		grpc.FailOnNonTempDialError(true),
		grpc.WithBackoffMaxDelay(maxDelay * time.Second),
	}...)

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, global.DefaultTimeout*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, address, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to dial %q : %v", address, err)
	}

	return NewClient(conn)
}
