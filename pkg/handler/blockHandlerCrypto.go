package handler

import (
	"github.com/n0rad/go-erlog/data"
	"github.com/n0rad/go-erlog/errs"
	"github.com/n0rad/go-erlog/logs"
	"github.com/n0rad/hard-disk-manager/pkg/system"
)

func init() {
	BlockHandlers["crypto"] = BlockHandlerBuilder{
		Filter: HandlerFilter{FSType: "crypto_LUKS"},
		New: func() BlockHandler {
			return &HandlerCrypto{}
		},
	}
}

type HandlerCrypto struct {
	CommonBlockHandler
}

//func (h *HandlerCrypto) Start() error {
//passwordSet := h.manager.hdm.PassService().Watch()
//defer h.manager.hdm.PassService().Unwatch(passwordSet)
//for {
//	select {
//	case <-passwordSet:
//		if err := h.open(); err != nil {
//			logs.WithE(err).Error("Failed to open crypto")
//		}
//	case <-h.stop:
//		return nil
//	}
//}
//}

func (h *HandlerCrypto) Remove() error {
	if len(h.manager.Block.Children) == 1 {
		if err := h.manager.Block.Children[0].LuksClose(); err != nil {
			return err
		}
	}
	return nil
}

func (h *HandlerCrypto) Add() error {
	if h.manager.Block.IsLuksOpen() {
		logs.WithF(h.GetFields()).Info("Already Open")
		return nil
	}

	b := []byte("aa")
	if err := h.manager.GetHDM().PassService().FromBytes(&b); err != nil {
		return errs.WithE(err, "Cannot get password")
	}

	//if err := h.manager.hdm.PassService().AskPassword(false); err != nil {
	//	return errs.WithE(err, "Cannot get password")
	//}

	used, err := h.manager.Block.IsLuksNameUsed()
	if err != nil {
		logs.WithEF(err, h.GetFields()).Warn("Cannot check is luks name is already in use, continuing but it may fail")
	}
	if used {
		if h.manager.Block.IsLuksUsed() {
			logs.WithF(h.GetFields()).Warn("Luks already open for same Block but not linked to device, trying to cleanup")

			h.cleanupRemovedBlockDevice(h.manager.Block.GetUsableLabel())
		} else {
			return errs.WithF(h.GetFields().WithField("name", h.manager.Block.GetUsableLabel()), "deviceMapper name is already used by another Block")
		}
	}

	if !h.manager.GetHDM().PassService().IsSet() {
		logs.WithF(h.GetFields()).Debug("Password is not set, cannot open")
		return nil
	}

	buffer, err := h.manager.GetHDM().PassService().Get()
	if err != nil {
		return errs.WithEF(err, h.GetFields(), "Failed to get password from password service")
	}
	defer buffer.Destroy()

	if err := h.manager.Block.LuksOpen(buffer); err != nil {
		return errs.WithEF(err, h.GetFields(), "Failed to luks open")
	}
	return nil
}

func (h *HandlerCrypto) cleanupRemovedBlockDevice(label string) {
	blockDevicePath := "/dev/mapper/" + label

	mapper := system.DeviceMapper{
		Exec: h.manager.Block.GetExec(),
	}

	blockName, err := mapper.BlockFromName(label)
	if err != nil {
		logs.WithEF(err, data.WithField("blockDevice", blockDevicePath)).Debug("Cannot get Block name from blockDevice")
	}

	mountPoint := ""
	if mount, err := system.MountFromBlockDevice(blockDevicePath); err != nil {
		logs.WithEF(err, data.WithField("blockDevice", blockDevicePath)).Debug("Cannot get mount from blockDevice")
	} else if mount != nil {
		mountPoint = mount.Path
	}

	// Block device
	fakeOpenedBlockDevice := system.BlockDevice{
		Name:       label,
		Path:       blockDevicePath,
		Mountpoint: mountPoint,
		Label:      label,
		Kname:      blockName,
	}
	fakeOpenedBlockDevice.Init(h.manager.Block.GetExec())

	// manager
	manager := BlockManager{}
	manager.Init(fakeOpenedBlockDevice)

	// mount handler
	handlerMount := BlockHandlers[handlerNameMount].New()
	handlerMount.Init(handlerNameMount, &manager)

	go handlerMount.Start()
	if err := handlerMount.Remove(); err != nil {
		logs.WithE(err).Error("Failed to cleanup removed device")
	}

	// remove device mapper
	if err := mapper.Remove(label); err != nil {
		logs.WithE(err).Error("Cannot cleanup removed device") // TODO
	}
}
