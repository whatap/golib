package fileutil

import (
	// "bytes"
	// "fmt"
	"os"
	"path/filepath"
	"testing"
	// "github.com/stretchr/testify/assert"
)

// ReadFile 함수에 대한 테스트를 수행하는 함수입니다.
func TestReadFile(t *testing.T) {
	// 테스트 케이스 구조체 정의
	testCases := []struct {
		name          string
		fileContent   string
		maxLength     int64
		expectContent string
		expectLength  int64
		expectError   bool
		setupFunc     func(string, string) error
	}{
		{
			name:          "ReadFullLength",
			fileContent:   "abcdefghijklmnopqrstuvwxyz",
			maxLength:     26,
			expectContent: "abcdefghijklmnopqrstuvwxyz",
			expectLength:  26,
			expectError:   false,
			setupFunc:     func(p, c string) error { return os.WriteFile(p, []byte(c), 0644) },
		},
		{
			name:          "ReadPartialLength",
			fileContent:   "abcdefghijklmnopqrstuvwxyz",
			maxLength:     10,
			expectContent: "abcdefghij",
			expectLength:  10,
			expectError:   false,
			setupFunc:     func(p, c string) error { return os.WriteFile(p, []byte(c), 0644) },
		},
		{
			name:          "ReadEmptyFile",
			fileContent:   "",
			maxLength:     10,
			expectContent: "",
			expectLength:  0,
			expectError:   false,
			setupFunc:     func(p, c string) error { return os.WriteFile(p, []byte(c), 0644) },
		},
		{
			name:          "FileDoesNotExist",
			fileContent:   "",
			maxLength:     10,
			expectContent: "",
			expectLength:  0,
			expectError:   true,
			setupFunc:     func(p, c string) error { return nil }, // 파일 생성하지 않음
		},
	}

	// 모든 테스트 케이스를 순회하며 테스트 실행
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 테스트를 위한 임시 파일 생성
			tempDir, err := os.MkdirTemp("", "test-readfile-*")
			if err != nil {
				t.Fatalf("Failed to create temp dir: %v", err)
			}
			defer os.RemoveAll(tempDir)

			filePath := filepath.Join(tempDir, "testfile.txt")

			// 테스트 환경 설정
			if err := tc.setupFunc(filePath, tc.fileContent); err != nil {
				t.Fatalf("Setup failed: %v", err)
			}

			// ReadFile 함수 호출
			actualBytes, actualLength, err := ReadFile(filePath, tc.maxLength)

			// 에러 검증
			if tc.expectError {
				if err == nil {
					t.Errorf("Expected an error but got none.")
				}
				// 파일이 존재하지 않는 경우 os.IsNotExist를 사용하여 정확한 오류 유형을 검증합니다.
				if !os.IsNotExist(err) {
					t.Errorf("Expected a 'file does not exist' error, but got: %v", err)
				}
			} else {
				if err != nil {
					t.Fatalf("Unexpected error: %v", err)
				}

				// 반환된 내용과 길이 검증
				if string(actualBytes) != tc.expectContent {
					t.Errorf("Content mismatch. Expected: %q, Got: %q", tc.expectContent, string(actualBytes))
				}
				if actualLength != tc.expectLength {
					t.Errorf("Length mismatch. Expected: %d, Got: %d", tc.expectLength, actualLength)
				}
			}
		})
	}
}

