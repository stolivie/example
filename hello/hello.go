package main

import (
	"bufio"
	"fmt"
	"github.com/stolivie/example/bubba"
	"github.com/stolivie/example/stringutil"
	"net"
	"os"
	"time"
)

var transactionId int32 = 1

type cmdLogin struct {
	namelen        byte
	name           string
	passwordLength byte
	password       string
}

const (
	//Asset single
	protectionStatus         int16 = 1
	isBlocked                      = 2
	paddingStart                   = 3
	paddingStop                    = 4
	serviceId                      = 5
	eTitle                         = 6
	timeStart                      = 7
	timeEnd                        = 8
	channelName                    = 9
	longDescription                = 10
	isClosedCaptionAvailable       = 11
	isHiDef                        = 12
	willSoonBeDeleted              = 13
	timeOfDeletion                 = 14
	dispChannelNumber              = 15
	isRepeating                    = 16
	seriesId                       = 17
	viewedStatus                   = 18
	isLocked                       = 19
	isManual                       = 20
	networkId                      = 21
	duration                       = 22
	protectionStatusExt            = 23
	//Series
	episodeDuration   = 101
	dayOfTheWeek      = 102
	isNewEpisodesOnly = 104
	isThisTimeOnly    = 105
	isThisDayOnly     = 106
	isThisChannelOnly = 107
	priority          = 108
	//Preference
	blockRatingThreshold      = 201
	isBlockContent            = 202
	isHideBlockTitle          = 203
	isBlockUnratedShowBlocked = 204
	isPcEnabled               = 205
	//Folder
	numBlocked               = 301
	numProtectionAutoDelSoon = 302
	numProtectionSaveUntilI  = 303
	numConflictLoser         = 304
	numRecorded              = 305
	numScheduled             = 306
	//Conflict Item
	serviceType       = 401
	isAssetCancelled  = 402
	isSeriesCancelled = 403
	isSeries          = 404
	seriesPriority    = 405
	//Conflict
	requiredBandwidth             = 501
	isBwShared                    = 502
	isDropPaddingOnConflict       = 503
	isSeriesConflict              = 504
	isEpisodeConflictDueToPadding = 505
	resourceType                  = 506
)

const (
	reserved int16 = iota
	sysGetCapabilities
	sysNoOp
	sysLogin
	sysLogout
	sysObtainPurchPermission
	sysRemovePurchPermission
	preferenceGetPcScheme
	preferenceGetPc
	preferenceGetPassword
	preferenceSetPassword
	preferenceChannelCreateQuery
	preferenceChannelDeleteQuery
	preferenceChannelGetListChunk
	preferenceSetChannel
	preferenceSetPc
	preferenceSysGetStbFriendlyName
	preferenceChangeNotify
	assetSchedByChanTimeCreateQuery
	assetDeleteQuery
	assetSchedByChanTimeGetListChunk
	assetCreateQuery
	assetRecBriefGetListChunk
	assetRecFolderDetailsGet
	assetRecAssetDetailsGet
	assetRecAssetDelete
	assetRecAssetChangeNotify
	assetRecAssetModify
	assetSchedBriefGetListChunk
	assetSchedFolderDetailsGet
	assetSchedSeriesDetailsGet
	assetSchedAssetDetailsGet
	assetSchedAssetCreate
	assetSchedAssetDelete
	assetSchedAssetModify
	assetSchedAssetChangeNotify
	assetSeriesDelete
	assetSeriesModify
	assetSeriesChangeNotify
	assetGetConflict
	conflictDetailsGetList
	conflictConflictItemDetailsGetList
	unused
	conflictCalculateConflictTestList
	conflictApplyResolutionChanges
	conflictDeleteTestList
	assetSchedSeriesBriefGetListChunk
	sysGetTimezone
	dvrGetStorageFree
	assetFolderDelete
	assetFolderChangeNotify
	conflictApplyDefaultResolutionPolicy
	preferenceChannelGetListChunkExt
	preferenceSetChannelExt
	assetSchedByChanTimeCreateQueryExt
	assetSchedByChanTimeGetListChunkExt
	assetCreateQueryExt
	assetRecBriefGetListChunkExt
	assetSchedAssetCreateExt
	assetSchedAssetListDetailsGet
	assetRecAssetListDetailsGet
	preferenceSysSetStbFriendlyName
	sysDbCheckpointSet
	sysDbUpdateGet
	sysDbConnect
	sysDbDisconnect
	assetFolderChangeNotifyExt
)

