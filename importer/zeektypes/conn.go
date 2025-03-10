package zeektypes

import (
	"errors"
	"strconv"
	"strings"
	"time"

	jsoniter "github.com/json-iterator/go"
)

// EntryTypeConn should be matched against zeekFile.EntryType()
// before using OpenZeekReader[ZeekConn](fs, zeekFile) to read from the file.
const EntryTypeConn = "conn"

// EntryTypeOpenConn should be matched against zeekFile.EntryType()
// before using OpenZeekReader[ZeekConn](fs, zeekFile) to read from the file.
// Zeek logs written to the open_conn file follow the same format normal conn logs.
const EntryTypeOpenConn = "open_conn"

// type Timestamp time.Time
type Timestamp int64

var ErrInvalidZeekTimestamp = errors.New("invalid zeek timestamp")

// Conn provides a data structure for zeek's connection data
type Conn struct {
	// TimeStamp of this connection
	TimeStamp Timestamp `zeek:"ts" zeektype:"time" json:"ts"`
	// UID is the Unique Id for this connection (generated by zeek)
	UID string `zeek:"uid" zeektype:"string" json:"uid"`
	// Source is the source address for this connection
	Source string `zeek:"id.orig_h" zeektype:"addr" json:"id.orig_h"`
	// SourcePort is the source port of this connection
	SourcePort int `zeek:"id.orig_p" zeektype:"port" json:"id.orig_p"`
	// Destination is the destination of the connection
	Destination string `zeek:"id.resp_h" zeektype:"addr" json:"id.resp_h"`
	// DestinationPort is the port at the destination host
	DestinationPort int `zeek:"id.resp_p" zeektype:"port" json:"id.resp_p"`
	// Proto is the string protocol identifier for this connection
	Proto string `zeek:"proto" zeektype:"enum" json:"proto"`
	// Service describes the service of this connection if there was one
	Service string `zeek:"service" zeektype:"string" json:"service"`
	// Duration is the floating point representation of connection length
	Duration float64 `zeek:"duration" zeektype:"interval" json:"duration"`
	// OrigBytes is the byte count coming from the origin
	OrigBytes int64 `zeek:"orig_bytes" zeektype:"count" json:"orig_bytes"`
	// RespBytes is the byte count coming in on response
	RespBytes int64 `zeek:"resp_bytes" zeektype:"count" json:"resp_bytes"`
	// ConnState has data describing the state of a connection
	ConnState string `zeek:"conn_state" zeektype:"string" json:"conn_state"`
	// LocalOrigin denotes that the connection originated locally
	LocalOrigin bool `zeek:"local_orig" zeektype:"bool" json:"local_orig"`
	// LocalResponse denote that the connection responded locally
	LocalResponse bool `zeek:"local_resp" zeektype:"bool" json:"local_resp"`
	// MissedBytes keeps a count of bytes missed
	MissedBytes int64 `zeek:"missed_bytes" zeektype:"count" json:"missed_bytes"`
	// History is a string containing historical information
	History string `zeek:"history" zeektype:"string" json:"history"`
	// OrigPkts is a count of origin packets
	OrigPackets int64 `zeek:"orig_pkts" zeektype:"count" json:"orig_pkts"`
	// OrigIpBytes is another origin data count
	OrigIPBytes int64 `zeek:"orig_ip_bytes" zeektype:"count" json:"orig_ip_bytes"`
	// RespPackets counts response packets
	RespPackets int64 `zeek:"resp_pkts" zeektype:"count" json:"resp_pkts"`
	// RespIpBytes gives the bytecount of response data
	RespIPBytes int64 `zeek:"resp_ip_bytes" zeektype:"count" json:"resp_ip_bytes"`
	// TunnelParents lists tunnel parents
	TunnelParents []string `zeek:"tunnel_parents" zeektype:"set[string]" json:"tunnel_parents"`
	// AgentHostname names which sensor recorded this event. Only set when combining logs from multiple sensors.
	AgentHostname string `zeek:"agent_hostname" zeektype:"string" json:"agent_hostname"`
	// AgentUUID identifies which sensor recorded this event. Only set when combining logs from multiple sensors.
	AgentUUID string `zeek:"agent_uuid" zeektype:"string" json:"agent_uuid"`
	// Path of log file containing this record
	LogPath string
}

func (c *Conn) SetLogPath(path string) { c.LogPath = path }

// Unmarshals JSON timestamps
func (ts *Timestamp) UnmarshalJSON(data []byte) error {
	var t interface{}
	if err := jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal(data, &t); err != nil {
		return err
	}

	switch input := t.(type) {
	// all number types are assumed to be in unix format, possibly with fractional seconds
	case int:
		tsVal, ok := t.(Timestamp)
		if !ok {
			return ErrInvalidZeekTimestamp
		}
		*ts = tsVal
	case int32:
		tsVal, ok := t.(Timestamp)
		if !ok {
			return ErrInvalidZeekTimestamp
		}
		*ts = tsVal
	case float32:
		tsVal, ok := t.(Timestamp)
		if !ok {
			return ErrInvalidZeekTimestamp
		}
		*ts = tsVal
	case int64:
		intVal, ok := t.(int64)
		if !ok {
			return ErrInvalidZeekTimestamp
		}
		*ts = Timestamp(intVal)
	case float64:
		floatVal, ok := t.(float64)
		if !ok {
			return ErrInvalidZeekTimestamp
		}
		*ts = Timestamp(floatVal)
	case string:
		// assumed to be in RFC8601 format, though other formats can be added as necessary
		// ex: 2019-11-13T09:00:01.932360Z
		t, err := time.Parse(time.RFC3339, input)
		if err != nil {
			// attempt to parse as a epoch
			tsVal, err := strconv.ParseFloat(strings.TrimSpace(input), 64)
			if err != nil {
				return errors.Join(ErrInvalidZeekTimestamp, err)
			}
			*ts = Timestamp(tsVal)
		}
		var unix Timestamp
		unix = Timestamp(t.UTC().Unix())
		*ts = unix
	default:
		return ErrInvalidZeekTimestamp
	}

	return nil
}
