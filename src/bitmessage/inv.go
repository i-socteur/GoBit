// Copyright 2013 msm595. All rights reserved.
// Copyright 2011 ThePiachu. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bitmessage

import (
	"fmt"
	"mymath"
)

//TODO: change inventory to be able to handle everything?
//TODO: test

const (
	ERROR     uint32 = iota //==0
	MSG_TX    uint32 = iota //==1
	MSG_BLOCK uint32 = iota //==2
)

type Inv struct {
	Count     mymath.VarInt //var_int
	Inventory []*InventoryVector
}

//getdata responce to inv
type GetData Inv

//TODO: delete
/*func (i *Inv)AddNetworkAddress(na *NetworkAddress){
	if i.Count<1388{ //so payload length would be less than 50000 bytes
		i.Count++
		i.Inventory=append(i.Inventory, NewInventoryVectorFromNetworkAddress(na))
	}
}*/

//TODO: support adding blocks and transactions

func InvMessageFromBytes(b []byte) *Inv {
	im := new(Inv)
	im.Count, b = mymath.DecodeVarIntGiveRest(b)

	im.Inventory = make([]*InventoryVector, len(b)/36)
	for i := 0; i < int(im.Count); i++ {
		im.Inventory[i] = InventoryVectorFromBytes(b[i*36:])
	}
	return im
}

func (im *Inv) String() string {
	s := ""

	s += fmt.Sprintf(
		"Inv message:\n"+
			"  %X\t\t\t\t- %d entries\n"+
			"  %s",
		mymath.VarInt2HexRev(im.Count), im.Count,
		im.Inventory)

	return s
}

func (i *Inv) Clear() {
	i.Count = 0
	i.Inventory = nil
}

//TODO: double check if 36 is really the answer
func (i *Inv) Compile() []byte {
	vi := mymath.VarInt2HexRev(i.Count) //TODO: check if Rev or not

	answer := make([]byte, len(vi)+36*len(i.Inventory))

	iterator := 0
	copy(answer[iterator:], vi)
	iterator += len(vi)
	for j := 0; j < len(i.Inventory); j++ {
		copy(answer[iterator:], i.Inventory[j].Compile())
		iterator += 34
	}

	return answer
}

//TODO: test
type InventoryVector struct {
	Vectortype uint32
	Hash       [32]byte
}

func InventoryVectorFromBytes(b []byte) *InventoryVector {
	iv := new(InventoryVector)
	iv.Vectortype = mymath.HexRev2Uint32(b[0:4])
	copy(iv.Hash[:], b[4:])
	return iv
}

//TODO: support adding blocks and transactions

//TODO: delete
/*
func NewInventoryVectorFromNetworkAddress(na *NetworkAddress) *InventoryVector{
	iv:=new(InventoryVector)
	//TODO: set proper type
	copy(iv.Hash[:], na.GetHash())

	return iv
}*/

func (iv *InventoryVector) SetType(newtype uint32) {
	iv.Vectortype = newtype
}

func (iv *InventoryVector) SetHash(newhash []byte) {
	if len(newhash) == 32 {
		copy(iv.Hash[:], newhash)
	}
}

func (iv *InventoryVector) Compile() []byte {
	answer := make([]byte, 4+len(iv.Hash))

	iterator := 0
	copy(answer[iterator:], mymath.Uint322HexRev(iv.Vectortype)) //TODO: check endianess
	iterator += 4
	copy(answer[iterator:], iv.Hash[:])

	return answer
}
