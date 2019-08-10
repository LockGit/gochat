/**
 * Created by lock
 * Date: 2019-08-10
 * Time: 18:35
 */
package connect

import "gochat/proto"

type Operator interface {
	Connect(conn *proto.ConnectRequest) (string, error)
	DisConnect(disConn *proto.DisConnectRequest) (err error)
}

type DefaultOperator struct {
}

//rpc call logic layer
func (o *DefaultOperator) Connect(conn *proto.ConnectRequest) (uid string, err error) {
	return "", nil
}

//rpc call logic layer
func (o *DefaultOperator) DisConnect(disConn *proto.DisConnectRequest) (err error) {
	return nil
}
