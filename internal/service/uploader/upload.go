/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package uploader

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/apache/incubator-answer/internal/base/reason"
	"github.com/apache/incubator-answer/internal/service/service_config"
	"github.com/apache/incubator-answer/internal/service/siteinfo_common"
	"github.com/apache/incubator-answer/pkg/checker"
	"github.com/apache/incubator-answer/pkg/dir"
	"github.com/apache/incubator-answer/pkg/uid"
	"github.com/apache/incubator-answer/plugin"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/disintegration/imaging"
	"github.com/gin-gonic/gin"
	exifremove "github.com/scottleedavis/go-exif-remove"
	"github.com/segmentfault/pacman/errors"
	"github.com/segmentfault/pacman/log"
)

type StorageConfig struct {
	Enable          bool   `json:"enable" mapstructure:"enable" yaml:"enable"`
	Bucket          string `json:"bucket" mapstructure:"bucket" yaml:"bucket"`
	BaseURL         string `json:"base_url" mapstructure:"base_url" yaml:"base_url"`
	Path            string `json:"path" mapstructure:"path" yaml:"path"`
	AccountID       string `json:"account_id" mapstructure:"account_id" yaml:"account_id"`
	AccessKeyID     string `json:"access_key_id" mapstructure:"access_key_id" yaml:"access_key_id"`
	AccessKeySecret string `json:"access_key_secret" mapstructure:"access_key_secret" yaml:"access_key_secret"`
}

const (
	avatarSubPath      = "avatar"
	avatarThumbSubPath = "avatar_thumb"
	postSubPath        = "post"
	brandingSubPath    = "branding"
)

var (
	subPathList = []string{
		avatarSubPath,
		avatarThumbSubPath,
		postSubPath,
		brandingSubPath,
	}
	supportedThumbFileExtMapping = map[string]imaging.Format{
		".jpg":  imaging.JPEG,
		".jpeg": imaging.JPEG,
		".png":  imaging.PNG,
		".gif":  imaging.GIF,
	}
)

type UploaderService interface {
	UploadAvatarFile(ctx *gin.Context) (url string, err error)
	UploadPostFile(ctx *gin.Context) (url string, err error)
	UploadBrandingFile(ctx *gin.Context) (url string, err error)
	AvatarThumbFile(ctx *gin.Context, fileName string, size int) (url string, err error)
}

// uploaderService uploader service
type uploaderService struct {
	serviceConfig   *service_config.ServiceConfig
	storageConfig   *StorageConfig
	siteInfoService siteinfo_common.SiteInfoCommonService
}

// NewUploaderService new upload service
func NewUploaderService(serviceConfig *service_config.ServiceConfig,
	storageConfig *StorageConfig,
	siteInfoService siteinfo_common.SiteInfoCommonService) UploaderService {
	for _, subPath := range subPathList {
		err := dir.CreateDirIfNotExist(filepath.Join(serviceConfig.UploadPath, subPath))
		if err != nil {
			panic(err)
		}
	}
	return &uploaderService{
		serviceConfig:   serviceConfig,
		storageConfig:   storageConfig,
		siteInfoService: siteInfoService,
	}
}

// UploadAvatarFile upload avatar file
func (us *uploaderService) UploadAvatarFile(ctx *gin.Context) (url string, err error) {
	url, err = us.tryToUploadByPlugin(ctx, plugin.UserAvatar)
	if err != nil {
		return "", err
	}
	if len(url) > 0 {
		return url, nil
	}

	// max size
	ctx.Request.Body = http.MaxBytesReader(ctx.Writer, ctx.Request.Body, 5*1024*1024)
	file, fileHeader, err := ctx.Request.FormFile("file")
	if err != nil {
		return "", errors.BadRequest(reason.RequestFormatError).WithError(err)
	}
	file.Close()
	fileExt := strings.ToLower(path.Ext(fileHeader.Filename))
	if _, ok := plugin.DefaultFileTypeCheckMapping[plugin.UserAvatar][fileExt]; !ok {
		return "", errors.BadRequest(reason.RequestFormatError).WithError(err)
	}

	newFilename := fmt.Sprintf("%s%s", uid.IDStr12(), fileExt)
	avatarFilePath := path.Join(avatarSubPath, newFilename)
	return us.uploadFile(ctx, fileHeader, avatarFilePath)
}

