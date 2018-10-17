package conv

import (
	"context"
	"errors"
	"reflect"
	"strconv"

	"github.com/alipay/sofa-mosn/pkg/protocol/rpc"
	"github.com/alipay/sofa-mosn/pkg/protocol/rpc/sofarpc"
)

// PropertyHeaders map the cmdkey and its data type
var (
	PropertyHeaders = make(map[string]reflect.Kind, 11)
	boltv1          = new(boltv1conv)
)

func init() {
	PropertyHeaders[sofarpc.HeaderProtocolCode] = reflect.Uint8
	PropertyHeaders[sofarpc.HeaderCmdType] = reflect.Uint8
	PropertyHeaders[sofarpc.HeaderCmdCode] = reflect.Int16
	PropertyHeaders[sofarpc.HeaderVersion] = reflect.Uint8
	PropertyHeaders[sofarpc.HeaderReqID] = reflect.Uint32
	PropertyHeaders[sofarpc.HeaderCodec] = reflect.Uint8
	PropertyHeaders[sofarpc.HeaderClassLen] = reflect.Int16
	PropertyHeaders[sofarpc.HeaderHeaderLen] = reflect.Int16
	PropertyHeaders[sofarpc.HeaderContentLen] = reflect.Int
	PropertyHeaders[sofarpc.HeaderTimeout] = reflect.Int
	PropertyHeaders[sofarpc.HeaderRespStatus] = reflect.Int16
	PropertyHeaders[sofarpc.HeaderRespTimeMills] = reflect.Int64

	sofarpc.RegisterConv(sofarpc.PROTOCOL_CODE_V1, boltv1)
}

type boltv1conv struct{}

func (b *boltv1conv) MapToCmd(ctx context.Context, headers map[string]string) (sofarpc.SofaRpcCmd, error) {
	if len(headers) < 8 {
		return nil, errors.New("headers count not enough")
	}

	value := sofarpc.GetPropertyValue1(PropertyHeaders, headers, sofarpc.HeaderProtocolCode)
	protocolCode := sofarpc.ConvertPropertyValueUint8(value)
	value = sofarpc.GetPropertyValue1(PropertyHeaders, headers, sofarpc.HeaderCmdType)
	cmdType := sofarpc.ConvertPropertyValueUint8(value)
	value = sofarpc.GetPropertyValue1(PropertyHeaders, headers, sofarpc.HeaderCmdCode)
	cmdCode := sofarpc.ConvertPropertyValueInt16(value)
	value = sofarpc.GetPropertyValue1(PropertyHeaders, headers, sofarpc.HeaderVersion)
	version := sofarpc.ConvertPropertyValueUint8(value)
	value = sofarpc.GetPropertyValue1(PropertyHeaders, headers, sofarpc.HeaderReqID)
	requestID := sofarpc.ConvertPropertyValueUint32(value)
	value = sofarpc.GetPropertyValue1(PropertyHeaders, headers, sofarpc.HeaderCodec)
	codec := sofarpc.ConvertPropertyValueUint8(value)
	//value = sofarpc.GetPropertyValue1(PropertyHeaders, headers, sofarpc.HeaderClassLen)
	//classLength := sofarpc.ConvertPropertyValueInt16(value)
	//value = sofarpc.GetPropertyValue1(PropertyHeaders, headers, sofarpc.HeaderHeaderLen)
	//headerLength := sofarpc.ConvertPropertyValueInt16(value)
	value = sofarpc.GetPropertyValue1(PropertyHeaders, headers, sofarpc.HeaderContentLen)
	contentLength := sofarpc.ConvertPropertyValueInt(value)

	//class
	className := sofarpc.GetPropertyValue1(PropertyHeaders, headers, sofarpc.HeaderClassName)

	//RPC Request
	if cmdType == sofarpc.REQUEST || cmdType == sofarpc.REQUEST_ONEWAY {
		value = sofarpc.GetPropertyValue1(PropertyHeaders, headers, sofarpc.HeaderTimeout)
		timeout := sofarpc.ConvertPropertyValueInt(value)

		//sofabuffers := sofarpc.SofaProtocolBuffersByContext(ctx)
		//request := &sofabuffers.BoltEncodeReq
		request := &sofarpc.BoltRequest{}
		request.Protocol = protocolCode
		request.CmdType = cmdType
		request.CmdCode = cmdCode
		request.Version = version
		request.ReqID = requestID
		request.Codec = codec
		request.Timeout = timeout
		//request.ClassLen = classLength
		//request.HeaderLen = headerLength
		request.ContentLen = contentLength
		request.RequestClass = className
		request.RequestHeader = headers
		return request, nil
	} else if cmdType == sofarpc.RESPONSE {
		value = sofarpc.GetPropertyValue1(PropertyHeaders, headers, sofarpc.HeaderRespStatus)
		responseStatus := sofarpc.ConvertPropertyValueInt16(value)
		value = sofarpc.GetPropertyValue1(PropertyHeaders, headers, sofarpc.HeaderRespTimeMills)
		responseTime := sofarpc.ConvertPropertyValueInt64(value)

		//sofabuffers := sofarpc.SofaProtocolBuffersByContext(ctx)
		//response := &sofabuffers.BoltEncodeRsp
		response := &sofarpc.BoltResponse{}
		response.Protocol = protocolCode
		response.CmdType = cmdType
		response.CmdCode = cmdCode
		response.Version = version
		response.ReqID = requestID
		response.Codec = codec
		response.ResponseStatus = responseStatus
		//response.ClassLen = classLength
		//response.HeaderLen = headerLength
		response.ContentLen = contentLength
		response.ResponseClass = className
		response.ResponseHeader = headers
		response.ResponseTimeMillis = responseTime
		return response, nil
	}

	return nil, rpc.ErrUnknownType
}

