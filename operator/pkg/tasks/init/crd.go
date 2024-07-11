/*
Copyright 2023 The Karmada Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package tasks

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	operatorv1alpha1 "github.com/karmada-io/karmada/operator/pkg/apis/operator/v1alpha1"
	"os"
	"path"
	"strings"

	"k8s.io/klog/v2"

	"github.com/karmada-io/karmada/operator/pkg/util"
	"github.com/karmada-io/karmada/operator/pkg/workflow"
)

var (
	crdsFileSuffix = "crds.tar.gz"
	crdPathSuffix  = "crds"
)

// NewPrepareCrdsTask init a prepare-crds task
func NewPrepareCrdsTask() workflow.Task {
	return workflow.Task{
		Name:        "prepare-crds",
		Run:         runPrepareCrds,
		RunSubTasks: true,
		Tasks: []workflow.Task{
			{
				Name: "download-crds",
				Skip: skipCrdsDownload,
				Run:  runCrdsDownload,
			},
			{
				Name: "Unpack",
				Run:  runUnpack,
			},
		},
	}
}

func runPrepareCrds(r workflow.RunData) error {
	data, ok := r.(InitData)
	if !ok {
		return errors.New("prepare-crds task invoked with an invalid data struct")
	}

	crdsDir := getCrdsDir(data)
	klog.V(4).InfoS("[prepare-crds] Running prepare-crds task", "karmada", klog.KObj(data))
	klog.V(2).InfoS("[prepare-crds] Using crd folder", "folder", crdsDir, "karmada", klog.KObj(data))

	return nil
}

func skipCrdsDownload(r workflow.RunData) (bool, error) {
	data, ok := r.(InitData)
	if !ok {
		return false, errors.New("prepare-crds task invoked with an invalid data struct")
	}

	if data.CrdDownloadPolicy() == operatorv1alpha1.DownloadAlways {
		klog.V(2).InfoS("[skipCrdsDownload] CrdDownloadPolicy is 'Always', skipping download check")
		return false, nil
	}

	crdsDir := getCrdsDir(data)
	klog.V(2).InfoS("[skipCrdsDownload] Checking if CRDs need to be downloaded", "folder", crdsDir)

	if exist, err := util.PathExists(crdsDir); !exist || err != nil {
		klog.V(2).InfoS("[skipCrdsDownload] CRDs directory does not exist or an error occurred", "folder", crdsDir, "error", err)
		return false, err
	}

	if !existCrdsTar(crdsDir) {
		klog.V(2).InfoS("[skipCrdsDownload] CRD tar file does not exist", "folder", crdsDir)
		return false, nil
	}

	klog.V(2).InfoS("[download-crds] Skip download CRD yaml files, the CRD tar exists on disk", "karmada", klog.KObj(data), "folder", crdsDir)
	return true, nil
}

func runCrdsDownload(r workflow.RunData) error {
	data, ok := r.(InitData)
	if !ok {
		return errors.New("download-crds task invoked with an invalid data struct")
	}

	crdsDir := getCrdsDir(data)
	crdsTarPath := path.Join(crdsDir, crdsFileSuffix)
	klog.V(2).InfoS("[runCrdsDownload] Starting CRDs download", "folder", crdsDir, "remoteURL", data.CrdsRemoteURL())

	// Check if the CRDs directory exists
	exist, err := util.PathExists(crdsDir)
	if err != nil {
		return err
	}

	// If the CRDs directory exists, delete and recreate it
	if exist {
		klog.V(2).InfoS("[runCrdsDownload] CRDs directory exists, deleting and recreating it", "folder", crdsDir)
		if err := os.RemoveAll(crdsDir); err != nil {
			return fmt.Errorf("failed to delete CRDs directory, err: %w", err)
		}
	}

	// Create the CRDs directory
	klog.V(2).InfoS("[runCrdsDownload] Creating CRDs directory", "folder", crdsDir)
	if err := os.MkdirAll(crdsDir, 0700); err != nil {
		return fmt.Errorf("failed to create CRDs directory, err: %w", err)
	}

	// Download the CRD tar file
	klog.V(2).InfoS("[runCrdsDownload] Downloading CRD tar file", "remoteURL", data.CrdsRemoteURL(), "tarPath", crdsTarPath)
	if err := util.DownloadFile(data.CrdsRemoteURL(), crdsTarPath); err != nil {
		return fmt.Errorf("failed to download CRD tar, err: %w", err)
	}

	klog.V(2).InfoS("[runCrdsDownload] Successfully downloaded CRD package from remote URL", "remoteURL", data.CrdsRemoteURL(), "folder", crdsDir)
	return nil
}

func runUnpack(r workflow.RunData) error {
	data, ok := r.(InitData)
	if !ok {
		return errors.New("unpack task invoked with an invalid data struct")
	}

	crdsDir := getCrdsDir(data)
	crdsTarPath := path.Join(crdsDir, crdsFileSuffix)
	crdsPath := path.Join(crdsDir, crdPathSuffix)
	klog.V(2).InfoS("[runUnpack] Starting to unpack CRDs", "tarPath", crdsTarPath, "unpackDir", crdsDir)

	exist, _ := util.PathExists(crdsPath)
	if !exist {
		klog.V(2).InfoS("[runUnpack] CRD yaml files do not exist, unpacking tar file", "unpackDir", crdsDir)
		if err := util.Unpack(crdsTarPath, crdsDir); err != nil {
			return fmt.Errorf("[unpack] failed to unpack CRD tar, err: %w", err)
		}
	} else {
		klog.V(2).InfoS("[unpack] These CRDs yaml files have been decompressed in the path", "path", crdsPath, "karmada", klog.KObj(data))
	}

	klog.V(2).InfoS("[unpack] Successfully unpacked CRD tar", "karmada", klog.KObj(data), "unpackDir", crdsDir)
	return nil
}

func existCrdsTar(crdsDir string) bool {
	files := util.ListFiles(crdsDir)
	klog.V(2).InfoS("[existCrdsTar] Checking for CRD tar file in directory", "directory", crdsDir)

	for _, file := range files {
		klog.V(2).InfoS("[existCrdsTar] Checking file", "fileName", file.Name(), "fileSize", file.Size())
		if strings.Contains(file.Name(), crdsFileSuffix) && file.Size() > 0 {
			klog.V(2).InfoS("[existCrdsTar] Found CRD tar file", "fileName", file.Name(), "fileSize", file.Size())
			return true
		}
	}
	return false
}

func getCrdsDir(data InitData) string {
	url := strings.TrimSpace(data.CrdsRemoteURL())
	hash := sha256.Sum256([]byte(url))
	hashStr := hex.EncodeToString(hash[:])
	return path.Join(data.DataDir(), "cache", hashStr)
}
