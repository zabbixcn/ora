// Copyright 2014 Rana Ian. All rights reserved.
// Use of this source code is governed by The MIT License
// found in the accompanying LICENSE file.

package ora

/*
#include <oci.h>
#include <stdlib.h>
*/
import "C"
import (
	"container/list"
	"fmt"
	"strings"
	"unsafe"
)

// Env is an Oracle environment.
type Env struct {
	id     uint64
	drv    *Drv
	ocienv *C.OCIEnv
	ocierr *C.OCIError

	srvId    uint64
	conId    uint64
	srvs     *list.List
	cons     *list.List
	elem     *list.Element
	stmtCfg  StmtCfg
	errBuf   [512]C.char
	isSqlPkg bool
}

// NumSrv returns the number of open Oracle servers.
func (env *Env) NumSrv() int {
	return env.srvs.Len()
}

// NumCon returns the number of open Oracle connections.
func (env *Env) NumCon() int {
	return env.cons.Len()
}

// checkIsOpen validates that the environment is open.
func (env *Env) checkIsOpen() error {
	if !env.IsOpen() {
		return errNewF("Env is closed (id %v)", env.id)
	}
	return nil
}

// IsOpen returns true when the environment is open; otherwise, false.
//
// Calling Close will cause IsOpen to return false.
// Once closed, the environment may be re-opened by
// calling Open.
func (env *Env) IsOpen() bool {
	return env.drv != nil
}

// Close disconnects from servers and resets optional fields.
func (env *Env) Close() (err error) {
	if err := env.checkIsOpen(); err != nil {
		return err
	}
	Log.Infof("E%v] Close", env.id)
	errs := env.drv.listPool.Get().(*list.List)
	defer func() {
		if value := recover(); value != nil {
			Log.Errorln(recoverMsg(value))
			errs.PushBack(errRecover(value))
		}

		drv := env.drv
		drv.envs.Remove(env.elem)
		env.srvs.Init()
		env.drv = nil
		env.ocienv = nil
		env.ocierr = nil
		env.elem = nil
		drv.envPool.Put(env)

		m := newMultiErrL(errs)
		if m != nil {
			err = *m
		}
		errs.Init()
		drv.listPool.Put(errs)
	}()

	// close connections
	for e := env.cons.Front(); e != nil; e = e.Next() {
		err0 := e.Value.(*Con).Close()
		errs.PushBack(err0)
	}
	// close servers
	for e := env.srvs.Front(); e != nil; e = e.Next() {
		err0 := e.Value.(*Srv).Close()
		errs.PushBack(err0)
	}

	// Free oci environment handle and all oci child handles
	// The oci error handle is released as a child of the environment handle
	err = env.freeOciHandle(unsafe.Pointer(env.ocienv), C.OCI_HTYPE_ENV)
	return err
}

// OpenSrv connects to an Oracle server returning a *Srv and possible error.
func (env *Env) OpenSrv(dbname string) (*Srv, error) {
	if err := env.checkIsOpen(); err != nil {
		return nil, err
	}
	Log.Infof("E%v] OpenSrv (dbname %v)", env.id, dbname)
	// allocate server handle
	ocisrv, err := env.allocOciHandle(C.OCI_HTYPE_SERVER)
	if err != nil {
		return nil, err
	}
	// attach to server
	cDbname := C.CString(dbname)
	defer C.free(unsafe.Pointer(cDbname))
	r := C.OCIServerAttach(
		(*C.OCIServer)(ocisrv),                //OCIServer     *srvhp,
		env.ocierr,                            //OCIError      *errhp,
		(*C.OraText)(unsafe.Pointer(cDbname)), //const OraText *dbname,
		C.sb4(len(dbname)),                    //sb4           dbname_len,
		C.OCI_DEFAULT)                         //ub4           mode);
	if r == C.OCI_ERROR {
		return nil, env.ociError()
	}
	// allocate service context handle
	ocisvcctx, err := env.allocOciHandle(C.OCI_HTYPE_SVCCTX)
	if err != nil {
		return nil, err
	}
	// set server handle onto service context handle
	err = env.setAttr(ocisvcctx, C.OCI_HTYPE_SVCCTX, ocisrv, C.ub4(0), C.OCI_ATTR_SERVER)
	if err != nil {
		return nil, err
	}

	// set srv struct
	srv := env.drv.srvPool.Get().(*Srv)
	if srv.id == 0 {
		env.srvId++
		srv.id = env.srvId
	}
	Log.Infof("E%v] OpenSrv (srvId %v)", env.id, srv.id)
	srv.env = env
	srv.ocisrv = (*C.OCIServer)(ocisrv)
	srv.ocisvcctx = (*C.OCISvcCtx)(ocisvcctx)
	srv.stmtCfg = env.stmtCfg
	srv.dbname = dbname
	srv.elem = env.srvs.PushBack(srv)
	return srv, nil
}

