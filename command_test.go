package command_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/gloo-foo/testable/assertion"
	"github.com/gloo-foo/testable/run"
	command "github.com/yupsh/basename"
)

// ==============================================================================
// Test Basic Functionality
// ==============================================================================

func TestBasename_SimplePath(t *testing.T) {
	result := run.Quick(command.Basename("/usr/local/bin/script.sh"))

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{"script.sh"})
}

func TestBasename_NoPath(t *testing.T) {
	result := run.Quick(command.Basename("script.sh"))

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{"script.sh"})
}

func TestBasename_WithSuffix(t *testing.T) {
	result := run.Quick(command.Basename("/usr/local/bin/script.sh", command.Suffix(".sh")))

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{"script"})
}

func TestBasename_SuffixNoMatch(t *testing.T) {
	result := run.Quick(command.Basename("/usr/local/bin/script.sh", command.Suffix(".txt")))

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{"script.sh"})
}

func TestBasename_Root(t *testing.T) {
	result := run.Quick(command.Basename("/"))

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{"/"})
}

func TestBasename_CurrentDir(t *testing.T) {
	result := run.Quick(command.Basename("."))

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{"."})
}

func TestBasename_ParentDir(t *testing.T) {
	result := run.Quick(command.Basename(".."))

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{".."})
}

// ==============================================================================
// Test Trailing Slashes
// ==============================================================================

func TestBasename_TrailingSlash(t *testing.T) {
	result := run.Quick(command.Basename("/path/to/dir/"))

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{"dir"})
}

func TestBasename_MultipleTrailingSlashes(t *testing.T) {
	result := run.Quick(command.Basename("/path/to/dir///"))

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{"dir"})
}

// ==============================================================================
// Test Multiple Paths
// ==============================================================================

func TestBasename_MultiplePaths(t *testing.T) {
	result := run.Quick(command.Basename(
		"/usr/bin/script.sh",
		"/home/user/doc.txt",
		"/tmp/temp",
	))

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{
		"script.sh",
		"doc.txt",
		"temp",
	})
}

func TestBasename_MultiplePathsWithSuffix(t *testing.T) {
	result := run.Quick(command.Basename(
		"/usr/bin/script.sh",
		"/home/user/doc.sh",
		command.Suffix(".sh"),
	))

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{
		"script",
		"doc",
	})
}

// ==============================================================================
// Test Different Path Styles
// ==============================================================================

func TestBasename_AbsolutePath(t *testing.T) {
	result := run.Quick(command.Basename("/usr/local/bin/myapp"))

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{"myapp"})
}

func TestBasename_RelativePath(t *testing.T) {
	result := run.Quick(command.Basename("relative/path/file.txt"))

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{"file.txt"})
}

func TestBasename_SingleLevel(t *testing.T) {
	result := run.Quick(command.Basename("filename"))

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{"filename"})
}

// ==============================================================================
// Test Special Characters
// ==============================================================================

func TestBasename_Spaces(t *testing.T) {
	result := run.Quick(command.Basename("/path/to/file with spaces.txt"))

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{"file with spaces.txt"})
}

func TestBasename_SpecialChars(t *testing.T) {
	result := run.Quick(command.Basename("/path/to/file-name_v2.txt"))

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{"file-name_v2.txt"})
}

func TestBasename_Dots(t *testing.T) {
	result := run.Quick(command.Basename("/path/to/file.tar.gz"))

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{"file.tar.gz"})
}

func TestBasename_DotsWithSuffix(t *testing.T) {
	result := run.Quick(command.Basename("/path/to/file.tar.gz", command.Suffix(".gz")))

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{"file.tar"})
}

// ==============================================================================
// Test Zero Terminator Flag
// ==============================================================================

func TestBasename_ZeroTerminator(t *testing.T) {
	result := run.Quick(command.Basename("file.txt", command.Zero))

	assertion.NoError(t, result.Err)
	// Should end with null byte instead of newline
	assertion.True(t, strings.HasSuffix(result.Stdout[0], "\x00"), "ends with null")
	assertion.True(t, strings.HasPrefix(result.Stdout[0], "file.txt"), "starts with filename")
}

func TestBasename_ZeroTerminatorMultiple(t *testing.T) {
	result := run.Quick(command.Basename("file1.txt", "file2.txt", command.Zero))

	assertion.NoError(t, result.Err)
	// Each output should be null-terminated
	fullOutput := strings.Join(result.Stdout, "")
	parts := strings.Split(fullOutput, "\x00")
	assertion.True(t, len(parts) >= 2, "at least 2 parts")
}

// ==============================================================================
// Test Suffix Variations
// ==============================================================================

func TestBasename_EmptySuffix(t *testing.T) {
	result := run.Quick(command.Basename("file.txt", command.Suffix("")))

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{"file.txt"})
}

func TestBasename_LongSuffix(t *testing.T) {
	result := run.Quick(command.Basename("document.backup.old", command.Suffix(".backup.old")))

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{"document"})
}

