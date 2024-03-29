// Code generated by go-bindata.
// sources:
// openconfigbeat.yml
// DO NOT EDIT!

package openconfigbeat

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func bindataRead(data []byte, name string) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	clErr := gz.Close()

	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}
	if clErr != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type asset struct {
	bytes []byte
	info  os.FileInfo
}

type bindataFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
}

func (fi bindataFileInfo) Name() string {
	return fi.name
}
func (fi bindataFileInfo) Size() int64 {
	return fi.size
}
func (fi bindataFileInfo) Mode() os.FileMode {
	return fi.mode
}
func (fi bindataFileInfo) ModTime() time.Time {
	return fi.modTime
}
func (fi bindataFileInfo) IsDir() bool {
	return false
}
func (fi bindataFileInfo) Sys() interface{} {
	return nil
}

var _openconfigbeatYml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x94\x57\x6d\x6f\xe3\x36\x12\xfe\xae\x5f\x31\x90\xbe\xb4\x07\x47\xf6\x2e\xb6\x8b\x3b\xa1\x5b\x60\x9b\xeb\x5d\x83\x1a\xe8\xe2\x92\x02\x77\x58\x2c\x6a\x9a\x1c\x49\x44\x68\x52\xc7\xa1\xec\xfa\xdf\x17\x43\x52\xb2\x9d\x97\x26\x2b\x20\x80\x23\xce\x3c\x33\xf3\x70\xde\x54\x3d\x7e\xe0\xd7\x01\xad\x74\xb6\xd5\xdd\x16\x45\x80\xeb\xf8\x73\xf4\x22\x68\x67\xe1\xa7\x3f\xc4\x6e\x30\x08\x4f\x28\xa6\xa7\x28\x9e\x3d\x7a\x0a\xfd\x2f\x85\xcf\x40\xdd\x85\x5a\x53\x14\x00\x15\xdc\xf5\x08\x42\x29\x8f\x44\x48\xe0\x5a\x08\x3d\x46\x03\xc9\x67\x50\xb8\xd7\x12\x09\x82\x03\xe9\xac\x45\x19\x20\xb8\xba\x80\x93\x52\x03\x9f\x4b\xe3\xa4\x30\xbd\xa3\x50\x7e\x39\xc1\x9e\xa1\x0c\x22\xf4\x11\x83\xc6\x2d\x49\xaf\xb7\x98\x51\xe2\x01\x23\x2c\xcf\x35\x15\xb6\x62\x34\x01\x06\xe7\xc3\xa5\x65\xd0\x2d\x58\x67\x11\x34\x81\xcc\xac\xa2\x62\xa4\xac\xf3\x3b\xeb\x34\xf0\x7e\xf5\xee\xed\x09\x6f\x24\xf4\x56\xec\x10\x9c\x8d\xf1\xa5\xa0\x58\xab\x9a\x8e\x1a\x28\x85\xda\x69\x5b\x9e\xb4\x06\x41\x74\x70\x5e\x41\xeb\x7c\x54\x63\xd9\x27\x20\x26\xb9\x06\xca\xe9\x67\x59\x14\x45\xf5\xe1\x85\x07\xfe\x8d\x16\xbd\x30\xf0\x92\x60\x7c\x8a\x22\x79\x95\xe2\x48\xf7\x44\xbd\x1e\x06\x64\xe7\x44\x80\x61\xdc\x1a\x4d\x3d\xdf\x15\x8b\x61\x38\x38\x7f\x0f\x4a\x04\x51\xc3\x4d\x00\x29\x2c\x6c\x63\x0c\x8a\x79\xec\xbc\x1b\x87\xa2\x02\x61\x4c\x94\x0f\x5e\x58\x12\x92\x13\x94\x80\xd0\x06\xd8\x1e\x41\x00\x69\xdb\x99\x93\x21\x9d\x82\x3f\xe0\x16\xb4\x0d\xe8\x5b\xc1\x1c\x54\x91\xc0\xc9\xc1\x20\x3a\x7a\xe8\xa0\xf0\x08\xda\x4a\x33\x2a\x54\x19\x44\x7b\x70\x07\x0b\xad\x46\xa3\xe0\xa0\x43\x0f\x28\x64\x5f\x54\xe7\x9e\xcc\x31\xa9\xba\xa8\x18\x97\x33\x85\xd0\x33\xf5\x57\xff\x2d\x17\x50\x1e\x70\x7b\x15\x34\x7a\xce\x1e\x2e\x0c\xd6\x12\x26\xa1\x52\xe2\xe5\xe8\xc6\x18\x3c\x0d\x28\x75\x7b\xe4\xe0\x85\x52\xfc\xa7\xb3\xb4\xb6\xad\xf3\xbb\x54\x9c\xc1\xb1\x73\x45\x05\x6e\x0c\xc3\x18\xea\xa2\x4a\x58\x4d\x51\x01\xa0\xdd\x37\x40\x41\x74\xda\x76\x2f\xdf\x30\xfc\x53\x50\xbf\x75\xc2\x2b\x7a\xe5\x15\x47\x02\x09\x81\x30\x04\x6d\xbb\x98\xe2\xc1\x3b\x03\xc6\x09\xa5\x6d\x97\x38\x4d\xdd\x43\x9d\xc0\x93\xcf\xf0\x8b\xde\x0a\x2b\x40\x5b\x85\x7f\xd4\xb0\x4e\x2a\x4c\x68\x7f\x21\xad\x09\x94\x26\xb1\x35\xa8\xf8\x8e\xa7\x5a\x13\x56\x4d\x39\x82\x36\x9d\xa2\x0e\x3d\x7a\x16\xca\x0e\x4d\xd4\x0c\x29\x4d\x7a\xf4\xb8\x00\x17\x25\x46\x9a\xfc\xdb\x5c\x11\x86\x71\xd8\xc0\xf5\xfa\x06\x5a\x23\x3a\xc8\xe5\xb3\xc9\xef\xa5\xdb\xed\x84\xe5\x2b\x8d\x2f\xea\x93\x6f\x75\xb6\xdc\x40\x2b\x0c\xe1\x94\x51\xbf\xfd\x67\x0d\xad\x77\x3b\x38\xb0\x45\x0e\x57\xb9\x83\x65\x4e\x1e\xc6\x26\xbc\xec\xf5\x1e\x6b\xf8\xf1\x14\x58\xe8\x35\x31\x44\x51\x41\x2f\x08\x04\xec\x85\x19\x11\x0e\xbd\x96\x7d\x6a\x23\xbb\x61\x0c\x4c\x86\xe0\xda\xc8\xf5\xfd\x23\x77\xd6\x58\x6b\x4c\xcc\x1e\x3d\x69\x67\x6b\xf8\x97\xf3\xe0\xd1\x20\x8b\x16\xd5\xf4\x9e\x16\xb3\x15\x18\x9c\xb6\x61\xbe\x93\xd9\xb7\xc9\xb5\x09\x5f\xf8\xa0\x5b\x21\x03\xd5\x68\x04\x05\x2d\x6b\xe9\x8a\x8a\x2b\x8b\x74\xc0\xa7\xb8\x19\xbd\x69\x5e\xce\xb9\x9c\x04\xaf\x6e\x29\xb7\x81\x3d\xb1\x5d\x2a\x41\x8e\x9a\xa6\xa8\xe0\x7d\xbd\xaa\x57\x8b\xc7\x1c\x63\xcc\x47\x54\xb0\xd7\xe2\x3c\xf3\x3e\x7e\xba\xa9\xe3\x95\x69\x02\x8f\xff\x1f\xb5\x47\xe6\x3b\x9f\xa2\x55\x91\x9b\xb9\x6f\xc7\x82\xab\x8b\x14\xe8\x7d\x14\xca\x43\x29\x6b\xfc\xec\x28\xc4\xff\x6f\x65\x8f\xf9\x26\xe2\x48\xc8\x79\x6a\xb0\x0d\x5c\xa7\xf1\xe0\xa0\x8d\xe1\x97\x84\x61\x26\x3f\x67\xc0\x37\x7d\x08\x43\x14\xfa\xee\xfd\xea\xcd\xb7\x11\xf2\xc6\x82\x14\x84\xb1\x3b\x4c\x9d\x81\x25\xce\xda\x02\xcf\xa6\x14\x3d\x25\xfb\xa7\xb0\x54\x03\x8c\xd9\x2c\x97\xf3\xe4\x6b\x18\x7b\xc9\x3a\x09\xff\xd3\xfe\xfd\xd9\x60\xa5\xde\x8d\x46\x81\x30\x07\x71\x24\x76\x53\x61\xab\x2d\x2a\x10\x94\xa0\xa8\x59\x2e\x3f\xbf\x5d\xad\xde\x34\x6a\xfb\xf7\xa6\x79\xf3\x25\xe2\x31\x54\x04\x87\xf2\xd2\x50\xf9\x52\x26\xc0\x4f\x29\xad\xe0\xda\xb8\x51\xbd\x22\x1f\x8a\xc7\xdd\x87\xf4\x6e\x30\x4c\x4c\x2a\xee\xcb\xfd\x21\x25\x0c\xb3\x73\x69\xe9\x9b\x29\x1a\xc9\xff\x9e\x65\xf7\xf2\xdb\x7a\x2a\xe8\x74\xa4\xd5\xdc\x57\xdc\x1e\xfd\xc1\xeb\x90\x87\xd7\x26\x37\xdf\xac\x4c\xc8\xc5\x53\x73\xf0\xb4\xe1\x5b\x2a\xaa\xdc\x4c\x72\xda\xc4\xa3\xcd\xd4\x97\x38\x07\xff\x97\x9b\x7e\xab\x6d\x6a\x12\x9b\xc9\xe6\x66\x1a\x62\x97\x6e\xf3\x48\xfb\x8d\xd3\x77\x92\x6b\x2e\x9d\x15\x63\xe8\xbf\xce\xdd\x69\xb3\x98\x3d\x7e\x52\x6a\xda\x18\x36\x33\xed\x75\xb4\x9a\x26\x12\xa7\xdc\xe6\x7b\x46\xfa\xa1\xf9\x9e\x45\x7f\xd8\xcc\x2e\xb2\x47\x2f\x37\x84\x0f\x1f\xe0\xd7\x68\xf7\xb5\x33\xa8\xa8\xe6\x75\x95\x9b\xa4\x08\x79\x12\x72\x51\x8d\xc4\xaf\xd0\xf2\x82\x30\xcf\x23\xde\x30\x40\x3a\x63\x50\x86\x34\x52\xf8\x2d\xa7\x08\xdf\xf7\xd5\xb3\xcf\x74\x01\x89\x88\xc9\xca\xf3\xf2\xfc\x14\x4f\x71\xd8\xc4\x82\xfb\xe8\xbd\x38\xf2\xde\x11\xd3\xe4\xf1\xbe\x1a\x5f\x5f\xec\xaa\xcd\x3f\xde\xae\x56\xd3\xda\x39\xaf\x0e\x83\x77\xc1\x49\x67\x62\x37\xd8\x0a\xd2\x12\xe2\xdd\x4b\x8f\x0a\x6d\xd0\xc2\x50\xda\xfa\xb2\x5c\x03\x65\xcc\xf8\xf2\xc1\x36\x99\x5d\x2c\x1f\x6c\x88\xb2\x17\xb6\xc3\x1d\x96\x7f\xc9\x0d\xd3\xb3\x76\x1d\x05\x41\xaf\x64\xe6\xea\xaa\xa8\x32\x39\x26\x2b\x36\xf3\x26\x3b\x43\x45\x12\xa6\xa6\xf2\x80\x8d\xef\x56\xef\xde\x3d\x62\xe3\xf6\x76\x7d\x31\x4c\x35\xaf\x76\x6d\x24\x00\xd6\x9a\x02\x13\xee\x9d\x0b\x20\x91\x07\x9a\x96\x82\x8b\x82\x17\xe6\x9f\xef\xee\x3e\xdd\x02\x2f\x6b\xe8\x79\xa2\xa4\x43\xae\x4f\x56\x26\x32\xf5\x99\xca\xef\x4c\xb1\xf3\x3a\xe8\xf4\x3d\xb1\xc4\x20\x97\xc3\xbd\x5e\x32\xf6\x52\x8a\x7a\xc0\xdd\xe4\xdb\xf5\x49\x2d\x1a\xba\xbd\x5d\x83\x34\x9a\x57\x56\x46\xe1\x3b\x4a\x86\x9e\xb0\xd3\xc0\x09\x3a\xe9\x2c\xf9\x34\xc2\x67\xf4\x84\x74\x6e\xe4\x17\x3c\x4e\x50\xf7\x78\x7c\x0e\xe2\x1e\x8f\x2f\x36\x65\xae\xc6\xb5\xeb\x78\x7d\xfc\x8a\x09\x8d\x81\xc0\xb8\x0e\x0c\xee\xd1\xd4\x17\x1f\x48\xf3\x6b\xbe\x18\x5e\x63\xb9\xf7\x7d\xdc\x0b\x6d\x78\x8f\x3a\x1d\xc7\xa9\xdd\x00\x7a\xef\xfc\x02\x0e\xc2\x5b\x6d\xbb\x45\xd4\x58\x80\xc2\xed\xd8\x15\x95\x49\x8e\xd5\x51\xa1\xc9\x6f\x19\x2e\xa4\xdf\x09\x69\x71\xda\xa7\x91\x2b\x5e\xef\xd1\x1c\xf3\xc6\x08\x19\x02\x9c\x35\xc7\x78\x37\xe4\x76\x18\x17\x2c\x67\xd1\x86\xd8\x99\xef\xdc\x24\xcd\x9f\x1e\x09\xc4\x79\x8a\xbd\xe5\x73\xf9\xb7\xf2\x4b\x3d\x7d\x20\xc7\xaf\x08\x17\x57\xd0\x93\x18\x6f\x1f\x25\x37\x97\x72\x51\x54\x50\xe6\x0f\x04\xfe\x14\xc8\xdf\x05\x65\x7d\x0a\x65\x56\x6b\x12\x74\xf1\x67\x00\x00\x00\xff\xff\x61\x39\xc9\x06\xac\x0f\x00\x00")

