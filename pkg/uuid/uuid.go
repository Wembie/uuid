package uuid

import (
    "crypto/rand"
    "database/sql/driver"
    "encoding/hex"
    "encoding/json"
    "fmt"
    "strings"
    "time"
)

// UUID represents a UUID value
type UUID [16]byte

// Version represents the UUID version
type Version byte

// Variant represents the UUID variant
type Variant byte

const (
    // UUID versions
    VersionUnknown Version = iota
    VersionTimeBased
    VersionDCESecurity
    VersionNameBasedMD5
    VersionRandom
    VersionNameBasedSHA1
)

const (
    // UUID variants
    VariantNCS Variant = iota
    VariantRFC4122
    VariantMicrosoft
    VariantFuture
)

// Nil is the nil UUID
var Nil = UUID{}

// Generator interface for UUID generation strategies
type Generator interface {
    Generate() (UUID, error)
    Version() Version
}

// UUIDGenerator is the default UUID generator
type UUIDGenerator struct {
    version Version
}

// NewGenerator creates a new UUID generator for the specified version
func NewGenerator(version Version) Generator {
    return &UUIDGenerator{version: version}
}

// Generate creates a new UUID based on the generator's version
func (g *UUIDGenerator) Generate() (UUID, error) {
    switch g.version {
    case VersionRandom:
        return generateV4()
    case VersionTimeBased:
        return generateV1()
    default:
        return generateV4() // Default to V4
    }
}

// Version returns the generator's version
func (g *UUIDGenerator) Version() Version {
    return g.version
}

// New generates a new random UUID (Version 4)
func New() UUID {
    uuid, _ := generateV4()
    return uuid
}

// NewV4 generates a new random UUID (Version 4)
func NewV4() (UUID, error) {
    return generateV4()
}

// NewV1 generates a new time-based UUID (Version 1)
func NewV1() (UUID, error) {
    return generateV1()
}

// Must is a helper that wraps a UUID generation function and panics if error occurs
func Must(uuid UUID, err error) UUID {
    if err != nil {
        panic(err)
    }
    return uuid
}

// MustNew generates a new UUID and panics if error occurs
func MustNew() UUID {
    return Must(NewV4())
}

// Parse parses a string into a UUID
func Parse(s string) (UUID, error) {
    var uuid UUID
    
    // Remove hyphens and braces
    s = strings.ReplaceAll(s, "-", "")
    s = strings.ReplaceAll(s, "{", "")
    s = strings.ReplaceAll(s, "}", "")
    
    if len(s) != 32 {
        return uuid, fmt.Errorf("invalid UUID length: %d", len(s))
    }
    
    decoded, err := hex.DecodeString(s)
    if err != nil {
        return uuid, fmt.Errorf("invalid UUID format: %v", err)
    }
    
    copy(uuid[:], decoded)
    return uuid, nil
}

// MustParse parses a string into a UUID and panics if error occurs
func MustParse(s string) UUID {
    uuid, err := Parse(s)
    if err != nil {
        panic(err)
    }
    return uuid
}

// ParseBytes parses a byte slice into a UUID
func ParseBytes(b []byte) (UUID, error) {
    var uuid UUID
    if len(b) != 16 {
        return uuid, fmt.Errorf("invalid UUID byte length: %d", len(b))
    }
    copy(uuid[:], b)
    return uuid, nil
}

// FromString is an alias for Parse
func FromString(s string) (UUID, error) {
    return Parse(s)
}

// String returns the string representation of the UUID
func (u UUID) String() string {
    return fmt.Sprintf("%x-%x-%x-%x-%x",
        u[:4], u[4:6], u[6:8], u[8:10], u[10:])
}

// URN returns the RFC 2141 URN form of the UUID
func (u UUID) URN() string {
    return "urn:uuid:" + u.String()
}

// Bytes returns the UUID as a byte slice
func (u UUID) Bytes() []byte {
    return u[:]
}