func main() {
	fmt.Printf(stringutil.Reverse("Hello Go!"))
	fmt.Printf("\n")
	fmt.Printf(stringutil.Reverse("!oG ,olleH"))
	fmt.Printf("\n")

	conn, err := net.Dial("tcp", "10.10.1.23:9994")
	//conn, err := net.Dial("tcp", "www.facebook.com:80")
	if err != nil {
		fmt.Printf("Connection to host failed - %s", err.Error())
		os.Exit(1)
	}
	defer conn.Close()

	err = conn.SetDeadline(time.Now().Add(60000 * time.Millisecond))
	if err != nil {
		fmt.Printf("SetDeadline failed: %s", err.Error())
		os.Exit(1)
	}

	bubba.Bubba()

	var cmd2 []byte
	var n int

	cmd2, err = SysDbConnect()
	if err != nil {
		fmt.Printf("SysDbConnect packet create failed: %s", err.Error())
		os.Exit(1)
	}
	fmt.Printf("length %d\n", len(cmd2))
	write(conn, cmd2)
	cmd2, n, err = read(conn)
	for i := 0; i < n; i++ {
		fmt.Printf(" %02x", cmd2[i])
	}
	fmt.Printf("\n\n")

	time.Sleep(1 * time.Second)

	cmd2, err = AssetRecAssetListDetailsGet()
	if err != nil {
		fmt.Printf("SysDbDisconnect packet create failed: %s", err.Error())
		os.Exit(1)
	}
	fmt.Printf("length %d message:%q\n", len(cmd2), cmd2)
	write(conn, cmd2)
	err = parse60(conn)
	if err != nil {
		fmt.Printf("Error in parse60() %q", err)
	}
	time.Sleep(3 * time.Second)
	cmd2, err = SysDbDisconnect()
	if err != nil {
		fmt.Printf("SysDbDisconnect packet create failed: %s", err.Error())
		os.Exit(1)
	}
	fmt.Printf("length %d\n", len(cmd2))
	write(conn, cmd2)
	cmd2, n, err = read(conn)
	for i := 0; i < n; i++ {
		fmt.Printf(" %02x", cmd2[i])
	}
	fmt.Printf("\n\ndone\n")

	time.Sleep(10 * time.Second)
	/*
		write(conn, cmd2)
		read(conn)

		time.Sleep(1 * time.Second)

		var numChar int
		numChar, err = conn.Write(cmd2)
		if err != nil {
			fmt.Printf("Client Write failed: %s", err.Error())
			os.Exit(1)
		}
		fmt.Printf("Client Wrote %d characters\n", numChar)

		time.Sleep(1 * time.Second)
		read(conn)

		var b []byte
		numChar, err = conn.Read(b)
		if err != nil {
			fmt.Printf("Client Read failed: %s", err.Error())
			os.Exit(1)
		}
		fmt.Printf("Client Read %d characters\n\n", numChar)
	*/
}
func parse60(conn net.Conn) error {
	var cmd []byte
	var n, pos int
	var err error
	var assetLoopLength int16
	var asset int16
	var assetId int32
	var descriptorTag int16
	var descriptorLength int8
	var descriptorPos int8
	var title []byte
	var assetOffset uint16
	var descriptorOffset uint16

	/*	
	daproHeader() {
   daproIdetnifier 32
   payloadLength   32
}

params() {
   assetLoopLength 16
   for (assetLoopLength) {
      assetId 32
      descriptorLoopLength 16
      for (descriptorLoopLength) {
         descriptor()
      }
   }
}

descriptor() {
			   descriptorTag 16
			   descriptorLength 8
			   for(descriptorLength) {
			      descriptorValue
			   }
			}
	44 61 50 72
	1a 0b 00 00
	00 00 00 00
	00
	02 02 01 0b 10 00 00 01 d2 00 08 00 06
	'D' 'a' 'P' 'r' 
	'\x1a' '\  v' '\x00' '\x00' 
	'\x00' '\x00' '\x00' '\x00' 
	'\x00' '\x02' '\x02' '\x01' 
	'\  v' '\x10' '\x00' '\x00' 
	'\x01' '   Ò' '\x00' '\  b'
	 
	'\x00' '\x06' '\x05' 'D' 'r' 'i' 'v' 'e'                                                      '\x00' '\x00' '\x01' 'Ó' '\x00' '\v' 
	'\x00' '\x06' '\  b' 'V' 'a' 'c' 'a' 't' 'i' 'o' 'n'                                          '\x00' '\x00' '\x01' 'Ô' '\x00' '\x0e' 
	'\x00' '\x06' '\  v' 'H' 'o' 'u' 's' 'e' 's' 'i' 't' 't' 'e' 'r'                              '\x00' '\x00' '\x01' 'Ø' '\x00' '\x13' 
	'\x00' '\x06' '\x10' 'S' 'o' 'u' 't' 'h' 'e' 'r' 'n' ' ' 'J' 'u' 's' 't' 'i' 'c' 'e'          '\x00' '\x00' '\x01' 'Ú' '\x00' '\x13' 
	'\x00' '\x06' '\x10' 'S' 'o' 'u' 't' 'h' 'e' 'r' 'n' ' ' 'J' 'u' 's' 't' 'i' 'c' 'e'          '\x00' '\x00' '\x01' 'Þ' '\x00' '\x13' 
	'\x00' '\x06' '\x10' 'S' 'o' 'u' 't' 'h' 'e' 'r' 'n' ' ' 'J' 'u' 's' 't' 'i' 'c' 'e'          '\x00' '\x00' '\x01' 'ì' '\x00' '\v' 
	'\x00' '\x06' '\  b' 'V' 'a' 'c' 'a' 't' 'i' 'o' 'n'                                          '\x00' '\x00' '\x01' 'í' '\x00' '\x15' 
	'\x00' '\x06' '\x12' 'C' 'o' 'l' 'l' 'e' 'g' 'e' ' ' 'B' 'a' 's' 'k' 'e' 't' 'b' 'a' 'l' 'l'  '\x00' '\x00' '\x01' 'î' '\x00' '\x15' 
	'\x00' '\x06' '\x12' 'C' 'o' 'l' 'l' 'e' 'g' 'e' ' ' 'B' 'a' 's' 'k' 'e' 't' 'b' 'a' 'l' 'l'  '\x00' '\x00' '\x04' 'h' '\x00' '\x15' 
	'\x00' '\x06' '\x12' 'C' 'o' 'l' 'l' 'e' 'g' 'e' ' ' 'B' 'a' 's' 'k' 'e' 't' 'b' 'a' 'l' 'l'  '\x00' '\x00' '\x04' 'i' '\x00' '\f' 
	'\x00' '\x06' '\  t' 'O' 'U' ' ' 'R' 'e' 'v' 'i' 'e' 'w'                                      '\x00' '\x00' '\x04' 'j' '\x00' '\x15' 
	'\x00' '\x06' '\x12' 'C' 'o' 'l' 'l' 'e' 'g' 'e' ' ' 'B' 'a' 's' 'k' 'e' 't' 'b' 'a' 'l' 'l' '\x00' '\x00' '\x04' 'k' '\x00' '\x1c' 
	'\x00' '\x06' '\x19' 'G' 'a' 'r' 'm' ' ' 'W' 'a' 'r' 's' ':' ' ' 'T' 'h' 'e' ' ' 'L' 'a' 's' 't' ' ' 'D' 'r' 'u' 'i' 'd' '\x00' '\x00' '\x04' 'm' '\x00' '\x1c' 
	'\x00' '\x06' '\x19' 'G' 'a' 'r' 'm' ' ' 'W' 'a' 'r' 's' ':' ' ' 'T' 'h' 'e' ' ' 'L' 'a' 's' 't' ' ' 'D' 'r' 'u' 'i' 'd' '\x00' '\x00' '\x04' 'n' '\x00' '\x1c' 
	'\x00' '\x06' '\x19' 'G' 'a' 'r' 'm' ' ' 'W' 'a' 'r' 's' ':' ' ' 'T' 'h' 'e' ' ' 'L' 'a' 's' 't' ' ' 'D' 'r' 'u' 'i' 'd' '\x00' '\x00' '\x04' 'q' '\x00' '&' 
	*/

	descriptorFinished := false
	done := false
	for !done {
		cmd, n, err = read(conn)
		if err != nil {
			return err
		}
		if pos == 0 {
			pos = 15
			fmt.Printf("ALL: %x %x\n", cmd[pos], cmd[pos+1])
			assetLoopLength = int16(cmd[pos]) << 8
			assetLoopLength += int16(cmd[pos+1])
			pos += 2
			assetOffset = 0
			fmt.Printf("AssetLoopLength: %d\n", assetLoopLength)
		}

		for ; asset < assetLoopLength || pos < n; asset++ {
			fmt.Printf("AID: %x %x %x %x\n", cmd[pos], cmd[pos+1], cmd[pos+2], cmd[pos+3])
			if assetOffset == 0 {
				assetId = int32(cmd[pos]) 
				pos++
				assetOffset += 8
				if pos == n {
					break
				}
			}
			if assetOffset == 8 {
				assetId += int32(cmd[pos]) << assetOffset
				assetOffset += 8
				pos++
				if pos == n {
					break
				}
			}
			if assetOffset == 16 {
				assetId += int32(cmd[pos]) << assetOffset
				assetOffset += 8
				pos++
				if pos == n {
					break
				}
			}
			if assetOffset == 24 {
				assetId += int32(cmd[pos]) << assetOffset
				assetOffset -= 8
				pos++
				if pos == n {
					break
				}
				fmt.Printf("AssetId: %d\n", assetId)
				descriptorOffset = 8
			}
			assetOffset = 0xFF

			if descriptorFinished {
				if descriptorOffset == 8 {
					descriptorTag = int16(cmd[pos]) << descriptorOffset
					descriptorOffset -= 8
					pos++
					if pos == n {
						break
					}
				}
				if descriptorOffset == 0 {
					descriptorTag += int16(cmd[pos]) << descriptorOffset
					descriptorOffset -= 8
					pos++
					if pos == n {
						break
					}
					descriptorLength = -1
				}
				fmt.Printf("DescriptorTag: %d\n", descriptorTag)
				descriptorOffset = 0x00
				if descriptorLength == -1 {
					descriptorLength = int8(cmd[pos])
					pos++
					if pos == n {
						break
					}
				}
				descriptorFinished = false
				descriptorPos = 0
			}
			for ; descriptorPos < descriptorLength && pos < n; asset++ {
				if descriptorTag == eTitle {
					if descriptorPos == 0 {
						title = make([]byte, descriptorLength)
					}
					for ; descriptorPos < descriptorLength && pos < n; descriptorPos++ {
						title[descriptorPos] = cmd[pos]
						pos++
					}
					if pos == n && descriptorPos < descriptorLength {
						break
					}
				}
				if descriptorPos == descriptorLength {
					descriptorFinished = true
					descriptorPos = 0
					if descriptorTag == eTitle {
						fmt.Printf("Title: %s\n", title)
					}
				}
			}
			assetOffset = 0
		}
		done = true
	}
	return nil
}
func DaProHeader(method int16) ([]byte, int, error) {
	var cmd []byte
	cmd = make([]byte, 1024)

	cmd[0] = 0x44
	cmd[1] = 0x61
	cmd[2] = 0x50
	cmd[3] = 0x72

	cmd[8] = byte(transactionId >> 24 & 0xff)
	cmd[9] = byte(transactionId >> 16 & 0xff)
	cmd[10] = byte(transactionId >> 8 & 0xff)
	cmd[11] = byte(transactionId & 0xff)
	transactionId++

	cmd[12] = 1

	cmd[13] = byte(method >> 8 & 0xff)
	cmd[14] = byte(method & 0xff)

	return cmd, 15, nil
}

