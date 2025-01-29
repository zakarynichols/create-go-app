package fsys

import (
	"embed"
	"errors"
	"io"
	"os"
	"testing"
)

//go:embed all:test
var mockEmbed embed.FS

type MockFileOps struct {
	ReadAllErr   error
	CreateErr    error
	WriteFileErr error
	MkdirErr     error
}

func (m MockFileOps) Create(name string) (*os.File, error) {
	if m.CreateErr != nil {
		return nil, m.CreateErr
	}
	return os.Create(name)
}

func (m MockFileOps) WriteFile(name string, data []byte, perm os.FileMode) error {
	if m.WriteFileErr != nil {
		return m.WriteFileErr
	}
	return os.WriteFile(name, data, perm)
}

func (m MockFileOps) ReadAll(r io.Reader) ([]byte, error) {
	if m.ReadAllErr != nil {
		return nil, m.ReadAllErr
	}
	return io.ReadAll(r)
}

func (m MockFileOps) Mkdir(name string, perm os.FileMode) error {
	if m.MkdirErr != nil {
		return m.MkdirErr
	}
	return os.Mkdir(name, perm)
}

type mockFSWrapper struct {
	fs embed.FS
}

func (w mockFSWrapper) Open(name string) (FileReaderCloser, error) {
	file, err := w.fs.Open(name)
	if err != nil {
		return nil, err
	}

	return file, nil
}

// TODO: Make tests idempotent by adding setup and teardown logic per test.
func TestWriteEmit(t *testing.T) {
	tests := []struct {
		appName string
		path    string
		title   string
		fileOps fileService
		wantErr bool
		isDir   bool
	}{
		{
			title:   "Open error",
			appName: "",
			path:    "",
			fileOps: MockFileOps{},
			wantErr: true,
		},
		{
			title:   "ReadAll error",
			appName: "test",
			path:    "test",
			fileOps: MockFileOps{
				ReadAllErr: errors.New("mock ReadAll error"),
			},
			wantErr: true,
		},
		{
			title:   "Create error",
			appName: "test",
			path:    "test/root.txt",
			fileOps: MockFileOps{
				CreateErr: errors.New("mock Create error"),
			},
			wantErr: true,
		},
		{
			title:   "WriteFile error",
			appName: "test",
			path:    "test/root.txt",
			fileOps: MockFileOps{
				WriteFileErr: errors.New("mock WriteFile error"),
			},
			wantErr: true,
		},
		{
			title:   "Mkdir error",
			appName: "test",
			path:    "test",
			fileOps: MockFileOps{
				MkdirErr: errors.New("mock Mkdir error"),
			},
			wantErr: true,
			isDir:   true,
		},
		{
			title:   "isExists error",
			appName: "test",
			path:    "test",
			fileOps: MockFileOps{
				MkdirErr: os.ErrExist,
			},
			wantErr: false,
			isDir:   true,
		},
		{
			title:   "Success",
			appName: "test",
			path:    "test/root.txt",
			fileOps: MockFileOps{},
			wantErr: false,
			isDir:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.title, func(t *testing.T) {
			err := Output(tt.appName, tt.path, tt.isDir, mockFSWrapper{fs: mockEmbed}, tt.fileOps)
			if (err != nil) != tt.wantErr {
				t.Errorf("got = %v, want = %v", err, tt.wantErr)
			}
		})
	}
}

type MockChangeImportsOpts struct {
	ReadFileErr  error
	WriteFileErr error
	ReadFileData string
}

func (m MockChangeImportsOpts) ReadFile(name string) ([]byte, error) {
	if m.ReadFileErr != nil {
		return nil, m.ReadFileErr
	}
	return os.ReadFile(name)
}

func (m MockChangeImportsOpts) WriteFile(name string, data []byte, perm os.FileMode) error {
	if m.WriteFileErr != nil {
		return m.WriteFileErr
	}
	return os.WriteFile(name, data, perm)
}

type MockFileDescriptor struct {
	isDir bool
	name  string
}

func (m MockFileDescriptor) IsDir() bool {
	return m.isDir
}

func (m MockFileDescriptor) Name() string {
	return m.name
}

func TestChangeImports(t *testing.T) {
	tests := []struct {
		title        string
		path         string
		fd           FileDescriptor
		prev         string
		new          string
		ops          MockChangeImportsOpts
		wantErr      bool
		expectedData string
		setup        func()
		teardown     func()
	}{
		{
			title:   "FileDescriptor is a directory",
			path:    "dir",
			fd:      MockFileDescriptor{isDir: true, name: "dir"},
			prev:    "old/import/path",
			new:     "new/import/path",
			ops:     MockChangeImportsOpts{},
			wantErr: false,
		},
		{
			title:   "File path does not match *.go",
			path:    "file.txt",
			fd:      MockFileDescriptor{isDir: false, name: "file.txt"},
			prev:    "old/import/path",
			new:     "new/import/path",
			ops:     MockChangeImportsOpts{},
			wantErr: false,
		},
		{
			title:   "ReadFile error",
			path:    "file.go",
			fd:      MockFileDescriptor{isDir: false, name: "file.go"},
			prev:    "old/import/path",
			new:     "new/import/path",
			ops:     MockChangeImportsOpts{ReadFileErr: errors.New("mock ReadFile error")},
			wantErr: true,
		},
		{
			title: "WriteFile error",
			path:  "file.go",
			fd:    MockFileDescriptor{isDir: false, name: "file.go"},
			prev:  "old/import/path",
			new:   "new/import/path",
			ops: MockChangeImportsOpts{
				ReadFileData: "package main\nimport \"old/import/path\"",
				WriteFileErr: errors.New("mock WriteFile error"),
			},
			wantErr: true,
			setup: func() {
				_, err := os.Create("file.go")
				if err != nil {
					t.Fatal(err)
				}
				err = os.WriteFile("file.go", []byte("package main\nimport \"old/import/path\""), os.FileMode(0777))
				if err != nil {
					t.Fatal(err)
				}
			},
			teardown: func() {
				err := os.Remove("file.go")
				if err != nil {
					t.Fatal(err)
				}
			},
		},
		{
			title:        "Successful import replacement",
			path:         "file.go",
			fd:           MockFileDescriptor{isDir: false, name: "file.go"},
			prev:         "old/import/path",
			new:          "new/import/path",
			ops:          MockChangeImportsOpts{},
			wantErr:      false,
			expectedData: "package main\nimport \"new/import/path\"",
			setup: func() {
				_, err := os.Create("file.go")
				if err != nil {
					t.Fatal(err)
				}
				err = os.WriteFile("file.go", []byte("package main\nimport \"old/import/path\""), os.FileMode(0777))
				if err != nil {
					t.Fatal(err)
				}
			},
			teardown: func() {
				err := os.Remove("file.go")
				if err != nil {
					t.Fatal(err)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.title, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup()
			}

			err := ReplaceImports("*.go", tt.path, tt.prev, tt.new, tt.fd, tt.ops)

			if (err != nil) != tt.wantErr {
				t.Errorf("got = %v, want = %v", err, tt.wantErr)
			}

			if tt.expectedData != "" && err == nil {
				readData, _ := tt.ops.ReadFile(tt.path)
				if string(readData) != tt.expectedData {
					t.Errorf("expected data = %v, got = %v", tt.expectedData, string(readData))
				}
			}

			if tt.teardown != nil {
				tt.teardown()
			}
		})
	}
}
