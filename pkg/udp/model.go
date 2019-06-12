package udp

import (
	"golang.org/x/text/encoding/unicode/utf32"
)

type SessionType uint8

func (s SessionType) String() string {
	switch s {
	case 0:
		return "Booking"
	case 1:
		return "Practice"
	case 2:
		return "Qualifying"
	case 3:
		return "Race"
	default:
		return "Unknown"
	}
}

type Event uint8

const (
	// Receive
	EventCollisionWithCar Event = 10
	EventCollisionWithEnv Event = 11
	EventNewSession       Event = 50
	EventNewConnection    Event = 51
	EventConnectionClosed Event = 52
	EventCarUpdate        Event = 53
	EventCarInfo          Event = 54
	EventEndSession       Event = 55
	EventVersion          Event = 56
	EventChat             Event = 57
	EventClientLoaded     Event = 58
	EventSessionInfo      Event = 59
	EventError            Event = 60
	EventLapCompleted     Event = 73
	EventClientEvent      Event = 130

	// Send
	EventRealtimeposInterval Event = 200
	EventGetCarInfo          Event = 201
	EventSendChat            Event = 202
	EventBroadcastChat       Event = 203
	EventGetSessionInfo      Event = 204
	EventSetSessionInfo      Event = 205
	EventKickUser            Event = 206
	EventNextSession         Event = 207
	EventRestartSession      Event = 208
	EventAdminCommand        Event = 209

	SessionTypeRace       SessionType = 3
	SessionTypeQualifying SessionType = 2
	SessionTypePractice   SessionType = 1
	SessionTypeBooking    SessionType = 0
)

type Message interface {
	Event() Event
}

type ServerError struct {
	error
}

func (ServerError) Event() Event {
	return EventError
}

type CarID uint8
type DriverGUID string

type lapCompletedInternal struct {
	CarID     CarID  `json:"CarID"`
	LapTime   uint32 `json:"LapTime"`
	Cuts      uint8  `json:"Cuts"`
	CarsCount uint8  `json:"CarsCount"`
}

func (LapCompleted) Event() Event {
	return EventLapCompleted
}

type LapCompleted struct {
	CarID     CarID  `json:"CarID"`
	LapTime   uint32 `json:"LapTime"`
	Cuts      uint8  `json:"Cuts"`
	CarsCount uint8  `json:"CarsCount"`

	Cars []*LapCompletedCar `json:"Cars"`
}

type LapCompletedCar struct {
	CarID     CarID  `json:"CarID"`
	LapTime   uint32 `json:"LapTime"`
	Laps      uint16 `json:"Laps"`
	Completed uint8  `json:"Completed"`
}

type Vec struct {
	X float32 `json:"X"`
	Y float32 `json:"Y"`
	Z float32 `json:"Z"`
}

type CollisionWithCar struct {
	CarID       CarID   `json:"CarID"`
	OtherCarID  CarID   `json:"OtherCarID"`
	ImpactSpeed float32 `json:"ImpactSpeed"`
	WorldPos    Vec     `json:"WorldPos"`
	RelPos      Vec     `json:"RelPos"`
}

func (CollisionWithCar) Event() Event {
	return EventCollisionWithCar
}

type CollisionWithEnvironment struct {
	CarID       CarID   `json:"CarID"`
	ImpactSpeed float32 `json:"ImpactSpeed"`
	WorldPos    Vec     `json:"WorldPos"`
	RelPos      Vec     `json:"RelPos"`
}

func (CollisionWithEnvironment) Event() Event {
	return EventCollisionWithEnv
}

type SessionCarInfo struct {
	CarID      CarID      `json:"CarID"`
	DriverName string     `json:"DriverName"`
	DriverGUID DriverGUID `json:"DriverGUID"`
	CarModel   string     `json:"CarModel"`
	CarSkin    string     `json:"CarSkin"`

	EventType Event `json:"EventType"`
}

func (s SessionCarInfo) Event() Event {
	return s.EventType
}

type Chat struct {
	CarID   CarID  `json:"CarID"`
	Message string `json:"Message"`
}

func (Chat) Event() Event {
	return EventChat
}

type CarInfo struct {
	CarID       CarID      `json:"CarID"`
	IsConnected bool       `json:"IsConnected"`
	CarModel    string     `json:"CarModel"`
	CarSkin     string     `json:"CarSkin"`
	DriverName  string     `json:"DriverName"`
	DriverTeam  string     `json:"DriverTeam"`
	DriverGUID  DriverGUID `json:"DriverGUID"`
}

func (CarInfo) Event() Event {
	return EventCarInfo
}

type CarUpdate struct {
	CarID               CarID   `json:"CarID"`
	Pos                 Vec     `json:"Pos"`
	Velocity            Vec     `json:"Velocity"`
	Gear                uint8   `json:"Gear"`
	EngineRPM           uint16  `json:"EngineRPM"`
	NormalisedSplinePos float32 `json:"NormalisedSplinePos"`
}

func (CarUpdate) Event() Event {
	return EventCarUpdate
}

type EndSession string

func (EndSession) Event() Event {
	return EventEndSession
}

type Version uint8

func (Version) Event() Event {
	return EventVersion
}

type ClientLoaded CarID

func (ClientLoaded) Event() Event {
	return EventClientLoaded
}

type SessionInfo struct {
	Version             uint8       `json:"Version"`
	SessionIndex        uint8       `json:"SessionIndex"`
	CurrentSessionIndex uint8       `json:"CurrentSessionIndex"`
	SessionCount        uint8       `json:"SessionCount"`
	ServerName          string      `json:"ServerName"`
	Track               string      `json:"Track"`
	TrackConfig         string      `json:"TrackConfig"`
	Name                string      `json:"Name"`
	SessionType         SessionType `json:"SessionType"`
	Time                uint16      `json:"Time"`
	Laps                uint16      `json:"Laps"`
	WaitTime            uint16      `json:"WaitTime"`
	AmbientTemp         uint8       `json:"AmbientTemp"`
	RoadTemp            uint8       `json:"RoadTemp"`
	WeatherGraphics     string      `json:"WeatherGraphics"`
	ElapsedMilliseconds int32       `json:"ElapsedMilliseconds"`

	EventType Event `json:"EventType"`
}

func (s SessionInfo) Event() Event {
	return s.EventType
}

type GetSessionInfo struct {
}

func (GetSessionInfo) Event() Event {
	return EventGetSessionInfo
}

type EnableRealtimePosInterval struct {
	Type     uint8
	Interval uint16
}

func (EnableRealtimePosInterval) Event() Event {
	return EventRealtimeposInterval
}

func NewEnableRealtimePosInterval(interval uint16) EnableRealtimePosInterval {
	return EnableRealtimePosInterval{
		Type:     uint8(EventRealtimeposInterval),
		Interval: interval,
	}
}

type SendChat struct {
	EventType    uint8
	CarID        uint8
	Len          uint8
	UTF32Encoded []byte
}

func (SendChat) Event() Event {
	return EventSendChat
}

func NewSendChat(carID CarID, data string) (*SendChat, error) {
	strlen := len(data)

	encoded, err := utf32.UTF32(utf32.LittleEndian, utf32.IgnoreBOM).NewEncoder().Bytes([]byte(data))

	if err != nil {
		return nil, err
	}

	return &SendChat{
		EventType:    uint8(EventSendChat),
		CarID:        uint8(carID),
		Len:          uint8(strlen),
		UTF32Encoded: encoded,
	}, nil
}
