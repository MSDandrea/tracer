package tracer

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

type mockLogWriter struct {
	mock.Mock
}

func (mw mockLogWriter) Write(entry Entry) {
	mw.Called(entry.Message, entry.TransactionId, entry.Level, entry.Args)
}

func TestLogger_Alert(t *testing.T) {
	t.Parallel()
	t.Run("when the logger should create implicit transactions", func(t *testing.T) {
		t.Parallel()
		is := assert.New(t)
		expected := Entry{
			Message:       "this is a test",
			TransactionId: "some-id",
			Level:         Alert,
			Args:          []interface{}{"some-id"},
		}
		lw1 := mockLogWriter{}
		lw1.On("Write", expected.Message, expected.TransactionId, expected.Level, expected.Args).Return().Once()
		lw2 := mockLogWriter{}
		lw2.On("Write", expected.Message, expected.TransactionId, expected.Level, expected.Args).Return().Once()

		subject := logger{
			createImplicitTransactions: true,
			writers:                    []Writer{lw1, lw2},
		}
		is.NotPanics(func() {
			subject.Alert("this is a test", "some-id")
		}, "it should not panics")
		lw1.AssertExpectations(t)
		lw2.AssertExpectations(t)
	})
	t.Run("when the logger should not create implicit transactions", func(t *testing.T) {
		t.Parallel()
		t.Run("but is not on trace mode", func(t *testing.T) {
			t.Parallel()
			is := assert.New(t)
			expected := Entry{
				Message:       "this is a test",
				TransactionId: "",
				Level:         Alert,
				Args:          []interface{}{"some-field"},
			}
			lw1 := mockLogWriter{}
			lw1.On("Write", expected.Message, expected.TransactionId, expected.Level, expected.Args).Return().Once()
			lw2 := mockLogWriter{}
			lw2.On("Write", expected.Message, expected.TransactionId, expected.Level, expected.Args).Return().Once()

			subject := logger{
				createImplicitTransactions: false,
				writers:                    []Writer{lw1, lw2},
			}
			is.NotPanics(func() {
				subject.Alert("this is a test", "some-field")
			}, "it should not panics")
			lw1.AssertExpectations(t)
			lw2.AssertExpectations(t)
		})
		t.Run("and is on trace mode", func(t *testing.T) {
			t.Parallel()
			is := assert.New(t)
			expected := Entry{
				Message:       "this is a test",
				TransactionId: "some-id",
				Level:         Alert,
				Args:          []interface{}{"some-field"},
			}
			lw1 := mockLogWriter{}
			lw1.On("Write", expected.Message, expected.TransactionId, expected.Level, expected.Args).Return().Once()
			lw2 := mockLogWriter{}
			lw2.On("Write", expected.Message, expected.TransactionId, expected.Level, expected.Args).Return().Once()

			subject := logger{
				createImplicitTransactions: false,
				transactionId:              "some-id",
				writers:                    []Writer{lw1, lw2},
			}
			is.NotPanics(func() {
				subject.Alert("this is a test", "some-field")
			}, "it should not panics")
			lw1.AssertExpectations(t)
			lw2.AssertExpectations(t)
		})
	})
}

func TestLogger_Trace(t *testing.T) {
	t.Parallel()
	is := assert.New(t)
	subject := &logger{}
	newLogger := subject.Trace("some-id")
	is.NotEqual(subject.transactionId, (newLogger.(*logger)).transactionId, "it should not have the same transaction id as the original logger")
	is.Equal("some-id", (newLogger.(*logger)).transactionId, "it should set the transaction id")
}

func TestGetLogger(t *testing.T) {
	t.Parallel()
	is := assert.New(t)
	loggerA := GetLogger("A")
	loggerB := GetLogger("B")
	loggerB2 := GetLogger("B")
	is.NotEqual(loggerA, loggerB, "it should return two different loggers")
	is.Equal(loggerB, loggerB2, "it should return the same logger")
}
