package txtartest

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"
	"testing"

	"github.com/grafana/codejen"
	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/envvars"
	"github.com/stretchr/testify/assert"
	"golang.org/x/tools/txtar"
)

// A TxTarTest represents a test run that processes all txtar-formatted files
// from a given directory. See the [Test] documentation for
// more details.
type TxTarTest struct {
	// Run TxTarTest on this directory.
	Root string

	// Name is a unique name for this test. The golden file for this test is
	// derived from the out/<name> file in the .txtar file.
	Name string

	// Skip is a map of tests to skip; the key is the test name; the value is the
	// skip message.
	Skip map[string]string

	// ToDo is a map of tests that should be skipped now, but should be fixed.
	ToDo map[string]string
}

// A Test represents a single test based on a .txtar file.
//
// A Test embeds *[testing.T] and should be used to report errors.
//
// Entries within the txtar file define a JSON-marshaled intermediate representation
// under the `ir.json` file. It can be accessed in tests with [Test.TypesIR].
//
// The rest of the .txtar file contains the results of running code generation jennies
// on that intermediate representation.
//
// These test cases (or "golden") files have names starting with "out/\(testname)". The "main" golden
// file is "out/\(testname)" itself, used when [Test] is used directly as an [io.Writer]
// and with [Test.WriteFile].
//
// When the test function has returned, output written with [Test.Write], [Test.Writer],
// [Test.WriteFile] and friends is checked against the expected output files.
//
// A txtar file can define test-specific tags and values in the comment section.
// These are available via the [Test.HasTag] and [Test.Value] methods.
// The #skip tag causes a [Test] to be skipped.
//
// If the output differs and $COG_UPDATE_GOLDEN is non-empty, the txtar file will be
// updated and written to disk with the actual output data replacing the
// out files.
type Test struct {
	// Allow Test to be used as a T.
	*testing.T

	prefix   string
	buf      *bytes.Buffer // the default buffer
	outFiles []file

	Archive *txtar.Archive

	// The absolute path of the current test directory.
	Dir string

	hasGold bool
}

// WriteFile writes a [codejen.File] to the main output,
// prefixed by a line of the form:
//
//	== name
//
// where name is the base name of f.RelativePath.
func (t *Test) WriteFile(f *codejen.File) {
	// TODO: use FileWriter instead in separate CL.
	fmt.Fprintln(t, "==", f.RelativePath)
	_, _ = t.Write(f.Data)
}

// WriteFiles writes a list of [codejen.File] to the main output.
func (t *Test) WriteFiles(files codejen.Files) {
	for i := range files {
		t.WriteFile(&files[i])
	}
}

// Write implements [io.Writer] by writing to the output for the test,
// which will be tested against the main golden file.
func (t *Test) Write(b []byte) (int, error) {
	if t.buf == nil {
		t.buf = &bytes.Buffer{}
		t.outFiles = append(t.outFiles, file{t.prefix, t.buf})
	}

	return t.buf.Write(b)
}

// HasTag reports whether the tag with the given key is defined
// for the current test. A tag x is defined by a line in the comment
// section of the txtar file like:
//
//	#x
func (t *Test) HasTag(key string) bool {
	prefix := []byte("#" + key)
	s := bufio.NewScanner(bytes.NewReader(t.Archive.Comment))
	for s.Scan() {
		b := s.Bytes()
		if bytes.Equal(bytes.TrimSpace(b), prefix) {
			return true
		}
	}

	return false
}

// Value returns the value for the given key for this test and
// reports whether it was defined.
//
// A value is defined by a line in the comment section of the txtar
// file like:
//
//	#key: value
//
// White space is trimmed from the value before returning.
func (t *Test) Value(key string) (string, bool) {
	prefix := []byte("#" + key + ":")
	s := bufio.NewScanner(bytes.NewReader(t.Archive.Comment))
	for s.Scan() {
		b := s.Bytes()
		if bytes.HasPrefix(b, prefix) {
			return string(bytes.TrimSpace(b[len(prefix):])), true
		}
	}

	return "", false
}

// Bool searches for a line starting with #key: value in the comment and
// reports whether the key exists and its value is true.
func (t *Test) Bool(key string) bool {
	s, ok := t.Value(key)

	return ok && s == "true"
}

// TypesIR locates and returns the raw types intermediate representation described
// within the txtar archive.
func (t *Test) TypesIR() *ast.Schema {
	for _, f := range t.Archive.Files {
		if f.Name != "ir.json" {
			continue
		}

		parsedIR := &ast.Schema{}
		if err := json.Unmarshal(f.Data, parsedIR); err != nil {
			t.Fatalf("could not load types IR: %s", err)
		}

		return parsedIR
	}

	// the ir.json file wasn't found, let's fail hard.
	t.Fatal("could not load types IR: file 'ir.json' not found in test archive")

	return nil
}

