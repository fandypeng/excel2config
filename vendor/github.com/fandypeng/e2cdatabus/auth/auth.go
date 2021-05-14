package auth

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"log"
	"math"
)

const _abortIndex int8 = math.MaxInt8 / 2

// implement from credentials.PerRPCCredentials

type Auth struct {
	appKey    string
	appSecret string
	handlers  []grpc.UnaryServerInterceptor
}

func New(appKey, appSecret string) *Auth {
	return &Auth{appKey: appKey, appSecret: appSecret}
}

func (a *Auth) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	reqInfo, ok := credentials.RequestInfoFromContext(ctx)
	if !ok {
		log.Fatal("no request info")
	}
	return map[string]string{
		"token": a.getToken(reqInfo.Method),
	}, nil
}

func (a *Auth) RequireTransportSecurity() bool {
	return false
}

func (a *Auth) getToken(method string) (token string) {
	source := method + a.appKey
	bytes := a.aesEncrypt([]byte(source), []byte(a.appSecret))
	token = base64.StdEncoding.EncodeToString(bytes)
	return
}

func (a *Auth) AccessControl() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			err = status.Errorf(codes.Unauthenticated, "no auth token")
			return
		}
		reqAuthKey, ok := md["token"]
		if !ok {
			err = status.Errorf(codes.Unauthenticated, "no auth token")
			return
		}
		if reqAuthKey[0] != a.getToken(info.FullMethod) {
			err = status.Errorf(codes.Unauthenticated, "auth failed")
			return
		}
		return handler(ctx, req)
	}
}

func (a *Auth) Interceptor(ctx context.Context, req interface{}, args *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	var (
		i     int
		chain grpc.UnaryHandler
	)
	n := len(a.handlers)
	if n == 0 {
		return handler(ctx, req)
	}
	chain = func(ictx context.Context, ireq interface{}) (interface{}, error) {
		if i == n-1 {
			return handler(ctx, req)
		}
		i++
		return a.handlers[i](ictx, ireq, args, chain)
	}
	return a.handlers[0](ctx, req, args, chain)
}

func (a *Auth) Use(handlers ...grpc.UnaryServerInterceptor) {
	finalSize := len(a.handlers) + len(handlers)
	if finalSize >= int(_abortIndex) {
		panic("rpc server use too many handlers")
	}
	mergedHandlers := make([]grpc.UnaryServerInterceptor, finalSize)
	copy(mergedHandlers, a.handlers)
	copy(mergedHandlers[len(a.handlers):], handlers)
	a.handlers = mergedHandlers
}

//aes加密 分组模式ctr
func (a *Auth) aesEncrypt(plaintext, key []byte) []byte {
	//1. 建立一个底层使用aes的密码接口
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	iv := []byte("12345678abcdefgh")
	stream := cipher.NewCTR(block, iv)
	cipherText := make([]byte, len(plaintext))
	stream.XORKeyStream(cipherText, plaintext)
	return cipherText
}

func (a *Auth) aesDecrypt(cipherText, key []byte) []byte {
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	iv := []byte("12345678abcdefgh")
	stream := cipher.NewCTR(block, iv)
	plainText := make([]byte, len(cipherText))
	stream.XORKeyStream(plainText, cipherText)
	return plainText
}
