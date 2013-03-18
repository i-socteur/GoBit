// Copyright 2013 msm595. All rights reserved.
// Copyright 2011 ThePiachu. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bitmessage

//Package for handling low-level Bitcoin messages

import (
	"io"
	"mymath"
)

const MainNetMagic uint32 = 0xD9B4BEF9
const TestNetMagic uint32 = 0xDAB5BFFA

const MainNetPort uint16 = 8333
const TestNetPort uint16 = 18333

const (
	WRONGMESSAGE = iota
	VERSION      = iota
	VERACK       = iota
	ADDR         = iota
	INV          = iota
	GETDATA      = iota
	GETBLOCKS    = iota
	GETHEADERS   = iota
	TX           = iota
	BLOCK        = iota
	HEADERS      = iota
	GETADDR      = iota
	CHECKORDER   = iota
	SUBMITORDER  = iota
	REPLY        = iota
	PING         = iota
	ALERT        = iota
)

/*
Messages

version - Information about program version and block count. Exchanged when first connecting.
verack - Sent in response to a version message to acknowledge that we are willing to connect.
addr - List of one or more IP addresses and ports.
inv - "I have these blocks/transactions: ..." Normally sent only when a new block or transaction is being relayed. This is only a list, not the actual data.
getdata - Request a single block or transaction by hash.
getblocks - Request an inv of all blocks in a range.
getheaders - Request a headers message containing all block headers in a range.
tx - Send a transaction. This is only sent in response to a getdata request.
block - Send a block. This is only sent in response to a getdata request.
headers - Send up to 2,000 block headers. Non-generators can download the headers of blocks instead of entire blocks.
getaddr - Request an addr message containing a bunch of known-active peers (for bootstrapping).
submitorder, checkorder, and reply - Used when performing an IP transaction.
alert - Send a network alert.
ping - Does nothing. Used to check that the connection is still online. A TCP error will occur if the connection has died.
*/

//getaddr has no data transmitted with the message

type BitMessage struct {
	Magic    [4]byte  //Magic 4 network bytes
	Command  [12]byte //12 ASCII characters identifying package content, NULL padded
	Length   [4]byte  //Length of the Payload
	Checksum [4]byte  //first 4 bytes of sha256(sha256(Payload))
	Payload  []byte   //the actual data
}

func (bm *BitMessage) SetMagic(newMagic uint32) {
	magchars := mymath.Uint322Hex(newMagic)

	bm.Magic[0] = magchars[3]
	bm.Magic[1] = magchars[2]
	bm.Magic[2] = magchars[1]
	bm.Magic[3] = magchars[0]
}

func (bm *BitMessage) SetCommand(newCommand string) {

	if len(newCommand) > 11 {
		return
	}

	for i := 0; i < len(bm.Command); i++ {
		bm.Command[i] = 0
	}

	for i := 0; i < len(newCommand); i++ {
		bm.Command[i] = newCommand[i]
	}
}

func (bm *BitMessage) GetMagic() []byte {
	return bm.Magic[:]
}

func (bm *BitMessage) SetPayloadVersion(vm *VersionMessage) {
	tmp := vm.Compile()
	bm.Payload = make([]byte, len(tmp))
	for i := 0; i < len(tmp); i++ {
		bm.Payload[i] = tmp[i]
	}
	bm.calculateChecksumAndLength()
}

func (bm *BitMessage) SetPayloadByte(pld []byte) {
	bm.Payload = make([]byte, len(pld))
	for i := 0; i < len(pld); i++ {
		bm.Payload[i] = pld[i]
	}
	bm.calculateChecksumAndLength()
}

func (bm *BitMessage) SetPayloadAddr(al *AddressList) {
	tmp := al.Compile()
	bm.Payload = make([]byte, len(tmp))
	for i := 0; i < len(tmp); i++ {
		bm.Payload[i] = tmp[i]
	}
	bm.calculateChecksumAndLength()
}

func (bm *BitMessage) calculateChecksumAndLength() {
	tmp := mymath.Uint322HexRev(uint32(len(bm.Payload)))

	for i := 0; i < 4; i++ {
		bm.Length[i] = tmp[i]
	}

	tmp = mymath.DoubleSHA(bm.Payload)

	for i := 0; i < 4; i++ {
		bm.Checksum[i] = tmp[i]
	}
}

