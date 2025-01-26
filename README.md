# julianephdate

A Go library to convert between standard UTC time and the Julian Ephemeris Date (JED).
Julian Ephemeris Date (JED) is a time scale used in astronomical calculations, similar to Julian Date (JD) but accounting for terrestrial time offsets (TT, TAI, leap seconds, etc.).

This repository provides:
• A function `Date(time.Time) float64` that returns the Julian Ephemeris Date of a given standard library time.Time (UTC).
• A function `StdTime(float64) time.Time` that returns the UTC time.Time from a JED float64.

## Table of Contents
1. [Installation](#installation)
2. [Usage](#usage)
3. [Leap Second Considerations](#leap-second-considerations)
4. [Testing Instructions](#testing-instructions)
5. [Contributing](#contributing)

---

### Installation

To install:

    go get github.com/dtimm/julianephdate

Then, in your code:

    import "github.com/dtimm/julianephdate"

---

### Usage

• Date(t time.Time) -> float64

  Pass in a time.Time; get back the Julian Ephemeris Date (JED), which is based on TT (Terrestrial Time) and accounts for leap seconds via a lookup table.

• StdTime(jed float64) -> time.Time

  Converts the JED float64 back into a time.Time in UTC. Leap-second boundaries are approximated by a single iteration looking up the correct TAI−UTC offset from a table of known leap seconds.

#### Example

    package main

    import (
        "fmt"
        "time"

        "github.com/dtimm/julianephdate"
    )

    func main() {
        // Example: J2000 epoch
        t := time.Date(2000, 1, 1, 11, 58, 55, 816e6, time.UTC)
        jed := julianephdate.Date(t)
        fmt.Printf("Time: %v -> JED: %f\n", t, jed)

        backToTime := julianephdate.StdTime(jed)
        fmt.Printf("JED: %f -> back to Time: %v\n", jed, backToTime)
    }

---

### Leap Second Considerations

A table of historical leap seconds is built into the package (`leapSecondsTable`). It starts at the first leap second in 1972 through the last announced leap second. For historical or future dates beyond the scope of this table, the conversion will be approximate or incomplete, and you may need to update the table with any future leap second announcements.

---

### Testing Instructions

This project uses [Ginkgo](https://github.com/onsi/ginkgo) and [Gomega](https://github.com/onsi/gomega) for Behavior-Driven Development (BDD)-style tests.

Commands to run the tests:

1. Make sure you have Go installed.
2. Make sure you have the Ginkgo testing framework and Gomega installed:

       go get github.com/onsi/ginkgo/v2
       go get github.com/onsi/gomega

3. Clone this repository and navigate to the project root.
4. Run tests with either:

       go test ./...

   Or, if you have Ginkgo CLI installed:

       ginkgo ./...

You'll see output indicating whether all tests have passed.

---

### Contributing

Feel free to open an issue or create a pull request if you have suggestions, bug fixes, or new features. We welcome any improvements, especially updates to the leap-second table if new announcements are made.
