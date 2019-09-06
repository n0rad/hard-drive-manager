package agent

import "github.com/n0rad/go-erlog/logs"

type HandlerAdd struct {
	CommonHandler
}

func (h *HandlerAdd) Start() {

	watch := h.manager.PassService.Watch()
	<-watch

	buffer, err := h.manager.PassService.Get()
	if err != nil {
		logs.WithEF(err, h.fields).Error("Cannot get password to add disk")
		return
	}

	if err := h.disk.Add(buffer); err != nil {
		logs.WithEF(err, h.fields).Error("Failed to add disk")
	}



	//disk, err := h.server.ScanDisk(h.path)
	//if err != nil {
	//	logs.WithE(err).Error("Failed to scan disk")
	//	return
	//}

	//h.disk.Add()
}