func SysDbConnect() ([]byte, error) {
	cmd, pos, _ := DaProHeader(sysDbConnect)

	var len uint32
	len = uint32(pos - 8)
	cmd[4] = byte(len >> 24 & 0xff)
	cmd[5] = byte(len >> 16 & 0xff)
	cmd[6] = byte(len >> 8 & 0xff)
	cmd[7] = byte(len & 0xff)

	var cmd2 []byte
	cmd2 = make([]byte, pos)

	fmt.Printf("SysDbConnect message: ")
	for i := 0; i < pos; i++ {
		cmd2[i] = cmd[i]
		fmt.Printf(" %02x", cmd2[i])
	}
	fmt.Printf("\n")

	//fmt.Printf("message: % q\n", cmd2)
	return cmd2, nil

}

func SysDbDisconnect() ([]byte, error) {
	cmd, pos, _ := DaProHeader(sysDbDisconnect)

	var len uint32
	len = uint32(pos - 8)
	cmd[4] = byte(len >> 24 & 0xff)
	cmd[5] = byte(len >> 16 & 0xff)
	cmd[6] = byte(len >> 8 & 0xff)
	cmd[7] = byte(len & 0xff)

	var cmd2 []byte
	cmd2 = make([]byte, pos)

	fmt.Printf("SysDbDisconnect message: ")
	for i := 0; i < pos; i++ {
		cmd2[i] = cmd[i]
		fmt.Printf(" %02x", cmd2[i])
	}
	fmt.Printf("\n")

	//fmt.Printf("message: % q\n", cmd2)
	return cmd2, nil

}

