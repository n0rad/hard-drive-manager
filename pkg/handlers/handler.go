package handlers

import (
	"github.com/n0rad/go-erlog/data"
	"github.com/n0rad/go-erlog/logs"
	"github.com/n0rad/hard-disk-manager/pkg/system"
)

type handler struct {
	filter HandlerFilter
	new func() Handler
}

var handlers []handler


type Handler interface {
	Init(manager *BlockDeviceManager)
	Start()
	Stop()
	Name() string
}

/////////////////////////

type HandlerFilter struct {
	Type   string
	FSType string
}

func (h HandlerFilter) Match(filter HandlerFilter) bool {
	typeRes := true
	if h.Type != "" {
		if filter.Type == h.Type {
			typeRes = true
		} else {
			typeRes = false
		}
	}

	fsTypeRes := true
	if h.FSType != "" {
		if filter.FSType == h.FSType {
			fsTypeRes = true
		} else {
			fsTypeRes = false
		}
	}
	return typeRes && fsTypeRes
}

///////////////////////////

type CommonHandler struct {
	handlerName string
	disk        *system.Disk
	server      system.Server
	fields      data.Fields
	manager     *BlockDeviceManager
	stop        chan struct{}
}

func (h *CommonHandler) Init(manager *BlockDeviceManager) {
	h.fields = data.WithField("path", manager.Path)
	h.server = system.Server{}
	h.manager = manager
	h.stop = make(chan struct{})

	if err := h.server.Init(); err != nil {
		logs.WithE(err).Error("fail")
	}
}

func (h *CommonHandler) Stop() {
	close(h.stop)
}

func (h *CommonHandler) Name() string {
	return h.handlerName
}