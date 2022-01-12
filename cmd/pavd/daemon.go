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

package main

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
	"google.golang.org/grpc"
	"net"
	"os"
	"os/signal"
	"path/filepath"
	"gitee.com/openeuler/secpaver/common/config"
	"gitee.com/openeuler/secpaver/common/global"
	"gitee.com/openeuler/secpaver/common/log"
	"gitee.com/openeuler/secpaver/common/pid"
	"gitee.com/openeuler/secpaver/common/server"
	"gitee.com/openeuler/secpaver/common/utils"
	"gitee.com/openeuler/secpaver/engine"
	"gitee.com/openeuler/secpaver/modules/server/service"
	"gitee.com/openeuler/secpaver/repository"
	"syscall"
	"time"
)

const signalProcInterval = 100

func listenSocket(sock string) (net.Listener, error) {
	if !utils.IsUnixFilePath(sock) {
		return nil, fmt.Errorf("invalid unix socket filepath")
	}

	if err := os.RemoveAll(sock); err != nil {
		return nil, fmt.Errorf("fail to remove old unix socket file")
	}

	if err := os.MkdirAll(filepath.Dir(sock), global.DefaultDirPerm); err != nil {
		return nil, fmt.Errorf("fail to create directory for unix socket %s", filepath.Base(sock))
	}

	lis, err := net.Listen("unix", sock)
	if err != nil {
		return nil, errors.Wrap(err, "fail to listen unix socket")
	}

	if err := os.Chmod(sock, global.SocketFilePerm); err != nil {
		return nil, fmt.Errorf("fail to chmod unix socket file %s", filepath.Base(sock))
	}

	log.Infof("grpc server is listening on %s", filepath.Base(sock))

	return lis, nil
}

func runPavd(ctx *cli.Context) error {
	// create pid file lock
	pl, err := pid.CreatePidLock(global.DefaultPidFile)
	if err != nil {
		return errors.Wrap(err, "fail to create pid lock")
	}

	if err := pl.Lock(); err != nil {
		return fmt.Errorf(
			"error starting daemon: pid file found, ensure Pavd is not running or delete %s", global.DefaultPidFile)
	}

	defer func() {
		_ = pl.Unlock()
		pl.Release()
	}()

	// parse config file
	confPath := ctx.String("config")
	if confPath == "" {
		confPath = global.DefaultConfigFile
	}

	conf, err := config.ParseConfig(confPath)
	if err != nil {
		return errors.Wrapf(err, "fail to parse config file %s", filepath.Base(confPath))
	}

	conf = config.MergeConfigInfo(conf, ctx)

	// initialize log
	if err := log.InitLogger(&conf.Log); err != nil {
		return errors.Wrap(err, "fail to init logger")
	}
	defer log.Close()

	// initialize repository
	if err := repository.InitRepo(&conf.Repository); err != nil {
		return errors.Wrap(err, "fail to init repository")
	}

	// load engine plugins
	engines, err := engine.LoadEnginePlugin(global.DefaultPluginRoot)
	if err != nil {
		return errors.Wrap(err, "fail to load engine plugins")
	}

	// load all services
	services, err := service.NewServerServices(repository.GetRepo(), engines)
	if err != nil {
		return errors.Wrap(err, "fail to create services")
	}

	return runServices(conf.Connect.Grpc.Socket, services)
}

func runServices(sock string, services []server.Service) error {
	svr := server.NewServer(grpc.NewServer())
	if err := svr.LoadService(services); err != nil {
		return errors.Wrap(err, "fail to load service")
	}

	defer func() {
		log.Info("shutdown grpc server")
		svr.Stop()
	}()

	lis, err := listenSocket(sock)
	if err != nil {
		log.Errorf("fail to listen grpc socket: %v", err)
		return errors.Wrap(err, "fail to listen grpc socket")
	}

	go signalProc()

	if err := svr.Serve(lis); err != nil {
		log.Errorf("fail to serve grpc server: %v", err)
		return errors.Wrap(err, "fail to serve grpc server")
	}

	return nil
}

func signalProc() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGSTOP, syscall.SIGKILL)

	for true {
		select {
		case s, ok := <-sigChan:
			if !ok {
				return
			}

			log.Infof("sig recv %s", s.String())
			os.Exit(0)

		default:
			time.Sleep(time.Millisecond * signalProcInterval)
		}
	}
}