func openconfigbeatYmlBytes() ([]byte, error) {
	return bindataRead(
		_openconfigbeatYml,
		"openconfigbeat.yml",
	)
}

func openconfigbeatYml() (*asset, error) {
	bytes, err := openconfigbeatYmlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "openconfigbeat.yml", size: 4012, mode: os.FileMode(420), modTime: time.Unix(1515793919, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("Asset %s can't read by error: %v", name, err)
		}
		return a.bytes, nil
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// MustAsset is like Asset but panics when Asset would return an error.
// It simplifies safe initialization of global variables.
func MustAsset(name string) []byte {
	a, err := Asset(name)
	if err != nil {
		panic("asset: Asset(" + name + "): " + err.Error())
	}

	return a
}

// AssetInfo loads and returns the asset info for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func AssetInfo(name string) (os.FileInfo, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("AssetInfo %s can't read by error: %v", name, err)
		}
		return a.info, nil
	}
	return nil, fmt.Errorf("AssetInfo %s not found", name)
}

// AssetNames returns the names of the assets.
func AssetNames() []string {
	names := make([]string, 0, len(_bindata))
	for name := range _bindata {
		names = append(names, name)
	}
	return names
}

// _bindata is a table, holding each asset generator, mapped to its name.
var _bindata = map[string]func() (*asset, error){
	"openconfigbeat.yml": openconfigbeatYml,
}