// OpenCon starts an Oracle session on a server returning a *Con and possible error.
//
// The connection string has the form username/password@dbname e.g., scott/tiger@orcl
// dbname is a connection identifier such as a net service name,
// full connection identifier, or a simple connection identifier.
// The dbname may be defined in the client machine's tnsnames.ora file.
func (env *Env) OpenCon(str string) (*Con, error) {
	if err := env.checkIsOpen(); err != nil {
		return nil, err
	}
	Log.Infof("E%v] OpenCon", env.id)
	// parse connection string
	var username string
	var password string
	var dbname string
	str = strings.TrimSpace(str)
	if strings.HasPrefix(str, "/@") {
		dbname = str[2:]
	} else {
		str = strings.Replace(str, "/", " / ", 1)
		str = strings.Replace(str, "@", " @ ", 1)
		_, err := fmt.Sscanf(str, "%s / %s @ %s", &username, &password, &dbname)
		Log.Infof("E%v] OpenCon (dbname %v, username %v)", env.id, dbname, username)
		if err != nil {
			return nil, fmt.Errorf("parse %q: %v", err)
		}
	}
	// connect to server
	srv, err := env.OpenSrv(dbname)
	if err != nil {
		return nil, err
	}
	// open a session on the server
	ses, err := srv.OpenSes(username, password)
	if err != nil {
		return nil, err
	}
	// set con struct
	con := env.drv.conPool.Get().(*Con)
	if con.id == 0 {
		env.conId++
		con.id = env.conId
	}
	Log.Infof("E%v] OpenCon (conId %v)", env.id, con.id)
	con.env = env
	con.srv = srv
	con.ses = ses
	con.elem = env.cons.PushBack(con)

	return con, nil
}

// allocateOciHandle allocates an oci handle.
func (env *Env) allocOciHandle(handleType C.ub4) (unsafe.Pointer, error) {
	// OCIHandleAlloc returns: OCI_SUCCESS, OCI_INVALID_HANDLE
	var handle unsafe.Pointer
	r := C.OCIHandleAlloc(
		unsafe.Pointer(env.ocienv), //const void    *parenth,
		&handle,                    //void          **hndlpp,
		handleType,                 //ub4           type,
		C.size_t(0),                //size_t        xtramem_sz,
		nil)                        //void          **usrmempp
	if r == C.OCI_INVALID_HANDLE {
		return nil, errNew("Unable to allocate handle")
	}
	return handle, nil
}

// freeOciHandle deallocates an oci handle.
func (env *Env) freeOciHandle(ociHandle unsafe.Pointer, handleType C.ub4) error {
	// OCIHandleFree returns: OCI_SUCCESS, OCI_INVALID_HANDLE, or OCI_ERROR
	r := C.OCIHandleFree(
		unsafe.Pointer(env.ocienv), //void      *hndlp,
		handleType)                 //ub4       type );
	if r == C.OCI_INVALID_HANDLE {
		return errNew("Unable to free handle")
	} else if r == C.OCI_ERROR {
		return env.ociError()
	}

	return nil
}

// setOciAttribute sets an attribute value on a handle or descriptor.
func (env *Env) setAttr(
	target unsafe.Pointer,
	targetType C.ub4,
	attribute unsafe.Pointer,
	attributeSize C.ub4,
	attributeType C.ub4) (err error) {

	r := C.OCIAttrSet(
		target,        //void        *trgthndlp,
		targetType,    //ub4         trghndltyp,
		attribute,     //void        *attributep,
		attributeSize, //ub4         size,
		attributeType, //ub4         attrtype,
		env.ocierr)    //OCIError    *errhp );
	if r == C.OCI_ERROR {
		return env.ociError()
	}
	return nil
}

// getOciError gets an error returned by an Oracle server.
func (env *Env) ociError() error {
	var errcode C.sb4
	C.OCIErrorGet(
		unsafe.Pointer(env.ocierr),
		1, nil,
		&errcode,
		(*C.OraText)(unsafe.Pointer(&env.errBuf[0])),
		512,
		C.OCI_HTYPE_ERROR)
	return errNew(C.GoString(&env.errBuf[0]))
}

// Sets the StmtCfg on the Environment and all open Environment Servers.
func (env *Env) SetStmtCfg(c StmtCfg) {
	env.stmtCfg = c
	for e := env.srvs.Front(); e != nil; e = e.Next() {
		e.Value.(*Srv).SetStmtCfg(c)
	}
}

// StmtCfg returns a *StmtCfg.
func (env *Env) StmtCfg() *StmtCfg {
	return &env.stmtCfg
}
