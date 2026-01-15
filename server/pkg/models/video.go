package models

import (
	"github.com/Mahaveer86619/FrameSense/pkg/consts"
	"gorm.io/gorm"
)

type Video struct {
	gorm.Model

	Title             string
	Description       string
	SourceFilePath    string
	ProcessedFilePath string
	ErrorMessage      string

	HLSPlaylist       string
	HLSMasterPlaylist string
	Duration          string

	OwnerUserId uint
	OwnerUser   User

	Status consts.VideoStatus
}
