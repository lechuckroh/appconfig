package appconfig

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func createTempFile(t *testing.T, name, content string) (string, func()) {
	t.Helper()

	dir, err := ioutil.TempDir("", "confita")
	So(err, ShouldBeNil)

	path := filepath.Join(dir, name)
	f, err := os.Create(path)
	So(err, ShouldBeNil)

	_, _ = fmt.Fprintf(f, content)

	So(f.Close(), ShouldBeNil)

	return path, func() {
		So(os.RemoveAll(dir), ShouldBeNil)
	}
}

func TestLoadConfig(t *testing.T) {
	type config struct {
		Name   string   `config:"name"`
		Age    int      `config:"age"`
		Tags   []string `config:"tags"`
		Nested struct {
			Name string `config:"nested.name"`
			Age  int    `config:"nested.age"`
		}
	}

	Convey("yaml", t, func() {
		path, cleanup := createTempFile(t, "config.yaml", `
name: "some name"
age: 10
tags:
   - foo
   - bar
nested:
   name: "nested name"
   age: 20
`)
		defer cleanup()

		cfg := config{}
		_, err := LoadConfig(path, &cfg)

		So(err, ShouldBeNil)
		So(cfg.Name, ShouldEqual, "some name")
		So(cfg.Age, ShouldEqual, 10)
		So(cfg.Tags, ShouldResemble, []string{"foo", "bar"})
		So(cfg.Nested.Name, ShouldEqual, "nested name")
		So(cfg.Nested.Age, ShouldEqual, 20)
	})

	Convey("env", t, func() {
		_ = os.Setenv("NAME", "some name")
		_ = os.Setenv("AGE", "10")
		_ = os.Setenv("TAGS", "foo,bar")
		_ = os.Setenv("NESTED.NAME", "nested name")
		_ = os.Setenv("NESTED.AGE", "20")

		cfg := config{}
		_, err := LoadConfig("", &cfg)

		So(err, ShouldBeNil)
		So(cfg.Name, ShouldEqual, "some name")
		So(cfg.Age, ShouldEqual, 10)
		So(cfg.Tags, ShouldResemble, []string{"foo", "bar"})
		So(cfg.Nested.Name, ShouldEqual, "nested name")
		So(cfg.Nested.Age, ShouldEqual, 20)
	})

	Convey("precedence", t, func() {
		path, cleanup := createTempFile(t, "config.yaml", `
name: "some name"
age: 10
tags:
    - foo
    - bar
nested:
    name: "nested name"
    age: 20
`)
		defer cleanup()

		_ = os.Setenv("NAME", "override name")

		cfg := config{}
		_, err := LoadConfig(path, &cfg)

		So(err, ShouldBeNil)
		So(cfg.Name, ShouldEqual, "override name")
		So(cfg.Age, ShouldEqual, 10)
		So(cfg.Tags, ShouldResemble, []string{"foo", "bar"})
		So(cfg.Nested.Name, ShouldEqual, "nested name")
		So(cfg.Nested.Age, ShouldEqual, 20)
	})

}