func SysLogin() ([]byte, error) {
	cmd, pos, _ := DaProHeader(sysLogin)
	cmd[pos] = 9
	pos++
	for _, ch := range "anonymous" {
		cmd[pos] = byte(ch)
		pos++
	}
	cmd[pos] = 4
	pos++
	for _, ch := range "0000" {
		cmd[pos] = byte(ch)
		pos++
	}

	var len uint32
	len = uint32(pos - 8)
	cmd[4] = byte(len >> 24 & 0xff)
	cmd[5] = byte(len >> 16 & 0xff)
	cmd[6] = byte(len >> 8 & 0xff)
	cmd[7] = byte(len & 0xff)

	var cmd2 []byte
	cmd2 = make([]byte, pos)

	fmt.Printf("SysDbDisconnect message: ")
	for i := 0; i < pos; i++ {
		cmd2[i] = cmd[i]
		fmt.Printf(" %02x", cmd2[i])
	}
	fmt.Printf("\n")

	//fmt.Printf("message: %q\n", cmd2)
	return cmd2, nil
}

func write(conn net.Conn, cmd2 []byte) {
	numChar, err := conn.Write(cmd2)
	if err != nil {
		fmt.Printf("write: Client Write failed: %s", err.Error())
		os.Exit(1)
	}
	fmt.Printf("write:Client Wrote %d characters\n", numChar)
}

