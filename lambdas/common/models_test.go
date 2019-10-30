package common

import (
	"testing"
)

func TestNewReleaseEvent_Validate(t *testing.T) {
	ev1 := NewReleaseEvent{
		Event:       "test",
		BuildId:     123,
		Branch:      "somebranch",
		DownloadUrl: "https://someurl.server.com/path",
		ProductName: "some product",
	}

	er1 := ev1.Validate()
	if er1 != nil {
		t.Errorf("Test failed to validate: got %s", er1)
	}

	ev2 := NewReleaseEvent{
		Event:       "test",
		BuildId:     123,
		Branch:      "",
		DownloadUrl: "https://someurl.server.com/path",
		ProductName: "some product",
	}

	er2 := ev2.Validate()
	if er2 == nil {
		t.Errorf("Validation on empty branch should have failed but it succeeded")
	}

	ev3 := NewReleaseEvent{
		Event:       "test",
		BuildId:     123,
		Branch:      "somebranch",
		DownloadUrl: "https://someurl.server.com/path",
		ProductName: "",
	}

	er3 := ev3.Validate()
	if er3 == nil {
		t.Errorf("Validation on product name should have failed but it succeeded")
	}

	ev4 := NewReleaseEvent{
		Event:       "test",
		BuildId:     123,
		Branch:      "somebranch",
		DownloadUrl: "",
		ProductName: "some product",
	}

	er4 := ev4.Validate()
	if er4 == nil {
		t.Errorf("Validation on product name should have failed but it succeeded")
	}

	ev5 := NewReleaseEvent{
		Event:       "test",
		BuildId:     123,
		Branch:      "somebranch",
		DownloadUrl: "malformedurl!",
		ProductName: "some product",
	}

	er5 := ev5.Validate()
	if er5 == nil {
		t.Errorf("Validation on malformed URL should have failed but it succeeded")
	}
}