// Version returns the version of the UUID
func (u UUID) Version() Version {
    return Version(u[6] >> 4)
}

// Variant returns the variant of the UUID
func (u UUID) Variant() Variant {
    switch {
    case (u[8] & 0x80) == 0x00:
        return VariantNCS
    case (u[8] & 0xc0) == 0x80:
        return VariantRFC4122
    case (u[8] & 0xe0) == 0xc0:
        return VariantMicrosoft
    default:
        return VariantFuture
    }
}

// IsNil returns true if the UUID is the nil UUID
func (u UUID) IsNil() bool {
    return u == Nil
}

// Equal returns true if u and other are equal
func (u UUID) Equal(other UUID) bool {
    return u == other
}

// Compare compares two UUIDs lexicographically
func (u UUID) Compare(other UUID) int {
    for i := 0; i < 16; i++ {
        if u[i] < other[i] {
            return -1
        }
        if u[i] > other[i] {
            return 1
        }
    }
    return 0
}

// MarshalJSON implements json.Marshaler
func (u UUID) MarshalJSON() ([]byte, error) {
    return json.Marshal(u.String())
}

// UnmarshalJSON implements json.Unmarshaler
func (u *UUID) UnmarshalJSON(data []byte) error {
    var s string
    if err := json.Unmarshal(data, &s); err != nil {
        return err
    }
    
    parsed, err := Parse(s)
    if err != nil {
        return err
    }
    
    *u = parsed
    return nil
}

// MarshalText implements encoding.TextMarshaler
func (u UUID) MarshalText() ([]byte, error) {
    return []byte(u.String()), nil
}

// UnmarshalText implements encoding.TextUnmarshaler
func (u *UUID) UnmarshalText(text []byte) error {
    parsed, err := Parse(string(text))
    if err != nil {
        return err
    }
    
    *u = parsed
    return nil
}

// Value implements driver.Valuer for database operations
func (u UUID) Value() (driver.Value, error) {
    return u.String(), nil
}

// Scan implements sql.Scanner for database operations
func (u *UUID) Scan(value interface{}) error {
    if value == nil {
        *u = Nil
        return nil
    }
    
    switch v := value.(type) {
    case string:
        parsed, err := Parse(v)
        if err != nil {
            return err
        }
        *u = parsed
    case []byte:
        if len(v) == 16 {
            copy(u[:], v)
        } else {
            parsed, err := Parse(string(v))
            if err != nil {
                return err
            }
            *u = parsed
        }
    default:
        return fmt.Errorf("cannot scan %T into UUID", value)
    }
    
    return nil
}

// Internal generation functions
func generateV4() (UUID, error) {
    var uuid UUID
    _, err := rand.Read(uuid[:])
    if err != nil {
        return uuid, err
    }
    
    // Set version (4) and variant bits
    uuid[6] = (uuid[6] & 0x0f) | 0x40 // Version 4
    uuid[8] = (uuid[8] & 0x3f) | 0x80 // Variant RFC4122
    
    return uuid, nil
}

func generateV1() (UUID, error) {
    var uuid UUID
    _, err := rand.Read(uuid[:])
    if err != nil {
        return uuid, err
    }
    
    // Simplified V1 generation (in real implementation, use proper timestamp and MAC)
    now := time.Now().UnixNano()
    
    // Time low
    uuid[0] = byte(now)
    uuid[1] = byte(now >> 8)
    uuid[2] = byte(now >> 16)
    uuid[3] = byte(now >> 24)
    
    // Time mid
    uuid[4] = byte(now >> 32)
    uuid[5] = byte(now >> 40)
    
    // Time high and version
    uuid[6] = byte(now>>48) & 0x0f
    uuid[6] |= 0x10 // Version 1
    
    // Clock sequence and variant
    uuid[8] = (uuid[8] & 0x3f) | 0x80 // Variant RFC4122
    
    return uuid, nil
}