func (bm *BitMessage) Compile() []byte {
	answer := make([]byte, len(bm.Magic)+len(bm.Command)+len(bm.Length)+len(bm.Checksum)+len(bm.Payload))

	iterator := 0
	copy(answer[iterator:], bm.Magic[:])
	iterator += 4
	copy(answer[iterator:], bm.Command[:])
	iterator += 12
	copy(answer[iterator:], bm.Length[:])
	iterator += 4
	copy(answer[iterator:], bm.Checksum[:])
	iterator += 4
	copy(answer[iterator:], bm.Payload[:])

	return answer
}

//TODO: test
func (bm *BitMessage) GiveMessageType() int {
	return MessageType(bm.Command[:])
}

func (bm *BitMessage) String() string {
	s := ""

	magic := mymath.Hex2Uint32(bm.Magic[:])
	if magic == TestNetMagic {
		s += "TestNet\n"
	} else if magic == MainNetMagic {
		s += "MainNet\n"
	} else {
		s += "UnknownNet\n"
	}

	

	return s
}

func MessageType(Command []byte) int {

	msglen := 0
	for i := 0; i < len(Command); i++ {
		if Command[i] == 0 {
			msglen = i
			break
		}
	}

	cmd := mymath.Byte2String(Command[0:msglen])

	//log.Printf("%i - len", msglen)
	//log.Printf("%i - len of Command", len(cmd))
	//log.Print(cmd)

	switch cmd {
	default:
		return WRONGMESSAGE
	case "version":
		return VERSION
	case "verack":
		return VERACK
	case "addr":
		return ADDR
	case "inv":
		return INV
	case "getdata":
		return GETDATA
	case "getblocks":
		return GETBLOCKS
	case "getheaders":
		return GETHEADERS
	case "tx":
		return TX
	case "block":
		return BLOCK
	case "headers":
		return HEADERS
	case "getaddr":
		return GETADDR
	case "checkorder":
		return CHECKORDER
	case "submitorder":
		return SUBMITORDER
	case "reply":
		return REPLY
	case "ping":
		return PING
	case "alert":
		return ALERT
	}
	return WRONGMESSAGE
}

func DeserializeMessages(msgs []byte) []*BitMessage {
	iterator := 0
	var vec []*BitMessage
	for len(msgs[iterator:]) >= 20 {
		msgtype := MessageType(msgs[iterator+4 : iterator+16])

		Payloadbyte := msgs[iterator+16 : iterator+20]
		Payload := int(mymath.HexRev2Big(Payloadbyte).Int64())

		if msgtype == VERSION || msgtype == VERACK {
			vec = append(vec, DecodeMessage(msgs[iterator:iterator+20+Payload]))
			iterator = iterator + 20 + Payload
		} else {
			vec = append(vec, DecodeMessage(msgs[iterator:iterator+24+Payload]))
			iterator = iterator + 24 + Payload
		}
	}
	/*answer:=make([]BitMessage, vec.Len())

	for i:=0;i<vec.Len();i++{
		answer[i]=vec.At(i).(BitMessage)
	}*/

	return vec
}

func DecodeMessage(msg []byte) *BitMessage {
	bm := new(BitMessage)

	if len(msg) < 24 {
		return nil
	}
	for i := 0; i < 4; i++ {
		bm.Magic[i] = msg[0+i]
		bm.Length[i] = msg[16+i]
		bm.Checksum[i] = msg[20+i]
	}
	for i := 0; i < 12; i++ {
		bm.Command[i] = msg[4+i]
	}

//	msgtype := MessageType(msg[4:16])

	Payloadbyte := msg[16:20]
	Payload := int(mymath.HexRev2Big(Payloadbyte).Int64())

	if len(msg) < 24+Payload {
		return nil
	} else {
		bm.SetPayloadByte(msg[24 : 24+Payload])
	}

	return bm
}

//TODO: seek magic nums, deal with faulty data, check checksum
func DecodeMessages(data io.Reader, msgChan chan<- *BitMessage) {
	for {
		buf := make([]byte, 24)
		n, err := io.ReadFull(data, buf)
		if err != nil || n < 24 {
			close(msgChan)
			return
		}

		bm := new(BitMessage)

		for i := 0; i < 4; i++ {
			bm.Magic[i] = buf[0+i]
			bm.Length[i] = buf[16+i]
			bm.Checksum[i] = buf[20+i]
		}
		for i := 0; i < 12; i++ {
			bm.Command[i] = buf[4+i]
		}

		payload := int(mymath.HexRev2Big(bm.Length[:]).Int64())

		buf = make([]byte, payload)
		n, err = io.ReadFull(data, buf)
		if err != nil || n < payload {
			close(msgChan)
			return
		}

		bm.SetPayloadByte(buf)
		msgChan <- bm
	}
}
