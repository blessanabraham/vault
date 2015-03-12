package vault

import (
	"reflect"
	"testing"
)

func testSystem(t *testing.T) *SystemBackend {
	c, _ := testUnsealedCore(t)
	return &SystemBackend{c}
}

func TestSystem_mounts(t *testing.T) {
	s := testSystem(t)

	req := &Request{
		Operation: ReadOperation,
		Path:      "mounts",
	}
	resp, err := s.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	exp := map[string]interface{}{
		"secret/": map[string]string{
			"type":        "generic",
			"description": "generic secret storage",
		},
		"sys/": map[string]string{
			"type":        "system",
			"description": "system endpoints used for control, policy and debugging",
		},
	}

	if !reflect.DeepEqual(resp.Data, exp) {
		t.Fatalf("got: %#v expect: %#v", resp.Data, exp)
	}

	req = &Request{
		Operation: HelpOperation,
		Path:      "mounts",
	}
	resp, err = s.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if resp.Data["help"] != "logical backend mount table" {
		t.Fatalf("got: %#v", resp.Data)
	}
}

func TestSystem_mount_help(t *testing.T) {
	s := testSystem(t)

	req := &Request{
		Operation: HelpOperation,
		Path:      "mount/prod/secret/",
	}
	resp, err := s.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if resp.Data["help"] != "used to mount or unmount a path" {
		t.Fatalf("got: %#v", resp.Data)
	}
}

func TestSystem_mount(t *testing.T) {
	s := testSystem(t)

	req := &Request{
		Operation: WriteOperation,
		Path:      "mount/prod/secret/",
		Data: map[string]interface{}{
			"type": "generic",
		},
	}
	resp, err := s.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if resp != nil {
		t.Fatalf("bad: %v", resp)
	}
}

func TestSystem_mount_invalid(t *testing.T) {
	s := testSystem(t)

	req := &Request{
		Operation: WriteOperation,
		Path:      "mount/prod/secret/",
		Data: map[string]interface{}{
			"type": "what",
		},
	}
	resp, err := s.HandleRequest(req)
	if err != ErrInvalidRequest {
		t.Fatalf("err: %v", err)
	}
	if resp.Data["error"] != "unknown logical backend type: what" {
		t.Fatalf("bad: %v", resp)
	}
}

func TestSystem_unmount(t *testing.T) {
	s := testSystem(t)

	req := &Request{
		Operation: DeleteOperation,
		Path:      "mount/secret/",
	}
	resp, err := s.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if resp != nil {
		t.Fatalf("bad: %v", resp)
	}
}

func TestSystem_unmount_invalid(t *testing.T) {
	s := testSystem(t)

	req := &Request{
		Operation: DeleteOperation,
		Path:      "mount/foo/",
	}
	resp, err := s.HandleRequest(req)
	if err != ErrInvalidRequest {
		t.Fatalf("err: %v", err)
	}
	if resp.Data["error"] != "no matching mount" {
		t.Fatalf("bad: %v", resp)
	}
}