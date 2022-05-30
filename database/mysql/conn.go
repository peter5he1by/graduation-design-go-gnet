package mysql

import (
	"errors"
	"go-gnet/database/mysql/model"
	"gorm.io/gorm"
)

type Handle struct {
	DB *gorm.DB
}

func (h Handle) SelectDeviceByUuid(uuid string) (*model.Device, error) {
	d := model.Device{}
	res := h.DB.Model(&model.Device{}).Where("uuid = ?", uuid).Limit(1).Take(&d)
	if res.Error != nil {
		if errors.Is(gorm.ErrRecordNotFound, res.Error) {
			return nil, nil
		}
		return nil, res.Error
	}
	return &d, nil
}

func (h Handle) InsertDeviceInfo(d *model.DeviceInfo) error {
	d2 := model.DeviceInfo{
		DeviceID:     d.DeviceID,
		SoftwareInfo: d.SoftwareInfo,
		HardwareInfo: d.HardwareInfo,
		Remark:       d.Remark,
	}
	result := h.DB.Create(&d2)
	if result.RowsAffected != 1 {
		return result.Error
	}
	return nil
}

func (h Handle) InsertDeviceEventLogin(deviceId uint, d *model.DeviceEventLogin) error {
	d2 := &model.DeviceEventLogin{
		Username:     d.Username,
		LoginPurpose: d.LoginPurpose,
	}
	d2event := &model.DeviceEvent{
		DeviceID: deviceId,
		Type:     d2.TableName(),
	}
	err := h.DB.Transaction(func(tx *gorm.DB) error {
		result := h.DB.Create(d2event)
		if result.RowsAffected != 1 {
			return result.Error
		}
		d2.EventID = d2event.ID
		result = h.DB.Create(d2)
		if result.RowsAffected != 1 {
			return result.Error
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (h Handle) InsertDeviceEventLogout(deviceId uint, d *model.DeviceEventLogout) error {
	d2 := &model.DeviceEventLogout{
		Username:     d.Username,
		WorkContents: d.WorkContents,
	}
	d2event := &model.DeviceEvent{
		DeviceID: deviceId,
		Type:     d2.TableName(),
	}
	err := h.DB.Transaction(func(tx *gorm.DB) error {
		result := h.DB.Create(d2event)
		if result.RowsAffected != 1 {
			return result.Error
		}
		d2.EventID = d2event.ID
		result = h.DB.Create(d2)
		if result.RowsAffected != 1 {
			return result.Error
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (h Handle) InsertDeviceStatusChange(deviceId uint, d *model.DeviceEventStatusChange) error {
	d2 := &model.DeviceEventStatusChange{
		Status: d.Status,
	}
	d2event := &model.DeviceEvent{
		DeviceID: deviceId,
		Type:     d2.TableName(),
	}
	err := h.DB.Transaction(func(tx *gorm.DB) error {
		result := h.DB.Create(d2event)
		if result.RowsAffected != 1 {
			return result.Error
		}
		d2.EventID = d2event.ID
		result = h.DB.Create(d2)
		if result.RowsAffected != 1 {
			return result.Error
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (h Handle) InsertDeviceEventOperation(deviceId uint, d *model.DeviceEventOperation) error {
	d2 := &model.DeviceEventOperation{
		Username:      d.Username,
		OperationType: d.OperationType,
		Detail:        d.Detail,
	}
	d2event := &model.DeviceEvent{
		DeviceID: deviceId,
		Type:     d2.TableName(),
	}
	err := h.DB.Transaction(func(tx *gorm.DB) error {
		result := h.DB.Create(d2event)
		if result.RowsAffected != 1 {
			return result.Error
		}
		d2.EventID = d2event.ID
		result = h.DB.Create(d2)
		if result.RowsAffected != 1 {
			return result.Error
		}
		return nil
	})
	return err
}

func (h Handle) InsertDeviceLog(l *model.DeviceLog) error {
	log := &model.DeviceLog{
		DeviceID:  l.DeviceID,
		Content:   l.Content,
		StartTime: l.StartTime,
		EndTime:   l.EndTime,
	}
	result := h.DB.Create(log)
	if result.RowsAffected != 1 {
		return result.Error
	}
	return nil
}

func (h Handle) SelectLatestDeviceConfig() (*model.DeviceConfig, error) {
	c := model.DeviceConfig{}
	result := h.DB.Model(&model.DeviceConfig{}).Order("updated_at DESC").Limit(1).Take(&c)
	if result.Error != nil {
		if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, result.Error
		}
		return nil, nil
	}
	return &c, nil
}

func (h Handle) InsertDeviceConfig(config *model.DeviceConfig) (uint, error) {
	d := &model.DeviceConfig{
		DeviceID: config.DeviceID,
		Content:  config.Content,
		Type:     config.Type,
	}
	result := h.DB.Create(d)
	if result.RowsAffected != 1 {
		return 0, result.Error
	}
	return d.ID, nil
}

func (h Handle) UpdateDeviceConfig(c *model.DeviceConfig) error {
	result := h.DB.UpdateColumns(
		// &model.DeviceConfig{
		// 	BaseModel: model.BaseModel{
		// 		ID:        c.ID,
		// 		CreatedAt: c.CreatedAt,
		// 		UpdatedAt: c.UpdatedAt,
		// 	},
		// 	DeviceID: c.DeviceID,
		// 	Content:  c.Content,
		// 	Type:     c.Type,
		// },
		c,
	)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (h Handle) InsertDeviceDataTemperature(d *model.DeviceDataTemperature) error {
	temperature := model.DeviceDataTemperature{
		DeviceID: d.DeviceID,
		Time:     d.Time,
		Data:     d.Data,
	}
	result := h.DB.Create(&temperature)
	if result.RowsAffected != 1 {
		return result.Error
	}
	return nil
}

func (h Handle) SelectLatestStatusChangeEvent(id uint) (*model.DeviceEventStatusChange, error) {
	e := model.DeviceEvent{}
	result := h.DB.Model(&model.DeviceEvent{}).Where(
		"device_id = ? and type = 'device_event_status_change'", id,
	).Order("updated_at DESC").Limit(1).Take(&e)
	if result.Error != nil {
		if errors.Is(gorm.ErrRecordNotFound, result.Error) {
			return nil, nil
		}
		return nil, result.Error
	}
	ec := model.DeviceEventStatusChange{}
	result = h.DB.Model(&model.DeviceEventStatusChange{}).Where("event_id = ?", e.ID).Limit(1).Take(&ec)
	if result.Error != nil {
		if errors.Is(gorm.ErrRecordNotFound, result.Error) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &ec, nil
}
