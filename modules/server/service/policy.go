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
	pb "gitee.com/openeuler/secpaver/api/proto/policy"
	"gitee.com/openeuler/secpaver/common/ack"
	"gitee.com/openeuler/secpaver/common/log"
	"gitee.com/openeuler/secpaver/common/policy"
	"gitee.com/openeuler/secpaver/engine"
	repo "gitee.com/openeuler/secpaver/repository"
	"sync"
	"time"
)

// PolicyService is the implement of ProjectMgr service
type PolicyService struct {
	sync.Mutex
	pb.UnimplementedPolicyMgrServer

	repo    repo.Repo
	engines map[string]engine.Engine
	Name    string
}

// NewPolicyService is the constructor of PolicyService
func NewPolicyService(repo repo.Repo, engines map[string]engine.Engine) (*PolicyService, error) {
	if repo == nil {
		return nil, fmt.Errorf("nil repository handle")
	}

	if engines == nil {
		return nil, fmt.Errorf("nil engines map")
	}

	return &PolicyService{
		Name:    "policy",
		repo:    repo,
		engines: engines,
	}, nil
}

// ServiceName returns the name of service
func (s *PolicyService) ServiceName() string {
	return s.Name
}

// RegisterServer registers service in the grpc server
func (s *PolicyService) RegisterServer(server *grpc.Server) error {
	pb.RegisterPolicyMgrServer(server, s)
	return nil
}

// ListPolicy is the remote implement of the ListPolicy rpc
func (s *PolicyService) ListPolicy(ctx context.Context, req *pb.Req) (*pb.ListPolicyAck, error) {
	s.Lock()
	defer s.Unlock()

	log.Infof("rpc request: ListPolicy")

	policies, err := s.repo.FindAllPolicies()
	if err != nil {
		return nil, errors.Wrap(err, "fail to find all policies")
	}

	ackInfo := &pb.ListPolicyAck{}

	for _, p := range policies {
		mgr := s.engines[p.Engine]
		if mgr == nil {
			continue
		}

		fName := p.Name + "_" + p.Engine

		status, err := mgr.GetPolicyStatus(p)
		if err != nil {
			log.Warnf("fail to get status of %s: %v", fName, err)
			status = "N/A"
		}

		ackInfo.PolicyInfos = append(ackInfo.PolicyInfos, &pb.PolicyInfo{
			Name:   fName,
			Status: status,
		})
	}

	return ackInfo, nil
}

// InstallPolicy is the remote implement of the InstallPolicy rpc
func (s *PolicyService) InstallPolicy(req *pb.Req, srv pb.PolicyMgr_InstallPolicyServer) error {
	s.Lock()
	defer s.Unlock()

	if err := checkPolicyReq(req); err != nil {
		log.Errorf("fail to check request args")
		return errors.Wrap(err, "fail to check request args")
	}

	log.Infof("rpc request: InstallPolicy, name: %s, engine: %s",
		req.GetName(), req.GetEngine())

	p, err := s.repo.FindPolicyByName(req.Name, req.Engine)
	if err != nil {
		return errors.Wrap(err, "fail to find policy")
	}

	mgr, ok := s.engines[p.Engine]
	if !ok || mgr == nil {
		return fmt.Errorf("invalid engine %s", req.Engine)
	}

	ch := make(chan *pb.Ack, 1)
	defer close(ch)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		for {
			select {
			case value, ok := <-ch:
				if !ok {
					return
				}

				_ = srv.Send(value)
			case <-ctx.Done():
				return
			}
		}
	}()

	if err := mgr.Install(p, ch); err != nil {
		log.Errorf("fail to install policy: %v", err)
		return errors.Wrap(err, "fail to install policy")
	}

	time.Sleep(time.Second)
	log.Infof("finish installing policy")

	return nil
}

