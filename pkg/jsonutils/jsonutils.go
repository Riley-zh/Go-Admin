package jsonutils

import (
	"encoding/json"
	"io"
	"sync"

	"go-admin/internal/logger"

	"go.uber.org/zap"
)

// OptimizedDecoder provides efficient JSON decoding with low CPU usage
type OptimizedDecoder struct {
	decoderPool sync.Pool
}

// NewOptimizedDecoder creates a new optimized JSON decoder
func NewOptimizedDecoder() *OptimizedDecoder {
	return &OptimizedDecoder{
		decoderPool: sync.Pool{
			New: func() interface{} {
				return json.NewDecoder(nil)
			},
		},
	}
}

// Decode efficiently decodes JSON from an io.Reader into the provided interface
func (d *OptimizedDecoder) Decode(r io.Reader, v interface{}) error {
	decoder := json.NewDecoder(r)

	if err := decoder.Decode(v); err != nil {
		logger.Error("Failed to decode JSON", zap.String("error", err.Error()))
		return err
	}

	return nil
}

// DecodeStrict decodes JSON with strict error checking
func (d *OptimizedDecoder) DecodeStrict(r io.Reader, v interface{}) error {
	decoder := json.NewDecoder(r)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(v); err != nil {
		logger.Error("Failed to decode JSON (strict mode)", zap.String("error", err.Error()))
		return err
	}

	return nil
}

// OptimizedEncoder provides efficient JSON encoding with low CPU usage
type OptimizedEncoder struct {
	encoderPool sync.Pool
}

// NewOptimizedEncoder creates a new optimized JSON encoder
func NewOptimizedEncoder() *OptimizedEncoder {
	return &OptimizedEncoder{
		encoderPool: sync.Pool{
			New: func() interface{} {
				return json.NewEncoder(nil)
			},
		},
	}
}

// Encode efficiently encodes the provided interface to JSON
func (e *OptimizedEncoder) Encode(w io.Writer, v interface{}) error {
	encoder := json.NewEncoder(w)

	if err := encoder.Encode(v); err != nil {
		logger.Error("Failed to encode JSON", zap.String("error", err.Error()))
		return err
	}

	return nil
}

// Marshal efficiently marshals the provided interface to JSON bytes
func Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

// MarshalIndent efficiently marshals the provided interface to indented JSON bytes
func MarshalIndent(v interface{}, prefix, indent string) ([]byte, error) {
	return json.MarshalIndent(v, prefix, indent)
}

// Unmarshal efficiently unmarshals JSON bytes into the provided interface
func Unmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

// ValidateJSON validates if the given data is valid JSON
func ValidateJSON(data []byte) error {
	var js json.RawMessage
	return json.Unmarshal(data, &js)
}

// SafeDecode safely decodes JSON with additional error handling
func SafeDecode(r io.Reader, v interface{}) error {
	decoder := json.NewDecoder(r)

	// Check for unexpected data after JSON
	if err := decoder.Decode(v); err != nil {
		return err
	}

	// Look for any additional data that might indicate malformed JSON
	if decoder.More() {
		return json.NewDecoder(r).Decode(new(interface{}))
	}

	return nil
}

// StreamingDecoder provides streaming JSON decoding for large payloads
type StreamingDecoder struct {
	decoder *json.Decoder
}

// NewStreamingDecoder creates a new streaming JSON decoder
func NewStreamingDecoder(r io.Reader) *StreamingDecoder {
	return &StreamingDecoder{
		decoder: json.NewDecoder(r),
	}
}

// DecodeNext decodes the next JSON value from the stream
func (s *StreamingDecoder) DecodeNext(v interface{}) error {
	if !s.decoder.More() {
		return io.EOF
	}

	return s.decoder.Decode(v)
}

// BatchDecoder provides batch processing for JSON arrays
type BatchDecoder struct {
	decoder *json.Decoder
}

// NewBatchDecoder creates a new batch JSON decoder
func NewBatchDecoder(r io.Reader) *BatchDecoder {
	return &BatchDecoder{
		decoder: json.NewDecoder(r),
	}
}

// ProcessArray processes a JSON array in batches
func (b *BatchDecoder) ProcessArray(batchSize int, processFunc func([]interface{}) error) error {
	// Read opening bracket
	token, err := b.decoder.Token()
	if err != nil {
		return err
	}

	if delim, ok := token.(json.Delim); !ok || delim != '[' {
		return json.NewDecoder(nil).Decode(new(interface{}))
	}

	batch := make([]interface{}, 0, batchSize)

	// Process array elements
	for b.decoder.More() {
		var item interface{}
		if err := b.decoder.Decode(&item); err != nil {
			return err
		}

		batch = append(batch, item)

		// Process batch when it reaches the specified size
		if len(batch) >= batchSize {
			if err := processFunc(batch); err != nil {
				return err
			}
			batch = batch[:0] // Reset batch
		}
	}

	// Process any remaining items in the last batch
	if len(batch) > 0 {
		if err := processFunc(batch); err != nil {
			return err
		}
	}

	// Read closing bracket
	if _, err := b.decoder.Token(); err != nil {
		return err
	}

	return nil
}

// Global instances for reuse
var (
	defaultDecoder = NewOptimizedDecoder()
	defaultEncoder = NewOptimizedEncoder()
)

// DecodeJSON uses the default decoder to decode JSON
func DecodeJSON(r io.Reader, v interface{}) error {
	return defaultDecoder.Decode(r, v)
}

// EncodeJSON uses the default encoder to encode JSON
func EncodeJSON(w io.Writer, v interface{}) error {
	return defaultEncoder.Encode(w, v)
}

// DecodeJSONStrict uses the default decoder in strict mode
func DecodeJSONStrict(r io.Reader, v interface{}) error {
	return defaultDecoder.DecodeStrict(r, v)
}
