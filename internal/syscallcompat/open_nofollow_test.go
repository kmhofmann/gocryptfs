package syscallcompat

import (
	"os"
	"syscall"
	"testing"
)

func TestOpenNofollow(t *testing.T) {
	err := os.MkdirAll(tmpDir+"/d1/d2/d3", 0700)
	if err != nil {
		t.Fatal(err)
	}
	// Create a file
	fd, err := OpenNofollow(tmpDir, "d1/d2/d3/f1", syscall.O_RDWR|syscall.O_CREAT|syscall.O_EXCL, 0600)
	if err != nil {
		t.Fatal(err)
	}
	syscall.Close(fd)
	_, err = os.Stat(tmpDir + "/d1/d2/d3/f1")
	if err != nil {
		t.Fatal(err)
	}
	// Replace "d1" with a symlink - open should fail with ELOOP
	err = os.Rename(tmpDir+"/d1", tmpDir+"/d1.renamed")
	if err != nil {
		t.Fatal(err)
	}
	os.Symlink(tmpDir+"/d1.renamed", tmpDir+"/d1")
	fd, err = OpenNofollow(tmpDir, "d1/d2/d3/f1", syscall.O_RDWR|syscall.O_CREAT, 0600)
	if err == nil {
		t.Fatalf("should have failed")
	}
	if err != syscall.ELOOP {
		t.Errorf("expected ELOOP, got %v", err)
	}
}
