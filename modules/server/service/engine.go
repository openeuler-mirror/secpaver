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

package service

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	pbEngine "secpaver/api/proto/engine"
	"secpaver/common/log"
	"secpaver/engine"
	"sync"
)

// EngineService is the implement of EngineMgr service
type EngineService struct {
	sync.Mutex

	engines map[string]engine.Engine

	Name string
	pbEngine.UnimplementedEngineMgrServer
}

// NewEngineService is the constructor of EngineService
func NewEngineService(engines map[string]engine.Engine) (*EngineService, error) {
	if engines == nil {
		return nil, fmt.Errorf("nil engines map")
	}

	return &EngineService{
		Name:    "engine",
		engines: engines,
	}, nil
}

// ServiceName returns the name of service
func (s *EngineService) ServiceName() string {
	return s.Name
}

// RegisterServer registers service in the grpc server
func (s *EngineService) RegisterServer(server *grpc.Server) error {
	pbEngine.RegisterEngineMgrServer(server, s)
	return nil
}

// ListEngine is the remote implement of the ListEngine rpc
func (s *EngineService) ListEngine(tx context.Context, req *pbEngine.Req) (*pbEngine.ListEngineAck, error) {
	s.Lock()
	defer s.Unlock()

	log.Infof("rpc request: ListEngine")

	ack := &pbEngine.ListEngineAck{}
	for _, eng := range s.engines {
		if eng == nil {
			continue
		}

		ack.EngineInfos = append(ack.EngineInfos, &pbEngine.EngineInfo{
			Name: eng.GetName(),
			Desc: eng.GetDescription(),
		})
	}

	return ack, nil
}

// InfoEngine is the remote implement of the InfoEngine rpc
func (s *EngineService) InfoEngine(ctx context.Context, req *pbEngine.Req) (*pbEngine.InfoEngineAck, error) {
	s.Lock()
	defer s.Unlock()

	if err := checkEngineReq(req); err != nil {
		log.Errorf("fail to check request args")
		return &pbEngine.InfoEngineAck{}, errors.Wrap(err, "fail to check request args")
	}

	log.Infof("rpc request: InfoEngine, name: %s", req.GetName())

	eng, ok := s.engines[req.GetName()]
	if !ok || eng == nil {
		return &pbEngine.InfoEngineAck{}, fmt.Errorf("invalid engine %s", req.GetName())
	}

	return &pbEngine.InfoEngineAck{
		BaseInfo: &pbEngine.EngineInfo{
			Name: eng.GetName(),
			Desc: eng.GetDescription(),
		},
	}, nil
}