func (us *uploaderService) AvatarThumbFile(ctx *gin.Context, fileName string, size int) (url string, err error) {
	fileSuffix := path.Ext(fileName)
	if _, ok := supportedThumbFileExtMapping[fileSuffix]; !ok {
		// if file type is not supported, return original file
		return path.Join(us.serviceConfig.UploadPath, avatarSubPath, fileName), nil
	}
	if size > 1024 {
		size = 1024
	}

	thumbFileName := fmt.Sprintf("%d_%d@%s", size, size, fileName)
	thumbFilePath := fmt.Sprintf("%s/%s/%s", us.serviceConfig.UploadPath, avatarThumbSubPath, thumbFileName)
	avatarFile, err := os.ReadFile(thumbFilePath)
	if err == nil {
		return thumbFilePath, nil
	}
	filePath := fmt.Sprintf("%s/%s/%s", us.serviceConfig.UploadPath, avatarSubPath, fileName)
	avatarFile, err = os.ReadFile(filePath)
	if err != nil {
		return "", errors.InternalServer(reason.UnknownError).WithError(err).WithStack()
	}
	reader := bytes.NewReader(avatarFile)
	img, err := imaging.Decode(reader)
	if err != nil {
		return "", errors.InternalServer(reason.UnknownError).WithError(err).WithStack()
	}

	var buf bytes.Buffer
	newImage := imaging.Fill(img, size, size, imaging.Center, imaging.Linear)
	if err = imaging.Encode(&buf, newImage, supportedThumbFileExtMapping[fileSuffix]); err != nil {
		return "", errors.InternalServer(reason.UnknownError).WithError(err).WithStack()
	}

	if err = dir.CreateDirIfNotExist(path.Join(us.serviceConfig.UploadPath, avatarThumbSubPath)); err != nil {
		return "", errors.InternalServer(reason.UnknownError).WithError(err).WithStack()
	}

	avatarFilePath := path.Join(avatarThumbSubPath, thumbFileName)
	saveFilePath := path.Join(us.serviceConfig.UploadPath, avatarFilePath)
	out, err := os.Create(saveFilePath)
	if err != nil {
		return "", errors.InternalServer(reason.UnknownError).WithError(err).WithStack()
	}
	defer out.Close()

	thumbReader := bytes.NewReader(buf.Bytes())
	if _, err = io.Copy(out, thumbReader); err != nil {
		return "", errors.InternalServer(reason.UnknownError).WithError(err).WithStack()
	}
	return saveFilePath, nil
}

func (us *uploaderService) UploadPostFile(ctx *gin.Context) (
	url string, err error) {
	url, err = us.tryToUploadByPlugin(ctx, plugin.UserPost)
	if err != nil {
		return "", err
	}
	if len(url) > 0 {
		return url, nil
	}

	// max size
	ctx.Request.Body = http.MaxBytesReader(ctx.Writer, ctx.Request.Body, 512*1024*1024)
	file, fileHeader, err := ctx.Request.FormFile("file")
	if err != nil {
		return "", errors.BadRequest(reason.RequestFormatError).WithError(err)
	}
	defer file.Close()
	fileExt := strings.ToLower(path.Ext(fileHeader.Filename))
	if _, ok := plugin.DefaultFileTypeCheckMapping[plugin.UserPost][fileExt]; !ok {
		return "", errors.BadRequest(reason.RequestFormatError).WithError(err)
	}

	newFilename := fmt.Sprintf("%s%s", uid.IDStr12(), fileExt)
	postFilePath := path.Join(postSubPath, newFilename)

	if us.storageConfig.Enable {
		return us.uploadFileByStorage(ctx, fileHeader, postFilePath)
	}

	return us.uploadFile(ctx, fileHeader, postFilePath)
}

func (us *uploaderService) UploadBrandingFile(ctx *gin.Context) (
	url string, err error) {
	url, err = us.tryToUploadByPlugin(ctx, plugin.AdminBranding)
	if err != nil {
		return "", err
	}
	if len(url) > 0 {
		return url, nil
	}

	// max size
	ctx.Request.Body = http.MaxBytesReader(ctx.Writer, ctx.Request.Body, 10*1024*1024)
	file, fileHeader, err := ctx.Request.FormFile("file")
	if err != nil {
		return "", errors.BadRequest(reason.RequestFormatError).WithError(err)
	}
	file.Close()
	fileExt := strings.ToLower(path.Ext(fileHeader.Filename))
	if _, ok := plugin.DefaultFileTypeCheckMapping[plugin.AdminBranding][fileExt]; !ok {
		return "", errors.BadRequest(reason.RequestFormatError).WithError(err)
	}

	newFilename := fmt.Sprintf("%s%s", uid.IDStr12(), fileExt)
	avatarFilePath := path.Join(brandingSubPath, newFilename)
	return us.uploadFile(ctx, fileHeader, avatarFilePath)
}

