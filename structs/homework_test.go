package main

import (
	"math"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

type Option func(*GamePerson)

func WithName(name string) func(*GamePerson) {
	return func(person *GamePerson) {
		person.name = name
	}
}

func WithCoordinates(x, y, z int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.coordinates = [3]int32{int32(x), int32(y), int32(z)}
	}
}

func WithGold(gold int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.gold = uint32(gold)
	}
}

func WithMana(mana int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.mana = uint16(mana)
	}
}

func WithHealth(health int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.health = uint16(health)
	}
}

func WithRespect(respect int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.rselOptions = person.rselOptions | (uint16(respect) << 12)
	}
}

func WithStrength(strength int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.rselOptions = person.rselOptions | (uint16(strength) << 8)
	}
}

func WithExperience(experience int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.rselOptions = person.rselOptions | (uint16(experience) << 4)
	}
}

func WithLevel(level int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.rselOptions = person.rselOptions | (uint16(level))
	}
}

func WithHouse() func(*GamePerson) {
	return func(person *GamePerson) {
		person.hgftOptions |= 1 << 0
	}
}

func WithGun() func(*GamePerson) {
	return func(person *GamePerson) {
		person.hgftOptions |= 1 << 1
	}
}

func WithFamily() func(*GamePerson) {
	return func(person *GamePerson) {
		person.hgftOptions |= 1 << 2
	}
}

func WithType(personType int) func(*GamePerson) {
	return func(person *GamePerson) {
		switch personType {
		case 0:
			person.hgftOptions |= 1 << 3
		case 1:
			person.hgftOptions |= 1 << 4
		case 2:
			person.hgftOptions |= 1 << 5
		}
	}
}

const (
	BuilderGamePersonType = iota
	BlacksmithGamePersonType
	WarriorGamePersonType
)

type byteOptions = byte

/*
	    withHouse  						// 0000 0001 (1)
		withGun    						// 0000 0010 (2)
		withFamily 						// 0000 0100 (4)
		withTypeWarrior   				// 0000 1000 (8)
		withTypeBuilder					// 0001 0000 (16)
		withTypeBlacksmith				// 0010 0000 (32)
		резерв							// 0100 0000 (64)
		резерв							// 1000 0000 (128)
*/

type GamePerson struct {
	coordinates [3]int32    //координаты X Y Z каждая от [2 000 000..2 000 000]
	gold        uint32      //золото [0..2 000 000 000]
	mana        uint16      //мана  [0..1000]
	health      uint16      //здоровье  [0..1000]
	rselOptions uint16      //уважение, сила, опыт, уровень все значения [0..10] в uint16 по 4 бита
	hgftOptions byteOptions //опции дом, оружие, семья, тип(воин/строитель/кузнец), bitmask
	name        string      //имя от 0-42 симловов латиницей
}

func NewGamePerson(options ...Option) GamePerson {
	gamePerson := GamePerson{}
	for _, option := range options {
		option(&gamePerson)
	}
	return gamePerson
}

func (p *GamePerson) Name() string {
	return p.name
}

func (p *GamePerson) X() int {
	return int(p.coordinates[0])
}

func (p *GamePerson) Y() int {
	return int(p.coordinates[1])
}

func (p *GamePerson) Z() int {
	return int(p.coordinates[2])
}

func (p *GamePerson) Gold() int {
	return int(p.gold)
}

func (p *GamePerson) Mana() int {
	return int(p.mana)
}

func (p *GamePerson) Health() int {
	return int(p.health)
}

func (p *GamePerson) Respect() int {
	return int((p.rselOptions >> 12) & 0xF)
}

func (p *GamePerson) Strength() int {
	return int((p.rselOptions >> 8) & 0xF)
}

func (p *GamePerson) Experience() int {
	return int((p.rselOptions >> 4) & 0xF)
}

func (p *GamePerson) Level() int {
	return int(p.rselOptions & 0xF)
}

func (p *GamePerson) HasHouse() bool {
	value := p.hgftOptions&(1<<0) != 0
	return value
}

func (p *GamePerson) HasGun() bool {
	value := p.hgftOptions&(1<<1) != 0
	return value
}

func (p *GamePerson) HasFamilty() bool {
	value := p.hgftOptions&(1<<2) != 0
	return value
}

func (p *GamePerson) Type() int {
	if p.hgftOptions&(1<<3) != 0 {
		return 0
	} else {
		if p.hgftOptions&(1<<4) != 0 {
			return 1
		} else {
			if p.hgftOptions&(1<<5) != 0 {
				return 2
			}
		}
	}
	return 0
}

func TestGamePerson(t *testing.T) {
	assert.LessOrEqual(t, unsafe.Sizeof(GamePerson{}), uintptr(64))

	const x, y, z = math.MinInt32, math.MaxInt32, 0
	const name = "aaaaaaaaaaaaa_bbbbbbbbbbbbb_cccccccccccccc"
	const personType = BuilderGamePersonType
	const gold = math.MaxInt32
	const mana = 1000
	const health = 1000
	const respect = 10
	const strength = 10
	const experience = 10
	const level = 10

	options := []Option{
		WithName(name),
		WithCoordinates(x, y, z),
		WithGold(gold),
		WithMana(mana),
		WithHealth(health),
		WithRespect(respect),
		WithStrength(strength),
		WithExperience(experience),
		WithLevel(level),
		WithHouse(),
		WithFamily(),
		WithType(personType),
	}

	person := NewGamePerson(options...)
	assert.Equal(t, name, person.Name())
	assert.Equal(t, x, person.X())
	assert.Equal(t, y, person.Y())
	assert.Equal(t, z, person.Z())
	assert.Equal(t, gold, person.Gold())
	assert.Equal(t, mana, person.Mana())
	assert.Equal(t, health, person.Health())
	assert.Equal(t, respect, person.Respect())
	assert.Equal(t, strength, person.Strength())
	assert.Equal(t, experience, person.Experience())
	assert.Equal(t, level, person.Level())
	assert.True(t, person.HasHouse())
	assert.True(t, person.HasFamilty())
	assert.False(t, person.HasGun())
	assert.Equal(t, personType, person.Type())
}
