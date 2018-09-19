// Copyright 2018 necomeshi, All Rights Reserved.
// fix bugs by frog
package rpmlib

/*
#include <stdlib.h>
#include <stdint.h>
#include <fcntl.h>
#include <rpm/header.h>
#include <rpm/rpmlib.h>
#include <rpm/rpmdb.h>
#include <rpm/rpmmacro.h>
#include <rpm/rpmts.h>
#include <rpm/rpmtag.h>

const char * api_headerGetString(Header h, unsigned int tag)
{
	return  headerGetString(h, tag);
}

int api_headerGet(Header h, unsigned int tag, rpmtd td, headerGetFlags flags)
{
	return headerGet(h,  tag, td, flags);
}


#cgo LDFLAGS: -lrpm -lrpmio -lpopt
*/
import "C"
import "fmt"
import "io"
import "os"
import "unsafe"

//////////////////////////////////////////
// Header
//////////////////////////////////////////
type Header struct {
	rpmheader C.Header
}

func (h *Header) GetString(tag RpmTagVal) (s string, err error) {

	cstring := C.api_headerGetString(h.rpmheader, C.uint(tag))
	if cstring == nil {
		err = fmt.Errorf("C.headerGetString: Cannot get tag value")
		return
	}

	if cstring == nil {
		err = fmt.Errorf("C.headerGetString: Cannot get tag value")
		return
	}

	s = C.GoString(cstring)

	// cstring is 'const char*'. Do not C.free it.

	return
}

func (h *Header) GetStringArray(tag RpmTagVal) (values []string, err error) {

	var sa C.struct_rpmtd_s

	if C.api_headerGet(h.rpmheader, C.uint(tag), &sa, C.HEADERGET_MINMEM) == 0 {
		err = fmt.Errorf("C.headerGet: Cannot get tag value")
		return
	}

	if C.rpmtdCount(&sa) == 0 {
		C.rpmtdFreeData(&sa)
		err = fmt.Errorf("C.rpmtdCount: Tag not contain data")
		return
	}

	for {
		v := C.rpmtdNextString(&sa)
		if v == nil {
			break
		}

		gv := C.GoString(v)
		values = append(values, gv)
	}

	C.rpmtdFreeData(&sa)

	return
}

func (h *Header) GetUint16(tag RpmTagVal) (value uint16, err error) {

	var sa C.struct_rpmtd_s

	if C.api_headerGet(h.rpmheader, C.uint(tag), &sa, C.HEADERGET_MINMEM) == 0 {
		err = fmt.Errorf("C.headerGet: Cannot get tag value")
		return
	}

	v := C.rpmtdGetUint16(&sa)
	if v == nil {
		C.rpmtdFreeData(&sa)
		err = fmt.Errorf("C.rpmtdGetUint16: Tag not contain data")
		return
	}

	value = uint16(*v)
	C.free(unsafe.Pointer(v))

	C.rpmtdFreeData(&sa)
	return
}

func (h *Header) GetUint32(tag RpmTagVal) (value uint32, err error) {

	var sa C.struct_rpmtd_s

	if C.api_headerGet(h.rpmheader, C.uint(tag), &sa, C.HEADERGET_MINMEM) == 0 {
		err = fmt.Errorf("C.headerGet: Cannot get tag value")
		return
	}

	v := C.rpmtdGetUint32(&sa)
	if v == nil {
		C.rpmtdFreeData(&sa)
		err = fmt.Errorf("C.rpmtdGetUint32: Tag not contain data")
		return
	}

	value = uint32(*v)
	C.free(unsafe.Pointer(v))

	C.rpmtdFreeData(&sa)
	return
}

func (h *Header) GetUint32Array(tag RpmTagVal) (values []uint32, err error) {

	var sa C.struct_rpmtd_s

	if C.api_headerGet(h.rpmheader, C.uint(tag), &sa, C.HEADERGET_MINMEM) == 0 {
		err = fmt.Errorf("C.headerGet: Cannot get tag value")
		return
	}

	if C.rpmtdCount(&sa) == 0 {
		C.rpmtdFreeData(&sa)
		err = fmt.Errorf("C.rpmtdCount: Tag not contain data")
		return
	}

	for {
		v := C.rpmtdNextUint32(&sa)
		if v == nil {
			break
		}

		gv := uint32(*v)
		values = append(values, gv)

		C.free(unsafe.Pointer(v))
	}

	C.rpmtdFreeData(&sa)

	return
}

