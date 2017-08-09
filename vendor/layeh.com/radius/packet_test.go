package radius_test

import (
	"bytes"
	"net"
	"testing"

	"layeh.com/radius"
	"layeh.com/radius/rfc2865"
)

func Test_RFC2865_7_1(t *testing.T) {
	// Source: https://tools.ietf.org/html/rfc2865#section-7.1

	secret := []byte("xyzzy5461")

	// Request
	request := []byte{
		0x01, 0x00, 0x00, 0x38, 0x0f, 0x40, 0x3f, 0x94, 0x73, 0x97, 0x80, 0x57, 0xbd, 0x83, 0xd5, 0xcb,
		0x98, 0xf4, 0x22, 0x7a, 0x01, 0x06, 0x6e, 0x65, 0x6d, 0x6f, 0x02, 0x12, 0x0d, 0xbe, 0x70, 0x8d,
		0x93, 0xd4, 0x13, 0xce, 0x31, 0x96, 0xe4, 0x3f, 0x78, 0x2a, 0x0a, 0xee, 0x04, 0x06, 0xc0, 0xa8,
		0x01, 0x10, 0x05, 0x06, 0x00, 0x00, 0x00, 0x03,
	}

	p, err := radius.Parse(request, secret)
	if err != nil {
		t.Fatal(err)
	}
	if p.Code != radius.CodeAccessRequest {
		t.Fatal("expecting Code = PacketCodeAccessRequest")
	}
	if p.Identifier != 0 {
		t.Fatal("expecting Identifier = 0")
	}
	if p.Len() != 4 {
		t.Fatal("expecting 4 attributes")
	}
	if rfc2865.UserName_GetString(p) != "nemo" {
		t.Fatal("expecting User-Name = nemo")
	}
	if rfc2865.UserPassword_GetString(p) != "arctangent" {
		t.Fatal("expecting User-Password = arctangent")
	}
	if !rfc2865.NASIPAddress_Get(p).Equal(net.ParseIP("192.168.1.16")) {
		t.Fatal("expecting NAS-IP-Address = 192.168.1.16")
	}
	if rfc2865.NASPort_Get(p) != 3 {
		t.Fatal("expecting NAS-Port = 3")
	}

	{
		wire, err := p.Encode()
		if err != nil {
			t.Fatal(err)
		}
		if !RADIUSPacketsEqual(wire, request) {
			t.Fatal("expecting q.Encode() and request to be equal")
		}
	}

	// Response
	response := []byte{
		0x02, 0x00, 0x00, 0x26, 0x86, 0xfe, 0x22, 0x0e, 0x76, 0x24, 0xba, 0x2a, 0x10, 0x05, 0xf6, 0xbf,
		0x9b, 0x55, 0xe0, 0xb2, 0x06, 0x06, 0x00, 0x00, 0x00, 0x01, 0x0f, 0x06, 0x00, 0x00, 0x00, 0x00,
		0x0e, 0x06, 0xc0, 0xa8, 0x01, 0x03,
	}

	q := radius.Packet{
		Code:          radius.CodeAccessAccept,
		Identifier:    p.Identifier,
		Authenticator: p.Authenticator,
		Secret:        secret,
		Attributes:    make(radius.Attributes),
	}
	rfc2865.ServiceType_Set(&q, rfc2865.ServiceType(1))
	rfc2865.LoginService_Set(&q, rfc2865.LoginService(0))
	if err := rfc2865.LoginIPHost_Set(&q, net.ParseIP("192.168.1.3")); err != nil {
		t.Fatal(err)
	}

	{
		wire, err := q.Encode()
		if err != nil {
			t.Fatal(err)
		}
		if !RADIUSPacketsEqual(wire, response) {
			t.Fatalf("expecting q.Encode() and response to be equal\n%v\n%v", wire, response)
		}
	}

	if !radius.IsAuthenticResponse(response, request, secret) {
		t.Fatal("expecting response to be valid")
	}
}

