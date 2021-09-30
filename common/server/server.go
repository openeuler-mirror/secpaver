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
Package server implements the grpc server and defines the services of server.
*/
package server

import (
	"fmt"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"net"
)

// Service is the server service interface
type Service interface {
	ServiceName() string
	RegisterServer(*grpc.Server) error
}

// Server is the grpc server package with the services
type Server struct {
	server   *grpc.Server
	services map[string]Service
}

// NewServer creates a new grpc server with blank service map
func NewServer(rpc *grpc.Server) *Server {
	return &Server{
		server:   rpc,
		services: make(map[string]Service),
	}
}

// LoadService load services to grpc server
func (s *Server) LoadService(services []Service) error {
	for _, svc := range services {
		if err := svc.RegisterServer(s.server); err != nil {
			return errors.Wrapf(err, "fail to register service %s", svc.ServiceName())
		}

		s.services[svc.ServiceName()] = svc
	}

	return nil
}

// Stop stops the server and does clean up work
func (s *Server) Stop() {
	if s.server != nil {
		s.server.Stop()
	}
}

// Serve runs the server
func (s *Server) Serve(lis net.Listener) error {
	if s.server == nil {
		return fmt.Errorf("nil grpc server")
	}

	return s.server.Serve(lis)
}
