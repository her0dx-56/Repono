package storage

import (
	"context"
	"strings"
	"testing"
	"io"
)

func TestLocalStorage_SaveGetDelete(t *testing.T) {

	ctx:=context.Background()

	dir:= t.TempDir()
	store:=NewLocalStorage(dir)

	key:= "test/file.txt"
	content:="hello go dfs"

	err:=store.Save(
		ctx,
		key,
		strings.NewReader(content),
	)
	if err!=nil{
		t.Fatalf("save() error = %v",err)
	}

	reader,err:=store.Get(ctx,key)
	if err!=nil{
		t.Fatalf("Get() error = %v",err)
	}
	

	got,err:=io.ReadAll(reader)
	if err!=nil{
		reader.Close()
		t.Fatalf("ReadAll() error = %v",err)
	}

	reader.Close()

	if string(got) != content {
		t.Fatalf(
			"Get() content = %q, wnt %q",
			string(got),
			content,
		)
	}

	err=store.Delete(ctx,key)
	if err!=nil{
		t.Fatalf("Delete() error = %v", err)
	}
}