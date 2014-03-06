// Package middlecoin is a package which contains structures for decoding
// Middlecoin JSON data in Go.
package middlecoin

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

// ReportFloat64 is a simple type on which I defined Unmarshal.
type ReportFloat64 float64

// ReportTime is a custom defined type so that I could easily unmarshal
// the json into the time according to a format.
type ReportTime time.Time

// AddressMap maps bitcoin addresses from middlecoin to reports on their
// progress.
type AddressMap map[string]*AddressReport

func (r *ReportTime) String() string {
	return (*time.Time)(r).String()
}

func (m *AddressMap) UnmarshalJSON(data []byte) error {
	*m = make(AddressMap)
	var elements []json.RawMessage
	json.Unmarshal(data, &elements)
	for _, raw := range elements {
		element := new(ReportElement)
		json.Unmarshal(raw, element)
		(*m)[element.Address] = element.Report
	}
	return nil
}

func (r *ReportElement) UnmarshalJSON(data []byte) error {
	var a []json.RawMessage
	json.Unmarshal(data, &a)
	json.Unmarshal(a[0], &r.Address)
	json.Unmarshal(a[1], &r.Report)
	return nil
}

func (r *ReportFloat64) UnmarshalJSON(data []byte) error {
	var s string
	json.Unmarshal(data, &s)
	f, err := strconv.ParseFloat(s, 64)
	*r = ReportFloat64(f)
	return err
}

func (r *ReportTime) UnmarshalJSON(data []byte) error {
	const layout = `2006-01-02 15:04:05`
	var s string
	json.Unmarshal(data, &s)
	t, err := time.Parse(layout, s)
	*r = ReportTime(t)
	return err
}

// OverviewReport is the top level report on middlecoin and the totals
// for all of the miners.
type OverviewReport struct {
	TotalPaidOut                     ReportFloat64
	TotalRejectedMegahashesPerSecond ReportFloat64
	TotalImmatureBalance             ReportFloat64
	TotalMegahashesPerSecond         ReportFloat64
	TotalBalance                     ReportFloat64
	Time                             ReportTime
	Report                           AddressMap
}

// ReportElement is a simple pair that contains a bitcoin address for a
// miner and the report on its progress.
type ReportElement struct {
	Address string
	Report  *AddressReport
}

// AddressReport the report of what the middlecoin miner has been up to
// for the last hour.
type AddressReport struct {
	LastHourShares              int
	ImmatureBalance             ReportFloat64
	LastHourRejectedShares      int
	PaidOut                     ReportFloat64
	UnexchangedBalance          ReportFloat64
	MegahashesPerSecond         ReportFloat64
	BitcoinBalance              ReportFloat64
	RejectedMegahashesPerSecond ReportFloat64
}

// Add is a simple addition operator for AddressReports for easier
// aggregation.
func (r *AddressReport) Add(o *AddressReport) *AddressReport {
	r.LastHourShares += o.LastHourShares
	r.ImmatureBalance += o.ImmatureBalance
	r.LastHourRejectedShares += o.LastHourRejectedShares
	r.PaidOut += o.PaidOut
	r.UnexchangedBalance += o.UnexchangedBalance
	r.MegahashesPerSecond += o.MegahashesPerSecond
	r.BitcoinBalance += o.BitcoinBalance
	r.RejectedMegahashesPerSecond += o.RejectedMegahashesPerSecond
	return r
}

func (r *AddressReport) Profit() float64 {
	return float64(r.BitcoinBalance + r.ImmatureBalance + r.PaidOut + r.UnexchangedBalance)
}

func (r *AddressReport) String() string {
	profit := r.Profit()
	bitcointousd := 650.0
	usdprofit := float64(profit) * bitcointousd
	return fmt.Sprintf(
		"Last Hour Shares: %d\n"+
			"Immature Balance: %g\n"+
			"Last Hour Rejected Shares: %d\n"+
			"Paid Out: %g\n"+
			"Unexchanged Balance: %g\n"+
			"Megahashes Per Second: %g\n"+
			"Bitcoin Balance: %g\n"+
			"Rejected Megahashes Per Second: %g\n"+
			"\nTotal Profit: %g\n"+
			"Total Profit(@$%g): %g\n",
		r.LastHourShares,
		r.ImmatureBalance,
		r.LastHourRejectedShares,
		r.PaidOut,
		r.UnexchangedBalance,
		r.MegahashesPerSecond,
		r.BitcoinBalance,
		r.RejectedMegahashesPerSecond,
		profit,
		bitcointousd, usdprofit,
	)
}