// AssetDir returns the file names below a certain
// directory embedded in the file by go-bindata.
// For example if you run go-bindata on data/... and data contains the
// following hierarchy:
//     data/
//       foo.txt
//       img/
//         a.png
//         b.png
// then AssetDir("data") would return []string{"foo.txt", "img"}
// AssetDir("data/img") would return []string{"a.png", "b.png"}
// AssetDir("foo.txt") and AssetDir("notexist") would return an error
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		cannonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(cannonicalName, "/")
		for _, p := range pathList {
			node = node.Children[p]
			if node == nil {
				return nil, fmt.Errorf("Asset %s not found", name)
			}
		}
	}
	if node.Func != nil {
		return nil, fmt.Errorf("Asset %s not found", name)
	}
	rv := make([]string, 0, len(node.Children))
	for childName := range node.Children {
		rv = append(rv, childName)
	}
	return rv, nil
}

type bintree struct {
	Func     func() (*asset, error)
	Children map[string]*bintree
}

var _bintree = &bintree{nil, map[string]*bintree{
	"openconfigbeat.yml": &bintree{openconfigbeatYml, map[string]*bintree{}},
}}

// RestoreAsset restores an asset under the given directory
func RestoreAsset(dir, name string) error {
	data, err := Asset(name)
	if err != nil {
		return err
	}
	info, err := AssetInfo(name)
	if err != nil {
		return err
	}
	err = os.MkdirAll(_filePath(dir, filepath.Dir(name)), os.FileMode(0755))
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(_filePath(dir, name), data, info.Mode())
	if err != nil {
		return err
	}
	err = os.Chtimes(_filePath(dir, name), info.ModTime(), info.ModTime())
	if err != nil {
		return err
	}
	return nil
}

// RestoreAssets restores an asset under the given directory recursively
func RestoreAssets(dir, name string) error {
	children, err := AssetDir(name)
	// File
	if err != nil {
		return RestoreAsset(dir, name)
	}
	// Dir
	for _, child := range children {
		err = RestoreAssets(dir, filepath.Join(name, child))
		if err != nil {
			return err
		}
	}
	return nil
}

func _filePath(dir, name string) string {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	return filepath.Join(append([]string{dir}, strings.Split(cannonicalName, "/")...)...)
}
