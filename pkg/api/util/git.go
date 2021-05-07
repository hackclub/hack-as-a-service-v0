package util

/*
#cgo pkg-config: libgit2
#include <git2.h>
#include <git2/global.h>
#include <git2/net.h>
#include <git2/oid.h>
#include <git2/remote.h>
*/
import "C"
import (
	"errors"
	"log"
	"reflect"
	"strings"
	"unsafe"
)

// Gets the default branch of a git repository without cloning it
func GetGitRepositoryDefaultBranch(url string) (string, error) {
	err := C.git_libgit2_init()
	if err < 0 {
		return "", errors.New("failed to initialize libgit2")
	}
	var odb *C.git_odb
	err = C.git_odb_new(&odb)
	if err != 0 {
		return "", errors.New("failed to create git_odb")
	}
	var repo *C.git_repository
	err = C.git_repository_wrap_odb(&repo, odb)
	if err != 0 {
		return "", errors.New("failed to wrap odb into repository")
	}
	var remote *C.git_remote
	err = C.git_remote_create_anonymous(&remote, repo, C.CString(url))
	if err != 0 {
		return "", errors.New("failed to create anonymous remote")
	}
	err = C.git_remote_connect(remote, C.GIT_DIRECTION_FETCH, nil, nil, nil)
	if err != 0 {
		return "", errors.New("failed to connect to remote (fetch)")
	}
	var remote_heads **C.git_remote_head
	var remote_heads_size C.ulong
	err = C.git_remote_ls(&remote_heads, &remote_heads_size, remote)
	if err != 0 {
		return "", errors.New("failed to list remote heads")
	}
	var remote_heads2 []*C.git_remote_head
	sh := (*reflect.SliceHeader)(unsafe.Pointer(&remote_heads2))
	sh.Data = uintptr(unsafe.Pointer(remote_heads))
	sh.Len = int(remote_heads_size)
	sh.Cap = int(remote_heads_size)
	found := ""
	for _, remote_head := range remote_heads2 {
		if remote_head.symref_target != nil {
			// s := C.GoString(C.git_oid_tostr_s(&remote_head.oid))
			h := C.GoString(remote_head.name)
			if h != "HEAD" {
				continue
			}
			sr := C.GoString(remote_head.symref_target)
			sr = strings.TrimPrefix(sr, "refs/heads/")
			log.Printf("[%s] Found default branch name = %s\n", h, sr)
			found = sr
		}
	}
	C.git_remote_free(remote)
	C.git_repository_free(repo)
	C.git_odb_free(odb)
	C.git_libgit2_shutdown()

	return found, nil
}
