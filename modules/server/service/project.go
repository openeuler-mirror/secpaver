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
	"archive/zip"
	"bytes"
	"context"
	"fmt"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"path/filepath"
	pb "gitee.com/openeuler/secpaver/api/proto/project"
	"gitee.com/openeuler/secpaver/common/ack"
	"gitee.com/openeuler/secpaver/common/errdefs"
	"gitee.com/openeuler/secpaver/common/log"
	"gitee.com/openeuler/secpaver/common/project"
	"gitee.com/openeuler/secpaver/engine"
	repo "gitee.com/openeuler/secpaver/repository"
	"sync"
	"time"
)

// ProjectService is the implement of ProjectMgr service
type ProjectService struct {
	sync.Mutex
	pb.UnimplementedProjectMgrServer

	repo    repo.Repo
	engines map[string]engine.Engine
	Name    string
}

// NewProjectService is the constructor of ProjectService
func NewProjectService(repo repo.Repo, engines map[string]engine.Engine) (*ProjectService, error) {
	if repo == nil {
		return nil, fmt.Errorf("nil repository handle")
	}

	if engines == nil {
		return nil, fmt.Errorf("nil engines map")
	}

	return &ProjectService{
		Name:    "project",
		repo:    repo,
		engines: engines,
	}, nil
}

// ServiceName returns the name of service
func (s *ProjectService) ServiceName() string {
	return s.Name
}

// RegisterServer registers service in the grpc server
func (s *ProjectService) RegisterServer(server *grpc.Server) error {
	pb.RegisterProjectMgrServer(server, s)
	return nil
}

// ListProject is the remote implement of the ListProject rpc
func (s *ProjectService) ListProject(ctx context.Context, req *pb.Req) (*pb.ListProjectAck, error) {
	s.Lock()
	defer s.Unlock()

	log.Infof("rpc request: ListProject")

	projects, err := s.repo.FindAllProjects()
	if err != nil {
		return &pb.ListProjectAck{}, errors.Wrap(err, "fail to find all projects")
	}

	ackInfo := &pb.ListProjectAck{}
	for _, prj := range projects {
		ver, err := project.ReadProjectMetaFromDir(prj.Path)
		if err != nil {
			log.Warnf("fail to get meta for project %s", prj.Name)
			continue
		}

		ackInfo.ProjectInfos = append(ackInfo.ProjectInfos, &pb.ListProjectAck_ProjectInfoSimple{
			Name:    prj.Name,
			Version: ver.Version,
		})
	}

	return ackInfo, nil
}

// ImportProject is the remote implement of the ImportProject rpc
func (s *ProjectService) ImportProject(ctx context.Context, req *pb.ImportProjectReq) (*pb.Ack, error) {
	s.Lock()
	defer s.Unlock()

	if err := checkImportProjectReq(req); err != nil {
		log.Errorf("fail to check request args")
		return &pb.Ack{}, errors.Wrap(err, "fail to check request args")
	}

	log.Infof("rpc request: ImportProject, filename: %s", req.GetFile().GetFilename())

	zipReader, err := zip.NewReader(
		bytes.NewReader(req.GetFile().GetData()), int64(len(req.GetFile().GetData())))
	if err != nil {
		return nil, errors.Wrap(err, "fail to read zip file")
	}

	meta, err := project.ReadProjectMetaFromZip(zipReader)
	if err != nil {
		return nil, errors.Wrap(err, "fail to parse project from zip")
	}

	// check if the project exist
	prj, err := s.repo.FindProjectByName(meta.GetName())
	if err == nil {
		if !req.GetForce() {
			return &pb.Ack{}, fmt.Errorf("project exists")
		}

		if err := s.repo.UpdateProjectByZip(prj.Name, zipReader); err != nil {
			log.Errorf("fail to update project: %v", err)
			return &pb.Ack{}, errors.Wrap(err, "fail to update project")
		}

	} else if errdefs.IsNotFoundError(errors.Cause(err)) { // project not exist
		if err := s.repo.AddProjectByZip(meta.GetName(), zipReader); err != nil {
			log.Errorf("fail to add project: %v", err)
			return &pb.Ack{}, errors.Wrap(err, "fail to add project")
		}

	} else if err != nil {
		return &pb.Ack{}, errors.Wrap(err, "fail to find project")
	}

	log.Infof("finish importing project")

	return &pb.Ack{
		Level:  ack.LevelInfo,
		Status: fmt.Sprintf("Finish importing project"),
	}, nil
}

