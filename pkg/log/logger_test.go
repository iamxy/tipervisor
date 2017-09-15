package log

import (
	"bytes"
	"encoding/json"
	"strconv"
	"strings"
	"sync"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func LogAndAssertJSON(t *testing.T, log func(), assertions func(map[string]string)) {
	var buffer bytes.Buffer
	var fields map[string]string

	SetOutput(&buffer)
	SetFormatter(new(logrus.JSONFormatter))
	SetLevel(logrus.DebugLevel)
	SetRootFields(Fields{"scope": "all", "module": "test"})

	log()

	err := json.Unmarshal(buffer.Bytes(), &fields)
	assert.Nil(t, err)

	assertions(fields)
}

func LogAndAssertText(t *testing.T, log func(), assertions func(map[string]string)) {
	var buffer bytes.Buffer

	SetOutput(&buffer)
	SetFormatter(&logrus.TextFormatter{
		DisableColors: true,
	})
	SetLevel(logrus.DebugLevel)
	SetRootFields(Fields{"scope": "all", "module": "test"})

	log()

	fields := make(map[string]string)
	for _, kv := range strings.Split(buffer.String(), " ") {
		if !strings.Contains(kv, "=") {
			continue
		}
		kvArr := strings.Split(kv, "=")
		key := strings.TrimSpace(kvArr[0])
		val := kvArr[1]
		if kvArr[1][0] == '"' {
			var err error
			val, err = strconv.Unquote(val)
			assert.NoError(t, err)
		}
		fields[key] = val
	}
	assertions(fields)
}

func TestPrint(t *testing.T) {
	LogAndAssertJSON(t, func() {
		Print("test print")
	}, func(fields map[string]string) {
		assert.Equal(t, fields["msg"], "test print")
		assert.Equal(t, fields["level"], "info")
		assert.Equal(t, fields["scope"], "all")
		assert.Equal(t, fields["module"], "test")
	})
}

func TestInfo(t *testing.T) {
	LogAndAssertJSON(t, func() {
		Info("test info")
	}, func(fields map[string]string) {
		assert.Equal(t, fields["msg"], "test info")
		assert.Equal(t, fields["level"], "info")
		assert.Equal(t, fields["scope"], "all")
		assert.Equal(t, fields["module"], "test")
	})
}

func TestWarn(t *testing.T) {
	LogAndAssertJSON(t, func() {
		Warn("test warning")
	}, func(fields map[string]string) {
		assert.Equal(t, fields["msg"], "test warning")
		assert.Equal(t, fields["level"], "warning")
		assert.Equal(t, fields["scope"], "all")
		assert.Equal(t, fields["module"], "test")
	})
}

func TestError(t *testing.T) {
	LogAndAssertJSON(t, func() {
		Error("test error")
	}, func(fields map[string]string) {
		assert.Equal(t, fields["msg"], "test error")
		assert.Equal(t, fields["level"], "error")
		assert.Equal(t, fields["scope"], "all")
		assert.Equal(t, fields["module"], "test")
	})
}

func TestInfolnShouldAddSpacesBetweenStrings(t *testing.T) {
	LogAndAssertJSON(t, func() {
		Infoln("test", "test")
	}, func(fields map[string]string) {
		assert.Equal(t, fields["msg"], "test test")
	})
}

func TestInfolnShouldAddSpacesBetweenStringAndNonstring(t *testing.T) {
	LogAndAssertJSON(t, func() {
		Infoln("test", 10)
	}, func(fields map[string]string) {
		assert.Equal(t, fields["msg"], "test 10")
	})
}

func TestInfolnShouldAddSpacesBetweenTwoNonStrings(t *testing.T) {
	LogAndAssertJSON(t, func() {
		Infoln(10, 10)
	}, func(fields map[string]string) {
		assert.Equal(t, fields["msg"], "10 10")
	})
}

func TestInfoShouldNotAddSpacesBetweenStringAndNonstring(t *testing.T) {
	LogAndAssertJSON(t, func() {
		Info("test", 10)
	}, func(fields map[string]string) {
		assert.Equal(t, fields["msg"], "test10")
	})
}

func TestInfoShouldNotAddSpacesBetweenStrings(t *testing.T) {
	LogAndAssertJSON(t, func() {
		Info("test", "test")
	}, func(fields map[string]string) {
		assert.Equal(t, fields["msg"], "testtest")
	})
}

func TestWithFieldsShouldAllowAssignments(t *testing.T) {
	LogAndAssertJSON(t, func() {
		f := WithField("key1", "value1")
		f.WithFields(Fields{
			"key2": "value2",
		}).Info("test withfields")
	}, func(fields map[string]string) {
		assert.Equal(t, "value2", fields["key2"])
		assert.Equal(t, "value1", fields["key1"])
		assert.Equal(t, "test withfields", fields["msg"])
		assert.Equal(t, "info", fields["level"])
		assert.Equal(t, "all", fields["scope"])
		assert.Equal(t, "test", fields["module"])
	})
}

func TestRootFieldsCouldBeOverwritten(t *testing.T) {
	LogAndAssertJSON(t, func() {
		WithField("module", "test2").Info("test field overwritten")
	}, func(fields map[string]string) {
		assert.Equal(t, "test2", fields["module"])
		assert.Equal(t, "test field overwritten", fields["msg"])
	})
}

func TestDefaultFieldAreNotPrefixed(t *testing.T) {
	LogAndAssertText(t, func() {
		ll := WithField("herp", "derp")
		ll.Info("hello")
		ll.Info("bye")
	}, func(fields map[string]string) {
		for _, fieldName := range []string{"fields.level", "fields.time", "fields.msg"} {
			if _, ok := fields[fieldName]; ok {
				t.Fatalf("should not have prefixed %q: %v", fieldName, fields)
			}
		}
	})
}

func TestGetSetLevelRace(t *testing.T) {
	var buffer bytes.Buffer
	SetOutput(&buffer)
	SetFormatter(new(logrus.JSONFormatter))
	SetRootFields(Fields{"scope": "all", "module": "test"})

	var wg sync.WaitGroup
	wg.Add(100)
	for i := 0; i < 100; i++ {
		go func(i int) {
			defer wg.Done()
			if i%2 == 0 {
				SetLevel(logrus.InfoLevel)
			} else {
				GetLevel()
			}
		}(i)
	}
	wg.Wait()
}

func TestLoggingRace(t *testing.T) {
	var buffer bytes.Buffer
	SetOutput(&buffer)
	SetFormatter(new(logrus.JSONFormatter))
	SetRootFields(Fields{"scope": "all", "module": "test"})

	var wg sync.WaitGroup
	wg.Add(100)
	for i := 0; i < 100; i++ {
		go func(i int) {
			defer wg.Done()
			if i%2 == 0 {
				Info("info")
			} else {
				WithFields(Fields{"key1": "value1"}).Error("error")
			}
		}(i)
	}
	wg.Wait()
}
