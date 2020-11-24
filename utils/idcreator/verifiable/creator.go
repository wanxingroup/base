package verifiable

import (
	"sync"
	"time"
)

// 参考 Snowflake 和 Sonyflake 项目进行开发
// 生成出来的 ID 更适用于有时间销毁周期的记录进行使用
// 也就是在一定周期之后，允许再次重复使用的 ID
// ID 的总长度为48 bit

const BitLengthTime = 32    // 时间存放长度
const BitLengthSequence = 8 // 顺序号存放长度
const BitLengthVerify = 8   // 验证码存放长度

// 设置配置信息
type Settings struct {
	// 其实时间
	StartTime time.Time
	// 密钥
	SecretKey string
}

type IDCreator struct {
	signature   signature
	mutex       sync.Mutex
	startTime   int64
	elapsedTime int64
	sequence    uint16
}

func NewIDCreator(settings Settings) *IDCreator {

	creator := &IDCreator{
		sequence:  uint16(1<<BitLengthSequence - 1),
		signature: newSignature(settings.SecretKey),
	}

	if settings.StartTime.IsZero() {
		creator.startTime = toTime(time.Date(2014, 9, 1, 0, 0, 0, 0, time.UTC))
	} else {
		creator.startTime = toTime(settings.StartTime)
	}

	return creator
}

func (creator *IDCreator) NextID() uint64 {

	creator.mutex.Lock()
	defer creator.mutex.Unlock()

	creator.buildSequence()

	return creator.toID()
}

func (creator *IDCreator) buildSequence() {

	const maskSequence = uint16(1<<BitLengthSequence - 1)

	current := currentElapsedTime(creator.startTime)
	if creator.elapsedTime < current {
		creator.elapsedTime = current
		creator.sequence = 0
	} else { // creator.elapsedTime >= current
		creator.sequence = (creator.sequence + 1) & maskSequence
		if creator.sequence == 0 {
			creator.elapsedTime++
			overtime := creator.elapsedTime - current
			time.Sleep(sleepTime(overtime))
		}
	}
}

func (creator *IDCreator) toID() uint64 {

	elapsedTime := creator.elapsedTime % (1 << BitLengthTime)

	id := uint64(elapsedTime)<<(BitLengthSequence+BitLengthVerify) |
		uint64(creator.sequence)<<BitLengthVerify
	id |= creator.signature.sign(id)

	return id
}

type signature struct {
	block uint8
}

func initBlock(secretKey string) uint8 {

	block := uint8(0)
	for _, c := range secretKey {

		r := rune(c)
		block ^= uint8(r) ^ uint8(r>>8) ^ uint8(r>>16) ^ uint8(r>>24)
	}

	return block
}

func newSignature(secretKey string) signature {

	return signature{block: initBlock(secretKey)}
}

const maskVerify = uint64(1<<BitLengthVerify - 1)

func (s signature) sign(id uint64) uint64 {

	id = id & (^maskVerify)

	verify := s.block
	for id > 0 {
		verify ^= uint8(id)
		id >>= 8
	}

	return uint64(verify)
}

const timeUnit = 1e7 // 单位：纳秒，表示10毫秒

func toTime(t time.Time) int64 {
	return t.UTC().UnixNano() / timeUnit
}

func currentElapsedTime(startTime int64) int64 {
	return toTime(time.Now()) - startTime
}

func sleepTime(overtime int64) time.Duration {
	return time.Duration(overtime)*10*time.Millisecond -
		time.Duration(time.Now().UTC().UnixNano()%timeUnit)*time.Nanosecond
}

func Verify(id uint64, secretKey string) bool {
	return newSignature(secretKey).sign(id) == (id & maskVerify)
}