func Test_RFC2865_7_2(t *testing.T) {
	// Source: https://tools.ietf.org/html/rfc2865#section-7.2

	secret := []byte("xyzzy5461")

	// Request
	request := []byte{
		0x01, 0x01, 0x00, 0x47, 0x2a, 0xee, 0x86, 0xf0, 0x8d, 0x0d, 0x55, 0x96, 0x9c, 0xa5, 0x97, 0x8e,
		0x0d, 0x33, 0x67, 0xa2, 0x01, 0x08, 0x66, 0x6c, 0x6f, 0x70, 0x73, 0x79, 0x03, 0x13, 0x16, 0xe9,
		0x75, 0x57, 0xc3, 0x16, 0x18, 0x58, 0x95, 0xf2, 0x93, 0xff, 0x63, 0x44, 0x07, 0x72, 0x75, 0x04,
		0x06, 0xc0, 0xa8, 0x01, 0x10, 0x05, 0x06, 0x00, 0x00, 0x00, 0x14, 0x06, 0x06, 0x00, 0x00, 0x00,
		0x02, 0x07, 0x06, 0x00, 0x00, 0x00, 0x01,
	}

	p, err := radius.Parse(request, secret)
	if err != nil {
		t.Fatal(err)
	}

	if p.Code != radius.CodeAccessRequest {
		t.Fatal("expecting code access request")
	}
	if p.Identifier != 1 {
		t.Fatal("expecting Identifier = 1")
	}
	if rfc2865.UserName_GetString(p) != "flopsy" {
		t.Fatal("expecting User-Name = flopsy")
	}
	if !rfc2865.NASIPAddress_Get(p).Equal(net.ParseIP("192.168.1.16")) {
		t.Fatal("expecting NAS-IP-Address = 192.168.1.16")
	}
	if rfc2865.NASPort_Get(p) != 20 {
		t.Fatal("expecting NAS-Port = 20")
	}
	if rfc2865.ServiceType_Get(p) != rfc2865.ServiceType_Value_FramedUser {
		t.Fatal("expecting Service-Type = Attr_ServiceType_FramedUser")
	}
	if rfc2865.FramedProtocol_Get(p) != rfc2865.FramedProtocol_Value_PPP {
		t.Fatal("expecting Framed-Protocol = Attr_FramedProtocol_PPP")
	}
}

func TestPasswords(t *testing.T) {
	passwords := []string{
		"",
		"qwerty",
		"helloworld1231231231231233489hegufudhsgdsfygdf8g",
	}

	for _, password := range passwords {
		secret := []byte("xyzzy5461")

		r := radius.New(radius.CodeAccessRequest, secret)
		if r == nil {
			t.Fatal("could not create new RADIUS packet")
		}
		if err := rfc2865.UserPassword_AddString(r, password); err != nil {
			t.Fatal(err)
		}

		b, err := r.Encode()
		if err != nil {
			t.Fatal(err)
		}

		q, err := radius.Parse(b, secret)
		if err != nil {
			t.Fatal(err)
		}

		if s := rfc2865.UserPassword_GetString(q); s != password {
			t.Fatalf("incorrect User-Password (expecting %q, got %q)", password, s)
		}
	}
}

// RADIUSPacketsEqual returns if two RADIUS packets are equal, ignoring the
// order of attributes of different types.
func RADIUSPacketsEqual(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	if !bytes.Equal(a[:4], b[:4]) {
		return false
	}

	// hash is going to be different, as the attribute order could change

	aa, err := radius.ParseAttributes(a[20:])
	if err != nil {
		panic(err)
	}
	ab, err := radius.ParseAttributes(b[20:])
	if err != nil {
		panic(err)
	}

	if len(aa) != len(ab) {
		return false
	}

	for typeA, attrsA := range aa {
		if len(attrsA) != len(ab[typeA]) {
			return false
		}
		for i, attrA := range attrsA {
			if !bytes.Equal([]byte(attrA), []byte(ab[typeA][i])) {
				return false
			}
		}
	}

	return true
}