func (us *uploaderService) uploadFile(ctx *gin.Context, file *multipart.FileHeader, fileSubPath string) (
	url string, err error) {
	siteGeneral, err := us.siteInfoService.GetSiteGeneral(ctx)
	if err != nil {
		return "", err
	}
	filePath := path.Join(us.serviceConfig.UploadPath, fileSubPath)
	if err := ctx.SaveUploadedFile(file, filePath); err != nil {
		return "", errors.InternalServer(reason.UnknownError).WithError(err).WithStack()
	}

	src, err := file.Open()
	if err != nil {
		return "", errors.InternalServer(reason.UnknownError).WithError(err).WithStack()
	}
	defer src.Close()

	if !checker.IsSupportedImageFile(filePath) {
		return "", errors.BadRequest(reason.UploadFileUnsupportedFileFormat)
	}

	if err := removeExif(filePath); err != nil {
		log.Error(err)
	}

	url = fmt.Sprintf("%s/uploads/%s", siteGeneral.SiteUrl, fileSubPath)
	return url, nil
}

func (us *uploaderService) uploadFileByStorage(ctx *gin.Context, file *multipart.FileHeader, fileSubPath string) (
	url string, err error) {

	fileName := fmt.Sprintf("%s/%s", us.storageConfig.Path, fileSubPath)
	go func() {
		session := us.newSession()
		client := s3manager.NewUploader(session)

		f, openError := file.Open()
		if openError != nil {
			log.Error("open file failed", openError)
		}
		defer f.Close() // 创建文件 defer 关闭

		input := &s3manager.UploadInput{
			Bucket: aws.String(us.storageConfig.Bucket),
			Key:    aws.String(fileName),
			Body:   f,
		}
		now := time.Now()
		_, err = client.Upload(input)
		if err != nil {
			log.Error("upload file to storage failed", err)
		}
		log.Info("spend time", time.Since(now))
	}()

	return fmt.Sprintf("%s/%s", us.storageConfig.BaseURL, fileName), nil
}

// TODO: delete file from storage
func (us *uploaderService) DeleteFile(key string) error {
	session := us.newSession()
	svc := s3.New(session)
	filename := fmt.Sprintf("%s/%s", us.storageConfig.Path, key)
	bucket := us.storageConfig.Bucket

	_, err := svc.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(filename),
	})
	if err != nil {
		log.Error("delete file from storage failed", err)
		return errors.InternalServer(reason.UnknownError).WithError(err).WithStack()
	}

	_ = svc.WaitUntilObjectNotExists(&s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(filename),
	})
	return nil
}

func (us *uploaderService) newSession() *session.Session {
	endpoint := fmt.Sprintf("%s.r2.cloudflarestorage.com", us.storageConfig.AccountID)

	return session.Must(session.NewSession(&aws.Config{
		Region:   aws.String("auto"),
		Endpoint: aws.String(endpoint),
		Credentials: credentials.NewStaticCredentials(
			us.storageConfig.AccessKeyID,
			us.storageConfig.AccessKeySecret,
			"",
		),
	}))
}

func (us *uploaderService) tryToUploadByPlugin(ctx *gin.Context, source plugin.UploadSource) (
	url string, err error) {
	_ = plugin.CallStorage(func(fn plugin.Storage) error {
		resp := fn.UploadFile(ctx, source)
		if resp.OriginalError != nil {
			log.Errorf("upload file by plugin failed, err: %v", resp.OriginalError)
			err = errors.BadRequest("").WithMsg(resp.DisplayErrorMsg.Translate(ctx)).WithError(err)
		} else {
			url = resp.FullURL
		}
		return nil
	})
	return url, err
}

// removeExif remove exif
// only support jpg/jpeg/png
func removeExif(path string) error {
	ext := strings.ToLower(strings.TrimPrefix(filepath.Ext(path), "."))
	if ext != "jpeg" && ext != "jpg" && ext != "png" {
		return nil
	}
	img, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	noExifBytes, err := exifremove.Remove(img)
	if err != nil {
		return err
	}
	return os.WriteFile(path, noExifBytes, 0644)
}