// Writer returns a Writer with the given name. Data written will
// be checked against the file with name "out/\(testName)/\(name)"
// in the txtar file. If name is empty, data will be written to the test
// output and checked against "out/\(testName)".
func (t *Test) Writer(name string) io.Writer {
	switch name {
	case "":
		name = t.prefix
	default:
		name = path.Join(t.prefix, name)
	}

	for _, f := range t.outFiles {
		if f.name == name {
			return f.buf
		}
	}

	w := &bytes.Buffer{}
	t.outFiles = append(t.outFiles, file{name, w})

	if name == t.prefix {
		t.buf = w
	}

	return w
}

// Run runs tests defined in txtar files in x.Root or its subdirectories.
//
// The function f is called for each such txtar file. See the [Test] documentation
// for more details.
func (x *TxTarTest) Run(t *testing.T, f func(tc *Test)) {
	t.Helper()

	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	err = filepath.WalkDir(x.Root, func(fullpath string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() || filepath.Ext(fullpath) != ".txtar" {
			return nil
		}

		str := filepath.ToSlash(fullpath)
		p := strings.Index(str, "/testdata/")
		testName := str[p+len("/testdata/") : len(str)-len(".txtar")]

		t.Run(testName, func(t *testing.T) {
			a, err := txtar.ParseFile(fullpath)
			if err != nil {
				t.Fatalf("error parsing txtar file: %v", err)
			}

			tc := &Test{
				T:       t,
				Archive: a,
				Dir:     filepath.Dir(filepath.Join(dir, fullpath)),
				prefix:  path.Join("out", x.Name),
			}

			if tc.HasTag("skip") {
				t.Skip()
			}
			if testing.Short() && tc.HasTag("slow") {
				tc.Skip("case is tagged #slow, skipping for -short")
			}

			if msg, ok := x.Skip[testName]; ok {
				t.Skip(msg)
			}
			if msg, ok := x.ToDo[testName]; ok {
				t.Skip(msg)
			}

			update := false
			for _, f := range a.Files {
				if strings.HasPrefix(f.Name, tc.prefix) && (f.Name == tc.prefix || f.Name[len(tc.prefix)] == '/') {
					// It's either "\(tc.prefix)" or "\(tc.prefix)/..." but not some other name
					// that happens to start with tc.prefix.
					tc.hasGold = true
				}
			}

			f(tc)

			// TODO we MAY need the below if trying to enable parallel tests
			//
			// Lock and re-parse the txtar file now that test execution is done. This does
			// make for some weird edge cases, but as long as underlying fs supports file
			// locking (windows? :scream:) it should make it safe to run multiple tests on same
			// txtar archive in parallel.
			// lock := flock.New(fullpath)
			// defer lock.Unlock()
			// a, err = txtar.ParseFile(fullpath)
			// if err != nil {
			// 	t.Fatalf("error parsing txtar file: %v", err)
			// }

			index := make(map[string]int, len(a.Files))
			for i, f := range a.Files {
				index[f.Name] = i
			}

			// Insert results of this test at first location of any existing
			// test or at end of list otherwise.
			k := len(a.Files)
			for _, sub := range tc.outFiles {
				if i, ok := index[sub.name]; ok {
					k = i

					break
				}
			}

			files := a.Files[:k:k]

			for _, sub := range tc.outFiles {
				result := sub.buf.Bytes()

				files = append(files, txtar.File{Name: sub.name})
				gold := &files[len(files)-1]

				if i, ok := index[sub.name]; ok {
					gold.Data = a.Files[i].Data
					delete(index, sub.name)

					if bytes.Equal(gold.Data, result) || bytes.Equal(bytes.TrimRight(gold.Data, "\n"), result) {
						continue
					}
				}

				if envvars.UpdateGoldenFiles {
					update = true
					gold.Data = result

					continue
				}

				assert.Equal(t, string(gold.Data), string(result), "result for %s differs", sub.name)
			}

			// Add remaining unrelated files, ignoring files that were already
			// added.
			for _, f := range a.Files[k:] {
				if _, ok := index[f.Name]; ok {
					files = append(files, f)
				}
			}
			a.Files = files

			if update {
				err = os.WriteFile(fullpath, txtar.Format(a), 0600)
				if err != nil {
					t.Fatal(err)
				}
			}
		})

		return nil
	})

	if err != nil {
		t.Fatal(err)
	}
}

func DumpTestInfo(tc *Test) {
	fmt.Println("=== TEST:", tc.Dir, tc.prefix)
	fmt.Println("=== Files")
	for _, f := range tc.Archive.Files {
		fmt.Printf("===   %s\n", f.Name)
	}
	fmt.Println("=== END TEST")
}

type file struct {
	name string
	buf  *bytes.Buffer
}
