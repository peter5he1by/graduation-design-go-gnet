package server_grpc

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"go-gnet/constant"
	"go-gnet/database"
	"go-gnet/database/mysql/model"
	"go-gnet/server/pool"
	pb "go-gnet/server_grpc/device"
	"go-gnet/util"
	"google.golang.org/grpc"
	"io"
	"net"
	"time"
)

type GrpcServer struct {
	pb.UnimplementedRouteGuideServer
}

func (s *GrpcServer) GetDeviceStatus(c pb.RouteGuide_GetDeviceStatusServer) error {
	for {
		in, err := c.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		id := uint(in.GetId())
		if status, ok := pool.CachedDevicesStatus[id]; ok {
			err = c.Send(&pb.DeviceStatus{Id: uint64(id), Status: uint64(status)})
		} else {
			// 取最近一次状态变更记录
			event, err := database.Mysql.SelectLatestStatusChangeEvent(id)
			if err != nil {
				return err
			}
			if event != nil {
				err = c.Send(&pb.DeviceStatus{Id: uint64(id), Status: uint64(event.Status)})
			} else {
				err = c.Send(&pb.DeviceStatus{Id: uint64(id), Status: uint64(constant.DeviceStatusUnknown)})
			}
		}
		if err != nil {
			return err
		}
	}
}

func (s *GrpcServer) GetDeviceConfig(context context.Context, d *pb.DeviceId) (*pb.DeviceConfig, error) {
	log.Debug("grpc call: GetDeviceConfig")
	id := uint(d.Id)
	pool.PoolLock.Lock()
	conn := pool.LoginConnectionPool[id]
	pool.PoolLock.Unlock()
	// 表示设备不在线
	if conn == nil {
		log.Debug("device not online")
		return &pb.DeviceConfig{}, nil
	}
	// 是否需要发包
	pool.WaitingForConfigLock.Lock()
	_, ok := pool.WaitingForConfig[id]
	if !ok {
		pool.WaitingForConfig[id] = 0
		pool.WaitingForConfigLock.Unlock()
		err := util.Send(conn.Conn, constant.ServerRequestUpdateDeviceConfig, nil)
		if err != nil {
			return nil, err
		}
	} else {
		pool.WaitingForConfigLock.Unlock()
	}
	// 轮询
	mutex := 5000
	for mutex > 0 {
		log.Debug("waiting for config...")
		time.Sleep(200)
		mutex -= 200
		pool.WaitingForConfigLock.Lock()
		if pool.WaitingForConfig[id] == 0 {
			if mutex <= 0 {
				delete(pool.WaitingForConfig, id)
				pool.WaitingForConfigLock.Unlock()
				log.Debug("waiting for config timeout")
				return &pb.DeviceConfig{}, nil
			}
			pool.WaitingForConfigLock.Unlock()
			continue
		}
		c := pool.WaitingForConfig[id]
		delete(pool.WaitingForConfig, id)
		pool.WaitingForConfigLock.Unlock()
		// 最好的结果
		log.Debugf("config found: %d", c)
		return &pb.DeviceConfig{
			Id: uint64(c),
		}, nil
	}
	return &pb.DeviceConfig{}, nil
}

func (s *GrpcServer) SetDeviceConfig(context context.Context, c *pb.DeviceConfig) (*pb.ReturnValue, error) {
	id := uint(c.Id)
	pool.PoolLock.Lock()
	if pool.LoginConnectionPool[id] == nil {
		pool.PoolLock.Unlock()
		// 1 设备不在线
		return &pb.ReturnValue{Ret: 1}, nil
	}
	err := util.Send(pool.LoginConnectionPool[id].Conn, constant.ServerIssueDeviceConfig, &model.DeviceConfig{
		DeviceID: id,
		Content:  c.Content,
	})
	pool.PoolLock.Unlock()
	if err != nil {
		// -1 下发失败
		return &pb.ReturnValue{Ret: -1}, nil
	}
	// 0 下发完成
	return &pb.ReturnValue{Ret: 0}, nil
}

func StartGrpcServer(port int) {
	lis, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterRouteGuideServer(s, &GrpcServer{})
	log.Infof("goroutine: grpc server is listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
