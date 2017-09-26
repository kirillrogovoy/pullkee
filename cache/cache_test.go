package cache

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWrite(t *testing.T) {
	t.Run("Works when successfully JSON-ed and wrote to the file", func(t *testing.T) {
		m := mockFS{cache: map[string]cacheEntity{}}
		c := FSCache{
			FS:        &m,
			CachePath: "/tmp/",
		}
		err := c.Set("key1", testStruct{"val1"})

		require.Nil(t, err)
		require.Equal(t, `{"x":"val1"}`, string(m.cache["/tmp/key1.json"].data))
		require.Contains(t, m.mkdirs, "/tmp/")
	})

	t.Run("Fails when couldn't do Mkdir", func(t *testing.T) {
		m := mockFS{
			cache:    map[string]cacheEntity{},
			mkdirErr: fmt.Errorf("Mkdir: fail"),
		}
		c := FSCache{
			FS:        &m,
			CachePath: "/tmp/",
		}
		err := c.Set("key1", testStruct{"val1"})

		require.EqualError(t, err, "Mkdir: fail")
		require.Nil(t, m.cache["/tmp/key1.json"].data)
		require.Contains(t, m.mkdirs, "/tmp/")
	})

	t.Run("Fails when couldn't json.Marshal() the input", func(t *testing.T) {
		m := mockFS{cache: map[string]cacheEntity{}}
		c := FSCache{
			FS:        &m,
			CachePath: "/tmp/",
		}
		err := c.Set("key1", func() {}) // functions are unmarshable

		log.Println(err)
		require.EqualError(t, err, "json: unsupported type: func()")
		require.Nil(t, m.cache["/tmp/key1.json"].data)
		require.Contains(t, m.mkdirs, "/tmp/")
	})

	t.Run("Fails when couldn't write to the file", func(t *testing.T) {
		m := mockFS{
			cache:        map[string]cacheEntity{},
			writeFileErr: fmt.Errorf("Failed to write to the file"),
		}
		c := FSCache{
			FS:        &m,
			CachePath: "/tmp/",
		}
		err := c.Set("key1", testStruct{"val1"})

		require.EqualError(t, err, "Failed to write to the file")
		require.Nil(t, m.cache["/tmp/key1.json"].data)
		require.Contains(t, m.mkdirs, "/tmp/")
	})
}

func TestRead(t *testing.T) {
	t.Run("Works when the file is readable and is correct JSON", func(t *testing.T) {
		s := testStruct{}
		m := mockFS{
			cache: map[string]cacheEntity{
				"/tmp/key1.json": {[]byte(`{"x":"val1"}`), 0777},
			},
		}
		c := FSCache{
			FS:        &m,
			CachePath: "/tmp/",
		}
		ok, err := c.Get("key1", &s)

		require.Nil(t, err)
		require.True(t, ok)
		require.Equal(t, testStruct{"val1"}, s)
	})

	t.Run("Fails when couldn't read the file", func(t *testing.T) {
		s := testStruct{}
		m := mockFS{
			readFileErr: fmt.Errorf("Some weird FS error"),
		}
		c := FSCache{
			FS:        &m,
			CachePath: "/tmp/",
		}
		ok, err := c.Get("key1", &s)

		require.False(t, ok)
		require.EqualError(t, err, "Some weird FS error")
	})

	t.Run("Fails when couldn't unmarshal the contents", func(t *testing.T) {
		s := testStruct{}
		m := mockFS{
			cache: map[string]cacheEntity{
				"/tmp/key1.json": {[]byte(`incorrect JSON`), 0777},
			},
		}
		c := FSCache{
			FS:        &m,
			CachePath: "/tmp/",
		}
		ok, err := c.Get("key1", &s)

		require.True(t, ok)
		require.Contains(t, err.Error(), "invalid character 'i'")
	})

	t.Run("Doesn't fail when the file didn't exist", func(t *testing.T) {
		s := testStruct{}
		m := mockFS{
			readFileErr: fmt.Errorf("/tmp/key1.json: no such file or directory"),
		}
		c := FSCache{
			FS:        &m,
			CachePath: "/tmp/",
		}
		ok, err := c.Get("key1", &s)

		require.False(t, ok)
		require.Nil(t, err)
	})
}

type cacheEntity struct {
	data []byte
	perm os.FileMode
}

type mockFS struct {
	mkdirs       []string
	cache        map[string]cacheEntity
	mkdirErr     error
	writeFileErr error
	readFileErr  error
}

func (m *mockFS) Mkdir(path string, perms os.FileMode) error {
	m.mkdirs = append(m.mkdirs, path)
	return m.mkdirErr
}

func (m *mockFS) WriteFile(key string, data []byte, perm os.FileMode) error {
	if m.writeFileErr == nil {
		m.cache[key] = cacheEntity{data, perm}
	}
	return m.writeFileErr
}

func (m *mockFS) ReadFile(key string) ([]byte, error) {
	entity, ok := m.cache[key]
	if !ok {
		return nil, m.readFileErr
	}
	return entity.data, m.readFileErr
}

type testStruct struct {
	X string `json:"x"`
}