//Convert BoltV1's Protocol Header  and Content Header to Map[string]string
func (b *boltv1conv) MapToFields(ctx context.Context, cmd sofarpc.SofaRpcCmd) (map[string]string, error) {
	switch c := cmd.(type) {
	case *sofarpc.BoltRequest:
		return mapReqToFields(ctx, c)
	case *sofarpc.BoltResponse:
		return mapRespToFields(ctx, c)
	}

	return nil, rpc.ErrUnknownType
}

func mapReqToFields(ctx context.Context, req *sofarpc.BoltRequest) (map[string]string, error) {
	// TODO: map reuse
	//protocolCtx := protocol.ProtocolBuffersByContext(ctx)
	//headers := make(map[string]string, 9+len(req.RequestHeader))
	headers := req.RequestHeader

	headers[sofarpc.SofaPropertyHeader(sofarpc.HeaderProtocolCode)] = strconv.FormatUint(uint64(req.Protocol), 10)
	headers[sofarpc.SofaPropertyHeader(sofarpc.HeaderCmdType)] = strconv.FormatUint(uint64(req.CmdType), 10)
	headers[sofarpc.SofaPropertyHeader(sofarpc.HeaderCmdCode)] = strconv.FormatUint(uint64(req.CmdCode), 10)
	headers[sofarpc.SofaPropertyHeader(sofarpc.HeaderVersion)] = strconv.FormatUint(uint64(req.Version), 10)
	headers[sofarpc.SofaPropertyHeader(sofarpc.HeaderReqID)] = strconv.FormatUint(uint64(req.ReqID), 10)
	headers[sofarpc.SofaPropertyHeader(sofarpc.HeaderCodec)] = strconv.FormatUint(uint64(req.Codec), 10)
	headers[sofarpc.SofaPropertyHeader(sofarpc.HeaderTimeout)] = strconv.FormatUint(uint64(req.Timeout), 10)

	// TODO: bypass length header
	//headers[sofarpc.SofaPropertyHeader(sofarpc.HeaderClassLen)] = strconv.FormatUint(uint64(req.ClassLen), 10)
	//headers[sofarpc.SofaPropertyHeader(sofarpc.HeaderHeaderLen)] = strconv.FormatUint(uint64(req.HeaderLen), 10)
	headers[sofarpc.SofaPropertyHeader(sofarpc.HeaderContentLen)] = strconv.FormatUint(uint64(req.ContentLen), 10)

	headers[sofarpc.SofaPropertyHeader(sofarpc.HeaderClassName)] = req.RequestClass

	return headers, nil
}

func mapRespToFields(ctx context.Context, resp *sofarpc.BoltResponse) (map[string]string, error) {
	// TODO: map reuse
	//protocolCtx := protocol.ProtocolBuffersByContext(ctx)
	//headers := make(map[string]string, 12)

	headers := resp.ResponseHeader

	headers[sofarpc.SofaPropertyHeader(sofarpc.HeaderProtocolCode)] = strconv.FormatUint(uint64(resp.Protocol), 10)
	headers[sofarpc.SofaPropertyHeader(sofarpc.HeaderCmdType)] = strconv.FormatUint(uint64(resp.CmdType), 10)
	headers[sofarpc.SofaPropertyHeader(sofarpc.HeaderCmdCode)] = strconv.FormatUint(uint64(resp.CmdCode), 10)
	headers[sofarpc.SofaPropertyHeader(sofarpc.HeaderVersion)] = strconv.FormatUint(uint64(resp.Version), 10)
	headers[sofarpc.SofaPropertyHeader(sofarpc.HeaderReqID)] = strconv.FormatUint(uint64(resp.ReqID), 10)
	headers[sofarpc.SofaPropertyHeader(sofarpc.HeaderCodec)] = strconv.FormatUint(uint64(resp.Codec), 10)

	// TODO: bypass length header
	//headers[sofarpc.SofaPropertyHeader(sofarpc.HeaderClassLen)] = strconv.FormatUint(uint64(resp.ClassLen), 10)
	//headers[sofarpc.SofaPropertyHeader(sofarpc.HeaderHeaderLen)] = strconv.FormatUint(uint64(resp.HeaderLen), 10)
	headers[sofarpc.SofaPropertyHeader(sofarpc.HeaderContentLen)] = strconv.FormatUint(uint64(resp.ContentLen), 10)

	// FOR RESPONSE,ENCODE RESPONSE STATUS and RESPONSE TIME
	headers[sofarpc.SofaPropertyHeader(sofarpc.HeaderRespStatus)] = strconv.FormatUint(uint64(resp.ResponseStatus), 10)
	headers[sofarpc.SofaPropertyHeader(sofarpc.HeaderRespTimeMills)] = strconv.FormatUint(uint64(resp.ResponseTimeMillis), 10)

	headers[sofarpc.SofaPropertyHeader(sofarpc.HeaderClassName)] = resp.ResponseClass

	return headers, nil
}
