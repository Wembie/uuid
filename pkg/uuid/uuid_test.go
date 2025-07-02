package uuid

import (
    "encoding/json"
    "strings"
    "testing"
    
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
    uuid := New()
    assert.NotEqual(t, Nil, uuid)
    assert.Equal(t, VersionRandom, uuid.Version())
    assert.Equal(t, VariantRFC4122, uuid.Variant())
}

func TestNewV4(t *testing.T) {
    uuid, err := NewV4()
    require.NoError(t, err)
    assert.NotEqual(t, Nil, uuid)
    assert.Equal(t, VersionRandom, uuid.Version())
}

func TestNewV1(t *testing.T) {
    uuid, err := NewV1()
    require.NoError(t, err)
    assert.NotEqual(t, Nil, uuid)
    assert.Equal(t, VersionTimeBased, uuid.Version())
}

func TestParse(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        wantErr bool
    }{
        {
            name:    "valid UUID with hyphens",
            input:   "550e8400-e29b-41d4-a716-446655440000",
            wantErr: false,
        },
        {
            name:    "valid UUID without hyphens",
            input:   "550e8400e29b41d4a716446655440000",
            wantErr: false,
        },
        {
            name:    "valid UUID with braces",
            input:   "{550e8400-e29b-41d4-a716-446655440000}",
            wantErr: false,
        },
        {
            name:    "invalid UUID length",
            input:   "550e8400-e29b-41d4-a716",
            wantErr: true,
        },
        {
            name:    "invalid UUID characters",
            input:   "550e8400-e29b-41d4-a716-44665544000g",
            wantErr: true,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            uuid, err := Parse(tt.input)
            if tt.wantErr {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
                assert.NotEqual(t, Nil, uuid)
            }
        })
    }
}

func TestUUIDString(t *testing.T) {
    uuid := New()
    s := uuid.String()
    
    // Check format: 8-4-4-4-12
    parts := strings.Split(s, "-")
    assert.Len(t, parts, 5)
    assert.Len(t, parts[0], 8)
    assert.Len(t, parts[1], 4)
    assert.Len(t, parts[2], 4)
    assert.Len(t, parts[3], 4)
    assert.Len(t, parts[4], 12)
}

func TestUUIDEqual(t *testing.T) {
    uuid1 := New()
    uuid2 := New()
    uuid3 := uuid1
    
    assert.False(t, uuid1.Equal(uuid2))
    assert.True(t, uuid1.Equal(uuid3))
}

func TestUUIDIsNil(t *testing.T) {
    assert.True(t, Nil.IsNil())
    assert.False(t, New().IsNil())
}

func TestUUIDJSON(t *testing.T) {
    uuid := New()
    
    // Marshal
    data, err := json.Marshal(uuid)
    require.NoError(t, err)
    
    // Unmarshal
    var unmarshaled UUID
    err = json.Unmarshal(data, &unmarshaled)
    require.NoError(t, err)
    
    assert.True(t, uuid.Equal(unmarshaled))
}

func TestGenerator(t *testing.T) {
    gen := NewGenerator(VersionRandom)
    assert.Equal(t, VersionRandom, gen.Version())
    
    uuid, err := gen.Generate()
    require.NoError(t, err)
    assert.Equal(t, VersionRandom, uuid.Version())
}

func TestMust(t *testing.T) {
    uuid := Must(NewV4())
    assert.NotEqual(t, Nil, uuid)
    
    assert.Panics(t, func() {
        Must(Nil, assert.AnError)
    })
}

func BenchmarkNew(b *testing.B) {
    for i := 0; i < b.N; i++ {
        New()
    }
}

func BenchmarkParse(b *testing.B) {
    uuid := New()
    s := uuid.String()
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        Parse(s)
    }
}