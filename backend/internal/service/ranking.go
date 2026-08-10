package service

import (
	"errors"
	"math"
	"time"

	"github.com/google/uuid"
)

// Ranking weights for Smart Dispatch Finder.
const (
	WeightAvailability = 0.50
	WeightProximity    = 0.30
	WeightSpecs        = 0.20
	TransportSpeedKmh  = 40.0
)

var (
	// ErrConflict indicates an overlapping assignment or maintenance window.
	ErrConflict = errors.New("conflict")
	// ErrNotFound indicates a missing entity.
	ErrNotFound = errors.New("not found")
	// ErrUnauthorized indicates bad credentials.
	ErrUnauthorized = errors.New("unauthorized")
	// ErrForbidden indicates insufficient role.
	ErrForbidden = errors.New("forbidden")
	// ErrValidation indicates invalid input.
	ErrValidation = errors.New("validation")
)

// IntervalsOverlap reports whether two half-open intervals overlap.
func IntervalsOverlap(aStart, aEnd, bStart, bEnd time.Time) bool {
	return aStart.Before(bEnd) && aEnd.After(bStart)
}

// AvailabilityScore scores how well availability matches the requested window.
func AvailabilityScore(available bool, availableFrom, requestedStart time.Time) float64 {
	if available {
		return 1.0
	}
	if !availableFrom.After(requestedStart) {
		return 0.8
	}
	delayHours := availableFrom.Sub(requestedStart).Hours()
	score := 1.0 / (1.0 + delayHours/24.0)
	if score < 0.05 {
		return 0.05
	}
	return score
}

// ProximityScore converts distance km into 0..1 score.
func ProximityScore(distanceKm float64) float64 {
	return 1.0 / (1.0 + distanceKm/10.0)
}

// SpecsScore compares requested min weight to asset weight.
func SpecsScore(assetWeightT *float64, minWeightT *float64) (score float64, underSpec bool) {
	if minWeightT == nil {
		return 1.0, false
	}
	if assetWeightT == nil {
		return 0.5, false
	}
	if *assetWeightT >= *minWeightT {
		return 1.0, false
	}
	ratio := *assetWeightT / *minWeightT
	if ratio < 0.1 {
		ratio = 0.1
	}
	return ratio, true
}

// RankScore combines availability, proximity, and specs.
func RankScore(availability, proximity, specs float64) float64 {
	return WeightAvailability*availability + WeightProximity*proximity + WeightSpecs*specs
}

// TransportMinutes estimates transport time from distance.
func TransportMinutes(distanceKm float64) float64 {
	return distanceKm / TransportSpeedKmh * 60.0
}

// HaversineKm returns great-circle distance in kilometers.
func HaversineKm(lat1, lng1, lat2, lng2 float64) float64 {
	const earthRadiusKm = 6371.0
	toRad := func(d float64) float64 { return d * math.Pi / 180 }
	dLat := toRad(lat2 - lat1)
	dLng := toRad(lng2 - lng1)
	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(toRad(lat1))*math.Cos(toRad(lat2))*math.Sin(dLng/2)*math.Sin(dLng/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return earthRadiusKm * c
}

// IsNilUUID reports whether id is the zero UUID.
func IsNilUUID(id uuid.UUID) bool {
	return id == uuid.Nil
}
