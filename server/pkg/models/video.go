package models

import (
	"github.com/Mahaveer86619/FrameSense/pkg/consts"
	"gorm.io/gorm"
)

type Video struct {
	gorm.Model

	// Basic Info
	Title       string `gorm:"not null"`
	Description string

	// File Paths
	SourceFilePath    string // Original uploaded video
	ProcessedFilePath string // Path to HLS master playlist

	// HLS Specific Fields
	HLSMasterPlaylist string // s3://bucket/hls/video_1/master.m3u8
	HLSDuration       int    // Duration in seconds
	HLSResolutions    string // JSON array of available resolutions, e.g., ["720p", "480p", "360p"]

	// Metadata
	Duration     string // Human-readable duration, e.g., "5:32"
	FileSize     int64  // Original file size in bytes
	ThumbnailURL string // Path to thumbnail image

	// Error Handling
	ErrorMessage string

	// Ownership
	OwnerUserId uint `gorm:"not null;index"`
	OwnerUser   User

	// Status
	Status consts.VideoStatus `gorm:"not null;index"`

	// Analytics
	ViewCount      int `gorm:"default:0"`
	LastStreamedAt *gorm.DeletedAt
}

func (v *Video) IsStreamable() bool {
	return v.Status == consts.VideoStatusReady && v.ProcessedFilePath != ""
}

func (v *Video) GetHLSPath() string {
	if v.HLSMasterPlaylist != "" {
		return v.HLSMasterPlaylist
	}
	return v.ProcessedFilePath
}