// 주어진 함수에 대한 테스트를 수행하는 함수입니다.
// Go 테스트 함수의 이름은 'Test'로 시작하고, *testing.T 타입을 인자로 받습니다.
func TestReplaceFileWithOldBackup(t *testing.T) {
	// 여러 테스트 케이스를 정의하는 구조체 슬라이스입니다.
	// 이 방식을 사용하면 코드를 반복하지 않고 여러 시나리오를 테스트할 수 있습니다.
	testCases := []struct {
		name           string                     // 테스트 케이스의 이름
		initialContent string                     // 테스트 시작 시 파일에 있을 내용
		newContent     string                     // 함수가 파일에 쓸 새로운 내용
		expectBackup   bool                       // 백업 파일이 생성될 것으로 예상하는지 여부
		setupFunc      func(string, string) error // 테스트를 위한 초기 파일 생성 함수
		cleanupFunc    func(string) error         // 테스트 후 파일 정리 함수
	}{
		{
			name:           "NoExistingFile",
			initialContent: "",
			newContent:     "This is the new content.",
			expectBackup:   false,
			// 기존 파일이 없는 시나리오를 설정합니다.
			setupFunc: func(p, c string) error { return nil },
			// 테스트 후 생성된 파일을 정리합니다.
			cleanupFunc: func(p string) error {
				return os.RemoveAll(filepath.Dir(p))
			},
		},
		{
			name:           "ExistingFile",
			initialContent: "This is the old content.",
			newContent:     "This is the new content.",
			expectBackup:   true,
			// 기존 파일이 있는 시나리오를 설정합니다.
			setupFunc: func(p, c string) error {
				return os.WriteFile(p, []byte(c), 0644)
			},
			// 테스트 후 생성된 파일과 백업 파일을 모두 정리합니다.
			cleanupFunc: func(p string) error {
				if err := os.Remove(p + ".old"); err != nil && !os.IsNotExist(err) {
					return err
				}
				return os.RemoveAll(filepath.Dir(p))
			},
		},
	}

	// 모든 테스트 케이스를 순회하며 테스트를 실행합니다.
	for _, tc := range testCases {
		// t.Run을 사용하면 각 테스트 케이스를 독립적으로 실행하고 결과를 명확하게 분리할 수 있습니다.
		t.Run(tc.name, func(t *testing.T) {
			// 테스트를 위해 임시 디렉토리를 생성합니다.
			// 이렇게 하면 실제 파일 시스템을 오염시키지 않고 테스트를 격리할 수 있습니다.
			tempDir, err := os.MkdirTemp("", "test-file-replace-*")
			if err != nil {
				t.Fatalf("Failed to create temp dir: %v", err)
			}
			defer os.RemoveAll(tempDir) // 테스트가 끝난 후 임시 디렉토리를 정리합니다.

			// 테스트할 파일 경로를 설정합니다.
			filePath := filepath.Join(tempDir, "testfile.txt")

			// 초기 설정 함수를 호출하여 테스트 시나리오를 준비합니다.
			if err := tc.setupFunc(filePath, tc.initialContent); err != nil {
				t.Fatalf("Setup failed: %v", err)
			}

			// 테스트할 함수를 호출합니다.
			err = ReplaceFileWithOldBackup(filePath, tc.newContent)
			if err != nil {
				t.Fatalf("Function returned an error: %v", err)
			}

			// --- 결과 검증 단계 ---

			// 1. 새 파일의 내용이 올바른지 확인합니다.
			actualContent, err := os.ReadFile(filePath)
			if err != nil {
				t.Fatalf("Failed to read new file: %v", err)
			}
			if string(actualContent) != tc.newContent {
				t.Errorf("New file content is incorrect. Expected: %q, Got: %q", tc.newContent, string(actualContent))
			}

			// 2. 백업 파일의 존재 여부와 내용이 올바른지 확인합니다.
			oldFilePath := filePath + ".old"
			_, err = os.Stat(oldFilePath)

			if tc.expectBackup {
				// 백업 파일이 있을 것으로 예상하는 경우
				if os.IsNotExist(err) {
					t.Error("Expected backup file but it was not found.")
				} else if err != nil {
					t.Errorf("Error checking for backup file: %v", err)
				} else {
					// 백업 파일이 있으면 내용을 확인합니다.
					backupContent, err := os.ReadFile(oldFilePath)
					if err != nil {
						t.Fatalf("Failed to read backup file: %v", err)
					}
					if string(backupContent) != tc.initialContent {
						t.Errorf("Backup content is incorrect. Expected: %q, Got: %q", tc.initialContent, string(backupContent))
					}
				}
			} else {
				// 백업 파일이 없을 것으로 예상하는 경우
				if !os.IsNotExist(err) {
					t.Errorf("Did not expect a backup file but it was found.")
				}
			}
		})
	}
}