func read(conn net.Conn) ([]byte, int, error) {
	cmdLine := make([]byte, (1024 * 4))
	readBuf := bufio.NewReaderSize(conn, len(cmdLine))

	n, err := readBuf.Read(cmdLine)

	if err != nil {
		fmt.Printf("Client Read Error: %s\n", err.Error())
		return nil, 0, err
	}
	if n == 0 {
		fmt.Printf("Read 0 bytes\n")
	}

	fmt.Printf("read message[%d]: ", n)
	for i := 0; i < n; i++ {
		fmt.Printf(" %02x", cmdLine[i])
	}
	fmt.Printf("\n")

	//	fmt.Printf("Client Recieved: %q\n\n", cmdLine[:n])

	return cmdLine, n, err
}

/*

params() {
   assetCount 16
   for (assetCount) {
     assetRecId 32
   }
   assetTagCount 16
   for (assetTagCount) {
      tag 16
   }
}
*/
func AssetRecAssetListDetailsGet() ([]byte, error) {
	cmd, pos, _ := DaProHeader(assetRecAssetListDetailsGet)
	//asset count
	cmd[pos+0] = byte(1 >> 8 & 0xff)
	cmd[pos+1] = byte(1 & 0xff)
	pos += 2

	cmd[pos+0] = byte(0 >> 24 & 0xff)
	cmd[pos+1] = byte(0 >> 16 & 0xff)
	cmd[pos+2] = byte(0 >> 8 & 0xff)
	cmd[pos+3] = byte(0 & 0xff)
	pos += 4
	/*
		6	title
		7	timeStart
		8	timeEnd
	*/
	//asset tag count
	cmd[pos+0] = byte(1 >> 8 & 0xff)
	cmd[pos+1] = byte(1 & 0xff)
	pos += 2

	cmd[pos+0] = byte(6 >> 8 & 0xff)
	cmd[pos+1] = byte(6 & 0xff)
	pos += 2
	//	cmd[pos+0] = byte(7 >> 8 & 0xff)
	//	cmd[pos+1] = byte(7 & 0xff)
	//	pos += 2
	//	cmd[pos+0] = byte(7 >> 8 & 0xff)
	//	cmd[pos+1] = byte(7 & 0xff)
	//	pos += 2

	var len uint32
	len = uint32(pos - 8)
	cmd[4] = byte(len >> 24 & 0xff)
	cmd[5] = byte(len >> 16 & 0xff)
	cmd[6] = byte(len >> 8 & 0xff)
	cmd[7] = byte(len & 0xff)

	var cmd2 []byte
	cmd2 = make([]byte, pos)

	fmt.Printf("AssetRecAssetListDetailsGet message: ")
	for i := 0; i < pos; i++ {
		cmd2[i] = cmd[i]
		fmt.Printf(" %02x", cmd2[i])
	}
	fmt.Printf("\n")

	//fmt.Printf("message: %q\n", cmd2)
	return cmd2, nil
}