// DeleteProject is the remote implement of the DeleteProject rpc
func (s *ProjectService) DeleteProject(ctx context.Context, req *pb.Req) (*pb.Ack, error) {
	s.Lock()
	defer s.Unlock()

	if err := checkProjectReq(req); err != nil {
		log.Errorf("fail to check request args")
		return &pb.Ack{}, errors.Wrap(err, "fail to check request args")
	}

	log.Infof("rpc request: DeleteProject, name: %s", req.GetName())

	prj, err := s.repo.FindProjectByName(req.GetName())
	if err != nil {
		return &pb.Ack{}, errors.Wrap(err, "fail to find project")
	}

	if err := s.repo.DeleteProject(prj); err != nil {
		log.Errorf("fail to delete project: %v", err)
		return &pb.Ack{}, errors.Wrap(err, "fail to delete project")
	}

	log.Infof("finish deleting project")

	return &pb.Ack{
		Level:  ack.LevelInfo,
		Status: fmt.Sprintf("Finish deleting project"),
	}, nil
}

// InfoProject is the remote implement of the InfoProject rpc
func (s *ProjectService) InfoProject(ctx context.Context, req *pb.Req) (*pb.InfoProjectAck, error) {
	s.Lock()
	defer s.Unlock()

	if err := checkProjectReq(req); err != nil {
		log.Errorf("fail to check request args")
		return &pb.InfoProjectAck{}, errors.Wrap(err, "fail to check request args")
	}

	log.Infof("rpc request: InfoProject, name: %s", req.GetName())

	// check if the project is existed
	prj, err := s.repo.FindProjectByName(req.GetName())
	if err != nil {
		return &pb.InfoProjectAck{}, errors.Wrap(err, "fail to find project")
	}

	// parse project struct and check project is valid
	meta, err := project.ReadProjectMetaFromDir(prj.Path)
	if err != nil {
		return &pb.InfoProjectAck{}, errors.Wrap(err, "fail to read project meta")
	}

	return &pb.InfoProjectAck{
		MetaInfo: meta,
	}, nil
}

// BuildProject is the remote implement of the BuildProject rpc
func (s *ProjectService) BuildProject(req *pb.BuildProjectReq, srv pb.ProjectMgr_BuildProjectServer) error {
	s.Lock()
	defer s.Unlock()

	if err := checkBuildProjectReq(req); err != nil {
		log.Errorf("fail to check request args")
		return errors.Wrap(err, "fail to check request args")
	}

	log.Infof("rpc request: BuildProject, local project: %s, remote project: %s, engine: %s",
		req.GetProject().GetMeta().GetName(), req.GetRemoteProject(), req.GetEngine())

	projectInfo := req.GetProject()
	if req.GetRemote() {
		prj, err := s.repo.FindProjectByName(req.GetRemoteProject())
		if err != nil {
			return errors.Wrap(err, "fail to find project")
		}

		projectInfo, err = project.ParseProjectFromDir(prj.Path)
		if err != nil {
			return errors.Wrap(err, "fail to parse project")
		}
	}

	builder, ok := s.engines[req.GetEngine()]
	if !ok || builder == nil {
		return fmt.Errorf("invalid engine %s", req.GetEngine())
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

	outDir := filepath.Join(s.repo.GetPolicyRoot(), req.GetEngine())

	if err := builder.Build(projectInfo, outDir, ch); err != nil {
		log.Errorf("fail to build project: %v", err)
		return errors.Wrap(err, "fail to build project")
	}

	ch <- &pb.Ack{Level: ack.LevelInfo, Status: "Finish building project"}

	time.Sleep(time.Second)
	log.Infof("finish building project")

	return nil
}

// ExportProject is the remote implement of the ExportProject rpc
func (s *ProjectService) ExportProject(ctx context.Context, req *pb.Req) (*pb.ExportProjectAck, error) {
	s.Lock()
	defer s.Unlock()

	if err := checkProjectReq(req); err != nil {
		log.Errorf("fail to check request args")
		return &pb.ExportProjectAck{}, errors.Wrap(err, "fail to check request args")
	}

	log.Infof("rpc request: ExportProject, name: %s", req.GetName())

	prj, err := s.repo.FindProjectByName(req.GetName())
	if err != nil {
		return &pb.ExportProjectAck{}, errors.Wrap(err, "fail to find project")
	}

	zipName := "export_project_" + req.GetName() + ".zip"
	data, err := s.repo.ExportProjectZip(prj, zipName)
	if err != nil {
		log.Errorf("fail to export project: %v", err)
		return &pb.ExportProjectAck{}, errors.Wrap(err, "fail to export project zip")
	}

	log.Infof("finish exporting project")

	return &pb.ExportProjectAck{
		File: &pb.ProjectZipFile{
			Filename: zipName,
			Data:     data,
		},
	}, nil
}