func (h *Header) GetUint64(tag RpmTagVal) (value uint64, err error) {

	var sa C.struct_rpmtd_s

	if C.api_headerGet(h.rpmheader, C.uint(tag), &sa, C.HEADERGET_MINMEM) == 0 {
		err = fmt.Errorf("C.headerGet: Cannot get tag value")
		return
	}

	v := C.rpmtdGetUint64(&sa)
	if v == nil {
		C.rpmtdFreeData(&sa)
		err = fmt.Errorf("C.rpmtdGetUint64: Tag not contain data")
		return
	}

	value = uint64(*v)
	C.free(unsafe.Pointer(v))

	C.rpmtdFreeData(&sa)
	return
}

func (h *Header) GetUint64Array(tag RpmTagVal) (values []uint64, err error) {

	var sa C.struct_rpmtd_s

	if C.api_headerGet(h.rpmheader, C.uint(tag), &sa, C.HEADERGET_MINMEM) == 0 {
		err = fmt.Errorf("C.headerGet: Cannot get tag value")
		return
	}

	if C.rpmtdCount(&sa) == 0 {
		C.rpmtdFreeData(&sa)
		err = fmt.Errorf("C.rpmtdCount: Tag not contain data")
		return
	}

	for {
		v := C.rpmtdNextUint64(&sa)
		if v == nil {
			break
		}

		gv := uint64(*v)
		values = append(values, gv)

		C.free(unsafe.Pointer(v))
	}

	C.rpmtdFreeData(&sa)

	return
}

func (h *Header) IsSource() (isSource bool) {
	if C.headerIsSource(h.rpmheader) == 1 {
		isSource = true
	} else {
		isSource = false
	}
	return
}

func (h *Header) Free() {
	C.headerFree(h.rpmheader)
}

//////////////////////////////////////////
// Database iterator
//////////////////////////////////////////
type Iterator struct {
	mi C.rpmdbMatchIterator
}

func (iter *Iterator) Next() (h *Header, err error) {
	h = new(Header)
	rpmheader := C.rpmdbNextIterator(iter.mi)

	if rpmheader == nil {
		err = io.EOF
	}

	h.rpmheader = C.headerLink(rpmheader)

	return
}

func (iter *Iterator) Free() {
	C.rpmdbFreeIterator(iter.mi)
}

//////////////////////////////////////////
// Transaction
//////////////////////////////////////////

type TransactionSet struct {
	ts C.rpmts
}

func NewTransactionSet() (ts *TransactionSet, err error) {

	ts = new(TransactionSet)

	ts.ts = C.rpmtsCreate()

	C.rpmtsSetRootDir(ts.ts, C.CString("/"))

	return
}

func (ts *TransactionSet) SequencialIterator() (iter *Iterator, err error) {
	iter = new(Iterator)

	iter.mi = C.rpmtsInitIterator(ts.ts, C.RPMDBI_PACKAGES, nil, 0)
	if iter.mi == nil {
		err = fmt.Errorf("C.rpmtsInitIterator: Cannot get iterator")
	}

	return
}

func (ts *TransactionSet) ReadPackageFile(name string, verifySignature bool) (header *Header, err error) {
	cname := C.CString(name)
	cmode := C.CString("r.ufdio")

	fd := C.Fopen(cname, cmode)

	if fd == nil {
		C.free(unsafe.Pointer(cname))
		C.free(unsafe.Pointer(cmode))
		return nil, fmt.Errorf("C.Fopen: Error")
	}

	header = new(Header)
	ret := C.rpmReadPackageFile(ts.ts, fd, cname, &header.rpmheader)

	if ret == C.RPMRC_OK ||
		(!verifySignature && (ret == C.RPMRC_NOTTRUSTED || ret == C.RPMRC_NOKEY)) {
		C.Fclose(fd)
		C.free(unsafe.Pointer(cname))
		C.free(unsafe.Pointer(cmode))
		return
	}

	// Else, we have a problem.
	C.free(unsafe.Pointer(cname))
	C.free(unsafe.Pointer(cmode))
	return nil, fmt.Errorf("C.rpmReadPackageFile: Error")
}

func (ts *TransactionSet) Free() {
	// Opened database will be closed in here
	C.rpmtsFree(ts.ts)
}

//////////////////////////////////////////
// Macros & helpers
//////////////////////////////////////////

// DefineGlobalMacro sets a priority 0 macro.
func DefineGlobalMacro(macro string) {
	// Defines a macro in the environment with the highest priority.
	C.rpmDefineMacro(nil, C.CString(macro), 0)
}

// SetDbPath changes the path the rpmLib will use to look for the DB.
// Note that this is the DIRECTORY where the Packages file resides.
func SetDbPath(path string) {
	DefineGlobalMacro("_dbpath " + path)
}

// ShowRC writes to a file the full RPM environment configuration.
func ShowRC(f *os.File) {
	C.rpmShowRC(C.fdopen(C.int(f.Fd()), C.CString("w*")))
}

//////////////////////////////////////////
// Macros
//////////////////////////////////////////

func init() {
	C.rpmReadConfigFiles(nil, nil)
}