// ExportPolicy is the remote implement of the ExportPolicy rpc
func (s *PolicyService) ExportPolicy(ctx context.Context, req *pb.Req) (*pb.ExportPolicyAck, error) {
	s.Lock()
	defer s.Unlock()

	if err := checkPolicyReq(req); err != nil {
		log.Errorf("fail to check request args")
		return &pb.ExportPolicyAck{}, errors.Wrap(err, "fail to check request args")
	}

	log.Infof("rpc request: ExportPolicy, name: %s, engine: %s",
		req.GetName(), req.GetEngine())

	p, err := s.repo.FindPolicyByName(req.Name, req.Engine)
	if err != nil {
		return &pb.ExportPolicyAck{}, errors.Wrap(err, "fail to find policy")
	}

	zipName := "export_policy_" + req.GetName() + "_" + req.GetEngine() + ".zip"
	data, err := s.repo.ExportPolicyZip(p, zipName)
	if err != nil {
		log.Errorf("fail to export policy: %v", err)
		return &pb.ExportPolicyAck{}, errors.Wrap(err, "fail to export policy zip")
	}

	log.Infof("finish exporting policy")

	return &pb.ExportPolicyAck{
		File: &pb.PolicyZipFile{
			Filename: zipName,
			Data:     data,
		},
	}, nil
}

// UninstallPolicy is the remote implement of the UninstallPolicy rpc
func (s *PolicyService) UninstallPolicy(req *pb.Req, srv pb.PolicyMgr_UninstallPolicyServer) error {
	s.Lock()
	defer s.Unlock()

	if err := checkPolicyReq(req); err != nil {
		log.Errorf("fail to check request args")
		return errors.Wrap(err, "fail to check request args")
	}

	log.Infof("rpc request: UninstallPolicy, name: %s, engine: %s",
		req.GetName(), req.GetEngine())

	p, err := s.repo.FindPolicyByName(req.Name, req.Engine)
	if err != nil {
		return errors.Wrap(err, "fail to find policy")
	}

	mgr, ok := s.engines[p.Engine]
	if !ok || mgr == nil {
		return fmt.Errorf("invalid engine %s", req.Engine)
	}

	ch := make(chan *pb.Ack, 1)
	defer close(ch)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		for {
			select {
			case value, ok := <-ch:
				if !ok {
					return
				}

				_ = srv.Send(value)
			case <-ctx.Done():
				return
			}
		}
	}()

	if err := mgr.Uninstall(p, ch); err != nil {
		log.Errorf("fail to uninstall policy: %v", err)
		return errors.Wrap(err, "fail to uninstall policy")
	}

	time.Sleep(time.Second)
	log.Infof("finish uninstalling policy")

	return nil
}

// DeletePolicy is the remote implement of the DeletePolicy rpc
func (s *PolicyService) DeletePolicy(ctx context.Context, req *pb.Req) (*pb.Ack, error) {
	s.Lock()
	defer s.Unlock()

	if err := checkPolicyReq(req); err != nil {
		log.Errorf("fail to check request args")
		return &pb.Ack{}, errors.Wrap(err, "fail to check request args")
	}

	log.Infof("rpc request: DeletePolicy, name: %s, engine: %s",
		req.GetName(), req.GetEngine())

	p, err := s.repo.FindPolicyByName(req.Name, req.Engine)
	if err != nil {
		return &pb.Ack{}, errors.Wrap(err, "fail to find policy")
	}

	mgr, ok := s.engines[p.Engine]
	if !ok || mgr == nil {
		return &pb.Ack{}, fmt.Errorf("invalid engine %s", req.Engine)
	}

	status, err := mgr.GetPolicyStatus(p)
	if err != nil {
		log.Errorf("fail to get policy status: %v", err)
		return &pb.Ack{}, errors.Wrap(err, "fail to get policy status")
	}

	if status == policy.StatusActive {
		return &pb.Ack{}, fmt.Errorf("please uninstall policy first")
	}

	if err := s.repo.DeletePolicy(p); err != nil {
		log.Errorf("fail to delete policy: %v", err)
		return &pb.Ack{}, err
	}

	log.Infof("finish deleting policy")

	return &pb.Ack{
		Level:  ack.LevelInfo,
		Status: "Finish deleting policy",
	}, nil
}
