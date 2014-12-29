package yelp

import (
	"errors"
	"fmt"
	"reflect"
)

type OptionProvider interface {
	GetParameters() (params map[string]string, err error)
}

/**
 * Top level search object used for doing searches.  You can define multiple
 * sets of options, and use them together.  Only one of LocationOptions,
 * CoordinateOptions, or BoundOptions can be used at the same time.
 */
type SearchOptions struct {
	GeneralOptions    *GeneralOptions    // standard general search options (filters, terms, etc)
	LocaleOptions     *LocaleOptions     // Results will be localized in the region format and language if supported.
	LocationOptions   *LocationOptions   // Use a location term and potentially coordinates to define the location
	CoordinateOptions *CoordinateOptions // Use coordinate options to define the location.
	BoundOptions      *BoundOptions      // Use bound options (an area) to define the location.
}

/**
 * Generate a map that contains the querystring parameters for
 * all of the defined options.
 */
func (o *SearchOptions) GetParameters() (params map[string]string, err error) {

	// ensure only one loc option provider is being used
	locOptionsCnt := 0
	if o.LocationOptions != nil {
		locOptionsCnt++
	}
	if o.CoordinateOptions != nil {
		locOptionsCnt++
	}
	if o.BoundOptions != nil {
		locOptionsCnt++
	}

	if locOptionsCnt == 0 {
		return params, errors.New("A single location search options type (Location, Coordinate, Bound) must be used.")
	}
	if locOptionsCnt > 1 {
		return params, errors.New("Only a single location search options type (Location, Coordinate, Bound) can be used at a time.")
	}
	fmt.Printf("There are %v location options defined\n", locOptionsCnt)

	// create an empty map of options
	params = make(map[string]string)

	// reflect over the properties in o, adding parameters to the global map
	val := reflect.ValueOf(o).Elem()
	for i := 0; i < val.NumField(); i++ {
		if !val.Field(i).IsNil() {
			o := val.Field(i).Interface().(OptionProvider)
			fieldParams, err := o.GetParameters()
			if err != nil {
				return params, err
			}
			for k, v := range fieldParams {
				params[k] = v
			}
		}
	}
	return params, nil
}