func TestBasename_SuffixLongerThanName(t *testing.T) {
	result := run.Quick(command.Basename("a.txt", command.Suffix(".txt.backup")))

	assertion.NoError(t, result.Err)
	// Suffix doesn't match, so name is unchanged
	assertion.Lines(t, result.Stdout, []string{"a.txt"})
}

// ==============================================================================
// Test Edge Cases
// ==============================================================================

func TestBasename_NoArguments(t *testing.T) {
	result := run.Quick(command.Basename())

	assertion.NoError(t, result.Err)
	// No output when no arguments
	assertion.Empty(t, result.Stdout)
}

// ==============================================================================
// Test Error Handling
// ==============================================================================

func TestBasename_OutputError(t *testing.T) {
	result := run.Command(command.Basename("test.txt")).
		WithStdoutError(errors.New("write failed")).
		Run()

	assertion.ErrorContains(t, result.Err, "write failed")
}

func TestBasename_OutputError_ZeroFlag(t *testing.T) {
	result := run.Command(command.Basename("test.txt", command.Zero)).
		WithStdoutError(errors.New("write failed")).
		Run()

	assertion.ErrorContains(t, result.Err, "write failed")
}

// ==============================================================================
// Test Multiple Flag
// ==============================================================================

func TestBasename_MultipleFlag(t *testing.T) {
	// Multiple flag is defined but not currently used in implementation
	result := run.Quick(command.Basename("file.txt", command.Multiple))

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{"file.txt"})
}

// ==============================================================================
// Table-Driven Tests
// ==============================================================================

func TestBasename_TableDriven(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		suffix   command.Suffix
		expected string
	}{
		{
			name:     "simple file",
			path:     "file.txt",
			suffix:   "",
			expected: "file.txt",
		},
		{
			name:     "with path",
			path:     "/usr/local/bin/app",
			suffix:   "",
			expected: "app",
		},
		{
			name:     "with suffix",
			path:     "document.pdf",
			suffix:   ".pdf",
			expected: "document",
		},
		{
			name:     "trailing slash",
			path:     "/path/to/dir/",
			suffix:   "",
			expected: "dir",
		},
		{
			name:     "root",
			path:     "/",
			suffix:   "",
			expected: "/",
		},
		{
			name:     "current dir",
			path:     ".",
			suffix:   "",
			expected: ".",
		},
		{
			name:     "relative",
			path:     "sub/dir/file",
			suffix:   "",
			expected: "file",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result *run.Result
			if tt.suffix != "" {
				result = run.Quick(command.Basename(tt.path, tt.suffix))
			} else {
				result = run.Quick(command.Basename(tt.path))
			}

			assertion.NoError(t, result.Err)
			assertion.Lines(t, result.Stdout, []string{tt.expected})
		})
	}
}

// ==============================================================================
// Test Real-World Scenarios
// ==============================================================================

func TestBasename_ScriptName(t *testing.T) {
	result := run.Quick(command.Basename("/usr/local/bin/deploy.sh"))

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{"deploy.sh"})
}

func TestBasename_LogFile(t *testing.T) {
	result := run.Quick(command.Basename("/var/log/app.log", command.Suffix(".log")))

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{"app"})
}

func TestBasename_BinaryName(t *testing.T) {
	result := run.Quick(command.Basename("/usr/bin/python3"))

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{"python3"})
}

func TestBasename_ConfigFile(t *testing.T) {
	result := run.Quick(command.Basename("/etc/nginx/nginx.conf"))

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{"nginx.conf"})
}

// ==============================================================================
// Test Edge Cases
// ==============================================================================

func TestBasename_EmptyPath(t *testing.T) {
	result := run.Quick(command.Basename(""))

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{"."})
}

func TestBasename_OnlySlashes(t *testing.T) {
	result := run.Quick(command.Basename("///"))

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{"/"})
}

func TestBasename_HiddenFile(t *testing.T) {
	result := run.Quick(command.Basename("/home/user/.bashrc"))

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{".bashrc"})
}

func TestBasename_HiddenFileWithSuffix(t *testing.T) {
	result := run.Quick(command.Basename("/home/user/.config.bak", command.Suffix(".bak")))

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{".config"})
}

func TestBasename_WindowsStylePath(t *testing.T) {
	// Even on Unix, basename should handle this reasonably
	result := run.Quick(command.Basename("C:\\Users\\name\\file.txt"))

	assertion.NoError(t, result.Err)
	// filepath.Base handles this platform-independently
	assertion.Count(t, result.Stdout, 1)
}

func TestBasename_DoubleExtension(t *testing.T) {
	result := run.Quick(command.Basename("archive.tar.gz", command.Suffix(".tar.gz")))

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{"archive"})
}

// ==============================================================================
// Test Unicode Paths
// ==============================================================================

func TestBasename_Unicode(t *testing.T) {
	result := run.Quick(command.Basename("/path/to/ファイル.txt"))

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{"ファイル.txt"})
}

func TestBasename_UnicodeWithSuffix(t *testing.T) {
	result := run.Quick(command.Basename("/path/文档.pdf", command.Suffix(".pdf")))

	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{"文档"})
}

