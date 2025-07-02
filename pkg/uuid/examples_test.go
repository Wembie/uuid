package uuid_test

import (
	"fmt"
	"log"

	"github.com/Wembie/uuid/pkg/uuid"
)

// ExampleNew demonstrates basic UUID generation
func ExampleNew() {
    id := uuid.New()
    fmt.Printf("Generated UUID: %s\n", id)
    fmt.Printf("Version: %d\n", id.Version())
    fmt.Printf("Variant: %d\n", id.Variant())
    // Output will vary, but format will be: xxxxxxxx-xxxx-4xxx-xxxx-xxxxxxxxxxxx
}

// ExampleParse demonstrates UUID parsing
func ExampleParse() {
    id, err := uuid.Parse("550e8400-e29b-41d4-a716-446655440000")
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Parsed UUID: %s\n", id)
    fmt.Printf("Is nil: %v\n", id.IsNil())
    // Output:
    // Parsed UUID: 550e8400-e29b-41d4-a716-446655440000
    // Is nil: false
}

// ExampleGenerator demonstrates using custom generators
func ExampleGenerator() {
    gen := uuid.NewGenerator(uuid.VersionRandom)
    
    id, err := gen.Generate()
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Generated UUID: %s\n", id)
    fmt.Printf("Generator version: %d\n", gen.Version())
    // Output will vary
}